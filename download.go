package gochallenge3

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(urls []ImageURL, downloadDir string) ([]string, error) {
	var filePaths = make([]string, len(urls))

	for i, url := range urls {
		CommonLog.Printf("download [%d] => \"%s\"", i, url)

		filePaths[i] = fmt.Sprintf("%s/thumb.%d", downloadDir, i)

		err := downloadToFile(string(url), filePaths[i])
		if err != nil {
			return nil, err
		}
		CommonLog.Printf("downloaded %s to %s\n", url, filePaths[i])
	}

	return filePaths, nil
}

func downloadToFile(url, toFile string) (error) {
	response, err := http.Get(string(url))
	if err != nil {
		return fmt.Errorf("error while downloading %s: %s", url, err)
	}
	defer response.Body.Close()

	output, err := os.Create(toFile)
	if err != nil {
		return fmt.Errorf("error while creating %s: %s", toFile, err)
	}
	defer output.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("error while downloading %s: %v", url, err)
	}

	return nil
}
