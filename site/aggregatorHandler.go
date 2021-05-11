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
	"sync/atomic"
	"time"

	"github.com/cosaques/algolia/indexer"
	"github.com/gorilla/websocket"
)

// aggregatorHandler is an http.Handler dealing with an index aggregator.
type aggregatorHandler struct {
	// aggregator of indexes.
	aggregator indexer.Aggregator
	// handledCount keeps number of handled logs.
	handledCount int32
}

// newAggregatorHandler creates a new instance of aggregatorHandler
func newAggregatorHandler() *aggregatorHandler {
	return &aggregatorHandler{
		aggregator: indexer.NewAggregator(),
	}
}

// ServeHTTP implements http.Handler interface.
func (h *aggregatorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Split the URL to defince which API action is called.
	segs := strings.Split(r.URL.Path, "/")
	action := segs[3]

	switch action {
	// /1/queries/count/<DATE_PREFIX>
	case "count":
		// URL should contain <DATE_PREFIX> as well.
		if len(segs) != 5 {
			http.Error(w, "Wrong URL", http.StatusNotFound)
			return
		}

		// Check if TimeRange (<DATE_PREFIX>) is valid.
		timeRange, err := indexer.ParseTimeRange(segs[4])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.handleCount(timeRange, w, r)
	// /1/queries/popular/<DATE_PREFIX>?size=<SIZE>
	case "popular":
		// URL should contain <DATE_PREFIX> as well.
		if len(segs) != 5 {
			http.Error(w, "Wrong URL", http.StatusNotFound)
			return
		}

		// Check if TimeRange (<DATE_PREFIX>) is valid.
		timeRange, err := indexer.ParseTimeRange(segs[4])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Try to get a "size" parameter.
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
	// /1/queries/monitoring
	case "monitoring":
		h.handleMonitor(w, r)
	}
}

// handleCount returns count of distinct queries for a given time range.
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

// handlePopular returns top queries for a given time range.
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

// some constants and vars allowing to handle a socket connection
const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// handleMonitor sends the actual number of indexed query traces via a socket connection.
func (h *aggregatorHandler) handleMonitor(w http.ResponseWriter, r *http.Request) {
	// Upgrade the request to a socket connection.
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("ServeHTTP (socket upgrade): ", err)
		return
	}
	defer socket.Close()

	// Send an actual index state to a socket with a given periodicity.
	go func() {
		msg := MonitoringMsg{}
		for {
			// Get a handledCount in a correct concurrent way.
			handledCount := int(atomic.LoadInt32(&h.handledCount))

			// If state not changed don't send it.
			if msg.Indexed != handledCount {
				msg.Indexed = handledCount
				err := socket.WriteJSON(msg)
				if err != nil {
					break
				}
			}

			// Periodicity.
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Keep socket connection alive till client is connected.
	for {
		err := socket.ReadJSON(nil)
		if err != nil {
			return
		}
	}
}

// uploadLogs uploads query traces from a log file.
func (h *aggregatorHandler) uploadLogs(filePath string) {
	file, _ := os.Open(filePath)
	defer file.Close()

	traceReader := indexer.NewTraceReader(file)

	// Index each query trace in a concurrent way
	var wg sync.WaitGroup
	for trace, err := traceReader.Read(); !errors.Is(err, io.EOF); trace, err = traceReader.Read() {
		if err != nil {
			log.Fatalf("aggregatorHandler.uploadLogs(): %v", err)
		}

		wg.Add(1)
		go func(trace indexer.Trace) {
			defer wg.Done()

			// Add query trace to an index aggregator containing all indexes.
			h.aggregator.Add(trace)

			// Increment the number of handles query traces.
			atomic.AddInt32(&h.handledCount, 1)
		}(trace)
	}

	wg.Wait()
}
