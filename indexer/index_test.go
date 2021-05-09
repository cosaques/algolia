package indexer_test

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/cosaques/algolia/indexer"
)

func TestListIndexAdd(t *testing.T) {
	idx := indexer.NewListIndex()
	testAdd(t, idx)
}

func TestListIndexTop(t *testing.T) {
	idx := indexer.NewListIndex()
	testTop(t, idx)
}

func BenchmarkListIndex(b *testing.B) {
	benchmarkIndex(b, indexer.NewListIndex)
}

func testAdd(t *testing.T, idx indexer.Index) {
	index(idx, 10, 10)
	if idx.Len() != 10 {
		t.Fatalf("Index len = %d, want %d (%T)", idx.Len(), 10, idx)
	}
}

func testTop(t *testing.T, idx indexer.Index) {
	index(idx, 10, 10)
	tops := idx.Top(3)
	for i, top := range tops {
		want := indexer.TopQuery{fmt.Sprintf("Query %d", 10-i), 10 * (10 - i)}
		if top != want {
			t.Errorf("Top %d = %v, want %v (%T)", i+1, top, want, idx)
		}
	}
}

func benchmarkIndex(b *testing.B, idxFactory func() indexer.Index) {
	file, _ := os.Open("testdata/bench.txt")
	defer file.Close()
	scanner := bufio.NewScanner(file)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		idx := idxFactory()
		var wg sync.WaitGroup
		for scanner.Scan() {
			wg.Add(1)
			go func(s string) {
				defer wg.Done()
				idx.Add(s)
			}(scanner.Text())
		}
		wg.Wait()
	}
}

func index(idx indexer.Index, queriesNb, countMultimply int) {
	var wg sync.WaitGroup
	for i := 1; i <= queriesNb; i++ {
		query := fmt.Sprintf("Query %d", i)
		for j := 0; j < countMultimply*i; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				idx.Add(query)
			}()
		}
	}
	wg.Wait()
}
