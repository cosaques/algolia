package indexer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
)

func TestTraceRead(t *testing.T) {
	want := []Trace{
		{time.Date(2015, 8, 1, 0, 4, 0, 0, time.UTC), "http%3A%2F%2Fquiltville.blogspot.com%2F2015%2F07%2Fa-little-stop-at-connies-quilt-shop.html"},
		{time.Date(2015, 8, 1, 0, 4, 1, 0, time.UTC), "%22http%3A%2F%2Fwww.metrowestdailynews.com%2Farticle%2F20150701%2FSPORTS%2F150709145%22"},
		{time.Date(2015, 8, 1, 0, 4, 3, 0, time.UTC), "%22http%3A%2F%2Fwww.nbcnews.com%2Fmeet-the-press%2Ffirst-read-hillary-clintons-keystone-problem-n400291%22"},
	}

	file, _ := os.Open("testdata/trace.tsv")
	defer file.Close()

	traceReader := NewTraceReader(file)
	for i := 0; i < 3; i++ {
		t.Run(fmt.Sprintf("Line %d", i), func(t *testing.T) {
			actual, err := traceReader.Read()
			if err != nil {
				t.Fatalf("Error %v occured", err)
			}
			if actual != want[i] {
				t.Fatalf("Get trace %v, want %v", actual, want[i])
			}
		})
	}

	_, err := traceReader.Read()
	if !errors.Is(err, io.EOF) {
		t.Fatalf("Get error %v, want EOF", err)
	}
}
