package metrics

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	updateMetricsFunctionDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "azure_exporter",
			Subsystem: "update_metrics_function",
			Name:      "duration_seconds",
			Help:      "Duration of update metrics functions (does not include run which returned an error)",
			Buckets:   []float64{1, 5, 10, 30, 60, 300, 600, 1800, 3600},
		},
		[]string{"function"},
	)

	updateMetricsFunctionLastDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure_exporter",
			Subsystem: "update_metrics_function",
			Name:      "last_duration_seconds",
			Help:      "Last duration of update metrics functions (does not include run which returned an error)",
		},
		[]string{"function"},
	)
)

var (
	mutex                  = sync.RWMutex{}
	updateMetricsInterval  = 30 * time.Second
	updateMetricsFunctions = make(map[time.Duration]map[string]func(context.Context) error)
)

func init() {
	prometheus.MustRegister(updateMetricsFunctionDurationHistogram)
	prometheus.MustRegister(updateMetricsFunctionLastDurationGauge)
}

// initUpdateMetricsFunctionsMap makes sure the map is initialized
func initUpdateMetricsFunctionsMap(interval time.Duration) {
	if updateMetricsFunctions[interval] == nil {
		updateMetricsFunctions[interval] = make(map[string]func(context.Context) error)
	}
}

// RegisterUpdateMetricsFunctions allows you to register a function
// that will update prometheus metrics
func RegisterUpdateMetricsFunctions(name string, f func(context.Context) error) {
	mutex.Lock()
	defer mutex.Unlock()

	initUpdateMetricsFunctionsMap(updateMetricsInterval)
	updateMetricsFunctions[updateMetricsInterval][name] = f
}

// RegisterUpdateMetricsFunctionsWithInterval allows you to register a function
// that will update prometheus metrics every interval
func RegisterUpdateMetricsFunctionsWithInterval(name string, f func(context.Context) error, interval time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()

	initUpdateMetricsFunctionsMap(interval)
	updateMetricsFunctions[interval][name] = f
}

// SetUpdateMetricsInterval sets interval
func SetUpdateMetricsInterval(interval time.Duration) {
	updateMetricsInterval = interval
}

// GetUpdateMetricsInterval gets interval
func GetUpdateMetricsInterval() time.Duration {
	return updateMetricsInterval
}

// UpdateMetrics main update metrics process
// This process loops forever so it needs to be detached
func UpdateMetrics(ctx context.Context) {
	wg := sync.WaitGroup{}

	for interval, _ := range updateMetricsFunctions {
		go updateMetricsWithInterval(ctx, interval)
		wg.Add(1)
	}

	wg.Wait()
}

func updateMetricsWithInterval(ctx context.Context, interval time.Duration) {
	// logger
	processLogger := log.WithFields(log.Fields{
		"_id":       "00000000",
		"_interval": interval,
	})

	processLogger.Infof("Start update metrics process with interval: %v", interval)

	// Aligning update metric processes with minute start
	sec := 60 - (time.Now().Unix() % 60)
	nsec := time.Now().UnixNano() - (time.Now().Unix() * int64(time.Second))
	processLogger.Debugf("Waiting %d seconds before starting to update metrics (ns: %d)", sec, nsec)
	time.Sleep(time.Second*time.Duration(sec) - time.Duration(nsec))

	ticker := time.NewTicker(interval)
	t := time.Now()

	for {
		// Loop over all update metrics functions
		for updateMetricsFuncName, updateMetricsFunc := range updateMetricsFunctions[interval] {
			// We detach the update process so if it takes more than the refresh
			// time it does not get blocked
			go func(ctx context.Context, updateMetricsFuncName string, updateMetricsFunc func(context.Context) error, t time.Time) {
				id := processHash(t, updateMetricsFuncName)
				functionLogger := processLogger.WithFields(log.Fields{
					"_id":       id,
					"_interval": interval,
					"_func":     updateMetricsFuncName,
				})

				ctx = context.WithValue(ctx, "id", id)

				functionLogger.Infof("Start update metrics function")

				// Run update metrics function
				t0 := time.Now()
				err := updateMetricsFunc(ctx)
				t1 := time.Since(t0)

				// metrics
				if err == nil {
					updateMetricsFunctionDurationHistogram.WithLabelValues(updateMetricsFuncName).Observe(t1.Seconds())
					updateMetricsFunctionLastDurationGauge.WithLabelValues(updateMetricsFuncName).Set(t1.Seconds())
				}

				functionLogger.Infof("End update metrics function in %v", t1)
			}(ctx, updateMetricsFuncName, updateMetricsFunc, t)
		}

		t = <-ticker.C
	}
}

func processHash(t time.Time, salt string) string {
	h := md5.New()
	io.WriteString(h, salt+":"+t.String())
	s := fmt.Sprintf("%x", h.Sum(nil))

	return s[0:8]
}
