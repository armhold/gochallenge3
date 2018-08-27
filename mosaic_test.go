package gosaic

import (
	"fmt"
	"testing"
)

func TestGenerate(t *testing.T) {
	thumbs := make([]string, 10)
	for i, _ := range thumbs {
		thumbs[i] = fmt.Sprintf("test_images/test_image%d.jpg", i)
	}

	mosaic := NewMosaic(50, 50, thumbs)
	outPath := "mosaic_out.jpg"

	err := mosaic.Generate("test_images/source_image.jpg", outPath, 1, 1)
	if err != nil {
		t.Error(err)
	}

	ExpectedImgBounds(3264, 2448, outPath, t)
}
