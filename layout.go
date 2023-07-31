package generator

import (
	"fmt"
	"image"
	"image/color"
)

const (
	buffer = 15 // buffer for cut cards
)

type (
	A4Arrange struct{}
)

// 300dpi
var sizeA4 = Size{
	Width:  2480,
	Height: 3508,
}

func NewA4Layout() Arrange {
	return A4Arrange{}
}

// Cardサイズが同じであること前提
func (A4Arrange) Arrange(cards Cards) (Canvases, error) {
	// サイズに合わせてA4印刷用のレイアウトを作成
	sizeCard := cards[0].Size
	cs, err := merge(cards, sizeA4, sizeCard)
	if err != nil {
		return nil, fmt.Errorf("failed to merge cards: %w", err)
	}
	return cs, nil
}

func merge(cards Cards, paperSize, cardSize Size) (Canvases, error) {
	row := paperSize.Width / (cardSize.Width + buffer)
	column := paperSize.Height / (cardSize.Height + buffer)
	canvases := make(Canvases, 1)
	cnt := 0
	page := 0
	canvases[page] = NewLayer(paperSize)
	for _, card := range cards {
		// out of 1 page max
		if cnt >= row*column {
			cnt = 0
			page++
			canvases = append(canvases, NewLayer(paperSize))
		}
		cr := (cnt % row)
		cc := (cnt / row)
		x := cr*(cardSize.Width+buffer) + buffer/2
		y := cc*(cardSize.Height+buffer) + buffer/2
		if err := canvases[page].OverwriteImage(card.Image, image.Point{X: x, Y: y}); err != nil {
			return nil, fmt.Errorf("failed to overwrite image page %v, row %v, column %v: %w", page, row, column, err)
		}
		cnt++
	}

	return canvases, nil
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
	return Canvas{Image: img}
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
