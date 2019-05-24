package azure

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-07-01/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools/cache"
)

var (
	cacheKeySubscriptionStorageAccounts          = `sub-%s-storageaccounts`
	cacheKeySubscriptionStorageAccountContainers = `sub-%s-rg-%s-storageaccount-%s-containers`
	cacheKeySubscriptionStorageAccountKeys       = `sub-%s-rg-%s-storageaccount-%s-keys`
)

var (
	// AzureAPIStorageCallsTotal Total number of Azure Storage API calls
	AzureAPIStorageCallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "storage",
			Name:      "calls_total",
			Help:      "Total number of calls to the Azure API",
		},
		[]string{"subscription", "resource_group", "account"},
	)

	// AzureAPIStorageCallsFailedTotal Total number of failed Azure Storage API calls
	AzureAPIStorageCallsFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "storage",
			Name:      "calls_failed_total",
			Help:      "Total number of failed calls to the Azure API",
		},
		[]string{"subscription", "resource_group", "account"},
	)

	// AzureAPIStorageCallsDurationSecondsBuckets Histograms of Azure Storage API calls durations in seconds
	AzureAPIStorageCallsDurationSecondsBuckets = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "azure_api",
			Subsystem: "storage",
			Name:      "calls_duration_seconds",
			Help:      "Histograms of Azure Storage API calls durations in seconds",
			Buckets:   []float64{0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 4.0, 5.0},
		},
		[]string{"subscription", "resource_group", "account"},
	)
)

func init() {
	prometheus.MustRegister(AzureAPIStorageCallsTotal)
	prometheus.MustRegister(AzureAPIStorageCallsFailedTotal)
	prometheus.MustRegister(AzureAPIStorageCallsDurationSecondsBuckets)
}

// ObserveAzureStorageAPICall ...
func ObserveAzureStorageAPICall(duration float64, labels ...string) {
	AzureAPIStorageCallsTotal.WithLabelValues(labels...).Inc()
	AzureAPIStorageCallsDurationSecondsBuckets.WithLabelValues(labels...).Observe(duration)
}

// ObserveAzureStorageAPICallFailed ...
func ObserveAzureStorageAPICallFailed(duration float64, labels ...string) {
	AzureAPIStorageCallsFailedTotal.WithLabelValues(labels...).Inc()
}

