package gochallenge3

import (
    "testing"
    "image"
    "os"
)


func TestScale(t *testing.T) {
    srcPath := "./sunrise.jpg"
    dstPath := "sunrise-scaled.png"

    expectW := 800
    expectH := 600

    err := Scale(srcPath, dstPath, image.Rect(0, 0, expectW, expectH))
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

    if w != expectW {
        t.Errorf("expected width %d, got %d", expectW, w)
    }

    if h != expectH {
        t.Errorf("expected width %d, got %d", expectH, h)
    }

    os.Remove(dstPath)
}
