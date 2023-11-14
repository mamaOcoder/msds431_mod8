package main

import (
	"fmt"
	"log"
	"os"
)

// Function to write the body text to file
func (wp *wikiPage) save() error {
	// save wikipedia page html to wikipages directory
	wikiDir := "./wikipages"

	// Check if the directory exists
	_, err := os.Stat(wikiDir)
	if os.IsNotExist(err) {
		// Directory does not exist, so create it
		err := os.MkdirAll(wikiDir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
		}
	}

	// Create the filename
	fileName := fmt.Sprintf("./wikipages/%s.txt", wp.Title)

	return os.WriteFile(fileName, wp.Body, 0644)
}

// Function to load pages that have been saved
func loadPage(title string) (*wikiPage, error) {
	filename := fmt.Sprintf("./wikipages/%s.txt", title)
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &wikiPage{Title: title, Body: body}, nil
}

// Function to find a wikiPage by title
func findWikiPageByTitle(title string, wikipages []wikiPage) (*wikiPage, error) {
	for _, page := range wikipages {
		if page.Title == title {
			return &page, nil
		}
	}

	return nil, fmt.Errorf("wikiPage with title %s not found", title)
}

// Function to get body of wikiPage
func getWikiPageBody(wp wikiPage) (*wikiPage, error) {
	// Check if the body of the wikiPage is empty or is already wriiten to file
	if len(wp.Body) != 0 {
		fmt.Printf("\nBody of %s has already been scraped and loaded.", wp.Title)
		return &wp, nil
	} else if _, err := os.Stat(fmt.Sprintf("./wikipages/%s.txt", wp.Title)); err == nil {
		fmt.Printf("\nBody of %s has already been scraped. Loading data.", wp.Title)
		loaded, err := loadPage(wp.Title)
		if err != nil {
			log.Printf("Error loading wikiPage from file: %v", err)
			return &wp, err
		}
		wp.Body = loaded.Body
		return &wp, nil
	} else {
		fmt.Printf("\nBody of %s is empty. Scraping wiki page.\n", wp.Title)
		scraped, err := scrapeData(wp)
		if err != nil {
			log.Printf("Error scraping wikiPage: %v", err)
			return &wp, err
		}
		wp.Body = scraped
		wp.save()
		return &wp, nil
	}
}
