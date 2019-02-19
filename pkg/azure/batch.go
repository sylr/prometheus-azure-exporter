package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/2018-08-01.7.0/batch"
	azurebatch "github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools/cache"
)

var (
	cacheKeySubscriptionBatchAccounts     = `sub-%s-batch-accounts`
	cacheKeySubscriptionBatchAccountPools = `sub-%s-batch-account-%s-pools`
	cacheKeySubscriptionBatchAccountJobs  = `sub-%s-batch-account-%s-jobs`
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

	// AzureAPIBatchCallsDurationSecondsBuckets Histograms of Azure Batch API calls durations in seconds
	AzureAPIBatchCallsDurationSecondsBuckets = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "azure_api",
			Subsystem: "batch",
			Name:      "calls_duration_seconds",
			Help:      "Histograms of Azure Batch API calls durations in seconds",
			Buckets:   []float64{0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09, 0.10, 0.15, 0.20, 0.30, 0.40, 0.50, 1.0},
		},
		[]string{"subscription", "resource_group", "account"},
	)
)

func init() {
	prometheus.MustRegister(AzureAPIBatchCallsTotal)
	prometheus.MustRegister(AzureAPIBatchCallsFailedTotal)
	prometheus.MustRegister(AzureAPIBatchCallsDurationSecondsBuckets)
}

// ObserveAzureBatchAPICall ...
func ObserveAzureBatchAPICall(duration float64, labels ...string) {
	AzureAPIBatchCallsTotal.WithLabelValues(labels...).Inc()
	AzureAPIBatchCallsDurationSecondsBuckets.WithLabelValues(labels...).Observe(duration)
}

// ObserveAzureBatchAPICallFailed ...
func ObserveAzureBatchAPICallFailed(duration float64, labels ...string) {
	AzureAPIBatchCallsFailedTotal.WithLabelValues(labels...).Inc()
}

// ListSubscriptionBatchAccounts List all subscription batch accounts
func ListSubscriptionBatchAccounts(ctx context.Context, clients *AzureClients, subscription *subscription.Model) (*[]azurebatch.Account, error) {
	c := cache.GetCache(5 * time.Minute)
	cacheKey := fmt.Sprintf(cacheKeySubscriptionBatchAccounts, *subscription.SubscriptionID)

	contextLogger := log.WithFields(log.Fields{
		"_id":          ctx.Value("id").(string),
		"subscription": *subscription.DisplayName,
	})

	if caccounts, ok := c.Get(cacheKey); ok {
		if accounts, ok := caccounts.(*[]azurebatch.Account); !ok {
			contextLogger.Errorf("Failed to cast object from cache back to []azurebatch.Account")
		} else {
			return accounts, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetBatchAccountClient(*subscription.SubscriptionID)

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

	vals := accounts.Values()
	c.SetDefault(cacheKey, &vals)

	return &vals, nil
}

// ListBatchAccountPools List all batch account's pools
func ListBatchAccountPools(ctx context.Context, clients *AzureClients, subscription *subscription.Model, account *azurebatch.Account) ([]azurebatch.Pool, error) {
	c := cache.GetCache(5 * time.Minute)

	accountDetails, _ := ParseResourceID(*account.ID)
	cacheKey := fmt.Sprintf(cacheKeySubscriptionBatchAccountPools, *subscription.SubscriptionID, *account.Name)

	contextLogger := log.WithFields(log.Fields{
		"_id":     ctx.Value("id").(string),
		"rg":      accountDetails.ResourceGroup,
		"account": *account.Name,
	})

	if cpools, ok := c.Get(cacheKey); ok {
		if pools, ok := cpools.([]azurebatch.Pool); !ok {
			contextLogger.Errorf("Failed to cast object from cache back to []azurebatch.Pool")
		} else {
			return pools, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetBatchPoolClient(accountDetails.SubscriptionID)

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	pools, err := client.ListByBatchAccount(ctx, accountDetails.ResourceGroup, *account.Name, nil, "", "")
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)
	ObserveAzureBatchAPICall(t1, *subscription.DisplayName, accountDetails.ResourceGroup, *account.Name)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		ObserveAzureBatchAPICallFailed(t1, *subscription.DisplayName, accountDetails.ResourceGroup, *account.Name)
		return nil, err
	}

	vals := pools.Values()
	c.SetDefault(cacheKey, vals)

	return vals, nil
}

// ListBatchAccountJobs list batch account jobs
func ListBatchAccountJobs(ctx context.Context, clients *AzureClients, subscription *subscription.Model, account *azurebatch.Account) ([]batch.CloudJob, error) {
	c := cache.GetCache(5 * time.Minute)

	accountDetails, _ := ParseResourceID(*account.ID)
	cacheKey := fmt.Sprintf(cacheKeySubscriptionBatchAccountJobs, *subscription.SubscriptionID, *account.Name)

	contextLogger := log.WithFields(log.Fields{
		"_id":     ctx.Value("id").(string),
		"rg":      accountDetails.ResourceGroup,
		"account": *account.Name,
	})

	if cjobs, ok := c.Get(cacheKey); ok {
		if jobs, ok := cjobs.([]batch.CloudJob); !ok {
			contextLogger.Errorf("Failed to cast object from cache back to []batch.CloudJob")
		} else {
			return jobs, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetBatchJobClientWithResource(*account.AccountEndpoint, "https://batch.core.windows.net/")

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	jobs, err := client.List(ctx, "", "", "", nil, nil, nil, nil, nil)
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)
	ObserveAzureBatchAPICall(t1, *subscription.DisplayName, accountDetails.ResourceGroup, *account.Name)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		ObserveAzureBatchAPICallFailed(t1, *subscription.DisplayName, accountDetails.ResourceGroup, *account.Name)
		return nil, err
	}

	vals := jobs.Values()
	c.SetDefault(cacheKey, vals)

	return jobs.Values(), nil
}

// GetBatchJobTaskCounts get job tasks metrics
func GetBatchJobTaskCounts(ctx context.Context, clients *AzureClients, subscription *subscription.Model, account *azurebatch.Account, job *batch.CloudJob) (*batch.TaskCounts, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	accountDetails, _ := ParseResourceID(*account.ID)
	client, err := clients.GetBatchJobClientWithResource(*account.AccountEndpoint, "https://batch.core.windows.net/")

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	taskCounts, err := client.GetTaskCounts(ctx, *job.ID, nil, nil, nil, nil)
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)
	ObserveAzureBatchAPICall(t1, *subscription.DisplayName, accountDetails.ResourceGroup, *account.Name)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		ObserveAzureBatchAPICallFailed(t1, *subscription.DisplayName, accountDetails.ResourceGroup, *account.Name)
		return nil, err
	}

	return &taskCounts, nil
}
