package generator

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"
)

var sizeA4 = Size{
	Width:  2480,
	Height: 3508,
}
var sizeCard = Size{
	Width:  widthCard,
	Height: heightCard,
}

func MergeCards(dir string) error {
	return mergeCards(dir, sizeA4, sizeCard)
}

// MergeCards merges cards into one image.
func mergeCards(dir string, paperSize, cardSize Size) error {
	row := paperSize.Width / cardSize.Width
	column := paperSize.Height / cardSize.Height
	img := clearA4(paperSize)
	imgs := make(map[int]image.Image, 1)
	cnt := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking the path %v: %w\n", dir, err)
		}
		if !info.IsDir() {
			// どこから書き出すか？
			x := (cnt % row) * cardSize.Width
			y := (cnt / row) * cardSize.Height
			// 次ここから開発
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", dir, err)
	}

	// ファイルに保存
	od := "./tmp1.png"
	if err := exportImage(od, img); err != nil {
		return fmt.Errorf("failed to export card(%s): %w", od, err)
	}

	return nil
}

func clearA4(size Size) image.Image {

	img := image.NewRGBA(image.Rect(0, 0, size.Width, size.Height))

	white := color.RGBA{255, 255, 255, 255}
	for y := 0; y < size.Height; y++ {
		for x := 0; x < size.Width; x++ {
			img.Set(x, y, white)
		}
	}
	return img
}

func mergeImage(imgs map[int]image.Image, size Size) image.Image {
	img := clearA4(size)
	for _, v := range imgs {
		// img = drawImage(img, v)
	}
	return img
}
