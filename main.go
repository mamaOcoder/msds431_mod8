package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// File type and fields for URL, Title and Body for each wiki page
type wikiPage struct {
	Url   string `json:"url"`
	Title string `json:"title"`
	Body  []byte `json:"body"`
}

// variables that are used in the home HTML template
// Note that variables need to be capitol letters in order to be accessed in the template
type homeVariables struct {
	SelectedPage string
	PageTitles   []string
}

func homeHandler(w http.ResponseWriter, r *http.Request, hv *homeVariables, wikipages []wikiPage) {

	// Check if the form is submitted
	// All code for submit response is located inside if statement to avoid it being run on render
	if r.Method == http.MethodPost && r.FormValue("form_submitted") == "1" {
		// Get the selected option value from the form
		selected := r.FormValue("pages")

		// Update the homeVariables with the selected option
		hv.SelectedPage = selected

		fmt.Printf("\nValue of title selected by user: %v", hv.SelectedPage)
		// Get wikipage associated with title
		wp, err := findWikiPageByTitle(hv.SelectedPage, wikipages)
		if err != nil {
			log.Printf("Error finding wikiPage by title: %v", err)
			http.Error(w, "Error finding wikiPage by title", http.StatusInternalServerError)
			return
		}

		// Update body of wikipage
		updatedWP, err := getWikiPageBody(*wp)
		if err != nil {
			log.Printf("Error updating body wikiPage: %v", err)
			http.Error(w, "Error updating body of wikiPage", http.StatusInternalServerError)
		}
		// Update wikipage body in the original wikipages slice
		for i := range wikipages {
			if wikipages[i].Title == hv.SelectedPage {
				wikipages[i].Body = updatedWP.Body
				break
			}
		}

		// Open view page with wikipage content
		http.Redirect(w, r, "/view/"+hv.SelectedPage, http.StatusSeeOther)
		return
	}

	homePath := filepath.Join("templates", "home.gohtml")
	t, err := template.ParseFiles(homePath)

	if err != nil {
		log.Printf("Unable to parse template file: %v", err)
		http.Error(w, "Unable to parse template file", http.StatusInternalServerError)
		return
	}

	t.Execute(w, hv)

}

// Function to allow users to view a wiki page
func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title)
	viewPath := filepath.Join("templates", "view.gohtml")
	t, err := template.ParseFiles(viewPath)
	if err != nil {
		log.Printf("Unable to parse template file: %v", err)
		http.Error(w, "Unable to parse template file", http.StatusInternalServerError)
		return
	}

	t.Execute(w, p)

}

func main() {

	var wikipages []wikiPage
	var titles []string

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
		titles = append(titles, wp.Title)
	}

	homevars := homeVariables{
		PageTitles:   titles,
		SelectedPage: "",
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		homeHandler(w, r, &homevars, wikipages)
	})
	http.HandleFunc("/view/", viewHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
