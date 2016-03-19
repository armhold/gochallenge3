package main

import (
	"flag"
	"github.com/armhold/gochallenge3"
	"log"
	"os"
)

var (
	devMode bool
)

func init() {
	flag.BoolVar(&devMode, "dev", false, "start the server in devmode")
	flag.Parse()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	uploadRootDir := os.Getenv("UPLOAD_DIR")
	if uploadRootDir == "" {
		uploadRootDir = "/tmp/upload_dir"
	}

	if err := os.MkdirAll(uploadRootDir, 0700); err != nil {
		log.Fatalf("unable to create temp dir %s: %s", uploadRootDir, err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port

	if devMode {
		// prevent OSX firewall popups by binding to localhost
		//
		// see: http://grokbase.com/t/gg/golang-nuts/15322dedhn/go-nuts-mac-firewall
		addr = "localhost" + addr
	}

	instagramClientID := os.Getenv("INSTAGRAM_CLIENT_ID")
	if instagramClientID == "" {
		log.Fatal("environment variable INSTAGRAM_CLIENT_ID not set")
	}

	imageSource := &gochallenge3.InstagramClient{instagramClientID}

	gochallenge3.Serve(addr, uploadRootDir, imageSource)
}
