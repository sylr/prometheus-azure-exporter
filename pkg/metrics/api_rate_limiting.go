package metrics

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/azure"
)

func init() {
	if GetUpdateMetricsFunctionInterval("api_rate_limiting") == nil {
		RegisterUpdateMetricsFunctionWithInterval("api_rate_limiting", UpdateAPIRateLimitingMetrics, 30*time.Second)
	}
}

// UpdateAPIRateLimitingMetrics updates api metrics.
func UpdateAPIRateLimitingMetrics(ctx context.Context) error {
	var err error

	contextLogger := log.WithFields(log.Fields{
		"_id":   ctx.Value("id").(string),
		"_func": "UpdateApiRateLimitingMetrics",
	})

	// In order to retrieve the remaining number of write API calls that Azure
	// allows for the subscription or the tenant id we need to do make a call to
	// an API which does a write operation.
	// The ListKeys API of storage accounts although really being a read API is
	// considered a write API by Azure so we are going to use that.

	// We do not need to register metrics in this module. All the work is done by
	// callbacks given to the Azure SDK. See:
	// - pkg/azure/clients.go -> respondInspect()
	// - pkg/azure/azure.go   -> SetReadRateLimitRemaining()
	//                           SetWriteRateLimitRemaining()

	azureClients := azure.NewAzureClients()
	sub, err := azure.GetSubscription(ctx, azureClients, os.Getenv("AZURE_SUBSCRIPTION_ID"))

	if err != nil {
		contextLogger.Errorf("Unable to get subscription: %s", err)
		return err
	}

	storageAccounts, err := azure.ListSubscriptionStorageAccounts(ctx, azureClients, sub)

	if err != nil {
		contextLogger.Errorf("Unable to list account azure storage accounts: %s", err)
		return err
	}

	// Loop over storage accounts.
	for accountKey := range *storageAccounts {
		_, err := azure.ListStorageAccountKeys(ctx, azureClients, sub, &(*storageAccounts)[accountKey])

		if err != nil {
			contextLogger.Error(err)
		} else {
			break
		}
	}

	return err
}
