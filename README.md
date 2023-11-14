# Week 8 Assignment: Application Development

## Project Summary
This assignment asked us to create an information application in the form of a web application using Go as the backend. To reduce the amount of front-end coding in JavaScript, we use server-side rendering in Go.

Data scientists are often asked to act a full-stack developers, as customers need a way to interact with the information product that has been developed. In my experience, I have had to prototype a front-end for projects to demonstrate capabilities before hiring a contract front-end developer onto the team for full implementation. This week's assignment excited me as utilizing Go has the potential to streamline the this process, allowing data scientists to concentrate on what they do best.

## Development Steps
1. **Identify Problem and User Needs**: Our firm has decided to create its own online library, a knowledge base focused on its current research and development efforts into intelligent systems and robotics (Week 5 assignment). We need to build an application that allows users to search and display this library.
    1. *Selection tool*: Display list of available documents and allow user to select one for display.
    2. *Search*: Search box allows user to search for terms and returns relevant documents. Requires documents to be indexed (Elasticsearch). *Not completed*
2. **Front-end Design and Development**: Design a simple interface using Figma.
3. **Back-end Development**: Server side rendering in Go.

## Current State
Web development is not easy and using server-side rendering in Go also requires a big learning curve. The current state of this application is not pretty, however, suffices as a simple information application. Aspirations for this assignment were high and I hope to revisit the project to improve and learn more about the server-side rendering in Go.  

This current version simply reads in a json lines file that contains the URLS and titles that make up the firm's library. The titles are displayed in the drop down menu of the selection tool. If a copy with the title selected is not currently in the firm's library (i.e. saved in the wikipages folder), then the wikipage will be scraped using the [Colly](https://github.com/gocolly/colly) framework and saved. The application will then load the content from the saved file and display it on a new page.

## Next Steps
1. Scrape body text, dropping HTML tags and file formatting to improve readability. This will also be valuable for indexing the documents for searching.
2. Add styling to the site with cascading style sheets (CSS). Start to make it look like the front-end designed in Figma.
3. Index the documents in Elasticsearch. 
4. Incoporate search capabilities into application.

## Files
### *main.go*
This file defines the handler functions for the homepage and view pages for our application. It also defines the data structure for the wikiPage variable and the homeVariables which holds data for the home page of the application. The main function reads in the jsonlines file and initiates a slice of wikiPages to hold the data. It also initiates a homeVariables instance containing a slice of the available document titles which is used to populate the drop down selection on the home page.

### *scrape_wiki.go*
This file defines the scrapeData() function. It takes in a wikiPage struct and creates a HTTP client using Colly to parse the wikiPage.Url. The output is the parsed html content.

### *utils.go*
This file defines the majority of the functions for handling the data and calls for the application.  
*save()* saves the Body of the selected wikiPage to a text file using the Title as the file name. 
*loadPage()* loades the Body text from the saved file in order to display the text in the view template.  
*findWikiPageByTitle()* is a helper function to find the wikipage associated with the title selected in the application.  
*getWikiPageBody()* loads the body of the document. It determines whether the body already exists in memory, needs to be loaded from library, or needs to be scraped. Returns a completed wikiPage object.

### *main_test.go*
This file runs a test for a number of the functions.

### *mod8.exe*
Executable file of cross-compiled Go code for Mac/Windows.

### *images/Week 8 Assignment.pdf*
This is the exported file from Figma of an initial front-end design.


