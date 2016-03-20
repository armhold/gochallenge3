package gochallenge3

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func Download(urls []ImageURL, downloadDir string) ([]string, error) {
	var filePaths = make([]string, len(urls))

	for i, url := range urls {
		log.Printf("download [%d] => \"%s\"", i, url)

		ext, err := url.guessImageExtension()
		if err != nil {
			return nil, err
		}

		filePaths[i] = fmt.Sprintf("%s/thumb%d%s", downloadDir, i, ext)

		err = downloadToFile(string(url), filePaths[i])
		if err != nil {
			return nil, err
		}
		log.Printf("downloaded %s to %s\n", url, filePaths[i])
	}

	return filePaths, nil
}

func downloadToFile(url, toFile string) error {
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
