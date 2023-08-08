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

	ImageRepository interface {
		// GenerateCardImages generates card images from card contents
		GenerateCardImages([]Card) ([]Image, error)
		// GeneratePrintImages generates card image for print from card contents
		GeneratePrintImages([]Image) ([]Image, error)
		// Save saves images to file
		Save([]Image, string) error
	}
)
