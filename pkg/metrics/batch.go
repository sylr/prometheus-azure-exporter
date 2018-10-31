package metrics

import (
	"context"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/azure"
)

var (
	batchPoolQuota = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "pool_quota",
			Help:      "Quota of pool for batch account",
		},
		[]string{"subscription", "resource_group", "account"},
	)

	batchDedicatedCoreQuota = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "dedicated_core_quota",
			Help:      "Quota of dedicated core for batch account",
		},
		[]string{"subscription", "resource_group", "account"},
	)

	batchPoolsDedicatedNodes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "pool_dedicated_nodes",
			Help:      "Number of dedicated nodes for batch pool",
		},
		[]string{"subscription", "resource_group", "account", "pool"},
	)

	batchJobsTasksActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "job_tasks_active",
			Help:      "Number of active batch job task",
		},
		[]string{"subscription", "resource_group", "account", "job_id", "job_name"},
	)

	batchJobsTasksRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "job_tasks_running",
			Help:      "Number of running batch job task",
		},
		[]string{"subscription", "resource_group", "account", "job_id", "job_name"},
	)

	batchJobsTasksCompleted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "job_tasks_completed_total",
			Help:      "Total number of completed batch job task",
		},
		[]string{"subscription", "resource_group", "account", "job_id", "job_name"},
	)

	batchJobsTasksSucceeded = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "job_tasks_succeeded_total",
			Help:      "Total number of succeeded batch job task",
		},
		[]string{"subscription", "resource_group", "account", "job_id", "job_name"},
	)

	batchJobsTasksFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "job_tasks_failed_total",
			Help:      "Total number of failed batch job task",
		},
		[]string{"subscription", "resource_group", "account", "job_id", "job_name"},
	)
)

func init() {
	prometheus.MustRegister(batchPoolQuota)
	prometheus.MustRegister(batchDedicatedCoreQuota)
	prometheus.MustRegister(batchPoolsDedicatedNodes)
	prometheus.MustRegister(batchJobsTasksActive)
	prometheus.MustRegister(batchJobsTasksRunning)
	prometheus.MustRegister(batchJobsTasksCompleted)
	prometheus.MustRegister(batchJobsTasksSucceeded)
	prometheus.MustRegister(batchJobsTasksFailed)

	RegisterUpdateMetricsFunctions("UpdateBatchMetrics", UpdateBatchMetrics)
}

// UpdateBatchMetrics updates batch metrics
func UpdateBatchMetrics(ctx context.Context) {
	contextLogger := log.WithFields(log.Fields{"_id": ctx.Value("id").(string)})
	azureClients := azure.GetNewAzureClients()
	batchAccounts, err := azure.ListSubscriptionBatchAccounts(ctx, azureClients, os.Getenv("AZURE_SUBSCRIPTION_ID"))

	if err != nil {
		contextLogger.Errorf("Unable to list account azure batch accounts: %s", err)
		return
	}

	for _, account := range batchAccounts {
		accountProperties, _ := azure.ParseResourceID(*account.ID)
		sub, err := azure.GetSubscription(ctx, azureClients, os.Getenv("AZURE_SUBSCRIPTION_ID"))

		// <!-- metrics
		batchPoolQuota.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name).Set(float64(*account.PoolQuota))
		batchDedicatedCoreQuota.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name).Set(float64(*account.DedicatedCoreQuota))
		// metrics -->

		// <!-- POOLS ----------------------------------------------------------
		pools, err := azure.ListBatchAccountPools(ctx, azureClients, &account)

		if err != nil {
			batchPoolsDedicatedNodes.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name)
			contextLogger.Errorf("Unable to list account `%s` pools: %s", *account.Name, err)
		} else {
			for _, pool := range pools {
				// <!-- metrics
				batchPoolsDedicatedNodes.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *pool.Name).Set(float64(*pool.PoolProperties.CurrentDedicatedNodes))
				// metrics -->

				contextLogger.WithFields(log.Fields{
					"_id":             ctx.Value("id").(string),
					"metric":          "pool",
					"rg":              accountProperties.ResourceGroup,
					"account":         *account.Name,
					"pool":            *pool.Name,
					"dedicated_nodes": *pool.PoolProperties.CurrentDedicatedNodes,
				}).Debug("")
			}
		}
		// ----------------------------------------------------------- POOLS -->

		// <!-- JOBS -----------------------------------------------------------
		jobs, err := azure.ListBatchAccountJobs(ctx, azureClients, &account)

		if err != nil {
			batchJobsTasksActive.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name)
			batchJobsTasksRunning.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name)
			batchJobsTasksCompleted.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name)
			batchJobsTasksSucceeded.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name)
			batchJobsTasksFailed.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name)

			contextLogger.Errorf("Unable to list account `%s` jobs: %s", *account.Name, err)
		} else {
			for _, job := range jobs {
				taskCounts, err := azure.GetBatchJobTaskCounts(ctx, azureClients, &account, &job)

				if err != nil {
					batchJobsTasksActive.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, *job.DisplayName)
					batchJobsTasksRunning.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, *job.DisplayName)
					batchJobsTasksCompleted.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, *job.DisplayName)
					batchJobsTasksSucceeded.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, *job.DisplayName)
					batchJobsTasksFailed.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, *job.DisplayName)

					contextLogger.Error(err)
					continue
				}

				// <!-- metrics
				batchJobsTasksActive.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, *job.DisplayName).Set(float64(*taskCounts.Active))
				batchJobsTasksRunning.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, *job.DisplayName).Set(float64(*taskCounts.Running))
				batchJobsTasksCompleted.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, *job.DisplayName).Set(float64(*taskCounts.Completed))
				batchJobsTasksSucceeded.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, *job.DisplayName).Set(float64(*taskCounts.Succeeded))
				batchJobsTasksFailed.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, *job.DisplayName).Set(float64(*taskCounts.Failed))
				// metrics -->

				contextLogger.WithFields(log.Fields{
					"_id":       ctx.Value("id").(string),
					"metric":    "job",
					"rg":        accountProperties.ResourceGroup,
					"account":   *account.Name,
					"job_id":    *job.ID,
					"job":       *job.DisplayName,
					"pool":      *job.PoolInfo.PoolID,
					"active":    *taskCounts.Active,
					"running":   *taskCounts.Running,
					"completed": *taskCounts.Completed,
					"succeeded": *taskCounts.Succeeded,
					"failed":    *taskCounts.Failed,
				}).Debug("")
			}
		}
		// ----------------------------------------------------------- JOBS --!>
	}
}
