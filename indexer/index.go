package indexer

import (
	"sort"
	"sync"
)

type (
	// Index introduces common methods for indexes.
	Index interface {
		// Add adds new query to the index.
		Add(string)
		// Len gets the count of distinct indexed queries.
		Len() int
		// Top returns most popular queries.
		Top(int) []TopQuery
	}

	// TopQuery represents response of Index.Top().
	TopQuery struct {
		Query string
		Count int
	}
)

type (
	// memoryIndex indexes data and stores it in memory.
	memoryIndex struct {
		// counts is a map containing distinct queries and their counts,
		// note that we operate with references to strings to optimize the memory usage.
		counts map[*string]int
		// order keeps order of map keys according to its counts.
		order []*string
		// mux allows to read/write maps and slices in concurrent way.
		mux sync.RWMutex
		// toIndex stores queries to be indexed.
		toIndex chan indexArgs
	}

	// indexArgs allows to track the completion of query's indexation.
	indexArgs struct {
		s         *string
		completed chan<- bool
	}
)

// NewMemoryIndex creates an instance of memoryIndex
func NewMemoryIndex() Index {
	idx := &memoryIndex{
		counts: make(map[*string]int),
		mux:    sync.RWMutex{},
		// toIndex is a buffered channel that allows to check
		// if there are other queries waiting to be indexed.
		toIndex: make(chan indexArgs, 1),
	}

	// Run indexation in parallel.
	go idx.run()

	return idx
}

// Add adds new query to the index.
func (idx *memoryIndex) Add(s string) {
	completed := make(chan bool)
	idx.toIndex <- indexArgs{LoadOrStoreStringPtr(s), completed}

	// Wait the end of indexation.
	<-completed
}

// Len gets the count of distinct indexed queries.
func (idx *memoryIndex) Len() int {
	idx.mux.RLock()
	defer idx.mux.RUnlock()

	return len(idx.counts)
}

// Top returns most popular queries.
func (idx *memoryIndex) Top(size int) []TopQuery {
	idx.mux.RLock()
	defer idx.mux.RUnlock()

	// Check that asked size is less than current count of distinct queries.
	if size > len(idx.order) {
		size = len(idx.order)
	}
	result := make([]TopQuery, size)

	// Create response.
	for i := 0; i < size; i++ {
		result[i] = TopQuery{*idx.order[i], idx.counts[idx.order[i]]}
	}

	return result
}

// run is listening the channel for new queries to be indexed and index them.
func (idx *memoryIndex) run() {
	for indexArgs := range idx.toIndex {
		s := indexArgs.s
		idx.mux.Lock()
		{
			// Add new query.
			if _, exists := idx.counts[s]; !exists {
				idx.order = append(idx.order, s)
			}
			// Increment query's count
			idx.counts[s]++

			// If no other queries wait for indexation -
			// order the queries by their counts.
			if len(idx.toIndex) == 0 {
				sort.Slice(idx.order, func(i, j int) bool {
					return idx.counts[idx.order[i]] > idx.counts[idx.order[j]]
				})
			}
		}
		idx.mux.Unlock()

		// Mark the indexation as completed.
		indexArgs.completed <- true
	}
}
