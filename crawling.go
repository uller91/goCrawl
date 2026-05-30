package main

import (
	"fmt"
	"io"
	"net/http"
)

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
	if contentType != "text/html" {
		return "", fmt.Errorf("Header Content-Type is not text/html: %v", contentType)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(body), nil
}
