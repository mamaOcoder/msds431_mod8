package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	invalidURL   = "http://random-invalid-url"
	validURL     = "https://en.wikipedia.org/wiki/Robotics"
	invalidTitle = "Invalid Title"
	validTitle   = "Robotics"
)

func TestScrape(t *testing.T) {
	assert := assert.New(t)

	invTestWP := wikiPage{
		Url: invalidURL,
	}
	// Test Scrape with invalid URL
	_, err := scrapeData(invTestWP)
	assert.Error(err, "scrapeData() should return an error for invalid URLs")

	valTestWP := wikiPage{
		Url: validURL,
	}
	// Test Scrape with valid URL
	_, err = scrapeData(valTestWP)
	assert.NoError(err, "scrapeData() should not return an error for valid URLs")

}

func TestWikiPageSave(t *testing.T) {
	assert := assert.New(t)

	// Create a wikiPage instance
	wp := wikiPage{
		Title: "TestPage",
		Body:  []byte("Test body content"),
	}

	// Call the save method
	err := wp.save()
	assert.NoError(err, "save() should not return error")

	// Check if the file was created successfully
	filePath := fmt.Sprintf("./wikipages/%s.txt", wp.Title)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		t.Errorf("File %s was not created", filePath)
	}

	// Clean up: delete the test file
	err = os.Remove(filePath)
	if err != nil {
		t.Errorf("Error deleting test file: %v", err)
	}
}

func TestFindWikiPageByTitle(t *testing.T) {
	assert := assert.New(t)

	// Load wikipages in library
	var wikipages []wikiPage

	// Open the JSON Lines file
	file, err := os.Open("pages.jl")
	if err != nil {
		fmt.Println("Error opening the file:", err)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate through each line in the file
	for scanner.Scan() {
		// Read the current line
		line := scanner.Text()

		// Create an instance of the struct to hold the data
		var wp wikiPage

		// Unmarshal the JSON data into the struct
		err := json.Unmarshal([]byte(line), &wp)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		wikipages = append(wikipages, wp)
	}

	// Test invalid title
	_, err = findWikiPageByTitle(invalidTitle, wikipages)
	assert.Error(err, "findWikiPageByTitle() should return an error for invalid titles")

	// Test valid title
	_, err = findWikiPageByTitle(validTitle, wikipages)
	assert.NoError(err, "findWikiPageByTitle() should not return an error for valid titles")
}

func TestHomeHandler(t *testing.T) {
	// Create a test server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var wikipages []wikiPage
		homeVariablesInstance := &homeVariables{
			PageTitles:   []string{"Page1", "Page2"},
			SelectedPage: "",
		}
		homeHandler(w, r, homeVariablesInstance, wikipages)
	}))
	defer testServer.Close()

	// Test case 1: GET request to render the home page
	req1, err := http.NewRequest("GET", testServer.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp1 := httptest.NewRecorder()
	testServer.Config.Handler.ServeHTTP(resp1, req1)

	// Check the response status code for the GET request
	if resp1.Code != http.StatusOK {
		t.Errorf("GET request returned wrong status code: got %v want %v", resp1.Code, http.StatusOK)
	}

	// Test case 2: POST request to simulate form submission
	formData := url.Values{}
	formData.Set("form_submitted", "1")
	formData.Set("pages", "Page1")

	req2, err := http.NewRequest("POST", testServer.URL, strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp2 := httptest.NewRecorder()
	testServer.Config.Handler.ServeHTTP(resp2, req2)

	// Check the response status code for the POST request
	if resp2.Code != http.StatusSeeOther {
		t.Errorf("POST request returned wrong status code: got %v want %v", resp2.Code, http.StatusSeeOther)
	}

	// Check if the redirect URL is as expected
	expectedRedirectURL := "/view/Page1"
	if got := resp2.Header().Get("Location"); got != expectedRedirectURL {
		t.Errorf("POST request redirected to wrong URL: got %v want %v", got, expectedRedirectURL)
	}
}
