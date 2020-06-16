package g

import (
	"sync"

	"github.com/flytd/urlooker/dataobj"
)

type DetectedItemSafeMap struct {
	sync.RWMutex
	M map[string][]*dataobj.DetectedItem
}

var (
	DetectedItemMap = &DetectedItemSafeMap{M: make(map[string][]*dataobj.DetectedItem)}
)

func (this *DetectedItemSafeMap) Get(key string) ([]*dataobj.DetectedItem, bool) {
	this.RLock()
	defer this.RUnlock()
	ipItem, exists := this.M[key]
	return ipItem, exists
}

func (this *DetectedItemSafeMap) GetAll() map[string][]*dataobj.DetectedItem {
	this.RLock()
	defer this.RUnlock()
	return this.M
}

func (this *DetectedItemSafeMap) Set(detectedItemMap map[string][]*dataobj.DetectedItem) {
	this.Lock()
	defer this.Unlock()
	this.M = detectedItemMap
}
