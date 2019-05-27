package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/config"
	"github.com/sylr/prometheus-azure-exporter/pkg/metrics"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools"
)

var (
	version   = "v0.5.0"
	goVersion = runtime.Version()
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
			fmt.Printf("prometheus-azure-exporter version %s\n", version)
			os.Exit(0)
		} else if val == "--" {
			break
		}
	}

	// Configuration
	err := setConfig()
	if err != nil {
		os.Exit(1)
	}
	go watchConfigFile()

	// Log options
	log.Debugf("Options: %+v", config.CurrentConfig)
	log.Infof("Version: %s", version)

	// Set build info
	azureExporterBuildInfo.WithLabelValues(version, goVersion).Set(1)

	// Configure metrics update interval
	metrics.SetDefaultUpdateMetricsInterval(config.CurrentConfig.UpdateInterval)

	// Update metrics process
	ctx := context.Background()
	go metrics.UpdateMetrics(ctx)

	// Prometheus http endpoint
	listeningAddress := fmt.Sprintf("%s:%d", config.CurrentConfig.ListeningAddress, config.CurrentConfig.ListeningPort)
	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(listeningAddress, nil)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
