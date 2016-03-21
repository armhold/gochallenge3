package gochallenge3

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	homeHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}

func TestInvalidUploadId(t *testing.T) {
	req, _ := http.NewRequest("GET", "/results/invalid_upload_id", nil)
	w := httptest.NewRecorder()

	expected := http.StatusBadRequest

	uploadRootDir := os.TempDir()

	resultsHandler(uploadRootDir).ServeHTTP(w, req)

	if w.Code != expected {
		t.Errorf("results page expected %v, got %v", expected, w.Code)
	}
}
