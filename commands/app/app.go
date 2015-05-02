package main

import (
	"errors"
	"fmt"
	"github.com/armhold/gochallenge3"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

var (
	templates     map[string]*template.Template
	uploadRootDir string
)

type Page struct {
	Title            string
	SearchResultRows [][]gochallenge3.ImageURL
	Error            error
	UploadID         string
	Project          *gochallenge3.Project
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Welcome"}
	renderTemplate(w, "welcome.html", p)
}

func init() {
	templates = make(map[string]*template.Template)
	templates["welcome.html"]   = template.Must(template.ParseFiles("../../templates/welcome.html",   "../../templates/layout.html"))
	templates["search.html"]    = template.Must(template.ParseFiles("../../templates/search.html",    "../../templates/layout.html"))
	templates["choose.html"]    = template.Must(template.ParseFiles("../../templates/choose.html",    "../../templates/layout.html"))
	templates["generated.html"] = template.Must(template.ParseFiles("../../templates/generated.html", "../../templates/layout.html"))
}

func searchHandler(imageSource gochallenge3.ImageSource) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := &Page{Title: "Search Results"}

		parts := gochallenge3.SplitPath(r.URL.Path)
		if len(parts) != 2 {
			err := errors.New("upload_id missing")
			gochallenge3.CommonLog.Println(err)
			p.Error = err
		} else {
			projectID := parts[1]
			project, err := gochallenge3.ReadProject(uploadRootDir, projectID)
			if err != nil {
				gochallenge3.CommonLog.Println(err)
				p.Error = err
			} else {
				p.Project = project
			}
		}

		searchTerm := r.FormValue("search_term")
		if searchTerm != "" {
			imageURLs, err := imageSource.Search(searchTerm)

			// save image URLs to disk so we can use them to render mosaic, if/when the user clicks "generate"
			p.Project.ToFile(imageURLs)

			filePaths, err := gochallenge3.Download(imageURLs, p.Project.ThumbnailsDir())
			for _, filePath := range filePaths {
				gochallenge3.CommonLog.Printf("filePath: %s\n", filePath)
			}

			if err != nil {
				gochallenge3.CommonLog.Printf("error searching for images: %v\n", err)
				p.Error = err
			} else {
				p.SearchResultRows = gochallenge3.ToRows(5, imageURLs)
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

	project, err := createProject(r)

	if err != nil {
		gochallenge3.CommonLog.Println(err)
		p.Error = err
		renderTemplate(w, "choose.html", p)
	} else {
		p.Project = project
		http.Redirect(w, r, fmt.Sprintf("/search/%s", filepath.Base(project.ID)), http.StatusFound)
	}
}

func createProject(r *http.Request) (*gochallenge3.Project, error) {
	project, err := gochallenge3.NewProject(uploadRootDir)
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

func generateMosaicHandler(w http.ResponseWriter, r *http.Request) {
	p := &Page{Title: "Generate Mosaic"}

	project, err := generateMosaic(w, r)
	if err != nil {
		gochallenge3.CommonLog.Println(err)
		p.Error = err
	}
	p.Project = project

	renderTemplate(w, "generated.html", p)
}

func generateMosaic(w http.ResponseWriter, r *http.Request) (*gochallenge3.Project, error) {
	parts := gochallenge3.SplitPath(r.URL.Path)
	if len(parts) != 2 {
		return nil, errors.New("upload_id missing")
	}

	projectID := parts[1]
	project, err := gochallenge3.ReadProject(uploadRootDir, projectID)
	if err != nil {
		return nil, err
	}

	thumbs, err := project.Thumbnails()
	if err != nil {
		return nil, err
	}

	m := gochallenge3.NewMosaic(50, 50, thumbs)
	err = m.Generate(project.UploadedImageFile(), project.GeneratedMosaicFile())
	if err != nil {
		return nil, err
	}

	return project, nil
}

func downloadMosaicHandler(w http.ResponseWriter, r *http.Request) {
	parts := gochallenge3.SplitPath(r.URL.Path)
	if len(parts) != 2 {
		err := "upload_id missing"
		gochallenge3.CommonLog.Println(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	projectID := parts[1]
	project, err := gochallenge3.ReadProject(uploadRootDir, projectID)
	if err != nil {
		gochallenge3.CommonLog.Println(err)
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, project.GeneratedMosaicFile())
}



func main() {
	if uploadRootDir = os.Getenv("UPLOAD_DIR"); uploadRootDir == "" {
		panic("environment variable UPLOAD_DIR not set")
	}

	http.HandleFunc("/", homeHandler)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("../../public"))))

	instagramClientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	if instagramClientID == "" {
		panic("environment variable INSTAGRAM_CLIENT_ID not set")
	}
	imageSource := gochallenge3.NewInstagramImageSource(instagramClientID)

	http.HandleFunc("/search/", searchHandler(imageSource))
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/choose", chooseFileHandler)
	http.HandleFunc("/generate/", generateMosaicHandler)
	http.HandleFunc("/download/", downloadMosaicHandler)

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
