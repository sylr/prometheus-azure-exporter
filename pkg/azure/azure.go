package azure

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// AzureAPICallsTotal Total number of Azure API calls
	AzureAPICallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "",
			Name:      "calls_total",
			Help:      "Total number of calls to the Azure API",
		},
		[]string{},
	)

	// AzureAPICallsFailedTotal Total number of failed Azure API calls
	AzureAPICallsFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "",
			Name:      "calls_failed_total",
			Help:      "Total number of failed calls to the Azure API",
		},
		[]string{},
	)

	// AzureAPICallsDurationSeconds Percentiles of Azure API calls durations in seconds over last 10 minutes
	AzureAPICallsDurationSeconds = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "azure_api",
			Subsystem:  "",
			Name:       "calls_duration_seconds",
			Help:       "Percentiles of Azure API calls durations in seconds over last 10 minutes",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
			BufCap:     50000,
			MaxAge:     10 * time.Minute,
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(AzureAPICallsTotal)
	prometheus.MustRegister(AzureAPICallsFailedTotal)
	prometheus.MustRegister(AzureAPICallsDurationSeconds)
}
