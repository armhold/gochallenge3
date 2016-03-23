// Fetches tagged images from Instagram
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/armhold/gochallenge3"
)

var (
	tag                 string
	clientID            string
	imagesDir           string
	maxResults          int
	downloadConcurrency int
)

func init() {
	flag.StringVar(&tag, "tag", "", "search tag")
	flag.StringVar(&clientID, "client", "", "instagram client id")
	flag.StringVar(&imagesDir, "images", "", "directory to save images")
	flag.IntVar(&maxResults, "max", 1000, "max result count")
	flag.IntVar(&downloadConcurrency, "dc", 10, "download concurrency")

	flag.Parse()

	if tag == "" || clientID == "" || imagesDir == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	if err := os.MkdirAll(imagesDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "unable to create images dir %s: %s\n", imagesDir, err)
		os.Exit(1)
	}

	client := &gochallenge3.InstagramClient{ClientID: clientID}
	urls, err := client.Search(tag, maxResults)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	savedFiles, err := gochallenge3.Download(urls, imagesDir, downloadConcurrency)

	if err != nil {
		log.Println(err)
	} else {
		for _, filepath := range savedFiles {
			fmt.Printf("saved to: %s\n", filepath)
		}
	}
}
