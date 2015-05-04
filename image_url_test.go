package gochallenge3

import (
	"fmt"
	"testing"
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
