package azure

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/2018-08-01.7.0/batch"
	azurebatch "github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools"
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

// ObserveAzureBatchAPICall
func ObserveAzureBatchAPICall(duration float64, labels ...string) {
	AzureAPIBatchCallsTotal.WithLabelValues(labels...).Inc()
	AzureAPIBatchCallsDurationSecondsBuckets.WithLabelValues(labels...).Observe(duration)
}

// ObserveAzureBatchAPICallFailed
func ObserveAzureBatchAPICallFailed(duration float64, labels ...string) {
	AzureAPIBatchCallsFailedTotal.WithLabelValues(labels...).Inc()
}

// ListSubscriptionBatchAccounts List all subscription batch accounts
func ListSubscriptionBatchAccounts(ctx context.Context, clients *AzureClients, subscriptionID string) (*[]azurebatch.Account, error) {
	c := tools.GetCache(5 * time.Minute)

	contextLogger := log.WithFields(log.Fields{
		"_id":          ctx.Value("id").(string),
		"subscription": subscriptionID,
	})

	if caccounts, ok := c.Get(subscriptionID + "-accounts"); ok {
		if accounts, ok := caccounts.(*[]azurebatch.Account); ok {
			contextLogger.Debugf("Got []azurebatch.Account from cache")
			return accounts, nil
		} else {
			contextLogger.Errorf("Failed to cast object from cache back to []azurebatch.Account")
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetBatchAccountClient(subscriptionID)

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
	c.SetDefault(subscriptionID+"-accounts", &vals)

	return &vals, nil
}

// ListBatchAccountPools List all batch account's pools
func ListBatchAccountPools(ctx context.Context, clients *AzureClients, account *azurebatch.Account) ([]azurebatch.Pool, error) {
	c := tools.GetCache(5 * time.Minute)

	accountResourceDetails, err := ParseResourceID(*account.ID)

	contextLogger := log.WithFields(log.Fields{
		"_id":     ctx.Value("id").(string),
		"rg":      accountResourceDetails.ResourceGroup,
		"account": *account.Name,
	})

	if cpools, ok := c.Get(*account.Name + "-pools"); ok {
		if pools, ok := cpools.([]azurebatch.Pool); ok {
			contextLogger.Debugf("Got []azurebatch.Pool from cache")
			return pools, nil
		} else {
			contextLogger.Errorf("Failed to cast object from cache back to []azurebatch.Pool")
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	sub, err := GetSubscription(ctx, clients, accountResourceDetails.SubscriptionID)
	client, err := clients.GetBatchPoolClient(accountResourceDetails.SubscriptionID)

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	pools, err := client.ListByBatchAccount(ctx, accountResourceDetails.ResourceGroup, *account.Name, nil, "", "")
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)
	ObserveAzureBatchAPICall(t1, *sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		ObserveAzureBatchAPICallFailed(t1, *sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name)
		return nil, err
	}

	vals := pools.Values()
	c.SetDefault(*account.Name+"-pools", vals)

	return vals, nil
}

// ListBatchAccountJobs list batch account jobs
func ListBatchAccountJobs(ctx context.Context, clients *AzureClients, account *azurebatch.Account) ([]batch.CloudJob, error) {
	c := tools.GetCache(5 * time.Minute)

	accountResourceDetails, err := ParseResourceID(*account.ID)

	contextLogger := log.WithFields(log.Fields{
		"_id":     ctx.Value("id").(string),
		"rg":      accountResourceDetails.ResourceGroup,
		"account": *account.Name,
	})

	if cjobs, ok := c.Get(*account.Name + "-jobs"); ok {
		if jobs, ok := cjobs.([]batch.CloudJob); ok {
			contextLogger.Debugf("Got []batch.CloudJob from cache")
			return jobs, nil
		} else {
			contextLogger.Errorf("Failed to cast object from cache back to []batch.CloudJob")
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	sub, err := GetSubscription(ctx, clients, accountResourceDetails.SubscriptionID)
	client, err := clients.GetBatchJobClientWithResource(*account.AccountEndpoint, "https://batch.core.windows.net/")

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	jobs, err := client.List(ctx, "", "", "", nil, nil, nil, nil, nil)
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)
	ObserveAzureBatchAPICall(t1, *sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		ObserveAzureBatchAPICallFailed(t1, *sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name)
		return nil, err
	}

	vals := jobs.Values()
	c.SetDefault(*account.Name+"-jobs", vals)

	return jobs.Values(), nil
}

// GetBatchJobTaskCounts get job tasks metrics
func GetBatchJobTaskCounts(ctx context.Context, clients *AzureClients, account *azurebatch.Account, job *batch.CloudJob) (*batch.TaskCounts, error) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	accountResourceDetails, err := ParseResourceID(*account.ID)
	sub, err := GetSubscription(ctx, clients, accountResourceDetails.SubscriptionID)
	client, err := clients.GetBatchJobClientWithResource(*account.AccountEndpoint, "https://batch.core.windows.net/")

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	taskCounts, err := client.GetTaskCounts(ctx, *job.ID, nil, nil, nil, nil)
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)
	ObserveAzureBatchAPICall(t1, *sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		ObserveAzureBatchAPICallFailed(t1, *sub.DisplayName, accountResourceDetails.ResourceGroup, *account.Name)
		return nil, err
	}

	return &taskCounts, nil
}
