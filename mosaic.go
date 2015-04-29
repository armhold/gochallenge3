package gochallenge3

import (
	"image"
	"image/color"
	"math"
	"os"
)

type Mosaic struct {
    W      int
    H      int
    thumbs []string
    Tiles  []Tile
}


func NewMosaic(width, height int, thumbs []string) Mosaic {
    return Mosaic{W: width, H: height, thumbs: thumbs}
}

func (m *Mosaic) Generate(infile, outfile string) error {
	// read in source image
	srcFile, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return err
	}

	// create tiles from thumbnails
	tiles := make([]*Tile, len(m.thumbs))
	rect := image.Rect(0, 0, m.W, m.H)

	for i, file := range m.thumbs {
		tile, err := NewTile(file, rect)
		tiles[i] = tile
		if err != nil {
			return err
		}
	}

	cols := srcImg.Bounds().Dx() / m.W
	rows := srcImg.Bounds().Dy() / m.H

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {

			x0 := col * m.W
			y0 := row * m.H

			x1 := x0 + m.W
			y1 := y0 + m.H

			rect := image.Rect(x0, y0, x1, y1)
			subImg := srcImg.(*image.RGBA).SubImage(rect)
			avgColor := ComputeAverageColor(subImg)
			CommonLog.Printf("average color for %d, %d => %v", row, col, avgColor)
		}
	}

	return nil
}

func colorDiff(c1, c2 color.RGBA) float64 {
	dR := math.Pow(float64(c1.R-c2.R), 2)
	dG := math.Pow(float64(c1.G-c2.G), 2)
	dB := math.Pow(float64(c1.B-c2.B), 2)
	dA := math.Pow(float64(c1.A-c2.A), 2)

	return math.Sqrt(dR + dG + dB + dA)
}
