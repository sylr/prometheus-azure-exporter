// vim: set tabstop=4 expandtab autoindent smartindent:

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/metrics"
	"github.com/sylr/prometheus-client-golang/prometheus/promhttp"
)

// PrometheusAzureExporterOptions options
type PrometheusAzureExporterOptions struct {
	Verbose          []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	Version          bool   `          long:"version" description:"Show version"`
	ListeningAddress string `short:"a" long:"address" description:"Listening address" env:"LISTENING_ADDRESS" default:"0.0.0.0"`
	ListeningPort    uint   `short:"p" long:"port" description:"Listening port" env:"LISTENING_PORT" default:"9000"`
}

var (
	opts    = PrometheusAzureExporterOptions{}
	parser  = flags.NewParser(&opts, flags.Default)
	version = "v0.0.1"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
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
	case len(opts.Verbose) >= 1:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// Loggin options
	log.Debugf("Options: %+v", opts)
	log.Infof("Version: %s", version)

	ctx := context.Background()

	metrics.UpdateMetrics(ctx)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("0.0.0.0:9000", nil))
}
