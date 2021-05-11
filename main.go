package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/cosaques/algolia/indexer"
)

func main() {
	ch := make(chan struct{})
	go func() {
		var i int
		for {
			if _, ok := <-ch; !ok {
				return
			}

			i++
			fmt.Printf("\rIndexed %d", i)
		}
	}()

	file, _ := os.Open("hn_logs.tsv")
	defer file.Close()
	traceReader := indexer.NewTraceReader(file)
	aggregator := indexer.NewAggregator()
	var wg sync.WaitGroup
	for trace, err := traceReader.Read(); !errors.Is(err, io.EOF); trace, err = traceReader.Read() {
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go func(trace indexer.Trace) {
			defer wg.Done()
			aggregator.Add(trace)
			ch <- struct{}{}
		}(trace)
	}
	wg.Wait()
	close(ch)

	fmt.Println("\nCompleted!")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Query type (L|T): ")
		scanner.Scan()
		query := scanner.Text()
		fmt.Printf("%q\n", query)
		switch query {
		case "L":
			fmt.Print("Date range: ")
			scanner.Scan()
			date := scanner.Text()
			tr, _ := indexer.ParseTimeRange(date)
			if idx := aggregator.GetIndex(tr); idx != nil {
				fmt.Println(idx.Len())
			} else {
				fmt.Println(0)
			}

		case "T":
			fmt.Print("Size: ")
			scanner.Scan()
			size, _ := strconv.Atoi(scanner.Text())
			fmt.Println(size)
			fmt.Print("Date range: ")
			scanner.Scan()
			date := scanner.Text()
			tr, _ := indexer.ParseTimeRange(date)
			if idx := aggregator.GetIndex(tr); idx != nil {
				fmt.Println(idx.Top(size))
			} else {
				fmt.Println("[]")
			}
		}
	}

}
