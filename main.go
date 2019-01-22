// vim: set tabstop=4 expandtab autoindent smartindent:

package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/metrics"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools"
)

// PrometheusAzureExporterOptions options
type PrometheusAzureExporterOptions struct {
	Verbose          []bool `short:"v"  long:"verbose"   description:"Show verbose debug information"`
	JSONOutput       bool   `short:"j"  long:"json"      description:"Use json format for output"`
	Version          bool   `           long:"version"   description:"Show version"`
	ListeningAddress string `short:"a"  long:"address"   description:"Listening address" env:"LISTENING_ADDRESS" default:"0.0.0.0"`
	ListeningPort    uint   `short:"p"  long:"port"      description:"Listening port" env:"LISTENING_PORT" default:"9000"`
	UpdateInterval   int    `short:"i"  long:"interval"  description:"Number of seconds between metrics updates" default:"60"`
	NoCache          bool   `           long:"no-cache"  description:"Disable internal caching"`

	// Env vars used for Azure Authent, see
	// https://github.com/Azure/go-autorest/blob/master/autorest/azure/auth/auth.go#L86-L94
	AzureTenantID            string `env:"AZURE_TENANT_ID" description:"Azure tenant id" required:"true"`
	AzureSubscriptionID      string `env:"AZURE_SUBSCRIPTION_ID" description:"Azure subscription id" required:"true"`
	AzureClientID            string `env:"AZURE_CLIENT_ID" description:"Azure client id"`
	AzureClientSecret        string `env:"AZURE_CLIENT_SECRET" description:"Azure client secret"`
	AzureCertificatePath     string `env:"AZURE_CERTIFICATE_PATH" description:"Azure certficate path"`
	AzureCertificatePassword string `env:"AZURE_CERTIFICATE_PASSWORD" description:"Azure certficate password"`
	AzureUsername            string `env:"AZURE_USERNAME" description:"Azure username"`
	AzurePassword            string `env:"AZURE_PASSWORD" description:"Azure password"`
	AzureEnvironment         string `env:"AZURE_ENVIRONMENT" description:"Azure environment"`
	AzureADResource          string `env:"AZURE_AD_RESOURCE" description:"Azure AD resource"`
}

var (
	// Options daemon options
	Options = PrometheusAzureExporterOptions{}
	// Version daemon version
	Version = "v0.1.2"
	// Go Version
	GoVersion = runtime.Version()
	// parser
	parser = flags.NewParser(&Options, flags.Default)
)

var (
	azureExporterBuildInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure_exporter",
			Subsystem: "",
			Name:      "build_info",
			Help:      "Prometheus azure exporter build info",
		},
		[]string{"version", "goversion"},
	)
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{
		DisableColors:  true,
		DisableSorting: false,
		SortingFunc:    tools.SortLogKeys,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the info severity or above.
	log.SetLevel(log.InfoLevel)

	// Register build info
	prometheus.MustRegister(azureExporterBuildInfo)
}

// main
func main() {
	// looping for --version in args
	for _, val := range os.Args {
		if val == "--version" {
			fmt.Printf("prometheus-azure-exporter version %s\n", Version)
			os.Exit(0)
		} else if val == "--" {
			break
		}
	}

	// parse flags
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	// Update logging level
	switch {
	case len(Options.Verbose) >= 1:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// JSON formatter
	if Options.JSONOutput {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	}

	// Loggin options
	log.Debugf("Options: %+v", Options)
	log.Infof("Version: %s", Version)

	// Turn on Noop caching
	if Options.NoCache {
		tools.NoopCaching = true
	}

	// Set build info
	azureExporterBuildInfo.WithLabelValues(Version, GoVersion).Set(1)

	// Configure metrics update interval
	metrics.SetUpdateMetricsInterval(time.Duration(Options.UpdateInterval) * time.Second)

	// Update metrics process
	ctx := context.Background()
	go metrics.UpdateMetrics(ctx)

	// Prometheus http endpoint
	listeningAddress := fmt.Sprintf("%s:%d", Options.ListeningAddress, Options.ListeningPort)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(listeningAddress, nil)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
