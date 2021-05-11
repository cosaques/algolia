package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":5000", "The addr of the application")
	file := flag.String("file", "", "The path to .tsv file containing logs")
	flag.Parse()

	aggregatorHandler := newAggregatorHandler()

	http.Handle("/1/queries/", aggregatorHandler)

	go aggregatorHandler.uploadLogs(*file)

	// start the web server
	log.Println("Starting the webserver on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
