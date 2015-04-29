package gochallenge3

import (
	"log"
	"os"
)

var (
	CommonLog = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
)
