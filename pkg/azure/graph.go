package azure

import (
	"context"
	"os"
	"time"

	graph "github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools/cache"
)

var (
	// AzureAPIGraphCallsTotal Total number of Azure Graph API calls
	AzureAPIGraphCallsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "graph",
			Name:      "calls_total",
			Help:      "Total number of calls to the Azure Graph API",
		},
		[]string{},
	)

	// AzureAPIGraphCallsFailedTotal Total number of failed Azure Graph API calls
	AzureAPIGraphCallsFailedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "azure_api",
			Subsystem: "graph",
			Name:      "calls_failed_total",
			Help:      "Total number of failed calls to the Azure Graph API",
		},
		[]string{},
	)

	// AzureAPIGraphCallsDurationSecondsBuckets Histograms of Azure Graph API calls durations in seconds
	AzureAPIGraphCallsDurationSecondsBuckets = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "azure_api",
			Subsystem: "graph",
			Name:      "calls_duration_seconds",
			Help:      "Histograms of Azure Graph API calls durations in seconds",
			Buckets:   []float64{0.50, 0.75, 1.0, 1.25, 1.5, 2.0},
		},
		[]string{},
	)
)

func init() {
	prometheus.MustRegister(AzureAPIGraphCallsTotal)
	prometheus.MustRegister(AzureAPIGraphCallsFailedTotal)
	prometheus.MustRegister(AzureAPIGraphCallsDurationSecondsBuckets)
}

// ObserveAzureGraphAPICall ...
func ObserveAzureGraphAPICall(duration float64, labels ...string) {
	AzureAPIGraphCallsTotal.WithLabelValues(labels...).Inc()
	AzureAPIGraphCallsDurationSecondsBuckets.WithLabelValues(labels...).Observe(duration)
}

// ObserveAzureGraphAPICallFailed ...
func ObserveAzureGraphAPICallFailed(duration float64, labels ...string) {
	AzureAPIGraphCallsFailedTotal.WithLabelValues(labels...).Inc()
}

// ListApplications list applications
func ListApplications(ctx context.Context, clients *AzureClients) (*[]graph.Application, error) {
	c := cache.GetCache(5 * time.Minute)

	contextLogger := log.WithFields(log.Fields{
		"_id": ctx.Value("id").(string),
	})

	cacheKey := os.Getenv("AZURE_TENANT_ID") + "-applications"

	if capplications, ok := c.Get(cacheKey); ok {
		if apps, ok := capplications.(*[]graph.Application); !ok {
			contextLogger.Errorf("Failed to cast object from cache back to *[]graph.Application")
		} else {
			return apps, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetApplicationsClient(os.Getenv("AZURE_TENANT_ID"))

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	apps, err := client.List(ctx, "")
	t1 := time.Since(t0).Seconds()

	if err != nil {
		if ctx.Err() != context.Canceled {
			ObserveAzureAPICallFailed(t1)
			ObserveAzureGraphAPICallFailed(t1)
		}
		return nil, err
	}

	ObserveAzureAPICall(t1)
	ObserveAzureGraphAPICall(t1)

	vals := apps.Values()
	c.SetDefault(cacheKey, &vals)

	return &vals, nil
}
