package gochallenge3

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// Downloads the images to the download directory, using goroutines for concurrency.
// NB: to avoid overwhelming image servers, best to keep maxConcurrency <= 10. This also
// prevents "dial tcp: no such host" errors on OSX due to file descriptors maxing out.
func Download(urls []ImageURL, downloadDir string, maxConcurrency int) ([]string, error) {
	var filePaths = make([]string, len(urls))

	// create target file paths
	for i, url := range urls {
		ext, err := url.guessImageExtension()
		if err != nil {
			return nil, err
		}
		filePaths[i] = fmt.Sprintf("%s/thumb%d%s", downloadDir, i, ext)
	}

	workerTokens := make(chan struct{}, maxConcurrency) // controls concurrency- only run N downloads at a time
	errC := make(chan error)

	doDownload := func(fromUrl ImageURL, toFile string, wg *sync.WaitGroup) {

		defer func() {
			wg.Done()

			// release token
			<-workerTokens
		}()

		// block here until we acquire a token
		workerTokens <- struct{}{}

		err := downloadToFile(string(fromUrl), toFile)
		if err != nil {
			errC <- err
			return
		}
		log.Printf("downloaded %s to %s\n", fromUrl, toFile)
	}

	// set up a WaitGroup to wait out the workers, so that we can close errC when they are done.
	// see explanation at https://blog.golang.org/pipelines
	var wg sync.WaitGroup

	for i, url := range urls {
		wg.Add(1)
		go doDownload(url, filePaths[i], &wg)
	}

	go func() {
		wg.Wait()
		close(errC)
	}()

	log.Printf("all download workers complete")


	// check to see if any workers had an error.
	// NB: this is necessary since errC is non-buffered, and workers will be blocked trying to write to it.
	var err error
	for err = range errC {
		log.Printf("error during download: %s", err)
	}

	if err == nil {
		log.Println("all downloads completed successfully")
	}

	return filePaths, err
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
