package gochallenge3

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
)

type Tile struct {
	AverageColor color.RGBA
	ScaledImage  image.Image
}

func (t *Tile) GoString() string {
	r, g, b, a := t.AverageColor.RGBA()
	return fmt.Sprintf("Tile: AverageColor: r: %d, g: %d, b: %d, a: %d", r, g, b, a)
}

func NewTile(filePath string, r image.Rectangle) (*Tile, error) {
	srcFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return nil, err
	}

	scaledImg := Scale(srcImg, r)
	avgColor := ComputeAverageColor(scaledImg)

	return &Tile{ScaledImage: scaledImg, AverageColor: avgColor}, nil
}

func ComputeAverageColor(img image.Image) color.RGBA {
	sumR := uint32(0)
	sumG := uint32(0)
	sumB := uint32(0)
	sumA := uint32(0)

	b := img.Bounds()
	log.Printf("ComputeAverageColor bounds: %v", b)
	count := uint32(0)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			sumR += uint32(r) >> 8
			sumG += uint32(g) >> 8
			sumB += uint32(b) >> 8
			sumA += uint32(a) >> 8
			count++
		}
	}

	return color.RGBA{
		R: uint8(sumR / count),
		G: uint8(sumG / count),
		B: uint8(sumB / count),
		A: uint8(sumA / count),
	}
}
