package gochallenge3

import (
    "testing"
    "image"
    "os"
)


func TestScale(t *testing.T) {
    srcPath := "./sunrise.jpg"
    dstPath := "sunrise-scaled.png"

    err := Scale(srcPath, dstPath, image.Rect(0, 0, 800, 600))
    if err != nil {
        t.Fatal(err)
    }

    dstFile, err := os.Open(dstPath)
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

    if w != 800 {
        t.Errorf("expected width %d, got %d", 800, w)
    }

    if h != 600 {
        t.Errorf("expected width %d, got %d", 600, h)
    }

    os.Remove(dstPath)
}
