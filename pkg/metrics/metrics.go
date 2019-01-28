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

	updateMetricsFunctionIntervalDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "azure_exporter",
			Subsystem: "update_metrics_function",
			Name:      "interval_duration_seconds",
			Help:      "Interval of update metrics functions",
		},
		[]string{"function"},
	)
)

func init() {
	prometheus.MustRegister(updateMetricsFunctionDurationHistogram)
	prometheus.MustRegister(updateMetricsFunctionLastDurationGauge)
	prometheus.MustRegister(updateMetricsFunctionIntervalDurationGauge)
}

var (
	// Mutex used to lock read/writes of updateMetricsFunctions/intervalUpdateMetricsFunctions.
	mutex = sync.RWMutex{}
	// Default update interval.
	defaultUpdateMetricsInterval = 30 * time.Second
	// This var holds all update functions that have been registered once.
	// It is used when we want to move a function from an update interval
	// to another.
	updateMetricsFunctions = make(map[string]UpdateMetricsFunction)
	// This var holds all update functions that are currently registered
	// and tied to an update interval.
	intervalUpdateMetricsFunctions = make(map[time.Duration]map[string]UpdateMetricsFunction)
	// This var holds all the cancel functions of the contexts used by
	// the interval processes.
	intervalCancelFunctions = make(map[time.Duration]context.CancelFunc)
)

// UpdateMetricsFunction is the function type which needs to respected to
// create update metrics functions.
type UpdateMetricsFunction func(context.Context) error

// initUpdateMetricsFunctionsMap makes sure the map is initialized.
func initUpdateMetricsFunctionsMap(interval time.Duration) {
	if intervalUpdateMetricsFunctions[interval] == nil {
		intervalUpdateMetricsFunctions[interval] = make(map[string]UpdateMetricsFunction)
	}
}

// RegisterUpdateMetricsFunction allows you to register a function
// that will update prometheus metrics.
func RegisterUpdateMetricsFunction(name string, f UpdateMetricsFunction) {
	RegisterUpdateMetricsFunctionWithInterval(name, f, defaultUpdateMetricsInterval)
}

// UnregisterUpdateMetricsFunctions allows you to unregister an update functions metrics
// which has been previously registered.
func UnregisterUpdateMetricsFunctions(name string) UpdateMetricsFunction {
	mutex.Lock()
	defer mutex.Unlock()

	initUpdateMetricsFunctionsMap(defaultUpdateMetricsInterval)

	for interval := range intervalUpdateMetricsFunctions {
		if f, ok := intervalUpdateMetricsFunctions[interval][name]; ok {
			delete(intervalUpdateMetricsFunctions[interval], name)
			return f
		}
	}

	return nil
}

// RegisterUpdateMetricsFunctionWithInterval allows you to register a function
// that will update prometheus metrics every interval.
func RegisterUpdateMetricsFunctionWithInterval(name string, f UpdateMetricsFunction, interval time.Duration) {
	mutex.Lock()
	defer mutex.Unlock()

	initUpdateMetricsFunctionsMap(interval)

	if _, ok := updateMetricsFunctions[name]; !ok {
		updateMetricsFunctions[name] = f
	}

	intervalUpdateMetricsFunctions[interval][name] = f
}

// GetUpdateMetricsFunction returns the update metrics function associated to `name`.
// It will only return a result if the function has previously been registered once.
// It does not matter if the function has been un-registered.
func GetUpdateMetricsFunction(name string) UpdateMetricsFunction {
	if f, ok := updateMetricsFunctions[name]; ok {
		return f
	}

	return nil
}

// GetUpdateMetricsFunctionInterval returns the interval the update metrics is
// currently registered at.
func GetUpdateMetricsFunctionInterval(name string) *time.Duration {
	initUpdateMetricsFunctionsMap(defaultUpdateMetricsInterval)

	for interval := range intervalUpdateMetricsFunctions {
		if _, ok := intervalUpdateMetricsFunctions[interval][name]; ok {
			return &interval
		}
	}

	return nil
}

