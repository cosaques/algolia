package indexer

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

type (
	// TraceReader reads the log file of queries.
	TraceReader interface {
		// Read reads the traces one by one.
		Read() (Trace, error)
	}

	// Trace represents a line in a log file of queries.
	Trace struct {
		Date  time.Time
		Query string
	}
)

// tsvTraceReader reads traces from a tsv file.
type tsvTraceReader struct {
	tsvReader *csv.Reader
}

// NewTraceReader creates a new instance of tsvTraceReader.
func NewTraceReader(tsvFile io.Reader) TraceReader {
	csvReader := csv.NewReader(tsvFile)
	csvReader.Comma = '\t'
	return &tsvTraceReader{
		tsvReader: csvReader,
	}
}

// Read reads the traces one by one.
func (t *tsvTraceReader) Read() (Trace, error) {
	// Read next record from a tsv file.
	csvRecord, err := t.tsvReader.Read()
	if err != nil {
		return Trace{}, fmt.Errorf("traceReader.Read(): %w.", err)
	}

	// Only to fields expected: date and query.
	if len(csvRecord) != 2 {
		return Trace{}, fmt.Errorf("traceReader.Read(): line should contain 2 args %v.", csvRecord)
	}

	// Parse date.
	date, err := time.Parse("2006-01-02 15:04:05", csvRecord[0])
	if err != nil {
		return Trace{}, fmt.Errorf("traceReader.Read(): %w.", err)
	}

	// Construct a Trace.
	return Trace{Date: date, Query: csvRecord[1]}, nil
}
