package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gocolly/colly"
)

func ping(w http.ResponseWriter, r *http.Request) {
	log.Println("Ping")
	w.Write([]byte("ping"))
}

func getData(w http.ResponseWriter, r *http.Request) {
	// Make sure the URL exists
	URL := r.URL.Query().Get("url")
	if URL == "" {
		log.Println("missing URL argument")
		return
	}
	log.Println("visiting", URL)

	// Create a new collector to collect the data from the HTML
	c := colly.NewCollector()

	//Array to store the data in
	var response []string

	// onHTML function lets the collector use a callback function when
	// the specific HTML tag is reached.
	// In this case whenever the collector finds an anchor tag with href
	// it will call the anonymous function specified below.
	// This function will get the info from the href and add it to the array
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link != "" {
			response = append(response, link)
		}
	})

	// Visit the website
	c.Visit(URL)

	// parse the response array into JSON
	b, err := json.Marshal(response)
	if err != nil {
		log.Println("failed to seralize response:", err)
		return
	}

	// Add some headers and write the body for the endpoint
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func main() {
	addr := ":7171"

	http.HandleFunc("/ping", ping)
	http.HandleFunc("/search", getData)

	log.Println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
