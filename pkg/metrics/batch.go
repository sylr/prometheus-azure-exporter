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

	for _, account := range *batchAccounts {
		accountProperties, _ := azure.ParseResourceID(*account.ID)

		// logger
		accountLogger := contextLogger.WithFields(log.Fields{
			"_id":     ctx.Value("id").(string),
			"rg":      accountProperties.ResourceGroup,
			"account": *account.Name,
		})

		// subscription
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

				accountLogger.WithFields(log.Fields{
					"metric":          "pool",
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

			accountLogger.Errorf("Unable to list account jobs: %s", err)
		} else {
			for _, job := range jobs {
				jobLogger := accountLogger.WithFields(log.Fields{
					"job_id": *job.ID,
				})

				// job task count
				taskCounts, err := azure.GetBatchJobTaskCounts(ctx, azureClients, &account, &job)

				// job.DisplayName can be nil but we don't want that
				displayName := *job.ID
				if job.DisplayName != nil {
					displayName = *job.DisplayName
				} else {
					jobLogger.Warnf("Job has no display name, defaulting to job.ID")
				}

				if err != nil {
					batchJobsTasksActive.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, displayName)
					batchJobsTasksRunning.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, displayName)
					batchJobsTasksCompleted.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, displayName)
					batchJobsTasksSucceeded.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, displayName)
					batchJobsTasksFailed.DeleteLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, displayName)

					contextLogger.Error(err)
					continue
				}

				// <!-- metrics
				batchJobsTasksActive.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, displayName).Set(float64(*taskCounts.Active))
				batchJobsTasksRunning.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, displayName).Set(float64(*taskCounts.Running))
				batchJobsTasksCompleted.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, displayName).Set(float64(*taskCounts.Completed))
				batchJobsTasksSucceeded.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, displayName).Set(float64(*taskCounts.Succeeded))
				batchJobsTasksFailed.WithLabelValues(*sub.DisplayName, accountProperties.ResourceGroup, *account.Name, *job.ID, displayName).Set(float64(*taskCounts.Failed))
				// metrics -->

				jobLogger.WithFields(log.Fields{
					"metric":    "job",
					"job":       displayName,
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
