package azure

import (
	"net/http"
	"strconv"
	"sync"

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

	// AzureAPICallsDurationSecondsBuckets Histograms of Azure API calls durations in seconds
	AzureAPICallsDurationSecondsBuckets = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "azure_api",
			Subsystem: "",
			Name:      "calls_duration_seconds",
			Help:      "Histograms of Azure API calls durations in seconds",
			Buckets:   []float64{0.02, 0.03, 0.04, 0.05, 0.06, 0.07, 0.08, 0.09, 0.10, 0.15, 0.20, 0.30, 0.40, 0.50, 1.0},
		},
		[]string{},
	)

	// AzureAPIReadRateLimitRemaining Gauge describing the current number of remaining read API calls
	AzureAPIReadRateLimitRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure_api",
			Subsystem: "",
			Name:      "read_rate_limit_remaining",
			Help:      "Gauge describing the current number of remaining read API calls allowed for the subscription",
		},
		[]string{"subscription"},
	)
)

var (
	azureAPIReadRateLimitRemainingMutex = sync.RWMutex{}
)

func init() {
	prometheus.MustRegister(AzureAPICallsTotal)
	prometheus.MustRegister(AzureAPICallsFailedTotal)
	prometheus.MustRegister(AzureAPICallsDurationSecondsBuckets)
	prometheus.MustRegister(AzureAPIReadRateLimitRemaining)
}

// ObserveAzureAPICall ...
func ObserveAzureAPICall(duration float64) {
	AzureAPICallsTotal.WithLabelValues().Inc()
	AzureAPICallsDurationSecondsBuckets.WithLabelValues().Observe(duration)
}

// ObserveAzureAPICallFailed ...
func ObserveAzureAPICallFailed(duration float64) {
	AzureAPICallsFailedTotal.WithLabelValues().Inc()
}

// SetReadRateLimitRemaining ...
func SetReadRateLimitRemaining(subscription string, reponse *http.Response) {
	remaining := reponse.Header.Get("x-ms-ratelimit-remaining-subscription-reads")

	if len(remaining) > 0 {
		f, err := strconv.ParseFloat(remaining, 64)

		if err == nil {
			azureAPIReadRateLimitRemainingMutex.Lock()
			AzureAPIReadRateLimitRemaining.WithLabelValues(subscription).Set(f)
			azureAPIReadRateLimitRemainingMutex.Unlock()
		}
	}
}
