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
	graphApplicationKeyExpire = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "graph",
			Name:      "application_key_expire_time",
			Help:      "Unix timestamp of application key expiration",
		},
		[]string{"application", "key"},
	)

	graphApplicationPasswordExpire = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure",
			Subsystem: "graph",
			Name:      "application_password_expire_time",
			Help:      "Unix timestamp of application password expiration",
		},
		[]string{"application", "password"},
	)
)

var (
	nameSanitationRegexp = regexp.MustCompile("[^a-zA-z0-9_./*-+ ]")
)

func init() {
	prometheus.MustRegister(graphApplicationKeyExpire)
	prometheus.MustRegister(graphApplicationPasswordExpire)

	RegisterUpdateMetricsFunctionsWithInterval("graph", UpdateGraphMetrics, 60*time.Second)
}

// UpdateGraphMetrics updates graph metrics
func UpdateGraphMetrics(ctx context.Context) error {
	contextLogger := log.WithFields(log.Fields{
		"_id":   ctx.Value("id").(string),
		"_func": "UpdateGraphMetrics",
	})
	azureClients := azure.NewAzureClients()

	// <!-- APPLICATIONS -------------------------------------------------------
	applications, err := azure.ListApplications(ctx, azureClients)

	if err != nil {
		contextLogger.Errorf("Unable to list applications: %s", err)

		graphApplicationKeyExpire.Reset()
		graphApplicationPasswordExpire.Reset()

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

			graphApplicationKeyExpire.WithLabelValues(*app.DisplayName, decodedName).Set(float64(key.EndDate.Unix()))
		}

		for _, password := range *app.PasswordCredentials {
			var decodedName string
			if password.CustomKeyIdentifier != nil {
				decodedName = string(*password.CustomKeyIdentifier)
				decodedName = nameSanitationRegexp.ReplaceAllString(decodedName, "")
			} else {
				decodedName = *password.KeyID
			}

			graphApplicationPasswordExpire.WithLabelValues(*app.DisplayName, decodedName).Set(float64(password.EndDate.Unix()))
		}
	}
	// -- APPLICATIONS -------------------------------------------------------!>

	return err
}
