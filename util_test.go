package gochallenge3

import (
	"reflect"
	"testing"
)

func TestUrlEncodeSpaces(t *testing.T) {
	input := "foo bar"
	got, err := UrlEncode(input)
	want := "foo%20bar"

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("UrlEncoded(%q) => %q, want %q", input, got, want)
	}
}

func TestUrlEncodeDotsAndSlashes(t *testing.T) {
	input := "foo/../etc/password"
	got, err := UrlEncode(input)
	want := "fooetcpassword"

	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Errorf("UrlEncoded(%q) => %q, want %q", input, got, want)
	}
}

func TestSplitPath(t *testing.T) {
	expected := []string{"search", "foo"}

	parts := SplitPath("/search/foo/")
	if !reflect.DeepEqual(expected, parts) {
		t.Errorf("expected: %v, got: %v", expected, parts)
	}

	parts = SplitPath("/search/foo")
	if !reflect.DeepEqual(expected, parts) {
		t.Errorf("expected: %v, got: %v", expected, parts)
	}

	parts = SplitPath("search/foo")
	if !reflect.DeepEqual(expected, parts) {
		t.Errorf("expected: %v, got: %v", expected, parts)
	}

	parts = SplitPath("search/foo/")
	if !reflect.DeepEqual(expected, parts) {
		t.Errorf("expected: %v, got: %v", expected, parts)
	}
}
