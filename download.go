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
        CommonLog.Printf("download [%d] => \"%s\"", i, url)

        response, err := http.Get(url)
        if err != nil {
            return nil, fmt.Errorf("error while downloading %s: %s", url, err)
        }
        defer response.Body.Close()

        filePaths[i] = fmt.Sprintf("%s/thumb.%d", tmpDir, i)

        output, err := os.Create(filePaths[i])
        if err != nil {
            return nil, fmt.Errorf("error while creating %s: %s", filePaths[i], err)
        }
        defer output.Close()

        _, err = io.Copy(output, response.Body)
        if err != nil {
            return nil, fmt.Errorf("error while downloading", url, "-", err)
        }

        CommonLog.Printf("downloaded %s to %s\n", url, filePaths[i])
    }


    return filePaths, nil
}
