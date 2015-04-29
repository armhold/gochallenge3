package gochallenge3

import "testing"

func testReadNonExistant(t *testing.T) {
	_, err := ReadProject("/tmp/foo", "nonexistant")

	if err != nil {
		t.Errorf("expected error for non-existant project ID")
	}
}
