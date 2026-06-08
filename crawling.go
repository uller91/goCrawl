package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl // release slot
		cfg.wg.Done()            // signal done
	}()

	if !strings.HasPrefix(rawCurrentURL, cfg.baseURL.String()) {
		return
	}

	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error during URL normalization: %v", err)
		return
	}

	isFirst := cfg.addPageVisit(normCurrentURL)
	if !isFirst {
		return
	}

	fmt.Printf("Fetching HTML from %v...\n", rawCurrentURL)
	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting HTML from %v URL: %v", rawCurrentURL, err)
		return
	}
	pageData := extractPageData(pageHTML, rawCurrentURL)
	cfg.mu.Lock()
	cfg.pages[normCurrentURL] = pageData
	cfg.mu.Unlock()

	for _, imageURL := range pageData.ImageURLs {
		cfg.wg.Add(1)
		go cfg.crawlPage(imageURL)
	}

	for _, outgoingLink := range pageData.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(outgoingLink)
	}

}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	_, exist := cfg.pages[normalizedURL]
	if !exist {
		cfg.pages[normalizedURL] = PageData{}
	}
	cfg.mu.Unlock()
	isFirst = !exist
	return isFirst
}

/*
func crawlPageNoConcurency(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	if !strings.HasPrefix(rawCurrentURL, rawBaseURL) {
		return
	}

	normCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("error during URL normalization: %v", err)
		return
	}

	count, exists := pages[normCurrentURL]
	if exists {
		count += 1
		pages[normCurrentURL] = count
		return
	} else {
		pages[normCurrentURL] = 1
	}

	fmt.Printf("Fetching HTML from %v...\n", rawCurrentURL)
	pageHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("error getting HTML from %v URL: %v", rawCurrentURL, err)
		return
	}

	pageData := extractPageData(pageHTML, rawCurrentURL)
	//fmt.Println(pageData)
	for _, imageURL := range pageData.ImageURLs {
		crawlPageNoConcurency(rawBaseURL, imageURL, pages)
	}

	for _, outgoingLink := range pageData.OutgoingLinks {
		crawlPageNoConcurency(rawBaseURL, outgoingLink, pages)
	}
}
*/

func getHTML(rawURL string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform the request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("Error-level status code: %v", res.StatusCode)
	}

	contentType := res.Header.Get("Content-Type")
	contentTypes := strings.Split(contentType, ";")
	//if contentType != "text/html" {
	if !slices.Contains(contentTypes, "text/html") {
		return "", fmt.Errorf("Header Content-Type is not text/html: %v", contentType)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}