// ListSubscriptionStorageAccounts ...
func ListSubscriptionStorageAccounts(ctx context.Context, clients *AzureClients, subscription *subscription.Model) (*[]storage.Account, error) {
	c := cache.GetCache(5 * time.Minute)
	cacheKey := fmt.Sprintf(cacheKeySubscriptionStorageAccounts, subscription.SubscriptionID)

	contextLogger := log.WithFields(log.Fields{
		"_id":          ctx.Value("id").(string),
		"subscription": subscription.DisplayName,
	})

	if caccounts, ok := c.Get(cacheKey); ok {
		if accounts, ok := caccounts.(*[]storage.Account); !ok {
			contextLogger.Errorf("Failed to cast object from cache back to *[]storage.Account")
		} else {
			return accounts, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetStorageAccountsClient(*subscription.SubscriptionID)

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	accounts, err := client.List(ctx)
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		return nil, err
	}

	vals := accounts.Value
	c.SetDefault(cacheKey, vals)

	return vals, nil
}

// ListStorageAccountContainers ...
func ListStorageAccountContainers(ctx context.Context, clients *AzureClients, subscription *subscription.Model, account *storage.Account) (*[]storage.ListContainerItem, error) {
	c := cache.GetCache(5 * time.Minute)

	accountDetails, _ := ParseResourceID(*account.ID)

	cacheKey := fmt.Sprintf(
		cacheKeySubscriptionStorageAccountContainers,
		*subscription.SubscriptionID,
		accountDetails.ResourceGroup,
		*account.Name,
	)

	contextLogger := log.WithFields(log.Fields{
		"_id":             ctx.Value("id").(string),
		"storage_account": *account.Name,
	})

	if ccontainers, ok := c.Get(cacheKey); ok {
		if containers, ok := ccontainers.(*[]storage.ListContainerItem); !ok {
			contextLogger.Errorf("Failed to cast object from cache back to *[]storage.ListContainerItem")
		} else {
			return containers, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetBlobContainersClient(*subscription.SubscriptionID)

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	containers, err := client.List(ctx, accountDetails.ResourceGroup, *account.Name)
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)
	ObserveAzureStorageAPICall(t1, *subscription.DisplayName, accountDetails.ResourceGroup, *account.Name)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		ObserveAzureStorageAPICallFailed(t1, *subscription.DisplayName, accountDetails.ResourceGroup, *account.Name)
		return nil, err
	}

	vals := *containers.Value
	c.SetDefault(cacheKey, &vals)

	return &vals, nil
}

// ListStorageAccountKeys ...
func ListStorageAccountKeys(ctx context.Context, clients *AzureClients, subscription *subscription.Model, account *storage.Account) (*[]storage.AccountKey, error) {
	c := cache.GetCache(30 * time.Second)

	accountDetails, _ := ParseResourceID(*account.ID)

	cacheKey := fmt.Sprintf(
		cacheKeySubscriptionStorageAccountKeys,
		*subscription.SubscriptionID,
		accountDetails.ResourceGroup,
		*account.Name,
	)

	contextLogger := log.WithFields(log.Fields{
		"_id":             ctx.Value("id").(string),
		"storage_account": *account.Name,
	})

	if ckeys, ok := c.Get(cacheKey); ok {
		if keys, ok := ckeys.(*[]storage.AccountKey); !ok {
			contextLogger.Errorf("Failed to cast object from cache back to *[]storage.AccountKey")
		} else {
			return keys, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetStorageAccountsClient(*subscription.SubscriptionID)

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	keys, err := client.ListKeys(ctx, accountDetails.ResourceGroup, *account.Name)
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)
	ObserveAzureStorageAPICall(t1, *subscription.DisplayName, accountDetails.ResourceGroup, *account.Name)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		ObserveAzureStorageAPICallFailed(t1, *subscription.DisplayName, accountDetails.ResourceGroup, *account.Name)
		return nil, err
	}

	vals := *keys.Keys
	c.SetDefault(cacheKey, &vals)

	return &vals, nil
}

// StorageAccountMetrics ...
type StorageAccountMetrics struct {
	ContainerBlobSizeHistogram *prometheus.HistogramVec
	mutex                      sync.RWMutex
}

// Lock is here to make sure several Walkers do not update ContainerBlobSizeHistogram
// at the same time.
func (s *StorageAccountMetrics) Lock() {
	s.mutex.Lock()
}

// Unlock releases the lock.
func (s *StorageAccountMetrics) Unlock() {
	s.mutex.Unlock()
}

// Reset resets the histogram data.
func (s *StorageAccountMetrics) Reset() {
	s.Lock()
	s.ContainerBlobSizeHistogram.Reset()
	s.Unlock()
}

// DeleteLabelValues deletes histogram's data associated with given labels.
func (s *StorageAccountMetrics) DeleteLabelValues(labels ...string) {
	s.Lock()
	s.ContainerBlobSizeHistogram.DeleteLabelValues(labels...)
	s.Unlock()
}

// WalkBlob is called over each blobs listed by the function walking the
// storage account container.
func (s *StorageAccountMetrics) WalkBlob(subscription *subscription.Model, group *resources.Group, account *storage.Account, container *storage.ListContainerItem, blob *azblob.BlobItem) {
	s.ContainerBlobSizeHistogram.
		WithLabelValues(*subscription.DisplayName, *group.Name, *account.Name, *container.Name).
		Observe(float64(*blob.Properties.ContentLength))
}
