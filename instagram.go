package gochallenge3

// simple wrapper for Instagram REST API
type Instagram interface {

	Search(s string) []string
}

// curl 'https://api.instagram.com/v1/tags/SEARCH-TAG/media/recent?client_id=CLIENT-ID&callback=YOUR-CALLBACK'
