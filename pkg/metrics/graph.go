package metrics

import (
	"context"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/azure"
)

var (
	nameSanitationRegexp = regexp.MustCompile("[^a-zA-z0-9_./*-+ ]")
)

var (
	graphApplicationKeyExpire      = newGraphApplicationKeyExpire()
	graphApplicationPasswordExpire = newGraphApplicationPasswordExpire()
)

// -----------------------------------------------------------------------------

func newGraphApplicationKeyExpire() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "graph",
			Name:      "application_key_expire_time",
			Help:      "Unix timestamp of application key expiration",
		},
		[]string{"application", "key"},
	)
}

func newGraphApplicationPasswordExpire() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "graph",
			Name:      "application_password_expire_time",
			Help:      "Unix timestamp of application password expiration",
		},
		[]string{"application", "password"},
	)
}

// -----------------------------------------------------------------------------

func init() {
	prometheus.MustRegister(graphApplicationKeyExpire)
	prometheus.MustRegister(graphApplicationPasswordExpire)

	if GetUpdateMetricsFunctionInterval("graph") == nil {
		RegisterUpdateMetricsFunctionWithInterval("graph", UpdateGraphMetrics, 60*time.Second)
	}
}

// UpdateGraphMetrics updates graph metrics
func UpdateGraphMetrics(ctx context.Context) error {
	var err error

	contextLogger := log.WithFields(log.Fields{
		"_id":   ctx.Value("id").(string),
		"_func": "UpdateGraphMetrics",
	})

	// Create new metric vectors
	nextGraphApplicationKeyExpire := newGraphApplicationKeyExpire()
	nextGraphApplicationPasswordExpire := newGraphApplicationPasswordExpire()

	// <!-- APPLICATIONS -------------------------------------------------------
	azureClients := azure.NewAzureClients()
	applications, err := azure.ListApplications(ctx, azureClients)

	if err != nil {
		contextLogger.Errorf("Unable to list applications: %s", err)
		return err
	}

	for _, app := range *applications {
		for _, key := range *app.KeyCredentials {
			var decodedName string
			if key.CustomKeyIdentifier != nil {
				decodedName = string(*key.CustomKeyIdentifier)
				decodedName = nameSanitationRegexp.ReplaceAllString(decodedName, "")
			} else {
				decodedName = *key.KeyID
			}

			nextGraphApplicationKeyExpire.WithLabelValues(*app.DisplayName, decodedName).Set(float64(key.EndDate.Unix()))
		}

		for _, password := range *app.PasswordCredentials {
			var decodedName string
			if password.CustomKeyIdentifier != nil {
				decodedName = string(*password.CustomKeyIdentifier)
				decodedName = nameSanitationRegexp.ReplaceAllString(decodedName, "")
			} else {
				decodedName = *password.KeyID
			}

			nextGraphApplicationPasswordExpire.WithLabelValues(*app.DisplayName, decodedName).Set(float64(password.EndDate.Unix()))
		}
	}
	// -- APPLICATIONS -------------------------------------------------------!>

	// swapping current registered metrics with updated copies
	*graphApplicationKeyExpire = *nextGraphApplicationKeyExpire
	*graphApplicationPasswordExpire = *nextGraphApplicationPasswordExpire

	return err
}
