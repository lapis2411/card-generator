package adapter

import (
	"fmt"

	"github.com/lapis2411/card-generator/domain"
	"github.com/lapis2411/card-generator/common"
)

type (
	// ImageAdapter is an interface for generating images
	imageAdapter struct {
		encoder  Encoder
		template Template
		size    common.Size
	}
	Encoder interface {
		EncodeImage(domain.Card, common.Size) (domain.Image, error)
	}
	Template interface {
		Arrange([]domain.Image) ([]domain.Image, error)
	}
)

func NewImageAdapter(e Encoder, t Template, common.Size) domain.ImageAdapter {
	return &imageAdapter{
		encoder:  e,
		template: t,
		size:     s,
	}
}

func (ia *imageAdapter) GenerateCardImages(cards []domain.Card) ([]domain.Image, error) {
	cs := make([]domain.Image, 0, len(cards))

	for i, c := range cards {
		img, err := ia.encoder.EncodeCard(c, ia.size)
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
