package azure

import (
	"context"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/batch/2018-08-01.7.0/batch"
	azurebatch "github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-client-golang/prometheus"
)

var (
	AzureAPIBatchCallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "batch",
			Name:      "calls_total",
			Help:      "Total number of calls to the Azure API",
		},
		[]string{"account", "job_id", "job_name"},
	)

	AzureAPIBatchCallsFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "batch",
			Name:      "calls_failed_total",
			Help:      "Total number of failed calls to the Azure API",
		},
		[]string{"account", "job_id", "job_name"},
	)
)

type AzureBatchJobMetrics struct {
	batch.TaskCounts
	Account        string
	JobID          string
	JobDisplayName string
}

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
	accountResourceDetails, err := ParseResourceID(*account.ID)
	client, err := clients.GetBatchPoolClient(accountResourceDetails.SubscriptionID)

	if err != nil {
		return nil, err
	}

	accountPools, err := client.ListByBatchAccount(ctx, accountResourceDetails.ResourceGroup, *account.Name, nil, "", "")
	AzureAPICallsTotal.WithLabelValues().Inc()
	AzureAPIBatchCallsTotal.WithLabelValues(*account.Name, "", "").Inc()

	if err != nil {
		AzureAPICallsFailedTotal.WithLabelValues().Inc()
		AzureAPIBatchCallsFailedTotal.WithLabelValues(*account.Name, "", "").Inc()
		return nil, err
	}

	return accountPools.Values(), nil
}

// ListBatchAccountJobs list batch account jobs
func ListBatchAccountJobs(ctx context.Context, clients *AzureClients, account *azurebatch.Account) ([]batch.CloudJob, error) {
	client, err := clients.GetBatchJobClientWithResource(*account.AccountEndpoint, "https://batch.core.windows.net/")

	if err != nil {
		return nil, err
	}

	jobs, err := client.List(ctx, "", "", "", nil, nil, nil, nil, nil)
	AzureAPICallsTotal.WithLabelValues().Inc()
	AzureAPIBatchCallsTotal.WithLabelValues(*account.Name, "", "").Inc()

	if err != nil {
		AzureAPICallsFailedTotal.WithLabelValues().Inc()
		AzureAPIBatchCallsFailedTotal.WithLabelValues(*account.Name, "", "").Inc()
		return nil, err
	}

	return jobs.Values(), nil
}

// FetchAzureBatchMetrics pwet
func FetchAzureBatchMetrics(ctx context.Context) (*[]AzureBatchJobMetrics, error) {
	azureClients := GetNewAzureClients()
	batchAccounts, err := ListSubscriptionBatchAccounts(ctx, azureClients, os.Getenv("AZURE_SUBSCRIPTION_ID"))

	if err != nil {
		log.Errorf("Unable to list account azure batch accounts: %s", err)
		return nil, err
	}

	jobsMetrics := make([]AzureBatchJobMetrics, 50)

	for i, account := range batchAccounts {
		jobClient, _ := azureClients.GetBatchJobClientWithResource(*account.AccountEndpoint, "https://batch.core.windows.net/")
		jobs, err := ListBatchAccountJobs(ctx, azureClients, &account)

		if err != nil {
			log.Errorf("Unable to list account `%s` jobs: %s", *account.Name, err)
			continue
		}

		for k, job := range jobs {
			taskCount, err := jobClient.GetTaskCounts(ctx, *job.ID, nil, nil, nil, nil)
			AzureAPICallsTotal.WithLabelValues().Inc()
			AzureAPIBatchCallsTotal.WithLabelValues(*account.Name, *job.ID, *job.DisplayName).Inc()

			if err != nil {
				log.Error(err)
				AzureAPICallsFailedTotal.WithLabelValues().Inc()
				AzureAPIBatchCallsFailedTotal.WithLabelValues(*account.Name, *job.ID, *job.DisplayName).Inc()
				continue
			}

			jobMetrics := AzureBatchJobMetrics{}
			jobMetrics.Active = taskCount.Active
			jobMetrics.Completed = taskCount.Completed
			jobMetrics.Failed = taskCount.Failed
			jobMetrics.Running = taskCount.Running
			jobMetrics.Account = *account.Name
			jobMetrics.JobID = *job.ID
			jobMetrics.JobDisplayName = *job.DisplayName

			index := (i+1)*(k+1) - 1

			if index < len(jobsMetrics) {
				jobsMetrics[index] = jobMetrics
			} else {
				jobsMetrics = append(jobsMetrics, jobMetrics)
			}
		}
	}

	return &jobsMetrics, nil
}
