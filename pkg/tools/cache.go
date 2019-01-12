package tools

import (
	"sync"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	// NoopCaching Use a Noop caching implementation
	NoopCaching = false
)

var (
	mutex  = sync.RWMutex{}
	caches = make(map[time.Duration]cache.Cacher)
)

// GetCache get a caching object
func GetCache(duration time.Duration) cache.Cacher {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := caches[duration]; !ok {
		if NoopCaching {
			caches[duration] = cache.NewNoopCacher(duration, 2*time.Minute)
		} else {
			caches[duration] = cache.NewCacher(duration, 2*time.Minute)
		}
	}

	return caches[duration]
}
