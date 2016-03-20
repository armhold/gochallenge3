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

func TestGuessImageExtension(t *testing.T) {
	cases := []struct {
		url, expected string
	} {
		{ "https://scontent.cdninstagram.com/t51.2885-15/s150x150/e35/1169705_1005749189478164_861869881_n.jpg?ig_cache_key=BOGUS123%3D%3D.2", ".jpg"},
		{ "https://scontent.cdninstagram.com/t51.2885-15/s150x150/e35/1169705_1005749189478164_861869881_n.gif?ig_cache_key=BOGUS123%3D%3D.2", ".gif"},
		{ "http://SCONTENT.CDNINSTAGRAM.COM/T51.2885-15/S150X150/E35/1169705_1005749189478164_861869881_N.PNG?IG_CACHE_KEY=BOGUS123%3D%3D.2", ".png"},
	}

	for _, c := range cases {
		url := ImageURL(c.url)
		actual, err := url.guessImageExtension()

		if err != nil {
			t.Fatal(err)
		}

		if actual != c.expected {
			t.Fatalf("expected %s, got %s", c.expected, actual)
		}
	}
}
