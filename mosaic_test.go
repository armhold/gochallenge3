package gochallenge3
import (
    "testing"
    "fmt"
)


func TestGenerate(t *testing.T) {

    thumbs := make([]string, 23)
    for i, _ := range thumbs {
        thumbs[i] = fmt.Sprintf("test_images/test_image%d.jpg", i)
    }


    mosaic := NewMosaic(3264, 2448, thumbs)
    outPath := "mosaic_out.jpg"

    err := mosaic.Generate("test_images/source_image.jpg", outPath)
    if err != nil {
        t.Error(err)
    }

    ExpectedImgBounds(3264, 2448, outPath, t)
}