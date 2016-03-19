package gochallenge3

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
	"bufio"
	"path/filepath"
	"log"
)

// Project represents a mosaic project- the uploaded file, selected tile images, and resulting mosaic image
type Project struct {
	ID            string
	uploadRootDir string
}

func ReadProject(uploadRootDir, id string) (*Project, error) {
	var result = &Project{ID: id, uploadRootDir: uploadRootDir}

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
	id, err := randomString()
	if err != nil {
		return nil, fmt.Errorf("error generating project ID: %v", err)
	}

	var result = &Project{ID: id, uploadRootDir: uploadRootDir}

	err = os.Mkdir(result.UploadedImageDir(), os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating upload dir: %v", result.UploadedImageDir())
	}

	err = os.Mkdir(result.ThumbnailsDir(), os.ModeDir|os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error creating thumbnails dir: %v", result.ThumbnailsDir())
	}

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

func randomString() (string, error) {
	rb := make([]byte, 8)
	_, err := rand.Read(rb)

	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(rb), nil
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

func (p *Project) FromFile() ([]string, error) {
	file, err := os.Open(p.ImageUrlsFile())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		result = append(result, url)
	}

	return result, scanner.Err()
}

