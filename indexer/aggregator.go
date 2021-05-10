package indexer

import "sync"

type Aggregator interface {
	Add(Trace)
	GetIndex(TimeRange) Index
}

type aggregator struct {
	indexes map[string]Index
	mux     sync.RWMutex
}

func NewAggregator() Aggregator {
	return &aggregator{
		indexes: make(map[string]Index),
	}
}

func (a *aggregator) Add(t Trace) {
	precisions := []TimePrecision{Year, Month, Day, Hour, Minute}
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

func (a *aggregator) GetIndex(r TimeRange) Index {
	idxKey := r.String()

	a.mux.RLock()
	defer a.mux.RUnlock()
	return a.indexes[idxKey]
}

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
			idx = NewIndex()
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
