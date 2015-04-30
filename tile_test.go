package gochallenge3

import (
    "testing"
    "image"
    "image/color"
    "image/draw"
)

func TestComputeAverageColor(t *testing.T) {

    img := image.NewRGBA(image.Rect(0, 0, 100, 100))

    expected := color.RGBA{
        R: uint8(127),
        G: uint8(127),
        B: uint8(127),
        A: uint8(255),
    }

    topR    := image.Rect(0, 0,          img.Bounds().Dx(), img.Bounds().Dy() / 2)
    bottomR := image.Rect(0, topR.Max.Y, img.Bounds().Dx(), img.Bounds().Dy())


    draw.Draw(img, topR,    &image.Uniform{color.White}, image.ZP, draw.Src)
    draw.Draw(img, bottomR, &image.Uniform{color.Black}, image.ZP, draw.Src)

    got := ComputeAverageColor(img)

    if expected != got {
        t.Errorf("expected %v, got %v", expected, got)
    }

}
