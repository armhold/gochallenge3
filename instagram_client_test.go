package gochallenge3

import (
	"testing"
)

func TestUrl(t *testing.T) {
	i := NewInstagramImageSource("client123")

	searchTerm := "dogs"
	got, err := i.instagramAPIUrl(searchTerm)
	want := "https://api.instagram.com/v1/tags/dogs?client_id=client123"

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("instagramAPIUrl(%q) => %q, want %q", searchTerm, got, want)
	}
}
