package main

import (
	"fmt"
	"sync"
	"time"
)

type imgResource struct {
	imgId             string
	path              string // path on the server
	webPath           string // path on the web server
	resourceTimestamp time.Time
	solverResponse    interface{}
}
type imageResources struct {
	sync.RWMutex
	mapIdResource map[uint64]*imgResource
	crtId         uint64
}

func (ir *imageResources) getResource(key uint64) (*imgResource, bool) {
	ir.RLock()
	r, ok := ir.mapIdResource[key]
	ir.RUnlock()
	return r, ok
}
func (ir *imageResources) getLastResourceId() uint64 {
	ir.RLock()
	crtId := ir.crtId
	ir.RUnlock()
	return crtId
}
func (ir *imageResources) getLastResource() (*imgResource, bool) {
	ir.RLock()
	key := ir.crtId
	ir.RUnlock()
	return ir.getResource(key)
}
func (ir *imageResources) getResourcesSince(pastDuration time.Duration) map[uint64]*imgResource {
	ret := make(map[uint64]*imgResource)
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
func (ir *imageResources) addResource(filePath string, webPath string) (key uint64) {
	ir.Lock()
	if ir.mapIdResource == nil {
		ir.mapIdResource = make(map[uint64]*imgResource)
	}

	ir.crtId++
	key = ir.crtId
	sKey := fmt.Sprintf("%d", ir.crtId)
	r := imgResource{imgId: sKey, path: filePath, webPath: webPath, resourceTimestamp: time.Now()}
	ir.mapIdResource[key] = &r
	ir.Unlock()
	return
}

var resources imageResources
