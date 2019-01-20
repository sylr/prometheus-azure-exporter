package metrics

import (
	"context"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-07-01/storage"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/azure"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools"
)

const (
	kiloBytes = 1000
	megaBytes = 1000 * 1000
	gigaBytes = 1000 * 1000 * 1000
)

var (
	storageAccountBlobSizeHistogram = newStorageAccountBlobSizeHistogram()
)

// -----------------------------------------------------------------------------

func newStorageAccountBlobSizeHistogram() *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "azure",
			Subsystem: "storage",
			Name:      "blob_size_bytes",
			Help:      "Histograms of Azure Storage blob size bytes",
			Buckets: []float64{
				1 * megaBytes, 50 * megaBytes, 100 * megaBytes,
				250 * megaBytes, 500 * megaBytes, 1 * gigaBytes,
			},
		},
		[]string{"account", "container"},
	)
}

// -----------------------------------------------------------------------------

func init() {
	prometheus.MustRegister(storageAccountBlobSizeHistogram)

	RegisterUpdateMetricsFunctionsWithInterval("UpdateStorageMetrics", UpdateStorageMetrics, 60*time.Minute)
}

// UpdateStorageMetrics updates storage metrics.
func UpdateStorageMetrics(ctx context.Context) {
	contextLogger := log.WithFields(log.Fields{
		"_id":       ctx.Value("id").(string),
		"_function": "UpdateStorageMetrics",
	})

	azureClients := azure.NewAzureClients()
	storageAccounts, err := azure.ListSubscriptionStorageAccounts(ctx, azureClients, os.Getenv("AZURE_SUBSCRIPTION_ID"))

	if err != nil {
		contextLogger.Errorf("Unable to list account azure storage accounts: %s", err)
		return
	}

	hist := newStorageAccountBlobSizeHistogram()
	blobsMetrics := StorageAccountContainerMetrics{BlobsSizeHistogram: hist}

	// Loop over storage accounts.
	for _, account := range *storageAccounts {
		accountLogger := contextLogger.WithFields(log.Fields{
			"account": *account.Name,
		})

		accountLogger.Debugf("Start updating storage account")
		containers, err := azure.ListStorageAccountContainers(ctx, azureClients, &account)

		if err != nil {
			contextLogger.Fatalf("%v", err)
			blobsMetrics.DeleteLabelValues(*account.Name)
			continue
		}

		// Create a bounded wait group which allows 4 concurrent processes for
		// updating account's container's metrics.
		wg := tools.NewBoundedWaitGroup(4)

		// Loop over storage accounts
		for key := range *containers {
			// wg needs to be incremented outside the goroutine otherwise we could
			// reach wg.Wait() before wg.Add(1) is hit if it is in the goroutine.
			wg.Add(1)

			go func(wg *tools.BoundedWaitGroup, account *storage.Account, container *storage.ListContainerItem, walker *StorageAccountContainerMetrics) {
				accountLogger.Debugf("Start updating container: %s", *container.Name)

				t0 := time.Now()
				err := azure.WalkStorageAccount(ctx, azureClients, account, container, walker)
				t1 := time.Since(t0)

				if err != nil {
					accountLogger.Error(err)
				}

				accountLogger.Debugf("Done updating container: %s (%v)", *container.Name, t1)

				wg.Done()
			}(&wg, &account, &(*containers)[key], &blobsMetrics)
			// --------------^^^^^^^^^^^^^^^^^^^----------------
			// https://play.golang.org/p/YRGEg4LS5jd
			// https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
			// -------------------------------------------------
		}

		wg.Wait()

		accountLogger.Debugf("Done updating storage account")
	}

	// swapping current registered histogram with an updated copy
	*storageAccountBlobSizeHistogram = *blobsMetrics.BlobsSizeHistogram
}
