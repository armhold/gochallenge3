package gochallenge3
import (
    "os"
    "image"
    "image/draw"
    "image/png"
    _ "image/jpeg"
)


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
    draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.Point{0,0}, draw.Src)

    toFile, err := os.Create(dstPath)
    if err != nil {
        return err
    }

    defer toFile.Close()

    return png.Encode(toFile, dstImg)
}
