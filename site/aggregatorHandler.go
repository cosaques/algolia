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

type aggregatorHandler struct {
	aggregator   indexer.Aggregator
	handledCount int
}

func newAggregatorHandler() *aggregatorHandler {
	return &aggregatorHandler{
		aggregator: indexer.NewAggregator(),
	}
}

func (h *aggregatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[3]

	switch action {
	case "count":
		if len(segs) != 5 {
			http.Error(w, "Wrong URL", http.StatusNotFound)
			return
		}

		timeRange, err := indexer.ParseTimeRange(segs[4])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.handleCount(timeRange, w, r)
	case "popular":
		if len(segs) != 5 {
			http.Error(w, "Wrong URL", http.StatusNotFound)
			return
		}

		timeRange, err := indexer.ParseTimeRange(segs[4])
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
			http.Error(w, "Query should contain a \"size\" parameter", http.StatusBadRequest)
			return
		}
		size, err := strconv.Atoi(queryValues["size"][0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.handlePopular(timeRange, size, w, r)
	}
}

func (h *aggregatorHandler) handleCount(timeRange indexer.TimeRange, w http.ResponseWriter, r *http.Request) {
	result := 0
	if idx := h.aggregator.GetIndex(timeRange); idx != nil {
		result = idx.Len()
	}

	resp := CountResponse{Count: result}

	w.WriteHeader(200)
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *aggregatorHandler) handlePopular(timeRange indexer.TimeRange, size int, w http.ResponseWriter, r *http.Request) {

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

func (h *aggregatorHandler) uploadLogs(filePath string) {
	file, _ := os.Open(filePath)
	defer file.Close()

	traceReader := indexer.NewTraceReader(file)
	var wg sync.WaitGroup
	for trace, err := traceReader.Read(); !errors.Is(err, io.EOF); trace, err = traceReader.Read() {
		if err != nil {
			log.Fatalf("aggregatorHandler.uploadLogs(): %v", err)
		}
		wg.Add(1)
		go func(trace indexer.Trace) {
			defer wg.Done()
			h.aggregator.Add(trace)
		}(trace)
	}

	wg.Wait()
}
