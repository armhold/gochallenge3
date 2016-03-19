package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/armhold/gochallenge3"
)

var (
	imagesDir             string
	sourceImage           string
	outFile               string
	tileWidth, tileHeight int
	widthMult, heightMult int
)

func init() {
	flag.StringVar(&imagesDir, "images", "", "directory containing source images")
	flag.StringVar(&sourceImage, "source", "", "source image file")
	flag.StringVar(&outFile, "out", "", "output image file")
	flag.IntVar(&tileWidth, "tw", 50, "tile width")
	flag.IntVar(&tileHeight, "th", 50, "tile height")
	flag.IntVar(&widthMult, "wm", 2, "width multiplier (image will be N times orig width)")
	flag.IntVar(&heightMult, "hm", 2, "height multiplier (image will be N times orig height")

	flag.Parse()

	if imagesDir == "" || sourceImage == "" || outFile == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	sourceFiles := getImageFiles(imagesDir)
	mosaic := gochallenge3.NewMosaic(tileWidth, tileHeight, sourceFiles)
	err := mosaic.Generate(sourceImage, outFile, widthMult, heightMult)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func getImageFiles(dir string) (result []string) {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Name()))
		if ext == ".jpg" || ext == ".png" || ext == ".gif" {
			fullPath := path.Join(dir, f.Name())
			result = append(result, fullPath)
		}
	}

	return
}
