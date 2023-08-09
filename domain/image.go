package domain

import (
	"image"
)

type (
	Size struct {
		Width  int
		Height int
	}

	Image struct {
		Image *image.RGBA
		Size  Size
	}

	ImageAdapter interface {
		// GenerateCardImages generates card images from card contents
		GenerateCardImages([]Card) (map[string]Image, error)
		// GeneratePrintImages generates card image for print from card contents
		GeneratePrintImages([]Image) ([]Image, error)
	}

	ExportImageAdapter interface {
		// Save saves images to file
		Save([]Image, string) error
	}
)
