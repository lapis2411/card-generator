package usecase

import (
	"fmt"

	"github.com/lapis2411/card-generator/domain"
)

type (
	// Export has a dependency on Generate
	Export struct {
		generate           Generate
		exportImageAdapter domain.ExportImageAdapter
	}
)

func NewExport(ca domain.CardAdapter, ia domain.ImageAdapter) Export {
	return Export{
		generate: NewGenerate(ca, ia),
	}
}

func (e Export) ExportCardImages(styles, cards []byte) error {
	ci, err := e.generate.GenerateCardImages(styles, cards)
	if err != nil {
		return fmt.Errorf("failed to generate card images in export image: %w", err)
	}
	return e.exportImageAdapter.Save(ci)
}

func (e Export) ExportCardImagesForPrint(styles, cards []byte) error {
	ci, err := e.generate.GeneratePrintImages(styles, cards)
	if err != nil {
		return fmt.Errorf("failed to generate card images in export image: %w", err)
	}
	return e.exportImageAdapter.Save(ci)
}