// SetDefaultUpdateMetricsInterval sets default interval
func SetDefaultUpdateMetricsInterval(interval time.Duration) {
	defaultUpdateMetricsInterval = interval
}

// GetDefaultUpdateMetricsInterval gets default interval
func GetDefaultUpdateMetricsInterval() time.Duration {
	return defaultUpdateMetricsInterval
}

// UpdateMetrics main update metrics process. It spawns goroutines of
// updateMetricsWithInterval() which are responsible for running the
// update metrics functions every desired update intervals.
// This method loops forever so it needs to be detached.
func UpdateMetrics(ctx context.Context) {
	wg := sync.WaitGroup{}

	for {
		if len(intervalUpdateMetricsFunctions) == 0 {
			time.Sleep(time.Second)
		}

		mutex.RLock()
		for interval := range intervalUpdateMetricsFunctions {
			if len(intervalUpdateMetricsFunctions[interval]) == 0 {
				continue
			}

			ctx, f := context.WithCancel(ctx)
			intervalCancelFunctions[interval] = f
			go updateMetricsWithInterval(ctx, &wg, interval)

			wg.Add(1)
		}
		mutex.RUnlock()

		wg.Wait()
	}
}

// CancelUpdateMetricsFunctions calls the contexts cancel() methods
// of all interval processes.
func CancelUpdateMetricsFunctions() {
	for _, f := range intervalCancelFunctions {
		f()
	}
}

// updateMetricsWithInterval is the method used to run update metrics functions
// at given intervals. It is spawned as goroutines by UpdateMetrics(), one for each interval.
func updateMetricsWithInterval(ctx context.Context, wg *sync.WaitGroup, interval time.Duration) {
	var t time.Time
	var ticker *time.Ticker
	// logger
	processLogger := log.WithFields(log.Fields{
		"_id":       "00000000",
		"_interval": interval,
	})

	processLogger.Infof("Start interval update metrics process: %s", interval)

	// Aligning update metric processes with minute start
	sec := int64(interval/time.Second) - (time.Now().Unix() % int64(interval/time.Second))
	nsec := time.Now().UnixNano() - (time.Now().Unix() * int64(time.Second))
	wait := time.Duration(sec)*time.Second - time.Duration(nsec)
	waiter := time.NewTicker(wait)
	processLogger.Infof("Waiting before starting to update metrics: %s", wait.Round(time.Second))

	// Wait for time sync or cancelation of context (reload).
	select {
	case <-waiter.C:
	case <-ctx.Done():
		processLogger.Infof("Interval process context has been canceled during initial time sync")
		goto done
	}

	ticker = time.NewTicker(interval)
	t = time.Now()

	for {
		// Loop over all update metrics functions
		for updateMetricsFuncName, updateMetricsFunc := range intervalUpdateMetricsFunctions[interval] {
			updateMetricsFunctionIntervalDurationGauge.WithLabelValues(updateMetricsFuncName).Set(float64(interval.Seconds()))

			// We detach the update process so that if it takes more than the refresh
			// time it does not get blocked
			go func(ctx context.Context, updateMetricsFuncName string, updateMetricsFunc UpdateMetricsFunction, t time.Time) {
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

				functionLogger.Infof("End update metrics function in %v", t1.Round(time.Millisecond))
			}(ctx, updateMetricsFuncName, updateMetricsFunc, t)
		}

		// wait for ticker or cancelation of context (reload).
		select {
		case t = <-ticker.C:
		case <-ctx.Done():
			processLogger.Infof("Interval process context has been canceled during waiting")
			goto done
		}
	}

done:
	wg.Done()
	return
}

// processHash generates a hash based on time and salt to be used
// as id in the logger.
func processHash(t time.Time, salt string) string {
	h := md5.New()
	io.WriteString(h, salt+":"+t.String())
	s := fmt.Sprintf("%x", h.Sum(nil))

	return s[0:8]
}
