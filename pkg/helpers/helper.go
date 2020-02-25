package helpers

import (
	"fmt"
	"sync"
	"time"
)

type Helper interface {
	Run(config string) string
	UpdateInterval() time.Duration
}

var regLock sync.Mutex
var registry = map[string]Helper{}

func Register(name string, helper Helper) {
	regLock.Lock()
	defer regLock.Unlock()
	if _, ok := registry[name]; ok {
		panic(fmt.Sprintf("Helper already exists with name '%s'", name))
	}
	registry[name] = helper
}

var ErrHelperNotFound = fmt.Errorf("helper not found")

var cacheLock sync.Mutex
var cache = map[string]helperRun{}

type helperRun struct {
	output  string
	runTime time.Time
}

func Run(name, config string) (string, error) {
	regLock.Lock()
	defer regLock.Unlock()
	helper, ok := registry[name]
	if !ok {
		return "", ErrHelperNotFound
	}

	cacheLock.Lock()
	defer cacheLock.Unlock()

	if lastRun, ok := cache[name]; ok {
		if time.Since(lastRun.runTime) < helper.UpdateInterval() {
			return lastRun.output, nil
		}
	}

	output := helper.Run(config)
	cache[name] = helperRun{
		output:  output,
		runTime: time.Now(),
	}
	return output, nil
}
