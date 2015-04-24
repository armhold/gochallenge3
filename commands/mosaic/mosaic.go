package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var templates = template.Must(template.ParseFiles("../../static/welcome.html", "../../static/upload.html"))

type Page struct {
	Title string
	Body  []byte
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Welcome"}
	renderTemplate(w, "welcome.html", p)
}

func main() {
	http.HandleFunc("/", homeHandler)
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
