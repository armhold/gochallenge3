package gochallenge3

import (
	"reflect"
	"testing"
)

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
