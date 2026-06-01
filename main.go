package main

import (
	"fmt"
	"os"
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

	url := args[0]
	fmt.Printf("starting crawl of: %s\n", url)

	pages := make(map[string]int)
	crawlPage(url, url, pages)

	fmt.Println("")
	fmt.Println("Total ctawled:\n")
	for key, val := range pages {
		fmt.Printf("%s: %v\n", key, val)
	}
	
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
