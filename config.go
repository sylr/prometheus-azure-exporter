package main

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/config"
	"github.com/sylr/prometheus-azure-exporter/pkg/metrics"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools/cache"
)

func setConfig() error {
	logger := log.WithFields(log.Fields{
		"_id": "00000000",
	})

	config.ParseOptions()
	conf, err := config.ParseConfigFile()

	if err != nil {
		logger.Errorf("Configuration not applied because parsing of config file failed: %s", err)
		return err
	}

	errs := config.ValidateConfig(conf)

	if len(errs) > 0 {
		for _, err := range errs {
			logger.Error(err)
		}

		err := errors.New("Configuration not applied because error(s) have been found")
		logger.Error(err)
		return err
	}

	// Apply configuration
	err = applyConfig(conf)

	return err
}

func applyConfig(conf *config.PrometheusAzureExporterConfig) error {
	config.CurrentConfig = conf

	// Update logging level
	switch {
	case len(config.CurrentConfig.Verbose) >= 1:
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	// JSON formatter
	if config.CurrentConfig.JSONOutput {
		log.SetFormatter(&log.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		})
	}

	// Turn on Noop caching
	if config.CurrentConfig.NoCache {
		cache.NoopCaching = true
	}

	// Update metrics functions interval
	needCancel := false
	for _, v := range config.CurrentConfig.UpdateMetricsFunctions {
		interval := metrics.GetUpdateMetricsFunctionInterval(v.Name)

		if v.Interval == time.Duration(0) {
			// New interval nil or unset.
			metrics.UnregisterUpdateMetricsFunctions(v.Name)
			needCancel = true
		} else if interval != nil && *interval != v.Interval {
			// Current interval set and different from new interval
			f := metrics.UnregisterUpdateMetricsFunctions(v.Name)
			metrics.RegisterUpdateMetricsFunctionWithInterval(v.Name, f, v.Interval)
			needCancel = true
		} else if interval == nil && v.Interval != time.Duration(0) {
			// Current interval nil and new interval set
			f := metrics.GetUpdateMetricsFunction(v.Name)
			metrics.RegisterUpdateMetricsFunctionWithInterval(v.Name, f, v.Interval)
			needCancel = true
		}
	}

	if needCancel {
		metrics.CancelUpdateMetricsFunctions()
	}

	return nil
}

func watchConfigFile() {
	logger := log.WithFields(log.Fields{
		"_id": "00000000",
	})

	watcher, err := fsnotify.NewWatcher()
	err = watcher.Add(config.CurrentConfig.ConfigFile)

	if err != nil {
		log.Fatal(err)
	}

	if len(os.Getenv("KUBERNETES_PORT")) > 0 {
		dir := filepath.Dir(config.CurrentConfig.ConfigFile)
		logger.Infof("In kubernetes context, adding %s to watch list", dir)
		watcher.Add(dir)
	}

	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				logger.Error("fsnotify: error")
				break
			}

			logger.Debugf("fsnotify: %s -> %s", event.Name, event.Op.String())

			if event.Op&fsnotify.Write == fsnotify.Write {
				if event.Name == config.CurrentConfig.ConfigFile {
					logger.Debugf("config: file changed")
				}
			} else if event.Op&fsnotify.Create == fsnotify.Create {
				if event.Name == config.CurrentConfig.ConfigFile {
					logger.Debugf("config: file created")
				} else if filepath.Base(event.Name) == "..data" {
					logger.Debugf("config: configmap volume updated")
				} else {
					break
				}
			} else {
				break
			}

			logger.Info("config: reloading config")
			setConfig()
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			logger.Errorf("fsnotify: %s", err)
		}
	}
}
