package gochallenge3

import (
	"image"
	"image/color"
	"log"
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


func NewMosaic(tileW, tileH int, thumbs []string) Mosaic {
	return Mosaic{
		TileW: tileW,
		TileH: tileH,
		thumbs: thumbs,
		Tiles: make([]*Tile, len(thumbs)),
	}
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
	targetImg := image.NewRGBA(image.Rect(0, 0, gridImg.Bounds().Dx(), gridImg.Bounds().Dy()))

	tileRect := image.Rect(0, 0, m.TileW, m.TileH)

	for i, file := range m.thumbs {
		log.Printf("loading tile: %s", file)

		tile, err := NewTile(file, tileRect)
		m.Tiles[i] = tile
		if err != nil {
			return err
		}
	}

	cols := targetImg.Bounds().Dx() / m.TileW
	rows := targetImg.Bounds().Dy() / m.TileH

	log.Printf("targetImg.Bounds().Dx(): %d, targetImg.Bounds().Dy(): %d, rows: %d, cols: %d, W: %d, H: %d", targetImg.Bounds().Dx(), targetImg.Bounds().Dy(), rows, cols, m.OutputW, m.OutputH)


	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {

			x0 := col * m.TileW
			y0 := row * m.TileH

			x1 := x0 + m.TileW
			y1 := y0 + m.TileH

			gridRect := image.Rect(x0, y0, x1, y1)
			log.Printf("processing grid: %d, %d, bounds: %v", row, col, gridRect)

			intR := gridImg.Bounds().Intersect(gridRect)
			log.Printf("gridImg.Bounds(): %v", gridImg.Bounds())
			log.Printf("intersection: %v", intR)

			subImg := gridImg.(interface {
				SubImage(r image.Rectangle) image.Image
			}).SubImage(gridRect)

			if subImg.Bounds().Dx() == 0 {
//				panic(fmt.Errorf("rows: %d, cols: %d, row: %d, col: %d, subImg.Bounds() == %v, gridRect = %v", rows, cols, row, col, subImg.Bounds(), gridRect))
				continue
			}

			if subImg.Bounds().Dy() == 0 {
//				panic(fmt.Errorf("subImg.Bounds().Dy() == %d, gridRect = %v", subImg.Bounds().Dy(), gridRect))
				continue
			}

            tile := m.bestMatch(subImg)

			log.Printf("tile bounds: %v", tile.ScaledImage.Bounds())


            draw.Draw(targetImg, gridRect, tile.ScaledImage, tile.ScaledImage.Bounds().Min, draw.Src)
		}
	}

    return SavePng(targetImg, outfile)
}

func (m *Mosaic) bestMatch(img image.Image) *Tile {
	bestDiff := uint32(math.MaxUint32)
	log.Printf("m.Tiles len: %d", len(m.Tiles))
	bestIndex := 0
	bestTile := m.Tiles[bestIndex]

	log.Printf("bestMatch img bounds: %v", img.Bounds())

	imgAvgColor := ComputeAverageColor(img)

	for i, tile := range m.Tiles {
		diff := colorDiff(imgAvgColor, tile.AverageColor)
		if diff <= bestDiff {
			bestDiff = diff
			bestTile = tile
			bestIndex = i
		}
	}

	log.Printf("best tile index: %d, colorDiff: %f, color: %v", bestIndex, bestDiff, bestTile.AverageColor)

	return bestTile
}


// borrowed from func (p Palette) Index(c Color) int
//
func colorDiff(c1, c2 color.RGBA) uint32 {
	delta := int32(c1.R) - int32(c2.R) >> 1
	ssd := uint32(delta * delta)

	delta = int32(c1.G) - int32(c2.G) >> 1
	ssd += uint32(delta * delta)

	delta = int32(c1.B) - int32(c2.B) >> 1
	ssd += uint32(delta * delta)

	return ssd
}
