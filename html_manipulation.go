package main

import (
	"fmt"
	"strings"
	"github.com/PuerkitoBio/goquery"
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