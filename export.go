package generator

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"
)

// func Export(cards CardContents) error {
// 	folder := filepath.Dir(g.output + "/")
// 	err := os.MkdirAll(folder, 0755)
// 	if err != nil {
// 		return fmt.Errorf("failed to create folder(%s): %w", folder, err)
// 	}
// 	for name, card := range cards {
// 		if err := g.generateCard(*card, name, widthCard, heightCard); err != nil {
// 			return fmt.Errorf("failed to generate card(%s): %w", name, err)
// 		}
// 	}
// 	return nil
// }

func (cs Canvases) ExportImages(dir string) error {
	d := dir + "/paper_%d.png"
	for i, img := range cs {
		p := fmt.Sprintf(d, i)
		if err := exportImage(p, img.Image); err != nil {
			return fmt.Errorf("failed to export paper(%s): %w", d, err)
		}
	}
	return nil
}

func (cs Cards) ExportImages(dir string) error {
	d := dir + "/card_%d.png"
	for i, img := range cs {
		p := fmt.Sprintf(d, i)
		if err := exportImage(p, img.Image); err != nil {
			return fmt.Errorf("failed to export card(%s): %w", d, err)
		}
	}
	return nil
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
