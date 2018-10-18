package azure

import (
	"github.com/sylr/prometheus-client-golang/prometheus"
)

var (
	AzureAPICallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "",
			Name:      "calls_total",
			Help:      "Total number of calls to the Azure API",
		},
		[]string{},
	)

	AzureAPICallsFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "",
			Name:      "calls_failed_total",
			Help:      "Total number of failed calls to the Azure API",
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(AzureAPICallsTotal)
	prometheus.MustRegister(AzureAPICallsFailedTotal)
}
