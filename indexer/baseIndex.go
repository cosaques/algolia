package indexer

import (
	"sort"
	"sync"
)

type baseIndex struct {
	counts  map[*string]int
	order   []*string
	mux     sync.RWMutex
	toIndex chan *string
}

func NewBaseIndex() Index {
	idx := &baseIndex{
		counts:  make(map[*string]int),
		mux:     sync.RWMutex{},
		toIndex: make(chan *string, 1),
	}
	go idx.run()
	return idx
}

func (idx *baseIndex) Add(s string) {
	p := LoadOrStoreStringPtr(s)
	idx.toIndex <- p
}

func (idx *baseIndex) Len() int {
	idx.mux.RLock()
	defer idx.mux.RUnlock()

	return len(idx.counts)
}

func (idx *baseIndex) Top(size int) []TopQuery {
	idx.mux.RLock()
	defer idx.mux.RUnlock()

	if size > len(idx.order) {
		size = len(idx.order)
	}
	result := make([]TopQuery, size)

	for i := 0; i < size; i++ {
		result[i] = TopQuery{*idx.order[i], idx.counts[idx.order[i]]}
	}

	return result
}

func (idx *baseIndex) run() {
	for p := range idx.toIndex {
		idx.mux.Lock()
		{
			if _, exists := idx.counts[p]; !exists {
				idx.order = append(idx.order, p)
			}
			idx.counts[p]++

			if len(idx.toIndex) == 0 {
				sort.Slice(idx.order, func(i, j int) bool {
					return idx.counts[idx.order[i]] > idx.counts[idx.order[j]]
				})
			}
		}
		idx.mux.Unlock()
	}
}
