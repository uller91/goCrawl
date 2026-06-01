package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

type PageData struct {
	URL            string
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

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
			urlLinkUrl, err := url.Parse(urlLink)
			if err != nil {
				fmt.Errorf("couldn't parse input URL: %v", err)
				return
			}
			urlLinkAbsolute := (baseURL.ResolveReference(urlLinkUrl)).String()
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
	//body := doc.Find("main")
	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		imgLink, _ := s.Attr("src")
		if strings.HasPrefix(imgLink, baseURLstring) {
			links = append(links, imgLink)
		} else {
			impLinkUrl, err := url.Parse(imgLink)
			if err != nil {
				fmt.Errorf("couldn't parse input URL: %v", err)
				return
			}
			imgLinkAbsolute := (baseURL.ResolveReference(impLinkUrl)).String()
			links = append(links, imgLinkAbsolute)
		}
	})

	return links, nil
}

func extractPageData(html, pageURL string) PageData {
	baseURL, err := url.Parse(pageURL)
	if err != nil {
		fmt.Errorf("couldn't parse input URL: %v", err)
		return PageData{}
	}

	heading := getHeadingFromHTML(html)
	paragraph := getFirstParagraphFromHTML(html)
	outgoingLinks, err := getURLsFromHTML(html, baseURL)
	if err != nil {
		fmt.Errorf("couldn't parse outgoing links: %v", err)
		return PageData{}
	}
	imageURLs, err := getImagesFromHTML(html, baseURL)
	if err != nil {
		fmt.Errorf("couldn't parse outgoing links: %v", err)
		return PageData{}
	}

	pageData := PageData{
		URL:            pageURL,
		Heading:        heading,
		FirstParagraph: paragraph,
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}

	return pageData
}
