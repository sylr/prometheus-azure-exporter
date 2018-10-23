package azure

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/2018-08-01.7.0/batch"
	azurebatch "github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	"github.com/sylr/prometheus-client-golang/prometheus"
)

var (
	// AzureAPIBatchCallsTotal Total number of Azure Batch API calls
	AzureAPIBatchCallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "batch",
			Name:      "calls_total",
			Help:      "Total number of calls to the Azure API",
		},
		[]string{"subscription", "resource_group", "account"},
	)

	// AzureAPIBatchCallsFailedTotal Total number of failed Azure Batch API calls
	AzureAPIBatchCallsFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "batch",
			Name:      "calls_failed_total",
			Help:      "Total number of failed calls to the Azure API",
		},
		[]string{"subscription", "resource_group", "account"},
	)
)

func init() {
	prometheus.MustRegister(AzureAPIBatchCallsTotal)
	prometheus.MustRegister(AzureAPIBatchCallsFailedTotal)
}

// ListSubscriptionBatchAccounts List all subscription batch accounts
func ListSubscriptionBatchAccounts(ctx context.Context, clients *AzureClients, subscriptionID string) ([]azurebatch.Account, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetBatchAccountClient(subscriptionID)

	if err != nil {
		return nil, err
	}

	batchAccounts, err := client.List(ctx)
	AzureAPICallsTotal.WithLabelValues().Inc()

	if err != nil {
		AzureAPICallsFailedTotal.WithLabelValues().Inc()
		return nil, err
	}

	return batchAccounts.Values(), nil
}

// ListBatchAccountPools List all batch account's pools
func ListBatchAccountPools(ctx context.Context, clients *AzureClients, account *azurebatch.Account) ([]azurebatch.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	accountResourceDetails, err := ParseResourceID(*account.ID)
	sub, err := GetSubscription(ctx, clients, accountResourceDetails.SubscriptionID)
	client, err := clients.GetBatchPoolClient(accountResourceDetails.SubscriptionID)

	if err != nil {
		return nil, err
	}

	accountPools, err := client.ListByBatchAccount(ctx, accountResourceDetails.ResourceGroup, *account.Name, nil, "", "")
	AzureAPICallsTotal.WithLabelValues().Inc()
	AzureAPIBatchCallsTotal.WithLabelValues(*sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name).Inc()

	if err != nil {
		AzureAPICallsFailedTotal.WithLabelValues().Inc()
		AzureAPIBatchCallsFailedTotal.WithLabelValues(*sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name).Inc()
		return nil, err
	}

	return accountPools.Values(), nil
}

// ListBatchAccountJobs list batch account jobs
func ListBatchAccountJobs(ctx context.Context, clients *AzureClients, account *azurebatch.Account) ([]batch.CloudJob, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	accountResourceDetails, err := ParseResourceID(*account.ID)
	sub, err := GetSubscription(ctx, clients, accountResourceDetails.SubscriptionID)
	client, err := clients.GetBatchJobClientWithResource(*account.AccountEndpoint, "https://batch.core.windows.net/")

	if err != nil {
		return nil, err
	}

	jobs, err := client.List(ctx, "", "", "", nil, nil, nil, nil, nil)
	AzureAPICallsTotal.WithLabelValues().Inc()
	AzureAPIBatchCallsTotal.WithLabelValues(*sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name).Inc()

	if err != nil {
		AzureAPICallsFailedTotal.WithLabelValues().Inc()
		AzureAPIBatchCallsFailedTotal.WithLabelValues(*sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name).Inc()
		return nil, err
	}

	return jobs.Values(), nil
}

// GetBatchJobTaskCounts
func GetBatchJobTaskCounts(ctx context.Context, clients *AzureClients, account *azurebatch.Account, job *batch.CloudJob) (*batch.TaskCounts, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	accountResourceDetails, err := ParseResourceID(*account.ID)
	sub, err := GetSubscription(ctx, clients, accountResourceDetails.SubscriptionID)
	client, err := clients.GetBatchJobClientWithResource(*account.AccountEndpoint, "https://batch.core.windows.net/")

	if err != nil {
		return nil, err
	}

	taskCounts, err := client.GetTaskCounts(ctx, *job.ID, nil, nil, nil, nil)
	AzureAPICallsTotal.WithLabelValues().Inc()
	AzureAPIBatchCallsTotal.WithLabelValues(*sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name).Inc()

	if err != nil {
		AzureAPICallsFailedTotal.WithLabelValues().Inc()
		AzureAPIBatchCallsFailedTotal.WithLabelValues(*sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name).Inc()
		return nil, err
	}

	return &taskCounts, nil
}
