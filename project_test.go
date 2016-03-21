package gochallenge3

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
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

func TestJSON(t *testing.T) {
	p := Project{ID: "abc123", Status: StatusCompleted}
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"ID":"abc123","Status":"completed"}` // NB: no "uploadRootDir" (intentionally unexported field)
	actual := string(b)
	if expected != actual {
		t.Fatalf("expected: %s, got: %s", expected, actual)
	}
}
