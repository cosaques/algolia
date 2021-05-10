package indexer

import "sync"

var stringCache map[string]*string = make(map[string]*string)
var mux = sync.RWMutex{}

func LoadOrStoreStringPtr(s string) *string {
	mux.RLock()
	if p, exist := stringCache[s]; !exist {
		// The string wasn't found, so we'll create it.
		mux.RUnlock()
		mux.Lock()
		defer mux.Unlock()
		if p, exist := stringCache[s]; !exist {
			// Insert the new string.
			stringCache[s] = &s
			return &s
		} else {
			return p
		}
	} else {
		mux.RUnlock()
		return p
	}
}
