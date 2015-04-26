package gochallenge3
import (
    "os"
    "image"
    "image/png"
    _ "image/jpeg"
    "math"
)


// Scale scales the src image to the given rectangle using Nearest Neighbor
func Scale(srcPath, dstPath string, r image.Rectangle) error {
    srcFile, err := os.Open(srcPath)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    srcImg, _, err := image.Decode(srcFile)
    if err != nil {
        return err
    }

    dstImg := image.NewRGBA(r)

    srcBounds := srcImg.Bounds()
    sw := srcBounds.Max.X - srcBounds.Min.X
    sh := srcBounds.Max.Y - srcBounds.Min.Y

    dstBounds := dstImg.Bounds()
    dw := dstBounds.Max.X - dstBounds.Min.X
    dh := dstBounds.Max.Y - dstBounds.Min.Y

    xAspect := float64(sw)/float64(dw)
    yAspect := float64(sh)/float64(dh)

    for y := 0; y < sh; y++ {
        for x := 0; x < sw; x++ {
            srcX := int(math.Floor(float64(x) * xAspect))
            srcY := int(math.Floor(float64(y) * yAspect))
            pix := srcImg.At(srcX, srcY)
            dstImg.Set(x, y, pix)
        }
    }

    toFile, err := os.Create(dstPath)
    if err != nil {
        return err
    }

    defer toFile.Close()

    return png.Encode(toFile, dstImg)
}
