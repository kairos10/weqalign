package main

import (
	"fmt"
	"sync"
	"time"
)

type imgResource struct {
	imgId             string
	path              string
	resourceTimestamp time.Time
	solverResponse    interface{}
}
type imageResources struct {
	sync.RWMutex
	mapIdResource map[string]*imgResource
	crtId         uint
}

func (ir *imageResources) getResource(key string) (*imgResource, bool) {
	ir.RLock()
	r, ok := ir.mapIdResource[key]
	ir.RUnlock()
	return r, ok
}
func (ir *imageResources) getLastResource() (*imgResource, bool) {
	ir.RLock()
	crtId := ir.crtId
	ir.RUnlock()
	key := fmt.Sprintf("%d", crtId)
	return ir.getResource(key)
}
func (ir *imageResources) getResourcesSince(pastDuration time.Duration) map[string]*imgResource {
	ret := make(map[string]*imgResource)
	now := time.Now()

	ir.RLock()
	for k, r := range ir.mapIdResource {
		if now.Sub(r.resourceTimestamp) < pastDuration {
			ret[k] = r
		}
	}
	ir.RUnlock()
	if len(ret) == 0 {
		ret = nil
	}
	return ret
}
func (ir *imageResources) addResource(filePath string) (key string) {
	ir.Lock()
	if ir.mapIdResource == nil {
		ir.mapIdResource = make(map[string]*imgResource)
	}

	ir.crtId++
	key = fmt.Sprintf("%d", ir.crtId)
	r := imgResource{imgId: key, path: filePath, resourceTimestamp: time.Now()}
	ir.mapIdResource[key] = &r
	ir.Unlock()
	return
}

var resources imageResources
