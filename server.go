package gochallenge3

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"log"
)

var (
	templates     map[string]*template.Template
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
	templates["search.html"] = template.Must(template.ParseFiles("./templates/search.html", "./templates/layout.html"))
	templates["choose.html"] = template.Must(template.ParseFiles("./templates/choose.html", "./templates/layout.html"))
	templates["results.html"] = template.Must(template.ParseFiles("./templates/results.html", "./templates/layout.html"))
	templates["404.html"] = template.Must(template.ParseFiles("./templates/404.html", "./templates/layout.html"))
	templates["500.html"] = template.Must(template.ParseFiles("./templates/500.html", "./templates/layout.html"))
	//	template.ParseGlob("./templates/*.html")
}

func homeHandler() (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := &Page{Title: "Welcome"}
		renderTemplate(w, "welcome.html", context)
	})
}

func searchHandler(uploadRootDir string, imageSource *InstagramImageSource) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := &Page{Title: "Search Results"}

		parts := SplitPath(r.URL.Path)
		if len(parts) != 2 {
			page.Error = errors.New("upload_id missing")
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate(w, "search.html", page)
			return
		}

		projectID := parts[1]
		project, err := ReadProject(uploadRootDir, projectID)
		if err != nil {
			page.Error = err
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate(w, "search.html", page)
			return
		}

		page.Project = project

		searchTerm := r.FormValue("search_term")
		if searchTerm == "" {
			page.Error = errors.New("enter a search_term")
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate(w, "search.html", page)
			return
		}

		imageURLs, err := imageSource.Search(searchTerm)

		// save image URLs to disk so we can use them to render mosaic, if/when the user clicks "generate"
		page.Project.ToFile(imageURLs)

		filePaths, err := Download(imageURLs, page.Project.ThumbnailsDir())
		for _, filePath := range filePaths {
			log.Printf("filePath: %s\n", filePath)
		}

		if err != nil {
			page.Error = fmt.Errorf("error searching for images: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			renderTemplate(w, "search.html", page)
			return
		}
		page.SearchResultRows = ToRows(5, imageURLs)

		renderTemplate(w, "search.html", page)
	})
}

func chooseFileHandler() (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		context := &Page{Title: "Image Upload"}
		renderTemplate(w, "choose.html", context)
	})
}

func resultsHandler(uploadRootDir string) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := &Page{Title:"Mosaic Results"}

		parts := SplitPath(r.URL.Path)
		if len(parts) != 2 {
			err := errors.New("upload_id_missing")
			page.Error = err
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate(w, "choose.html", page)
			return
		}

		projectID := parts[1]
		project, err := ReadProject(uploadRootDir, projectID)
		if err != nil {
			page.Error = err
			w.WriteHeader(http.StatusBadRequest)
			renderTemplate(w, "choose.html", page)
			return
		}

		page.Project = project
		w.WriteHeader(http.StatusOK)
		renderTemplate(w, "results.html", page)
	})
}

func uploadHandler(uploadRootDir string) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project, err := createProject(uploadRootDir, r)

		if err != nil {
			context := &Page{Title: "Receive Upload"}
			context.Error = err
			renderTemplate(w, "choose.html", context)
		} else {
			redirectTo := fmt.Sprintf("/search/%s", filepath.Base(project.ID))
			log.Printf("redirect to %s", redirectTo)
			http.Redirect(w, r, redirectTo, http.StatusFound)
		}
	})
}

func createProject(uploadRootDir string, r *http.Request) (*Project, error) {
	project, err := NewProject(uploadRootDir)
	if err != nil {
		return nil, err
	}

	file, _, err := r.FormFile("file")

	if err != nil {
		return nil, err
	}
	defer file.Close()

	return project, project.ReceiveUpload(file)
}

func generateMosaicHandler(uploadRootDir string) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project, err := generateMosaic(uploadRootDir, r)
		if err != nil {
			context := &Page{Title: "Generate Mosaic"}
			context.Error = err
			renderTemplate(w, "choose.html", context)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/results/%s", filepath.Base(project.ID)), http.StatusFound)
		}
	})
}

func generateMosaic(uploadRootDir string, r *http.Request) (*Project, error) {
	parts := SplitPath(r.URL.Path)
	if len(parts) != 2 {
		return nil, errors.New("upload_id missing")
	}

	projectID := parts[1]
	project, err := ReadProject(uploadRootDir, projectID)
	if err != nil {
		return nil, err
	}

	thumbs, err := project.Thumbnails()
	if err != nil {
		return nil, err
	}

	m := NewMosaic(50, 50, thumbs)
	err = m.Generate(project.UploadedImageFile(), project.GeneratedMosaicFile())
	if err != nil {
		return nil, err
	}

	return project, nil
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

func Serve(addr, uploadRootDir string, imageSource *InstagramImageSource) {
	log.Printf("start server on: %s\n", addr)

	http.Handle("/", homeHandler())
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	http.Handle("/choose", chooseFileHandler())
	http.Handle("/upload", uploadHandler(uploadRootDir))
	http.Handle("/search/", searchHandler(uploadRootDir, imageSource))
	http.Handle("/generate/", generateMosaicHandler(uploadRootDir))
	http.Handle("/results/", resultsHandler(uploadRootDir))
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

