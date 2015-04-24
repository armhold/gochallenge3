package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"github.com/armhold/gochallenge3"
	"log"
	"github.com/juju/errgo/errors"
)

var (
	templates = template.Must(template.ParseFiles("../../static/welcome.html", "../../static/upload.html", "../../static/search.html"))
)

type Page struct {
	Title string
	SearchResults []gochallenge3.InstagramImageSet
	Error error
	Body  []byte
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Welcome"}
	renderTemplate(w, "welcome.html", p)
}


func searchHandler(imageSource gochallenge3.ImageSource) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := &Page{Title: "Search Results"}

		searchTerm := r.FormValue("search_term")
		if searchTerm == "" {
			p.Error = errors.New("search_term required")
			// TODO: support multiple errors in p.Error
		} else {
			imageSets, err := imageSource.Search(searchTerm)
			if err != nil {
				log.Printf("error searching for images: %v\n", err)
				p.Error = err
			} else {
				p.SearchResults = imageSets
			}
		}

		renderTemplate(w, "search.html", p)
	})
}

func main() {
	http.HandleFunc("/", homeHandler)

	instagramClientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	if instagramClientID == "" {
		panic("environment variable INSTAGRAM_CLIENT_ID not set")
	}
	imageSource := gochallenge3.NewInstagramImageSource(instagramClientID)

    http.HandleFunc("/search", searchHandler(imageSource))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("using port: %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func renderTemplate(w http.ResponseWriter, templatePath string, p *Page) {
	err := templates.ExecuteTemplate(w, templatePath, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
