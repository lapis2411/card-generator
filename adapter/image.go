package adapter

import (
	"fmt"

	"github.com/lapis2411/card-generator/common"
	"github.com/lapis2411/card-generator/domain"
)

type (
	// ImageAdapter is an interface for generating images
	imageAdapter struct {
		encoder  ImageDriver
		template Template
		cardSize common.Size
	}
	ImageDriver interface {
		ImageEncode(domain.Card) (domain.Image, error)
	}
	Template interface {
		Arrange([]domain.Image) ([]domain.Image, error)
	}
)

func NewImageAdapter(imgd ImageDriver, t Template) domain.ImageAdapter {
	return &imageAdapter{
		encoder:  imgd,
		template: t,
	}
}

func (ia *imageAdapter) GenerateCardImages(cards []domain.Card) ([]domain.Image, error) {
	cs := make([]domain.Image, 0, len(cards))

	for i, c := range cards {
		img, err := ia.encoder.ImageEncode(c)
		if err != nil {
			return nil, fmt.Errorf("failed to generate %v th card: %w", i, err)
		}
		cs = append(cs, img)
	}
	return cs, nil
}

func (ia *imageAdapter) GeneratePrintImages(imgs []domain.Image) ([]domain.Image, error) {
	return ia.template.Arrange(imgs)
}
