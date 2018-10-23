package metrics

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/azure"
	"github.com/sylr/prometheus-client-golang/prometheus"
)

var (
	batchPoolsDedicatedNodes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "pools_dedicated_nodes",
			Help:      "Number of dedicated nodes for batch account",
		},
		[]string{"account", "pool_name"},
	)

	batchJobsTasksActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "jobs_tasks_active",
			Help:      "Number of active batch jobs",
		},
		[]string{"account", "job_id", "job_name"},
	)

	batchJobsTasksRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "jobs_tasks_running",
			Help:      "Number of running batch jobs",
		},
		[]string{"account", "job_id", "job_name"},
	)

	batchJobsTasksCompleted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "jobs_tasks_completed_total",
			Help:      "Total number of completed batch jobs",
		},
		[]string{"account", "job_id", "job_name"},
	)

	batchJobsTasksSucceeded = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "jobs_tasks_succeeded_total",
			Help:      "Total number of succeeded batch jobs",
		},
		[]string{"account", "job_id", "job_name"},
	)

	batchJobsTasksFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "jobs_tasks_failed_total",
			Help:      "Total number of failed batch jobs",
		},
		[]string{"account", "job_id", "job_name"},
	)
)

func init() {
	prometheus.MustRegister(batchPoolsDedicatedNodes)
	prometheus.MustRegister(batchJobsTasksActive)
	prometheus.MustRegister(batchJobsTasksRunning)
	prometheus.MustRegister(batchJobsTasksCompleted)
	prometheus.MustRegister(batchJobsTasksSucceeded)
	prometheus.MustRegister(batchJobsTasksFailed)

	RegisterUpdateMetricsFunctions("UpdateBatchMetrics", UpdateBatchMetrics)
}

// UpdateBatchMetrics updates batch metrics
func UpdateBatchMetrics(ctx context.Context, id string) {
	contextLogger := log.WithFields(log.Fields{"_id": id})
	azureClients := azure.GetNewAzureClients()
	batchAccounts, err := azure.ListSubscriptionBatchAccounts(ctx, azureClients, os.Getenv("AZURE_SUBSCRIPTION_ID"))

	if err != nil {
		contextLogger.Errorf("Unable to list account azure batch accounts: %s", err)
		return
	}

	for _, account := range batchAccounts {
		// <!-- POOLS ----------------------------------------------------------
		pools, err := azure.ListBatchAccountPools(ctx, azureClients, &account)

		if err != nil {
			contextLogger.Errorf("Unable to list account `%s` pools: %s", *account.Name, err)
		} else {
			for _, pool := range pools {
				batchPoolsDedicatedNodes.WithLabelValues(*account.Name, *pool.Name).Set(float64(*pool.PoolProperties.CurrentDedicatedNodes))

				contextLogger.WithFields(log.Fields{
					"_id":             id,
					"account":         *account.Name,
					"pool":            *pool.Name,
					"dedicated_nodes": *pool.PoolProperties.CurrentDedicatedNodes,
				}).Debug("Batch pool")
			}
		}
		// ---------------------------------------------------------- POOLS --!>

		// <!-- JOBS -----------------------------------------------------------
		jobs, err := azure.ListBatchAccountJobs(ctx, azureClients, &account)

		if err != nil {
			contextLogger.Errorf("Unable to list account `%s` jobs: %s", *account.Name, err)
		} else {
			for _, job := range jobs {
				taskCounts, err := azure.GetBatchJobTaskCounts(ctx, azureClients, &account, &job)

				if err != nil {
					contextLogger.Error(err)
					continue
				}

				batchJobsTasksActive.WithLabelValues(*account.Name, *job.ID, *job.DisplayName).Set(float64(*taskCounts.Active))
				batchJobsTasksRunning.WithLabelValues(*account.Name, *job.ID, *job.DisplayName).Set(float64(*taskCounts.Running))
				batchJobsTasksCompleted.WithLabelValues(*account.Name, *job.ID, *job.DisplayName).Set(float64(*taskCounts.Completed))
				batchJobsTasksSucceeded.WithLabelValues(*account.Name, *job.ID, *job.DisplayName).Set(float64(*taskCounts.Succeeded))
				batchJobsTasksFailed.WithLabelValues(*account.Name, *job.ID, *job.DisplayName).Set(float64(*taskCounts.Failed))

				contextLogger.WithFields(log.Fields{
					"_id":       id,
					"account":   *account.Name,
					"job":       *job.DisplayName,
					"active":    *taskCounts.Active,
					"running":   *taskCounts.Running,
					"completed": *taskCounts.Completed,
					"succeeded": *taskCounts.Succeeded,
					"failed":    *taskCounts.Failed,
				}).Debug("Batch job")
			}
		}
		// ----------------------------------------------------------- JOBS --!>
	}
}
