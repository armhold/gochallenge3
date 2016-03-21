package gochallenge3

import (
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
	SearchResultRows [][]ImageURL
	Error            error
	UploadID         string
	Project          *Project
}

func init() {
	templates = make(map[string]*template.Template)
	templates["welcome.html"] = template.Must(template.ParseFiles("./templates/welcome.html", "./templates/layout.html"))
	templates["results.html"] = template.Must(template.ParseFiles("./templates/results.html", "./templates/layout.html"))
	templates["404.html"] = template.Must(template.ParseFiles("./templates/404.html", "./templates/layout.html"))
	templates["500.html"] = template.Must(template.ParseFiles("./templates/500.html", "./templates/layout.html"))
	//	template.ParseGlob("./templates/*.html")
}

func homeHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := &Page{Title: "Welcome"}
		renderTemplate(w, "welcome.html", context)
	})
}

func resultsHandler(uploadRootDir string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := &Page{Title: "Mosaic Results"}

		handleErr := func(err error) {
			page.Error = err
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate(w, "welcome.html", page)
			return
		}

		parts := SplitPath(r.URL.Path)
		if len(parts) != 2 {
			handleErr(errors.New("upload_id_missing"))
			return
		}

		projectID := parts[1]
		project, err := ReadProject(uploadRootDir, projectID)
		if err != nil {
			handleErr(err)
			return
		}

		page.Project = project
		w.WriteHeader(http.StatusOK)
		renderTemplate(w, "results.html", page)
	})
}

func uploadHandler(uploadRootDir string, imageSource *InstagramClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("uploadHandler()...")

		handleErr := func(err error) {
			context := &Page{Title: "Welcome"}
			context.Error = err
			renderTemplate(w, "welcome.html", context)
		}

		project, err := createProject(uploadRootDir, r)
		if err != nil {
			handleErr(err)
			return
		}

		searchTerm := r.FormValue("search_term")
		if searchTerm == "" {
			handleErr(errors.New("enter a search_term"))
			return
		}

		imageURLs, err := imageSource.Search(searchTerm, maxInstagramSearchImages)

		// save image URLs to disk so we can use them to render mosaic, if/when the user clicks "generate"
		project.ToFile(imageURLs)

		filePaths, err := Download(imageURLs, project.ThumbnailsDir())
		for _, filePath := range filePaths {
			log.Printf("filePath: %s\n", filePath)
		}

		if err != nil {
			handleErr(err)
			return
		}

		thumbs, err := project.Thumbnails()
		if err != nil {
			handleErr(err)
			return
		}

		m := NewMosaic(50, 50, thumbs)
		err = m.Generate(project.UploadedImageFile(), project.GeneratedMosaicFile(), 10, 10)
		if err != nil {
			handleErr(err)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/results/%s", filepath.Base(project.ID)), http.StatusFound)
	})
}

func createProject(uploadRootDir string, r *http.Request) (*Project, error) {
	project, err := NewProject(uploadRootDir)
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("file")

	if err != nil {
		if err == http.ErrMissingFile {
			err = errors.New("no image file selected")
		}
		return nil, err
	}
	defer file.Close()

	return project, project.ReceiveUpload(file)
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

// TODO
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

		log.Printf("jobStatus: %s: %s", project.ID, project.Status)
	})
}

func Serve(addr, uploadRootDir string, imageSource *InstagramClient) {
	log.Printf("start server on: %s\n", addr)

	http.Handle("/", homeHandler())
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.Handle("/upload", uploadHandler(uploadRootDir, imageSource))
	http.Handle("/results/", resultsHandler(uploadRootDir))
	http.Handle("/status/", jobStatusHandler(uploadRootDir))
	http.Handle("/download/", downloadMosaicHandler(uploadRootDir))

	http.ListenAndServe(addr, nil)
}

func renderTemplate(w http.ResponseWriter, templatePath string, page *Page) {
	if page.Error != nil {
		log.Println(page.Error)
	}

	err := templates[templatePath].ExecuteTemplate(w, "layout", page)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
