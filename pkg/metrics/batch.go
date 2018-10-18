package metrics

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/azure"
	"github.com/sylr/prometheus-client-golang/prometheus"
)

var (
	batchJobsActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "jobs_active",
			Help:      "Number of active batch jobs",
		},
		[]string{"account", "job_id", "job_name"},
	)

	batchJobsRunning = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "jobs_running",
			Help:      "Number of running batch jobs",
		},
		[]string{"account", "job_id", "job_name"},
	)

	batchJobsCompleted = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "jobs_completed_total",
			Help:      "Total number of completed batch jobs",
		},
		[]string{"account", "job_id", "job_name"},
	)

	batchJobsFailed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure",
			Subsystem: "batch",
			Name:      "jobs_failed_total",
			Help:      "Total number of failed batch jobs",
		},
		[]string{"account", "job_id", "job_name"},
	)
)

func init() {
	prometheus.MustRegister(batchJobsActive)
	prometheus.MustRegister(batchJobsRunning)
	prometheus.MustRegister(batchJobsCompleted)
	prometheus.MustRegister(batchJobsFailed)
}

// UpdateMetrics
func UpdateMetrics(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)

	// Process which updates metrics
	go func(ctx context.Context) {
		for t := range ticker.C {
			log.Debugf("Start metrics update: %s", t)

			// We detach the update process so if it takes more than the refresh
			// time it does not get blocked
			go func(ctx context.Context, t time.Time) {
				log.Debugf("Start Batch metrics update: %s", t)

				jobsMetrics, err := azure.FetchAzureBatchMetrics(ctx)

				if err != nil {
					log.Error(err)
					return
				}

				for _, jobMetrics := range *jobsMetrics {
					batchJobsActive.WithLabelValues(jobMetrics.Account, jobMetrics.JobID, jobMetrics.JobDisplayName).Set(float64(*jobMetrics.Active))
					batchJobsRunning.WithLabelValues(jobMetrics.Account, jobMetrics.JobID, jobMetrics.JobDisplayName).Set(float64(*jobMetrics.Running))
					batchJobsCompleted.WithLabelValues(jobMetrics.Account, jobMetrics.JobID, jobMetrics.JobDisplayName).Set(float64(*jobMetrics.Completed))
					batchJobsFailed.WithLabelValues(jobMetrics.Account, jobMetrics.JobID, jobMetrics.JobDisplayName).Set(float64(*jobMetrics.Failed))

					log.WithFields(log.Fields{
						"account":   jobMetrics.Account,
						"job":       jobMetrics.JobDisplayName,
						"active":    *jobMetrics.Active,
						"running":   *jobMetrics.Running,
						"completed": *jobMetrics.Completed,
						"failed":    *jobMetrics.Failed,
					}).Debug("Batch job")
				}

				log.Debugf("End Batch metrics update: %s", t)
			}(ctx, t)
		}
	}(ctx)
}
