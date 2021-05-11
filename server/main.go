package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":5000", "The addr of the application")
	flag.Parse()

	h := NewQueriesHandler("../indexer/testdata/bench_aggr.tsv")
	http.HandleFunc("/1/queries/count/", h.Distinct)
	http.HandleFunc("/1/queries/popular/", h.Popular)

	// start the web server
	log.Println("Starting the webserver on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
