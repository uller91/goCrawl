package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]
	switch len(args) {
	case 0:
		fmt.Println("No website provided")
		os.Exit(1)
	case 1, 2:
		fmt.Println("Required arguments: URL maxConcurrency maxPages")
		os.Exit(1)
	case 3:
		fmt.Printf("starting crawl of: %s\n", args[0])
	default:
		fmt.Println("Too many arguments provided! Required arguments: URL maxConcurrency maxPages")
		os.Exit(1)
	}

	argURL := args[0]
	argMaxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Println("maxConcurrency is not int!")
		os.Exit(1)
	}
	argMaxPages, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("maxPages is not int!")
		os.Exit(1)
	}

	pages := make(map[string]PageData)
	baseURL, err := url.Parse(argURL)
	if err != nil {
		fmt.Errorf("couldn't parse input URL: %v", err)
		return
	}
	concurrencyControl := make(chan struct{}, argMaxConcurrency)
	var mu sync.Mutex
	var wg sync.WaitGroup

	cfg := &config{
		pages:              pages,
		baseURL:            baseURL,
		mu:                 &mu,
		concurrencyControl: concurrencyControl,
		wg:                 &wg,
		maxPages:           argMaxPages,
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(argURL)
	cfg.wg.Wait()

	fmt.Println("")
	fmt.Println("Total ctawled:\n")
	for pageURL, page := range cfg.pages {
		fmt.Printf("%s: %s\n", pageURL, page.Heading)
	}

	err = writeJSONReport(cfg.pages, "report.json")
	if err != nil {
		fmt.Errorf("couldn't write JSON report: %v", err)
		return
	}
}
