package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(urlStr string) (string, error) {
	urlFull, err := url.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("couldn't parse URL: %w", err)
	}

	fullPath := urlFull.Host + urlFull.Path
	fullPath = strings.ToLower(fullPath)
	fullPath = strings.TrimSuffix(fullPath, "/")

	return fullPath, nil
}
