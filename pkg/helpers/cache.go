package helpers

import (
	"strings"
	"sync"
	"time"
)

// HelperRun represents a single run of a helper
type HelperRun struct {
	output  string
	runTime time.Time
}

// Output returns the output of the run
func (hr HelperRun) Output() string {
	return hr.output
}

// Time returns the time of the run
func (hr HelperRun) Time() time.Time {
	return hr.runTime
}

// HelperRunCache stores the last run of all unique helpers.
// A helper is defined by its name and config
type HelperRunCache struct {
	lock  sync.Mutex
	store map[string]HelperRun
}

// NewHelperRunCache creates a new helper run cache
func NewHelperRunCache() HelperRunCache {
	return HelperRunCache{
		lock:  sync.Mutex{},
		store: map[string]HelperRun{},
	}
}

// GetOrAdd tries to find an item in the cache and creates one if not found
func (c *HelperRunCache) GetOrAdd(key string) HelperRun {
	c.lock.Lock()
	defer c.lock.Unlock()

	item, found := c.store[key]
	if !found {
		item = HelperRun{
			runTime: time.Now(),
		}
		c.store[key] = item
	}

	return item
}

// Put updates an existing item in the cache or creates one if not found
func (c *HelperRunCache) Put(key, output string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	item, found := c.store[key]
	if !found {
		item = HelperRun{}
	}

	item.output = output
	item.runTime = time.Now()
	c.store[key] = item
}

// Key generates a cache key from a Helper's attributes
func (c *HelperRunCache) Key(parts ...string) string {
	return strings.Join(parts, ":")
}
