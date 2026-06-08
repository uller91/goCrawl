package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	argURL := args[0]
	fmt.Printf("starting crawl of: %s\n", argURL)

	pages := make(map[string]PageData)
	baseURL, err := url.Parse(argURL)
	if err != nil {
		fmt.Errorf("couldn't parse input URL: %v", err)
		return
	}
	concurrencyControl := make(chan struct{}, 20)
	var mu sync.Mutex
	var wg sync.WaitGroup

	cfg := &config{
		pages:              pages,
		baseURL:            baseURL,
		mu:                 &mu,
		concurrencyControl: concurrencyControl,
		wg:                 &wg,
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(argURL)
	cfg.wg.Wait()

	fmt.Println("")
	fmt.Println("Total ctawled:\n")
	for pageURL, page := range cfg.pages {
		fmt.Printf("%s: %s\n", pageURL, page.Heading)
	}

	/*
		crawlPage(url, url, pages)

		fmt.Println("")
		fmt.Println("Total ctawled:\n")
		for key, val := range pages {
			fmt.Printf("%s: %v\n", key, val)
		}
	*/

	/*
		html, err := getHTML(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		//fmt.Println(html)

		pageData := extractPageData(html, url)
		fmt.Println(pageData)
	*/

}
