package gochallenge3

import (
	"io/ioutil"
	"testing"
)

func TestParseInstagramJSON(t *testing.T) {
	jsonBytes, err := ioutil.ReadFile("instagram_sample.json")
	if err != nil {
		t.Fatal(err)
	}

	got, err := ParseInstagramJSON(jsonBytes)

	if err != nil {
		t.Fatal(err)
	}

	wantPagination := InstagramPagination{MinTagID: "969806771123671476", NextURL: "https://api.instagram.com/v1/tags/dogs/media/recent?callback=YOUR-CALLBACK\u0026client_id=client123\u0026max_tag_id=969806141804481738"}

	if got.Pagination != wantPagination {
		t.Fatalf("Pagination => %q, want %q", got.Pagination, wantPagination)
	}

	wantFirstImageSet := InstagramImageSet{
		LowRes:      InstagramImage{Url: "https://scontent.cdninstagram.com/hphotos-xfa1/t51.2885-15/s306x306/e15/11191310_409598512535011_1831056149_n.jpg", Width: 306, Height: 306},
		Thumb:       InstagramImage{Url: "https://scontent.cdninstagram.com/hphotos-xfa1/t51.2885-15/s150x150/e15/11191310_409598512535011_1831056149_n.jpg", Width: 150, Height: 150},
		StandardRes: InstagramImage{Url: "https://scontent.cdninstagram.com/hphotos-xfa1/t51.2885-15/e15/11191310_409598512535011_1831056149_n.jpg", Width: 640, Height: 640},
	}

	if got.Data[0].ImageSet != wantFirstImageSet {
		t.Fatalf("ImageSet => %q, want %q", got.Data[0].ImageSet, wantFirstImageSet)
	}
}

func TestByRows(t *testing.T) {
	imageSets := make([]InstagramImageSet, 19)

	for i, _ :=  range imageSets {
		imageSets[i] = InstagramImageSet{}
	}

	rows := ToRows(5, imageSets)

	if len(rows) != 4 {
		t.Errorf("expected 4 rows, got %d", len(rows))
	}


	expect := func (s []InstagramImageSet, expectedLen int) {
		if (len(s) != expectedLen) {
			t.Errorf("expected %d cols, got %d", expectedLen, len(rows))
		}
	}

	expect(rows[0], 5)
	expect(rows[1], 5)
	expect(rows[2], 5)
	expect(rows[3], 4)
}
