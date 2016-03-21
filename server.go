package gochallenge3

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

const (
	maxInstagramSearchImages = 100
)

var (
	templates map[string]*template.Template
)

type Page struct {
	Title            string
	UploadID         string
	Project          *Project
	template         string
	redirectTo       string
	Error            error
	httpStatusCode   int
}

func init() {
	templates = make(map[string]*template.Template)
	templates["welcome.html"] = template.Must(template.ParseFiles("./templates/welcome.html", "./templates/layout.html"))
	templates["results.html"] = template.Must(template.ParseFiles("./templates/results.html", "./templates/layout.html"))
	templates["404.html"] = template.Must(template.ParseFiles("./templates/404.html", "./templates/layout.html"))
	templates["500.html"] = template.Must(template.ParseFiles("./templates/500.html", "./templates/layout.html"))
}

// simplify error handling, see: http://blog.golang.org/error-handling-and-go
type appHandler func(http.ResponseWriter, *http.Request) (*Page)

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	page := fn(w, r)

	if page.Error != nil {
		log.Println(page.Error)
		w.WriteHeader(page.httpStatusCode)
	}

	if page.redirectTo != "" {
		http.Redirect(w, r, page.redirectTo, http.StatusFound)
	} else if page.template != "" {
		err := templates[page.template].ExecuteTemplate(w, "layout", page)
		if err != nil {
			log.Printf("error rendering template %s: %s", page.template, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) (page *Page) {
	return &Page{Title: "Welcome", template: "welcome.html"}
}

func resultsHandler(uploadRootDir string) appHandler {
	return func(w http.ResponseWriter, r *http.Request) (page *Page) {
		page = &Page{Title: "Mosaic Results", template: "welcome.html"}

		parts := SplitPath(r.URL.Path)
		if len(parts) != 2 {
			page.Error = errors.New("upload_id_missing")
			page.httpStatusCode = http.StatusBadRequest
			return
		}

		projectID := parts[1]
		project, err := ReadProject(uploadRootDir, projectID)
		if err != nil {
			page.Error = err
			page.httpStatusCode = http.StatusBadRequest
			return
		}

		page.Project = project
		page.httpStatusCode = http.StatusOK
		page.template = "results.html"
		return
	}
}

func uploadHandler(uploadRootDir string, imageSource *InstagramClient) appHandler {
	return func(w http.ResponseWriter, r *http.Request) (page *Page) {
		log.Printf("uploadHandler...")

		page = &Page{Title: "Welcome"}

		project, err := NewProject(uploadRootDir)
		if err != nil {
			page.Error = err
			page.template = "welcome.html"
			page.httpStatusCode = http.StatusInternalServerError
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			page.template = "welcome.html"

			if err == http.ErrMissingFile {
				page.Error = errors.New("no image file selected")
				page.httpStatusCode = http.StatusBadRequest
			} else {
				page.Error = err
				page.httpStatusCode = http.StatusInternalServerError
			}

			return
		}
		defer file.Close()
		project.ReceiveUpload(file)

		searchTerm := r.FormValue("search_term")
		if searchTerm == "" {
			page.Error = errors.New("enter a search_term")
			page.template = "welcome.html"
			page.httpStatusCode = http.StatusBadRequest
			return
		}

		project.SetAndSaveStatus(StatusSearching)

		go processMosaic(searchTerm, project, imageSource)
		page.redirectTo = fmt.Sprintf("/results/%s", filepath.Base(project.ID))

		return
	}
}

// Once we have everything from the user, we can process the search, download & generation offline.
// This is intended to be run from a goroutine; we report success/error via project status.
func processMosaic(searchTerm string, project *Project, imageSource *InstagramClient) {
	imageURLs, err := imageSource.Search(searchTerm, maxInstagramSearchImages)
	if err != nil {
		log.Println(err)
		project.SetAndSaveStatus(StatusError)
		return
	}

	// save image URLs to disk so we can use them to render mosaic, if/when the user clicks "generate"
	err = project.ToFile(imageURLs)
	if err != nil {
		log.Println(err)
		project.SetAndSaveStatus(StatusError)
		return
	}

	project.SetAndSaveStatus(StatusDownloading)
	_, err = Download(imageURLs, project.ThumbnailsDir())
	if err != nil {
		log.Println(err)
		project.SetAndSaveStatus(StatusError)
		return
	}

	thumbs, err := project.Thumbnails()
	if err != nil {
		log.Println(err)
		project.SetAndSaveStatus(StatusError)
		return
	}

	project.SetAndSaveStatus(StatusGenerating)
	m := NewMosaic(50, 50, thumbs)
	err = m.Generate(project.UploadedImageFile(), project.GeneratedMosaicFile(), 10, 10)
	if err != nil {
		log.Println(err)
		project.SetAndSaveStatus(StatusError)
		return
	}

	project.SetAndSaveStatus(StatusCompleted)
	log.Printf("project: %s completed successfully", project.ID)
}

func downloadMosaicHandler(uploadRootDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := SplitPath(r.URL.Path)
		if len(parts) != 2 {
			err := "upload_id missing"
			log.Println(err)
			http.Error(w, err, http.StatusBadRequest)
			return
		}

		projectID := parts[1]
		project, err := ReadProject(uploadRootDir, projectID)
		if err != nil {
			log.Println(err)
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, project.GeneratedMosaicFile())
	})
}

func jobStatusHandler(uploadRootDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		parts := SplitPath(r.URL.Path)
		if len(parts) != 2 {
			err := "upload_id missing"
			log.Println(err)
			http.Error(w, err, http.StatusBadRequest)
			return
		}

		projectID := parts[1]
		project, err := ReadProject(uploadRootDir, projectID)
		if err != nil {
			log.Println(err)
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json, err := json.Marshal(project)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonString := string(json)
		fmt.Fprint(w, jsonString)
		log.Printf("jobStatus: %s:", jsonString)
	})
}

func Serve(addr, uploadRootDir string, imageSource *InstagramClient) {
	log.Printf("start server on: %s\n", addr)

	http.Handle("/", appHandler(homeHandler))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.Handle("/upload", uploadHandler(uploadRootDir, imageSource))
	http.Handle("/results/", resultsHandler(uploadRootDir))
	http.Handle("/status/", jobStatusHandler(uploadRootDir))
	http.Handle("/download/", downloadMosaicHandler(uploadRootDir))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("error from http:ListenAndServe(): %s", err)
	}
}
