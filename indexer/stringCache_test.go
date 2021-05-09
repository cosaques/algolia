package indexer_test

import (
	"fmt"
	"testing"

	"github.com/cosaques/algolia/indexer"
)

func TestLoadOrStoreStringPtr(t *testing.T) {
	p1 := indexer.LoadOrStoreStringPtr("query1")
	p2 := indexer.LoadOrStoreStringPtr("query1")
	p3 := indexer.LoadOrStoreStringPtr("query2")

	if p1 != p2 {
		t.Fatalf("Different pointers for same strings.")
	}

	if *p1 != "query1" || *p3 != "query2" {
		t.Fatalf("Actual %q, %q, want %q, %q", *p1, *p3, "query1", "query2")
	}
}

func BenchmarkStoreStringPtr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := fmt.Sprintf("Query %d", i)
		indexer.LoadOrStoreStringPtr(s)
	}
}

func BenchmarkLoadStringPtr(b *testing.B) {
	indexer.LoadOrStoreStringPtr("Query")
	for i := 0; i < b.N; i++ {
		indexer.LoadOrStoreStringPtr("Query")
	}
}
