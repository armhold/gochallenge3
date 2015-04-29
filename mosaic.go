package gochallenge3

import (
	"image"
	"image/color"
	"math"
	"os"
    "image/draw"
)

type Mosaic struct {
    W      int
    H      int
    thumbs []string
    Tiles  []*Tile
}


func NewMosaic(width, height int, thumbs []string) Mosaic {
	return Mosaic{W: width, H: height, thumbs: thumbs, Tiles: make([]*Tile, len(thumbs))}
}

func (m *Mosaic) Generate(infile, outfile string) error {
	// read in source image
	srcFile, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	gridImg, _, err := image.Decode(srcFile)
	if err != nil {
		return err
	}

	// create tiles from thumbnails
	rect := image.Rect(0, 0, m.W, m.H)

	for i, file := range m.thumbs {
		CommonLog.Printf("loading tile: %s", file)

		tile, err := NewTile(file, rect)
		m.Tiles[i] = tile
		if err != nil {
			return err
		}
	}

	cols := gridImg.Bounds().Dx() / m.W
	rows := gridImg.Bounds().Dy() / m.H

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {

			x0 := col * m.W
			y0 := row * m.H

			x1 := x0 + m.W
			y1 := y0 + m.H

			CommonLog.Printf("processing grid: %d, %d", row, col)

			rect := image.Rect(x0, y0, x1, y1)
			subImg := gridImg.(*image.YCbCr).SubImage(rect)

            tile := m.bestMatch(subImg)
            r := tile.ScaledImage.Bounds()
            outImg := image.NewRGBA(gridImg.Bounds())

            draw.Draw(outImg, r, tile.ScaledImage, r.Min, draw.Src)

            CommonLog.Printf("best tile match for %v => %v", subImg, tile)
		}
	}

    return SavePng(gridImg, outfile)
}

func (m *Mosaic) bestMatch(img image.Image) *Tile {
    bestDiff := math.MaxFloat64
	CommonLog.Printf("m.Tiles len: %d", len(m.Tiles))
    bestTile := m.Tiles[0]
    imgAvgColor := ComputeAverageColor(img)

    for _, tile := range m.Tiles {
        diff := colorDiff(imgAvgColor, tile.AverageColor)
        if diff <= bestDiff {
            bestDiff = diff
            bestTile = tile
        }
    }

    return bestTile
}


func colorDiff(c1, c2 color.RGBA) float64 {
	dR := math.Pow(math.Abs(float64(c1.R-c2.R)), 2)
	dG := math.Pow(math.Abs(float64(c1.G-c2.G)), 2)
	dB := math.Pow(math.Abs(float64(c1.B-c2.B)), 2)
	dA := math.Pow(math.Abs(float64(c1.A-c2.A)), 2)

	return math.Sqrt(dR + dG + dB + dA)
}
