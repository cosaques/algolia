package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// Allowed flags.
	addr := flag.String("addr", ":5000", "The addr of the application")
	file := flag.String("file", "", "The path to .tsv file containing logs")
	flag.Parse()

	aggregatorHandler := newAggregatorHandler()

	// Add possible routes and their handlers.
	http.Handle("/", &templateHandler{fileName: "index.html"})
	http.Handle("/1/queries/", aggregatorHandler)

	// Upload and handle log file in parallel.
	go aggregatorHandler.uploadLogs(*file)

	// Start the web server.
	log.Println("Starting the webserver on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
