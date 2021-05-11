package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/cosaques/algolia/indexer"
)

type QueriesHandler interface {
	Distinct(w http.ResponseWriter, r *http.Request)
	Popular(w http.ResponseWriter, r *http.Request)
}

type queriesHandler struct {
	aggregator indexer.Aggregator
}

func NewQueriesHandler(filePath string) QueriesHandler {
	h := queriesHandler{
		aggregator: indexer.NewAggregator(),
	}
	go h.uploadLogs(filePath)
	return &h
}

func (h *queriesHandler) Distinct(w http.ResponseWriter, r *http.Request) {
	datePrefix := strings.Trim(r.URL.Path, "/1/queries/count/")

	timeRange, err := indexer.ParseTimeRange(datePrefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := 0
	if idx := h.aggregator.GetIndex(timeRange); idx != nil {
		result = idx.Len()
	}

	resp := CountResponse{Count: result}

	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *queriesHandler) Popular(w http.ResponseWriter, r *http.Request) {
	datePrefix := strings.Trim(r.URL.Path, "/1/queries/popular/")

	timeRange, err := indexer.ParseTimeRange(datePrefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryValues, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(queryValues["size"]) == 0 {
		http.Error(w, "Query should contain a \"size\" parameter", http.StatusInternalServerError)
		return
	}
	size, err := strconv.Atoi(queryValues["size"][0])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result []indexer.TopQuery
	if idx := h.aggregator.GetIndex(timeRange); idx != nil {
		result = idx.Top(size)
	}

	resp := PopularResponse{Queries: make([]QueryCountResponse, len(result))}
	for i, r := range result {
		resp.Queries[i] = QueryCountResponse{Query: r.Query, Count: r.Count}
	}

	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *queriesHandler) uploadLogs(filePath string) {
	file, _ := os.Open(filePath)
	defer file.Close()
	traceReader := indexer.NewTraceReader(file)
	var wg sync.WaitGroup
	for trace, err := traceReader.Read(); !errors.Is(err, io.EOF); trace, err = traceReader.Read() {
		if err != nil {
			log.Printf("queriesHandler.uploadLogs: %v.", err)
			continue
		}
		wg.Add(1)
		go func(trace indexer.Trace) {
			defer wg.Done()
			h.aggregator.Add(trace)
		}(trace)
	}
	wg.Wait()
}
