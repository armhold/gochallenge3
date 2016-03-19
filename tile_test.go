package gochallenge3

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestComputeAverageColor(t *testing.T) {

	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	topR := image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()/2)
	bottomR := image.Rect(0, topR.Max.Y, img.Bounds().Dx(), img.Bounds().Dy())

	// top half white
	draw.Draw(img, topR, &image.Uniform{color.White}, image.ZP, draw.Src)

	// bottom half black
	draw.Draw(img, bottomR, &image.Uniform{color.Black}, image.ZP, draw.Src)

	// expect it to average out to gray
	expected := color.RGBA{
		R: uint8(127),
		G: uint8(127),
		B: uint8(127),
		A: uint8(255),
	}

	got := ComputeAverageColor(img)

	if expected != got {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

//func TestColorRoundTrip(t *testing.T) {
//    initialR := uint32(127)
//    initialG := uint32(127)
//    initialB := uint32(127)
//    initialA := uint32(255)
//
//
//    c := color.RGBA{
//        R: uint8(initialR),
//        G: uint8(initialG),
//        B: uint8(initialB),
//        A: uint8(initialA),
//    }
//
//    r, g, b, a := c.RGBA()
//
//    if r != initialR {
//        t.Errorf("expected %d, got %d", initialR, r)
//    }
//
//    if g != initialG {
//        t.Errorf("expected %d, got %d", initialG, g)
//    }
//
//    if r != initialB {
//        t.Errorf("expected %d, got %d", initialB, b)
//    }
//
//    if r != initialA {
//        t.Errorf("expected %d, got %d", initialA, a)
//    }
//}
