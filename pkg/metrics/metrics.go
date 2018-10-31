package metrics

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	updateMetricsInterval  = 30
	updateMetricsFunctions = make(map[string]func(context.Context))
)

// RegisterUpdateMetricsFunctions allows you to register a function
// that will update prometheus metrics
func RegisterUpdateMetricsFunctions(name string, f func(context.Context)) {
	updateMetricsFunctions[name] = f
}

// SetUpdateMetricsInterval sets interval
func SetUpdateMetricsInterval(interval int) {
	updateMetricsInterval = interval
}

// UpdateMetrics main update metrics process
// This process loops forever so it needs to be detached
func UpdateMetrics(ctx context.Context) {
	// Aligning udate metric processes with minute start
	sec := 60 - (time.Now().Unix() % 60)
	nsec := time.Now().UnixNano() - (time.Now().Unix() * 1000000000)
	log.WithField("_id", "000000000000").Debugf("Waiting %d seconds before starting to update metrics (ns: %d)", sec, nsec)
	time.Sleep(time.Second*time.Duration(sec) - time.Duration(nsec))

	ticker := time.NewTicker(time.Duration(updateMetricsInterval) * time.Second)
	t := time.Now()

	for {
		log.WithField("_id", "000000000000").Debugf("Start metrics update: %s", t)

		// Loop over all update functions metrics
		for updateMetricsFuncName, updateMetricsFunc := range updateMetricsFunctions {
			// We detach the update process so if it takes more than the refresh
			// time it does not get blocked
			go func(ctx context.Context, updateMetricsFuncName string, updateMetricsFunc func(context.Context), t time.Time) {
				id := hashTime(t)
				fields := log.Fields{
					"_id":      id,
					"function": updateMetricsFuncName,
				}

				ctx = context.WithValue(ctx, "id", id)

				log.WithFields(fields).Debug("Start update metrics function")
				updateMetricsFunc(ctx)
				log.WithFields(fields).Debug("End update metrics function")
			}(ctx, updateMetricsFuncName, updateMetricsFunc, t)
		}

		t = <-ticker.C
	}
}

func hashTime(t time.Time) string {
	h := md5.New()
	io.WriteString(h, t.String())
	s := fmt.Sprintf("%x", h.Sum(nil))

	return s[0:12]
}
