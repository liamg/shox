package helpers

import (
	"fmt"
	"sync"
	"time"
)

// Helper is a dynamic operation which creates output for the status bar
type Helper interface {
	Run(config string) string      // Run starts the helper and returns the output
	UpdateInterval() time.Duration // UpdateInterval returns the minimum time period before the helper should run again
}

var regLock = sync.Mutex{}
var registry = map[string]Helper{}
var runCache = NewHelperRunCache()

// Register registers a new helper by name
func Register(name string, helper Helper) {
	regLock.Lock()
	defer regLock.Unlock()
	if _, ok := registry[name]; ok {
		panic(fmt.Sprintf("Helper already exists with name '%s'", name))
	}
	registry[name] = helper
}

// ErrHelperNotFound means no helper exists by the specified name
var ErrHelperNotFound = fmt.Errorf("helper not found")

// Run executes a helper with the provided config string
func Run(name, config string) (string, error) {
	regLock.Lock()
	defer regLock.Unlock()
	helper, ok := registry[name]
	if !ok {
		return "", ErrHelperNotFound
	}

	lastRun := runCache.GetOrAdd(runCache.Key(name, config))
	if time.Since(lastRun.Time()) < helper.UpdateInterval() {
		return lastRun.Output(), nil
	}

	output := helper.Run(config)
	runCache.Put(runCache.Key(name, config), output)
	return output, nil
}
