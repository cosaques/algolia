package indexer

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"
)

type TraceReader interface {
	Read() (Trace, error)
}

type Trace struct {
	Date  time.Time
	Query string
}

type traceReader struct {
	csvReader *csv.Reader
}

func NewTraceReader(tsvFile io.Reader) TraceReader {
	csvReader := csv.NewReader(tsvFile)
	csvReader.Comma = '\t'
	return &traceReader{
		csvReader: csvReader,
	}
}

func (t *traceReader) Read() (Trace, error) {
	csvRecord, err := t.csvReader.Read()
	if err != nil {
		return Trace{}, fmt.Errorf("traceReader.Read(): %w", err)
	}

	if len(csvRecord) != 2 {
		return Trace{}, fmt.Errorf("traceReader.Read(): line should contain 2 args %v.", csvRecord)
	}

	date, err := time.Parse("2006-01-02 15:04:05", csvRecord[0])
	if err != nil {
		return Trace{}, fmt.Errorf("traceReader.Read(): %w", err)
	}

	return Trace{Date: date, Query: csvRecord[1]}, nil
}
