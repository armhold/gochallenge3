package gochallenge3
import (
    "testing"
    "fmt"
)

func TestByRows(t *testing.T) {
    imageURLs := make([]ImageURL, 19)

    for i, _ := range imageURLs {
        imageURLs[i] = ImageURL(fmt.Sprintf("http://example.com/foo%d", i))
    }

    rows := ToRows(5, imageURLs)

    if len(rows) != 4 {
        t.Errorf("expected 4 rows, got %d", len(rows))
    }

    expect := func(s []ImageURL, expectedLen int) {
        if len(s) != expectedLen {
            t.Errorf("expected %d cols, got %d", expectedLen, len(rows))
        }
    }

    expect(rows[0], 5)
    expect(rows[1], 5)
    expect(rows[2], 5)
    expect(rows[3], 4)
}

func TestSerialize(t *testing.T) {
    imageURLs := make([]ImageURL, 5)

    for i, _ := range imageURLs {
        imageURLs[i] = ImageURL(fmt.Sprintf("http://example.com/foo%d", i))
    }

    // TODO: use tmp file
    file := "image_url_output.txt"
    ToFile(imageURLs, file)

    urlsFromFile, err := FromFile(file)
    if err != nil {
        t.Fatal(err)
    }

    if len(urlsFromFile) != 5 {
        t.Error("expected 5 urls, got %d", len(urlsFromFile))
    }

    for i, url := range urlsFromFile {
        expected := fmt.Sprintf("http://example.com/foo%d", i)
        if string(url) != expected {
            t.Errorf("expected %s, got %s", expected, url)
        }
    }
}
