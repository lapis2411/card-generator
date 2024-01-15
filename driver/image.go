package driver

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/lapis2411/card-generator/adapter"
	"github.com/lapis2411/card-generator/common"
	"github.com/lapis2411/card-generator/domain"
	"golang.org/x/image/font"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"image"
	"image/color"
	"io"
)

const (
	WidthRate = 0.05
)

type (
	imageDriver struct {
		size common.Size
		font []byte
	}
)

var border = color.RGBA{R: 0x55, G: 0x3a, B: 0xed, A: 0xff}
var white = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

func NewImageDriver(s common.Size, f []byte) adapter.ImageDriver {
	return imageDriver{
		size: s,
		font: f,
	}
}

func (i imageDriver) ImageEncode(c domain.Card) (domain.Image, error) {
	bw := int(float64(i.size.Width()) * WidthRate)
	img, err := templateCard(border, i.size, bw)
	if err != nil {
		return domain.Image{}, fmt.Errorf("failed to generate card(%s): %w", c.Name(), err)
	}
	// Cardにしたがって、画像に文字を書き込む
	for _, ft := range c.FormattedTexts() {
		points := ft.Point26_6()
		st := ft.Style()
		ff, err := fontFace(i.font, truetype.Options{Size: st.FontSize()})
		if err != nil {
			return domain.Image{}, fmt.Errorf("failed to generate fontface(%s): %w", c.Name(), err)
		}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(st.RGBA()),
			Face: ff,
			Dot:  points,
		}
		u8, err := SJIStoUTF8(ft.Text())
		if err != nil {
			return domain.Image{}, fmt.Errorf("failed to convert sjis to utf8 %v: %w", []byte(ft.Text()), err)
		}
		d.DrawString(u8)
	}

	return domain.NewImage(img, i.size, c.Name()), nil
}

// 処理時間かかるようなら毎回生成するのではなく、テンプレートを用意しておく
func templateCard(borderCol color.RGBA, size common.Size, width int) (*image.RGBA, error) {
	if size.Width() < width*2 || size.Height() < width*2 || size.Width() < 0 || size.Height() < 0 {
		return nil, fmt.Errorf("invalid size: %d, %d", size.Width(), size.Height())
	}
	img := image.NewRGBA(image.Rect(0, 0, size.Width(), size.Height()))
	for x := 0; x < size.Width(); x++ {
		for y := 0; y < size.Height(); y++ {
			if isBorder(size, x, y, width) {
				img.Set(x, y, borderCol)
			} else {
				img.Set(x, y, white)
			}
		}
	}
	return img, nil
}

func isBorder(size common.Size, x, y, width int) bool {
	return x < width || x >= size.Width()-width || y < width || y >= size.Height()-width
}

func fontFace(font []byte, opt truetype.Options) (font.Face, error) {
	// need japanese font to write japanese
	ft, err := truetype.Parse(font)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}
	face := truetype.NewFace(ft, &opt)
	return face, nil
}

func SJIStoUTF8(s string) (string, error) {
	u8, err := io.ReadAll(transform.NewReader(
		bytes.NewReader([]byte(s)),
		japanese.ShiftJIS.NewDecoder(),
	))
	if err != nil {
		return "", fmt.Errorf("failed to convert sjis to utf8: %w", err)
	}

	return string(u8), nil
}
