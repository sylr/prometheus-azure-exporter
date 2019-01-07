package tools

import (
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var (
	mutex  = sync.RWMutex{}
	caches = make(map[time.Duration]*gocache.Cache)
)

// GetCache get a cache object
func GetCache(duration time.Duration) *gocache.Cache {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := caches[duration]; !ok {
		caches[duration] = gocache.New(duration, 2*time.Minute)
	}

	return caches[duration]
}
