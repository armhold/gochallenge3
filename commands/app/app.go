package main

import (
	"errors"
	"fmt"
	"github.com/armhold/gochallenge3"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"flag"
	"log"
)

var (
	templates     map[string]*template.Template
	uploadRootDir string
	devMode bool
)

type appContext struct {
	imageSource gochallenge3.ImageSource
	templates map[string]*template.Template

	Title            string
	SearchResultRows [][]gochallenge3.ImageURL
	Error            error
	UploadID         string
	Project          *gochallenge3.Project
}


type appHandler struct {
	*appContext
	h func(*appContext, http.ResponseWriter, *http.Request) (status int, template string, error error)
}



func homeHandler(w http.ResponseWriter, r *http.Request) {
	context := &appContext{Title: "Welcome"}
	renderTemplate(w, "welcome.html", context)
}

func init() {
	flag.BoolVar(&devMode, "dev", false, "start the server in devmode")
	flag.Parse()

	uploadRootDir = os.Getenv("UPLOAD_DIR")
	if uploadRootDir == "" {
		uploadRootDir = "/tmp/upload_dir"
	}
	if err := os.MkdirAll(uploadRootDir, 0700) ; err != nil {
		log.Fatalf("unable to create temp dir %s: %s", uploadRootDir, err)
	}

	templates = make(map[string]*template.Template)
	templates["welcome.html"]   = template.Must(template.ParseFiles("../../templates/welcome.html",   "../../templates/layout.html"))
	templates["search.html"]    = template.Must(template.ParseFiles("../../templates/search.html",    "../../templates/layout.html"))
	templates["choose.html"]    = template.Must(template.ParseFiles("../../templates/choose.html",    "../../templates/layout.html"))
	templates["results.html"]   = template.Must(template.ParseFiles("../../templates/results.html",   "../../templates/layout.html"))
	templates["404.html"]       = template.Must(template.ParseFiles("../../templates/404.html",       "../../templates/layout.html"))
	templates["500.html"]       = template.Must(template.ParseFiles("../../templates/500.html",       "../../templates/layout.html"))
	template.ParseGlob("../../*.html")
}

func searchHandler(context *appContext, w http.ResponseWriter, r *http.Request) (int, string, error) {
	context.Title = "Search Results"

	parts := gochallenge3.SplitPath(r.URL.Path)
	if len(parts) != 2 {
		return http.StatusBadRequest, "search.html", errors.New("upload_id missing")
	}

	projectID := parts[1]
	project, err := gochallenge3.ReadProject(uploadRootDir, projectID)
	if err != nil {
		return http.StatusInternalServerError, "search.html", err
	}

	context.Project = project

	searchTerm := r.FormValue("search_term")
	if searchTerm != "" {
		imageURLs, err := context.imageSource.Search(searchTerm)

		// save image URLs to disk so we can use them to render mosaic, if/when the user clicks "generate"
		context.Project.ToFile(imageURLs)

		filePaths, err := gochallenge3.Download(imageURLs, context.Project.ThumbnailsDir())
		for _, filePath := range filePaths {
			gochallenge3.CommonLog.Printf("filePath: %s\n", filePath)
		}
		if err != nil {
			return http.StatusInternalServerError, "search.html", fmt.Errorf("error searching for images: %v\n", err)
		}
		context.SearchResultRows = gochallenge3.ToRows(5, imageURLs)
	}

	return http.StatusOK, "search.html", nil
}

func chooseFileHandler(w http.ResponseWriter, r *http.Request) {
	context := &appContext{Title: "Image Upload"}
	renderTemplate(w, "choose.html", context)
}

func resultsHandler(context *appContext, w http.ResponseWriter, r *http.Request) (int, string, error) {
	context.Title = "Mosaic Results"

	parts := gochallenge3.SplitPath(r.URL.Path)
	if len(parts) != 2 {
		return http.StatusBadRequest, "results.html", errors.New("upload_id_missing")
	}

	projectID := parts[1]
	project, err := gochallenge3.ReadProject(uploadRootDir, projectID)
	if err != nil {
		return http.StatusBadRequest, "results.html", err
	}

	context.Project = project

	return http.StatusOK, "results.html", nil
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	context := &appContext{Title: "Receive Upload"}

	project, err := createProject(r)

	if err != nil {
		gochallenge3.CommonLog.Println(err)
		context.Error = err
		renderTemplate(w, "choose.html", context)
	} else {
		context.Project = project

		redirectTo := fmt.Sprintf("/search/%s", filepath.Base(project.ID))

		gochallenge3.CommonLog.Printf("redirect to %s", redirectTo)
		http.Redirect(w, r, redirectTo, http.StatusFound)
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
	context := &appContext{Title: "Generate Mosaic"}

	project, err := generateMosaic(w, r)
	if err != nil {
		gochallenge3.CommonLog.Println(err)
		context.Error = err
	}
	context.Project = project

	http.Redirect(w, r, fmt.Sprintf("/results/%s", filepath.Base(project.ID)), http.StatusFound)
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

// implement http.Handler
func (ah appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, template, err := ah.h(ah.appContext, w, r)

	if err != nil {
		gochallenge3.CommonLog.Println(err)
		w.WriteHeader(status)
	}

	if template != "" {
		renderTemplate(w, template, ah.appContext)
	}
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


	context := &appContext{imageSource: imageSource}

	http.HandleFunc("/choose",    chooseFileHandler)
	http.HandleFunc("/upload",    uploadHandler)
	http.Handle("/search/",       appHandler{context, searchHandler})
	http.HandleFunc("/generate/", generateMosaicHandler)
	http.Handle("/results/",      appHandler{context, resultsHandler})
	http.HandleFunc("/download/", downloadMosaicHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("using port: %s\n", port)

	addr := ":" + port
	log.Printf("devMode: %t", devMode)

	if devMode {
		// prevent OSX firewall popups by binding to localhost
		//
		// see: http://grokbase.com/t/gg/golang-nuts/15322dedhn/go-nuts-mac-firewall
		addr = "localhost" + addr
	}

	http.ListenAndServe(addr, nil)
}

func renderTemplate(w http.ResponseWriter, templatePath string, context *appContext) {
	err := templates[templatePath].ExecuteTemplate(w, "layout", context)
	if err != nil {
		gochallenge3.CommonLog.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
