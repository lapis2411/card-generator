package driver

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
	Template struct {
		size common.Size
	}
)

func NewA4Template() adapter.Template {
	return Template{
		size: common.NewSize(2480, 3508),
	}
}

func NewJISB5Template() adapter.Template {
	return Template{
		size: common.NewSize(1900, 2500),
	}
}

// Cardサイズが同じであること前提
func (t Template) Arrange(imgs []domain.Image) ([]domain.Image, error) {
	// サイズに合わせて印刷用のレイアウトを作成
	sizeCard := imgs[0].Size()
	cs, err := merge(imgs, t.size, sizeCard)
	if err != nil {
		return nil, fmt.Errorf("failed to merge images: %w", err)
	}
	return cs, nil
}

func merge(imgs []domain.Image, pageSize, cardSize common.Size) ([]domain.Image, error) {
	row := pageSize.Width() / (cardSize.Width() + Buffer)
	column := pageSize.Height() / (cardSize.Height() + Buffer)
	cnt := 0
	bi := newPage(pageSize)
	pages := []*image.RGBA{bi}
	for _, img := range imgs {
		// out of 1 page max size
		if cnt >= row*column {
			cnt = 0
			bi = newPage(pageSize)
			pages = append(pages, bi)
		}
		cr := (cnt % row)
		cc := (cnt / row)
		x := cr*(cardSize.Width()+Buffer) + Buffer/2
		y := cc*(cardSize.Height()+Buffer) + Buffer/2
		if err := overwriteImage(bi, *img.Image(), image.Point{X: x, Y: y}); err != nil {
			page := len(pages) + 1
			return nil, fmt.Errorf("failed to overwrite image page %v, row %v, column %v: %w", page, row, column, err)
		}
		cnt++
	}
	pimgs := make([]domain.Image, 0, len(pages))
	for i, p := range pages {
		pimgs = append(pimgs, domain.NewImage(p, pageSize, fmt.Sprintf("%04d", i+1)))
	}
	return pimgs, nil
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
