package gosaic

import (
	"image"
	"os"
	"testing"
)

func TestScale(t *testing.T) {
	srcPath := "./sunrise.jpg"
	dstPath := "sunrise-scaled.png"

	expectW := 800
	expectH := 600

	err := ScaleToFile(srcPath, dstPath, image.Rect(0, 0, expectW, expectH))
	if err != nil {
		t.Fatal(err)
	}
	ExpectedImgBounds(expectW, expectH, dstPath, t)

	os.Remove(dstPath)
}

func ExpectedImgBounds(expectW, expectH int, imgPath string, t *testing.T) {
	dstFile, err := os.Open(imgPath)
	if err != nil {
		t.Fatal(err)
	}
	defer dstFile.Close()

	scaledImg, _, err := image.Decode(dstFile)
	if err != nil {
		t.Fatal(err)
	}

	bounds := scaledImg.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	if w != expectW {
		t.Errorf("expected width %d, got %d", expectW, w)
	}

	if h != expectH {
		t.Errorf("expected width %d, got %d", expectH, h)
	}
}
