package generator

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	defaultWidth = 30
	// output       = "../temporary/"
	// fontPath     = "../static/fonts/AP.ttf"
)

var border = color.RGBA{R: 0x55, G: 0x3a, B: 0xed, A: 0xff}
var white = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

type (
	Size struct {
		Width  int
		Height int
	}
	Generator struct {
		fontPath string
		cardSize Size
	}
	Canvas struct {
		Image *image.RGBA
		Size  Size
	}
	Canvases []Canvas
	Card     struct {
		Image *image.RGBA
		Size  Size
	}
	Cards   []Card
	Arrange interface {
		// Arrange arranges some images into specified size pages
		Arrange(Cards) (Canvases, error)
	}
)

var normalCardSize = Size{Width: 600, Height: 800}

func NormalSizeCardGenerator(fontPath string) Generator {
	return buildGenerator(fontPath, normalCardSize)
}

func buildGenerator(fontPath string, cardSize Size) Generator {
	if fontPath == "" {
		fontPath = "./static/fonts/AP.ttf"
	}
	return Generator{
		fontPath: fontPath,
		cardSize: cardSize,
	}
}

func (g Generator) Generate(contents CardContents) (Cards, error) {
	cs := make(Cards, len(contents))
	i := 0
	for name, c := range contents {
		card, err := g.generateCard(*c, name)
		if err != nil {
			return nil, fmt.Errorf("failed to generate card(%s): %w", name, err)
		}
		cs[i] = card
		i++
	}
	return cs, nil
}

// generateCard is a function for exporting a card image
// c is a Card struct, name is a file name, rgba is a border color
// TODO: FontSizeと文字の位置に応じてフォントサイズや改行を調整する
// 例えばタイトルが長くなりそうな時は2行にしてフォントを小さくする、など
// 描画領域を指定する形のほうがよさげか？いずれにしてもサイトがざっとできてから着手
func (g Generator) generateCard(c CardContent, name string) (Card, error) {
	img, err := templateCard(border, g.cardSize, defaultWidth)
	if err != nil {
		return Card{}, fmt.Errorf("failed to generate card(%s): %w", name, err)
	}
	// Cardにしたがって、画像に文字を書き込む
	for _, s := range c.styledTexts {
		points := s.Point26_6()
		ff, err := fontFace(g.fontPath, truetype.Options{Size: s.style.fontsize})
		if err != nil {
			return Card{}, fmt.Errorf("failed to generate fontface(%s): %w", name, err)
		}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(s.style.rgba),
			Face: ff,
			Dot:  points,
		}
		u8, err := SJIStoUTF8(s.text)
		if err != nil {
			return Card{}, fmt.Errorf("failed to convert sjis to utf8 %v: %w", []byte(s.text), err)
		}
		d.DrawString(u8)
	}

	return Card{Image: img, Size: g.cardSize}, nil
}

// 処理時間かかるようなら毎回生成するのではなく、テンプレートを用意しておく
func templateCard(borderCol color.RGBA, size Size, width int) (*image.RGBA, error) {
	if size.Width < width*2 || size.Height < width*2 || size.Width < 0 || size.Height < 0 {
		return nil, fmt.Errorf("invalid size: %d, %d", size.Width, size.Height)
	}
	img := image.NewRGBA(image.Rect(0, 0, size.Width, size.Height))
	for x := 0; x < size.Width; x++ {
		for y := 0; y < size.Height; y++ {
			if isBorder(size, x, y, width) {
				img.Set(x, y, borderCol)
			} else {
				img.Set(x, y, white)
			}
		}
	}
	return img, nil
}

func isBorder(size Size, x, y, width int) bool {
	return x < width || x >= size.Width-width || y < width || y >= size.Height-width
}

func fontFace(fontPath string, opt truetype.Options) (font.Face, error) {
	// need japanese font to write japanese
	ftBin, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load font(%s): %w", fontPath, err)
	}
	ft, err := truetype.Parse(ftBin)
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
