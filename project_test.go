package gochallenge3

import (
	"testing"
	"fmt"
)

func TestReadNonExistant(t *testing.T) {
	_, err := ReadProject("/tmp/foo", "nonexistant")

	if err == nil {
		t.Errorf("expected error for non-existant project ID")
	}
}


func TestSerialize(t *testing.T) {
	imageURLs := make([]ImageURL, 5)

	for i, _ := range imageURLs {
		imageURLs[i] = ImageURL(fmt.Sprintf("http://example.com/foo%d", i))
	}

	project, _ := NewProject("/tmp")

	project.ToFile(imageURLs)

	urlsFromFile, err := project.FromFile()
	if err != nil {
		t.Fatal(err)
	}

	if len(urlsFromFile) != 5 {
		t.Errorf("expected 5 urls, got %d", len(urlsFromFile))
	}

	for i, url := range urlsFromFile {
		expected := fmt.Sprintf("http://example.com/foo%d", i)
		if string(url) != expected {
			t.Errorf("expected %s, got %s", expected, url)
		}
	}
}
