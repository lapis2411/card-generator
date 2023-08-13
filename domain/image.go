package domain

import (
	"image"

	"github.com/lapis2411/card-generator/common"
)

type (
	Image struct {
		image *image.RGBA
		size  common.Size
		name  string
	}

	ImageAdapter interface {
		// GenerateCardImages generates card images from card contents
		GenerateCardImages([]Card) ([]Image, error)
		// GeneratePrintImages generates card image for print from card contents
		GeneratePrintImages([]Image) ([]Image, error)
	}
	ExportImageAdapter interface {
		// Save saves images to file
		Save([]Image) error
	}
)

// NewImage returns new image
func NewImage(i *image.RGBA, s common.Size, n string) Image {
	return Image{
		image: i,
		size:  s,
		name:  n,
	}
}

// Image returns image
// @todo deepcopy?
func (i Image) Image() *image.RGBA {
	return i.image
}

// Name returns name
func (i Image) Name() string {
	return i.name
}

// Size returns size
func (i Image) Size() common.Size {
	return i.size
}
