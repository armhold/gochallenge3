package gochallenge3
import (
    "path"
    "os"
    "encoding/base64"
    "crypto/rand"
    "fmt"
    "io"
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
        err = fmt.Errorf("%s is not a directory", result.UploadedImageDir())
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
    CommonLog.Printf("ReceiveUpload to: %s", p.UploadedImageFile())

    out, err := os.Create(p.UploadedImageFile())
    if err != nil {
        return fmt.Errorf("error creating upload file: %v", p.UploadedImageFile())
    }
    defer out.Close()

    _, err = io.Copy(out, r)
    if err != nil {
        return err
    }

    CommonLog.Printf("Image successfully uploaded to %s", p.UploadedImageFile())

    return nil
}


func (p *Project) ThumbnailsDir() string {
    return path.Join(p.uploadRootDir, p.ID, "thumbs")
}

func (p *Project) ImageUrlsFile() string {
    return path.Join(p.uploadRootDir, p.ID, "image_urls.txt")
}

func (p *Project) UploadedImageFile() string {
    fmt.Printf("UploadedImageFile: %s, %s, %s", p.uploadRootDir, p.ID, "uploaded_image")

    return path.Join(p.uploadRootDir, p.ID, "uploaded_image")
}

func (p *Project) UploadedImageDir() string {
    return path.Join(p.uploadRootDir, p.ID)
}

func randomString() (string, error) {
    rb := make([]byte, 8)
    _, err := rand.Read(rb)

    if err != nil {
        return "", err
    }

    return base64.URLEncoding.EncodeToString(rb), nil
}
