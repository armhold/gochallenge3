package gochallenge3

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"
)

type Status string

const (
	StatusNew         Status = "new"
	StatusSearching   Status = "searching images"
	StatusDownloading Status = "downloading images"
	StatusGenerating  Status = "generating mosaic"
	StatusError       Status = "error"
	StatusCompleted   Status = "completed"
)

const (
	idCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Project represents a mosaic project- the uploaded file, selected tile images, and resulting mosaic image
type Project struct {
	ID            string
	uploadRootDir string
	Status
}

func ReadProject(uploadRootDir, id string) (*Project, error) {
	var result = &Project{ID: id, uploadRootDir: uploadRootDir}
	err := result.LoadStatus()
	if err != nil {
		return nil, err
	}

	// check if the image dir exists
	fileInfo, err := os.Stat(result.UploadedImageDir())
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", result.UploadedImageDir())
	}

	return result, nil
}

func NewProject(uploadRootDir string) (*Project, error) {
	id := randomString()
	var result = &Project{ID: id, uploadRootDir: uploadRootDir}

	err := os.Mkdir(result.UploadedImageDir(), os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating upload dir: %v", result.UploadedImageDir())
	}

	err = os.Mkdir(result.ThumbnailsDir(), os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating thumbnails dir: %v", result.ThumbnailsDir())
	}

	result.SetAndSaveStatus(StatusNew)

	return result, nil
}

func (p *Project) ReceiveUpload(r io.Reader) error {
	log.Printf("ReceiveUpload to: %s", p.UploadedImageFile())

	out, err := os.Create(p.UploadedImageFile())
	if err != nil {
		return fmt.Errorf("error creating upload file: %v", p.UploadedImageFile())
	}
	defer out.Close()

	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}

	log.Printf("Image successfully uploaded to %s", p.UploadedImageFile())

	return nil
}

func (p *Project) ThumbnailsDir() string {
	return path.Join(p.uploadRootDir, p.ID, "thumbs")
}

func (p *Project) ImageUrlsFile() string {
	return path.Join(p.uploadRootDir, p.ID, "image_urls.txt")
}

func (p *Project) UploadedImageFile() string {
	log.Printf("UploadedImageFile: %s, %s, %s", p.uploadRootDir, p.ID, "uploaded_image")

	return path.Join(p.uploadRootDir, p.ID, "uploaded_image")
}

func (p *Project) GeneratedMosaicFile() string {
	return path.Join(p.uploadRootDir, p.ID, "mosaic")
}

func (p *Project) UploadedImageDir() string {
	return path.Join(p.uploadRootDir, p.ID)
}

func (p *Project) Thumbnails() ([]string, error) {
	return filepath.Glob(p.ThumbnailsDir() + "/thumb*")
}

// adapted from http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randomString() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = idCharacters[rand.Intn(len(idCharacters))]
	}
	return string(b)
}

func (p *Project) ToFile(urls []ImageURL) error {
	file, err := os.Create(p.ImageUrlsFile())

	if err != nil {
		return err
	}
	defer file.Close()

	for _, url := range urls {
		line := string(url) + "\n"
		_, err := io.WriteString(file, line)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Project) SetAndSaveStatus(status Status) error {
	p.Status = status
	return ioutil.WriteFile(p.statusFilePath(), []byte(p.Status), 0644)
}

func (p *Project) LoadStatus() error {
	s, err := ioutil.ReadFile(p.statusFilePath())
	if err != nil {
		return err
	}

	p.Status = Status(s)
	return nil
}

func (p *Project) statusFilePath() string {
	return path.Join(p.UploadedImageDir(), "status")
}
