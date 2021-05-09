package indexer

import "sync"

var stringCache sync.Map

func LoadOrStoreStringPtr(s string) *string {
	p, _ := stringCache.LoadOrStore(s, &s)
	return p.(*string)
}
