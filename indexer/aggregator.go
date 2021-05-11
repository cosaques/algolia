package indexer

import "sync"

// Aggregator is an abstraction of index aggregation.
type Aggregator interface {
	// Add adds a new Trace to all the concerned indexes.
	Add(Trace)
	// GetIndex returns an index for a given TimeRange.
	GetIndex(TimeRange) Index
}

// aggregator contains indexes for each possible TimeRange.
type aggregator struct {
	// indexes is a map of indexes for presented TimeRanges.
	indexes map[string]Index
	// mux allows to read/write maps and slices in concurrent way.
	mux sync.RWMutex
}

// NewAggregator creates an instance of aggregator.
func NewAggregator() Aggregator {
	return &aggregator{
		indexes: make(map[string]Index),
	}
}

// Add adds a new Trace to all the concerned indexes.
func (a *aggregator) Add(t Trace) {
	// All possible time precisions.
	precisions := []TimePrecision{Year, Month, Day, Hour, Minute}

	// Index Trace with all possible time precisions.
	var wg sync.WaitGroup
	for _, precision := range precisions {
		wg.Add(1)
		go func(precision TimePrecision) {
			defer wg.Done()
			index := a.getOrCreateIndex(TimeRange{t.Date, precision})
			index.Add(t.Query)
		}(precision)
	}
	wg.Wait()
}

// GetIndex returns an index for a given TimeRange.
func (a *aggregator) GetIndex(r TimeRange) Index {
	idxKey := r.String()

	a.mux.RLock()
	defer a.mux.RUnlock()

	return a.indexes[idxKey]
}

// getOrCreateIndex returns either existing index or
// a newly created for a given TimeRange.
func (a *aggregator) getOrCreateIndex(r TimeRange) Index {
	idxKey := r.String()

	a.mux.RLock()
	if idx, exist := a.indexes[idxKey]; !exist {
		// The index wasn't found, so we'll create it.
		a.mux.RUnlock()
		a.mux.Lock()
		defer a.mux.Unlock()
		if idx, exist := a.indexes[idxKey]; !exist {
			// Insert the new string.
			idx = NewMemoryIndex()
			a.indexes[idxKey] = idx
			return idx
		} else {
			return idx
		}
	} else {
		a.mux.RUnlock()
		return idx
	}
}
