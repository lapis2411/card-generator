package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	widthCard    = 600
	heightCard   = 800
	defaultWidth = 30
	output       = "../temporary/"
	fontPath     = "../static/fonts/AP.ttf"
)

var border = color.RGBA{R: 0x55, G: 0x3a, B: 0xed, A: 0xff}
var white = color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

type (
	Size struct {
		Width  int
		Height int
	}
)

func Generate(cards Cards) error {
	for name, card := range cards {
		if err := generateCard(*card, name, widthCard, heightCard); err != nil {
			return err
		}
	}
	return nil
}

// generateCard is a function for exporting a card image
// c is a Card struct, name is a file name, rgba is a border color
// TODO: FontSizeと文字の位置に応じてフォントサイズや改行を調整する
// 例えばタイトルが長くなりそうな時は2行にしてフォントを小さくする、など
// 描画領域を指定する形のほうがよさげか？いずれにしてもサイトがざっとできてから着手
func generateCard(c Card, name string, width, height int) error {
	img, err := templateCard(border, Size{width, height}, defaultWidth)
	if err != nil {
		return fmt.Errorf("failed to generate card(%s): %w", name, err)
	}
	// Cardにしたがって、画像に文字を書き込む
	for _, s := range c.styledTexts {
		points := s.Point26_6()
		ff, err := fontFace(truetype.Options{Size: s.style.fontsize})
		if err != nil {
			return fmt.Errorf("failed to generate fontface(%s): %w", name, err)
		}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(s.style.rgba),
			Face: ff,
			Dot:  points,
		}
		d.DrawString(s.text)
	}

	fn := "../temporary/" + name + ".png"
	if err := exportImage(fn, img); err != nil {
		return fmt.Errorf("failed to export card(%s): %w", name, err)
	}
	return nil
}

// 処理時間かかるようなら毎回生成するのではなく、テンプレートを用意しておく
func templateCard(borderCol color.RGBA, size Size, width int) (*image.RGBA, error) {
	if size.Width < width*2 || size.Height < width*2 || size.Width < 0 || size.Height < 0 {
		return nil, fmt.Errorf("invalid size: %d, %d", size.Width, size.Height)
	}
	img := image.NewRGBA(image.Rect(0, 0, size.Width, size.Height))
	for x := 0; x < size.Width; x++ {
		for y := 0; y < heightCard; y++ {
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

func fontFace(opt truetype.Options) (font.Face, error) {
	// need japanese font to write japanese
	ftBin, err := ioutil.ReadFile(fontPath)
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

func exportImage(path string, image image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file(%s): %w", path, err)
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	if err := png.Encode(b, image); err != nil {
		return fmt.Errorf("failed to encode card image: %w", err)
	}
	if err := b.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}
	return nil
}
