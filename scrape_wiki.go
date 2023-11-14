package main

import (
	"log"

	"github.com/gocolly/colly"
)

func scrapeData(wpage wikiPage) ([]byte, error) {

	// HTTP client creation using colly
	// Collector manages the network communication and responsible for the execution of the attached callbacks while a collector job is running.
	c := colly.NewCollector(
		colly.AllowedDomains("en.wikipedia.org"),
	)

	// Check status
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	// Get HTML content
	c.OnResponse(func(r *colly.Response) {
		wpage.Body = r.Body
	})

	err := c.Visit(wpage.Url)
	if err != nil {
		return wpage.Body, err
	}

	return wpage.Body, nil
}
