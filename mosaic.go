package gochallenge3

import (
	"image"
	"image/color"
	"math"
	"os"
    "image/draw"
)

type Mosaic struct {
    OutputW      int
    OutputH      int
	TileW        int
	TileH        int
    thumbs       []string
    Tiles        []*Tile
}


func NewMosaic(outputW, outputH, tileW, tileH int, thumbs []string) Mosaic {
	return Mosaic{OutputW: outputW, OutputH: outputH, TileW: tileW, TileH: tileH, thumbs: thumbs, Tiles: make([]*Tile, len(thumbs))}
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
	rect := image.Rect(0, 0, m.OutputW, m.OutputH)

	for i, file := range m.thumbs {
		CommonLog.Printf("loading tile: %s", file)

		tile, err := NewTile(file, rect)
		m.Tiles[i] = tile
		if err != nil {
			return err
		}
	}

	cols := gridImg.Bounds().Dx() / m.TileW
	rows := gridImg.Bounds().Dy() / m.TileH

	CommonLog.Printf("gridImg.Bounds().Dx(): %d, gridImg.Bounds().Dy(): %d, rows: %d, cols: %d, W: %d, H: %d", gridImg.Bounds().Dx(), gridImg.Bounds().Dy(), rows, cols, m.OutputW, m.OutputH)

	outImg := image.NewRGBA(rect)

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {

			x0 := col * m.TileW
			y0 := row * m.TileH

			x1 := x0 + m.TileW
			y1 := y0 + m.TileH

			rect := image.Rect(x0, y0, x1, y1)
			CommonLog.Printf("processing grid: %d, %d, bounds: %v", row, col, rect)

			// TODO: detect proper image type for this type assertion?
			subImg := gridImg.(*image.YCbCr).SubImage(rect)

            tile := m.bestMatch(subImg)
            r := tile.ScaledImage.Bounds()
            draw.Draw(outImg, r, tile.ScaledImage, image.Point{X: x0, Y: y0}, draw.Src)
		}
	}

    return SavePng(outImg, outfile)
}

func (m *Mosaic) bestMatch(img image.Image) *Tile {
    bestDiff := math.MaxFloat64
	CommonLog.Printf("m.Tiles len: %d", len(m.Tiles))
	bestIndex := 0
    bestTile := m.Tiles[bestIndex]
    imgAvgColor := ComputeAverageColor(img)


    for i, tile := range m.Tiles {
        diff := colorDiff(imgAvgColor, tile.AverageColor)
        if diff <= bestDiff {
            bestDiff = diff
            bestTile = tile
			bestIndex = i
        }
    }

	CommonLog.Printf("best tile index: %d", bestIndex)

    return bestTile
}


func colorDiff(c1, c2 color.RGBA) float64 {
	dR := math.Pow(math.Abs(float64(c1.R-c2.R)), 2)
	dG := math.Pow(math.Abs(float64(c1.G-c2.G)), 2)
	dB := math.Pow(math.Abs(float64(c1.B-c2.B)), 2)
	dA := math.Pow(math.Abs(float64(c1.A-c2.A)), 2)

	return math.Sqrt(dR + dG + dB + dA)
}
