package indexer_test

import (
	"errors"
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/cosaques/algolia/indexer"
)

func TestAggregatorAdd(t *testing.T) {
	aggregator := indexer.NewAggregator()
	traces := []indexer.Trace{
		{time.Date(2015, 8, 1, 0, 3, 43, 0, time.UTC), "q1"},
		{time.Date(2015, 8, 2, 0, 3, 43, 0, time.UTC), "q1"},
		{time.Date(2015, 8, 2, 0, 3, 44, 0, time.UTC), "q2"},
		{time.Date(2015, 8, 2, 0, 5, 45, 0, time.UTC), "q3"},
	}

	var wg sync.WaitGroup
	for _, trace := range traces {
		wg.Add(1)
		go func(trace indexer.Trace) {
			defer wg.Done()
			aggregator.Add(trace)
		}(trace)
	}
	wg.Wait()

	tests := []struct {
		name      string
		timeRange string
		wantNil   bool
		len       int
	}{
		{
			"Year",
			"2015",
			false,
			3,
		},
		{
			"Day1",
			"2015-08-01",
			false,
			1,
		},
		{
			"Day2",
			"2015-08-02",
			false,
			3,
		},
		{
			"Minute",
			"2015-08-02 00:03",
			false,
			2,
		},
		{
			"NotExisting",
			"2015-09",
			true,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeRange, _ := indexer.ParseTimeRange(tt.timeRange)
			idx := aggregator.GetIndex(timeRange)

			if (idx == nil) != tt.wantNil {
				t.Fatalf("GetIndexFor(%q) is nil = %v, want nil %v", tt.timeRange, idx == nil, tt.wantNil)
			} else if idx == nil {
				return
			}

			if idx.Len() != tt.len {
				t.Errorf("GetIndexFor(%q).Len() = %d, want %d", tt.timeRange, idx.Len(), tt.len)
			}
		})
	}
}

func BenchmarkAggregator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		file, _ := os.Open("testdata/bench_aggr.tsv")
		defer file.Close()

		traceReader := indexer.NewTraceReader(file)
		aggregator := indexer.NewAggregator()
		var wg sync.WaitGroup
		for trace, err := traceReader.Read(); !errors.Is(err, io.EOF); trace, err = traceReader.Read() {
			if err != nil {
				b.Fatalf("Error %v occured", err)
			}

			wg.Add(1)
			go func(trace indexer.Trace) {
				defer wg.Done()
				aggregator.Add(trace)
			}(trace)
		}
		wg.Wait()
	}
}
