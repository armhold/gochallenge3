package gochallenge3
import (
    "io/ioutil"
    "net/http"
    "fmt"
    "io"
    "os"
)


func Download(urls []string) ([]string, error) {
    tmpDir, err := ioutil.TempDir("", "thumbnails")
    if err != nil {
        return nil, err
    }

    CommonLog.Printf("created tempDir: %s\n", tmpDir)

    var filePaths = make([]string, len(urls))

    for i, url := range urls {
        response, err := http.Get(url)
        if err != nil {

            return nil, fmt.Errorf("error while downloading", url, "-", err)
        }
        defer response.Body.Close()

        fileName := fmt.Sprintf("%s/thumb.%d", tmpDir, i)
        output, err := os.Create(fileName)
        if err != nil {
            return nil, fmt.Errorf("error while creating", fileName, "-", err)
        }
        defer output.Close()

        _, err = io.Copy(output, response.Body)
        if err != nil {
            return nil, fmt.Errorf("error while downloading", url, "-", err)
        }

        CommonLog.Printf("downloaded %s to %s\n", url, fileName)
    }


    return filePaths, nil
}
