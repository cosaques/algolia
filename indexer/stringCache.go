package indexer

import "sync"

var stringCache map[string]*string = make(map[string]*string)
var lock = sync.RWMutex{}

func LoadOrStoreStringPtr(s string) *string {
	lock.RLock()
	if p, exist := stringCache[s]; !exist {
		// The string wasn't found, so we'll create it.
		lock.RUnlock()
		lock.Lock()
		defer lock.Unlock()
		if p, exist := stringCache[s]; !exist {
			// Insert the new string.
			stringCache[s] = &s
			return &s
		} else {
			return p
		}
	} else {
		lock.RUnlock()
		return p
	}
}
