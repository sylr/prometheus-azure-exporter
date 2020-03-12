package azure

import (
	"net/http"
	"strconv"
	"sync"
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
			Help:      "Total number of successful calls to the Azure API",
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
			Help:      "Histograms of successful Azure API calls durations in seconds",
			Buckets:   []float64{0.02, 0.03, 0.04, 0.05, 0.10, 0.20, 0.30, 0.40, 0.50, 0.75, 1.0, 2.0},
		},
		[]string{},
	)

	// AzureAPITenantReadRateLimitRemaining Gauge describing the current number of remaining read API calls
	AzureAPITenantReadRateLimitRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure_api",
			Subsystem: "tenant",
			Name:      "read_rate_limit_remaining",
			Help:      "Gauge describing the current number of remaining read API calls allowed for the tenant",
		},
		[]string{"tenant"},
	)

	// AzureAPITenantWriteRateLimitRemaining Gauge describing the current number of remaining write API calls
	AzureAPITenantWriteRateLimitRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure_api",
			Subsystem: "tenant",
			Name:      "write_rate_limit_remaining",
			Help:      "Gauge describing the current number of remaining write API calls allowed for the tenant",
		},
		[]string{"tenant"},
	)

	// AzureAPISubscriptionReadRateLimitRemaining Gauge describing the current number of remaining read API calls
	AzureAPISubscriptionReadRateLimitRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure_api",
			Subsystem: "subscription",
			Name:      "read_rate_limit_remaining",
			Help:      "Gauge describing the current number of remaining read API calls allowed for the subscription",
		},
		[]string{"subscription"},
	)

	// AzureAPISubscriptionReadRateLimitLastUpdateTime Gauge describing the current number of remaining read API calls
	AzureAPISubscriptionReadRateLimitLastUpdateTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure_api",
			Subsystem: "subscription",
			Name:      "read_rate_limit_remaining_last_update_time",
			Help:      "Time of the last update of azure_api_subscription_read_rate_limit_remaining",
		},
		[]string{"subscription"},
	)

	// AzureAPISubscriptionWriteRateLimitRemaining Gauge describing the current number of remaining write API calls
	AzureAPISubscriptionWriteRateLimitRemaining = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure_api",
			Subsystem: "subscription",
			Name:      "write_rate_limit_remaining",
			Help:      "Gauge describing the current number of remaining write API calls allowed for the subscription",
		},
		[]string{"subscription"},
	)

	// AzureAPISubscriptionWriteRateLimitLastUpdateTime Gauge describing the current number of remaining write API calls
	AzureAPISubscriptionWriteRateLimitLastUpdateTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure_api",
			Subsystem: "subscription",
			Name:      "write_rate_limit_remaining_last_update_time",
			Help:      "Time of the last update of azure_api_subscription_write_rate_limit_remaining",
		},
		[]string{"subscription"},
	)
)

var (
	azureAPIRateMutex = sync.RWMutex{}
)

func init() {
	prometheus.MustRegister(AzureAPICallsTotal)
	prometheus.MustRegister(AzureAPICallsFailedTotal)
	prometheus.MustRegister(AzureAPICallsDurationSecondsBuckets)
	prometheus.MustRegister(AzureAPITenantReadRateLimitRemaining)
	prometheus.MustRegister(AzureAPITenantWriteRateLimitRemaining)
	prometheus.MustRegister(AzureAPISubscriptionReadRateLimitRemaining)
	prometheus.MustRegister(AzureAPISubscriptionReadRateLimitLastUpdateTime)
	prometheus.MustRegister(AzureAPISubscriptionWriteRateLimitRemaining)
	prometheus.MustRegister(AzureAPISubscriptionWriteRateLimitLastUpdateTime)
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
func SetReadRateLimitRemaining(tenant string, subscription string, response *http.Response) {
	remaining := response.Header.Get("x-ms-ratelimit-remaining-tenant-reads")

	if len(remaining) > 0 {
		f, err := strconv.ParseFloat(remaining, 64)

		if err == nil {
			azureAPIRateMutex.Lock()
			AzureAPITenantReadRateLimitRemaining.WithLabelValues(tenant).Set(f)
			azureAPIRateMutex.Unlock()
		}
	}

	remaining = response.Header.Get("x-ms-ratelimit-remaining-subscription-reads")

	if len(remaining) > 0 {
		f, err := strconv.ParseFloat(remaining, 64)

		if err == nil {
			azureAPIRateMutex.Lock()
			AzureAPISubscriptionReadRateLimitRemaining.WithLabelValues(subscription).Set(f)
			AzureAPISubscriptionReadRateLimitLastUpdateTime.WithLabelValues(subscription).Set(float64(time.Now().Unix()))
			azureAPIRateMutex.Unlock()
		}
	}
}

// SetWriteRateLimitRemaining ...
func SetWriteRateLimitRemaining(tenant string, subscription string, response *http.Response) {
	remaining := response.Header.Get("x-ms-ratelimit-remaining-tenant-writes")

	if len(remaining) > 0 {
		f, err := strconv.ParseFloat(remaining, 64)

		if err == nil {
			azureAPIRateMutex.Lock()
			AzureAPITenantWriteRateLimitRemaining.WithLabelValues(tenant).Set(f)
			azureAPIRateMutex.Unlock()
		}
	}

	remaining = response.Header.Get("x-ms-ratelimit-remaining-subscription-writes")

	if len(remaining) > 0 {
		f, err := strconv.ParseFloat(remaining, 64)

		if err == nil {
			azureAPIRateMutex.Lock()
			AzureAPISubscriptionWriteRateLimitRemaining.WithLabelValues(subscription).Set(f)
			AzureAPISubscriptionWriteRateLimitLastUpdateTime.WithLabelValues(subscription).Set(float64(time.Now().Unix()))
			azureAPIRateMutex.Unlock()
		}
	}
}
