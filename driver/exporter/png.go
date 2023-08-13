package driver

import (
	"bufio"
	"fmt"
	"image/png"
	"os"

	"github.com/lapis2411/card-generator/adapter"
	"github.com/lapis2411/card-generator/domain"
)

type (
	PngExporter struct{}
)

func NewPngExporter() adapter.ExportDriver {
	return PngExporter{}
}

func (p PngExporter) Save(img domain.Image, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file(%s): %w", path, err)
	}
	defer f.Close()

	b := bufio.NewWriter(f)
	if err := png.Encode(b, img.Image()); err != nil {
		return fmt.Errorf("failed to encode card image: %w", err)
	}
	if err := b.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}
	return nil
}
