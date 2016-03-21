package gochallenge3

import (
	"testing"
	"os"
	"io/ioutil"
)

func TestReadNonExistant(t *testing.T) {
	_, err := ReadProject("/tmp/foo", "nonexistant")

	if err == nil {
		t.Errorf("expected error for non-existant project ID")
	}
}

func TestSetLoadStatus(t *testing.T) {
	uploadRootDir, err := ioutil.TempDir("", "project_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(uploadRootDir)

	p, err := NewProject(uploadRootDir)
	if err != nil {
		t.Fatal(err)
	}

	// re-read it- status should be "new"
	p, err = ReadProject(uploadRootDir, p.ID)
	if err != nil {
		t.Fatal(err)
	}

	if p.Status != StatusNew {
		t.Fatalf("expected %s, got %s", StatusNew, p.Status)
	}

	p.SetAndSaveStatus(StatusCompleted)
	p, err = ReadProject(uploadRootDir, p.ID)
	if err != nil {
		t.Fatal(err)
	}

	if p.Status != StatusCompleted {
		t.Fatalf("expected %s, got %s", StatusCompleted, p.Status)
	}
}
