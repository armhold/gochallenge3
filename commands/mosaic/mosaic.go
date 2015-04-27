package main

import (
	"fmt"
	"github.com/armhold/gochallenge3"
	"html/template"
	"net/http"
	"os"
	"io"
	"crypto/rand"
	"encoding/base64"
	"path/filepath"
	"errors"
)

var (
	templates map[string]*template.Template
	uploadDir string
)

type Page struct {
	Title string
	SearchResultRows [][]gochallenge3.InstagramImageSet
	Error            error
	Body             []byte
	UploadID         string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Welcome"}
	renderTemplate(w, "welcome.html", p)
}

func init() {
	templates = make(map[string]*template.Template)
	templates["welcome.html"] = template.Must(template.ParseFiles("../../templates/welcome.html", "../../templates/layout.html"))
	templates["search.html"]  = template.Must(template.ParseFiles("../../templates/search.html", "../../templates/layout.html"))
	templates["choose.html"]  = template.Must(template.ParseFiles("../../templates/choose.html", "../../templates/layout.html"))
}

func searchHandler(imageSource gochallenge3.ImageSource) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := &Page{Title: "Search Results"}

		parts := gochallenge3.SplitPath(r.URL.Path)
		gochallenge3.CommonLog.Printf("r.URL.Path: %s, parts: %+v, len: %d", r.URL.Path, parts, len(parts))
		if len(parts) != 2 {
			err := errors.New("upload_id missing")
			gochallenge3.CommonLog.Println(err)
			p.Error = err
		} else {
			p.UploadID = parts[1]
		}

		searchTerm := r.FormValue("search_term")
		if searchTerm != "" {
			imageSets, err := imageSource.Search(searchTerm)

			urls := make([]string, len(imageSets))
			for i, imageSet := range imageSets {
				gochallenge3.CommonLog.Printf("processing thumbnail: %s", imageSet.Thumb.Url)
				urls[i] = imageSet.Thumb.Url
			}

			filePaths, err := gochallenge3.Download(urls)
			for _, filePath := range filePaths {
				gochallenge3.CommonLog.Printf("filePath: %s\n", filePath)
			}

			if err != nil {
				gochallenge3.CommonLog.Printf("error searching for images: %v\n", err)
				p.Error = err
			} else {
				p.SearchResultRows = gochallenge3.ToRows(5, imageSets)
			}
		}

		renderTemplate(w, "search.html", p)
	})
}

func chooseFileHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Image Upload"}
	renderTemplate(w, "choose.html", p)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Receive Upload"}

	file, err := receiveUploadFile(r)

	if err != nil {
		gochallenge3.CommonLog.Println(err)
		p.Error = err
		renderTemplate(w, "choose.html", p)
	} else {
		gochallenge3.CommonLog.Printf("Image successfully uploaded to %s", file.Name())
		http.Redirect(w, r, fmt.Sprintf("/search/%s", filepath.Base(file.Name())), http.StatusFound)
	}
}

func randomString() (string, error) {
	rb := make([]byte, 8)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
}

func receiveUploadFile(r *http.Request) (*os.File, error) {
	file, _, err := r.FormFile("file")

	if err != nil {
		return nil, err
	}
	defer file.Close()

	rs, err := randomString()
	if err != nil {
		return nil, fmt.Errorf("error generating upload filename: %v", err)
	}

	uploadFile := uploadDir + "/" + rs
	out, err := os.Create(uploadFile)
	if err != nil {
		return nil, fmt.Errorf("error creating upload file: %v", uploadFile)
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return nil, err
	}

	return out, nil
}



func main() {
	http.HandleFunc("/", homeHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("../../public"))))

	instagramClientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	if instagramClientID == "" {
		panic("environment variable INSTAGRAM_CLIENT_ID not set")
	}
	imageSource := gochallenge3.NewInstagramImageSource(instagramClientID)

	uploadDir = os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		panic("environment variable UPLOAD_DIR not set")
	}


	http.HandleFunc("/search/", searchHandler(imageSource))
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/choose", chooseFileHandler)


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("using port: %s\n", port)
	http.ListenAndServe(":"+port, nil)
}

func renderTemplate(w http.ResponseWriter, templatePath string, p *Page) {
	err := templates[templatePath].ExecuteTemplate(w, "layout", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
