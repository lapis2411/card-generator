package generator

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
)

const (
	buffer = 15 // buffer for cut cards
)

type (
	A4Arrange struct{}
	Canvas    struct {
		Image *image.RGBA
	}
	Canvases []Canvas
)

// 350dpi
var sizeA4 = Size{
	Width:  2480,
	Height: 3508,
}
var sizeCard = Size{
	Width:  widthCard,
	Height: heightCard,
}

func NewA4Layout() Arrange {
	return A4Arrange{}
}

func (A4Arrange) Arrange(images []*image.RGBA) ([]*image.RGBA, error) {
	return images, nil
}

// MergeCards merges cards into A4 size(350dpi) images of under the dir.
func MergeCards(dir string, outDir string) error {
	if outDir == "" {
		outDir = "./out"
	}
	return mergeCards(dir, outDir, sizeA4, sizeCard)
}

func mergeCards(dir, outDir string, paperSize, cardSize Size) error {
	row := paperSize.Width / (cardSize.Width + buffer)
	column := paperSize.Height / (cardSize.Height + buffer)
	canvases := make(Canvases, 1)
	cnt := 0
	page := 0
	canvases[page] = NewLayer(paperSize)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the path %v: %w\n", dir, err)
		}
		if !info.IsDir() {
			// どこから書き出すか？
			if cnt >= row*column {
				cnt = 0
				page++
				canvases = append(canvases, NewLayer(paperSize))
			}
			img, err := DecodeImage(path)
			if err != nil {
				return fmt.Errorf("failed to decode image(%s): %w", path, err)
			}
			cr := (cnt % row)
			cc := (cnt / row)
			x := cr*(cardSize.Width+buffer) + buffer/2
			y := cc*(cardSize.Height+buffer) + buffer/2
			if err = canvases[page].OverwriteImage(img, image.Point{X: x, Y: y}); err != nil {
				return fmt.Errorf("failed to overwrite image(%s): %w", path, err)
			}

			cnt++
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", dir, err)
	}

	if err := canvases.ExportImages(outDir); err != nil {
		return fmt.Errorf("failed to export images: %w", err)
	}

	return nil
}

// NewLayer creates a new layer of the specified size.
func NewLayer(size Size) Canvas {

	img := image.NewRGBA(image.Rect(0, 0, size.Width, size.Height))

	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < size.Height; y++ {
		for x := 0; x < size.Width; x++ {
			img.Set(x, y, white)
		}
	}
	return Canvas{
		Image: img,
	}
}

// OverwriteImage overwrites the image at the specified position by the writeImg.
func (c *Canvas) OverwriteImage(writeImg image.Image, start image.Point) error {
	bnd := writeImg.Bounds()
	w := bnd.Max.X - bnd.Min.X
	h := bnd.Max.Y - bnd.Min.Y
	if c.Image.Bounds().Max.X < start.X+w || c.Image.Bounds().Max.Y < start.Y+h {
		return fmt.Errorf("out of bounds: %d, %d", start.X+w, start.Y+h)
	}

	for x := start.X; x < start.X+w; x++ {
		for y := start.Y; y < start.Y+h; y++ {
			px := x - start.X
			py := y - start.Y
			c.Image.Set(x, y, writeImg.At(px, py))
		}
	}
	return nil
}
