package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Errorf("couldn't make new document with goquery: %w", err)
	}

	header := doc.Find("h1").Text()
	if len(header) == 0 {
		header = doc.Find("h2").Text()
	}
	if len(header) != 0 {
		return header
	}

	return ""
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Errorf("couldn't make new document with goquery: %w", err)
	}

	paragraph := (((doc.Find("main")).Find("p")).First()).Text()
	if len(paragraph) != 0 {
		return paragraph
	}

	paragraph = ((doc.Find("p")).First()).Text()
	if len(paragraph) != 0 {
		return paragraph
	}

	return ""
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	baseURLstring := baseURL.String()
	links := []string{}
	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		urlLink, _ := s.Attr("href")
		if strings.HasPrefix(urlLink, baseURLstring) {
			links = append(links, urlLink)
		} else {
			urlLinkAbsolute := baseURLstring + urlLink
			links = append(links, urlLinkAbsolute)
		}
	})

	return links, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	baseURLstring := baseURL.String()
	links := []string{}
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		imgLink, _ := s.Attr("src")
		if strings.HasPrefix(imgLink, baseURLstring) {
			links = append(links, imgLink)
		} else {
			imgLinkAbsolute := baseURLstring + imgLink
			links = append(links, imgLinkAbsolute)
		}
	})

	fmt.Println(links)
	return links, nil
}

func extractPageData(html, pageURL string) PageData {
	
}