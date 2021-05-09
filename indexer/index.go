package indexer

import (
	"sort"
	"sync"
)

// Abstractions
type (
	Index interface {
		Add(string)
		Len() int
		Top(int) []TopQuery
	}

	TopQuery struct {
		Query string
		Count int
	}
)

// Implementation
type (
	index struct {
		counts  map[*string]int
		order   []*string
		mux     sync.RWMutex
		toIndex chan indexArgs
	}

	indexArgs struct {
		s         *string
		completed chan<- bool
	}
)

func NewIndex() Index {
	idx := &index{
		counts:  make(map[*string]int),
		mux:     sync.RWMutex{},
		toIndex: make(chan indexArgs, 1),
	}
	go idx.run()
	return idx
}

func (idx *index) Add(s string) {
	completed := make(chan bool)
	idx.toIndex <- indexArgs{LoadOrStoreStringPtr(s), completed}
	<-completed
}

func (idx *index) Len() int {
	idx.mux.RLock()
	defer idx.mux.RUnlock()

	return len(idx.counts)
}

func (idx *index) Top(size int) []TopQuery {
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

func (idx *index) run() {
	for indexArgs := range idx.toIndex {
		s := indexArgs.s
		idx.mux.Lock()
		{
			if _, exists := idx.counts[s]; !exists {
				idx.order = append(idx.order, s)
			}
			idx.counts[s]++

			if len(idx.toIndex) == 0 {
				sort.Slice(idx.order, func(i, j int) bool {
					return idx.counts[idx.order[i]] > idx.counts[idx.order[j]]
				})
			}
		}
		idx.mux.Unlock()
		indexArgs.completed <- true
	}
}
