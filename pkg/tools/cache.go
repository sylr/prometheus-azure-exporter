package tools

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var (
	caches = make(map[time.Duration]*gocache.Cache)
)

// GetCache get a caching object
func GetCache(duration time.Duration) *gocache.Cache {
	if _, ok := caches[duration]; !ok {
		caches[duration] = gocache.New(duration, 2*time.Minute)
	}

	return caches[duration]
}
