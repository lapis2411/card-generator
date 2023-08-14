package template

import (
	"fmt"
	"image"
	"image/color"

	"github.com/lapis2411/card-generator/adapter"
	"github.com/lapis2411/card-generator/common"
	"github.com/lapis2411/card-generator/domain"
)

const (
	Buffer = 15 // buffer for cut cards
)

type (
	a4Template struct{}
)

// 300dpi
var sizeA4 = common.NewSize(2480, 3508)

func NewA4Template() adapter.Template {
	return a4Template{}
}

// Cardサイズが同じであること前提
func (a4Template) Arrange(imgs []domain.Image) ([]domain.Image, error) {
	// サイズに合わせてA4印刷用のレイアウトを作成
	sizeCard := imgs[0].Size()
	cs, err := merge(imgs, sizeA4, sizeCard)
	if err != nil {
		return nil, fmt.Errorf("failed to merge images: %w", err)
	}
	return cs, nil
}

func merge(imgs []domain.Image, pageSize, cardSize common.Size) ([]domain.Image, error) {
	row := pageSize.Width() / (cardSize.Width() + Buffer)
	column := pageSize.Height() / (cardSize.Height() + Buffer)
	canvases := make([]domain.Image, 1)
	cnt := 0
	bi := newPage(pageSize)
	page := 1
	for _, img := range imgs {
		// out of 1 page max size
		if cnt >= row*column {
			cnt = 0
			canvases = append(canvases, domain.NewImage(bi, pageSize, fmt.Sprint("%4d", page)))
			page++
			bi = newPage(pageSize)
		}
		cr := (cnt % row)
		cc := (cnt / row)
		x := cr*(cardSize.Width()+Buffer) + Buffer/2
		y := cc*(cardSize.Height()+Buffer) + Buffer/2
		if err := overwriteImage(bi, *img.Image(), image.Point{X: x, Y: y}); err != nil {
			return nil, fmt.Errorf("failed to overwrite image page %v, row %v, column %v: %w", page, row, column, err)
		}
		cnt++
	}

	return canvases, nil
}

// newPage creates a new layer of the specified size.
func newPage(s common.Size) *image.RGBA {

	img := image.NewRGBA(image.Rect(0, 0, s.Width(), s.Height()))

	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < s.Height(); y++ {
		for x := 0; x < s.Width(); x++ {
			img.Set(x, y, white)
		}
	}
	return img
}

// overwriteImage overwrites the image at the specified position by the writeImg.
func overwriteImage(baseImg *image.RGBA, writeImg image.RGBA, start image.Point) error {
	bnd := writeImg.Bounds()
	w := bnd.Max.X - bnd.Min.X
	h := bnd.Max.Y - bnd.Min.Y
	if baseImg.Bounds().Max.X < start.X+w || baseImg.Bounds().Max.Y < start.Y+h {
		return fmt.Errorf("out of bounds: %d, %d", start.X+w, start.Y+h)
	}

	for x := start.X; x < start.X+w; x++ {
		for y := start.Y; y < start.Y+h; y++ {
			px := x - start.X
			py := y - start.Y
			baseImg.Set(x, y, writeImg.At(px, py))
		}
	}
	return nil
}
