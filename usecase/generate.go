package usecase

import (
	"fmt"

	"github.com/lapis2411/card-generator/domain"
)

type (
	Generate struct {
		cardAdapter  domain.CardAdapter
		imageAdapter domain.ImageAdapter
	}
)

func NewGenerate(ca domain.CardAdapter, ir domain.ImageAdapter) *Generate {
	return &Generate{
		cardAdapter:  ca,
		imageAdapter: ir,
	}
}

func (g Generate) GenerateCardImages(styles, cards []byte) ([]domain.Image, error) {
	c, err := g.cardAdapter.Fetch(styles, cards)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}
	ci, err := g.imageAdapter.GenerateCardImages(c)
	if err != nil {
		return nil, fmt.Errorf("failed to generate card images: %w", err)
	}
	return ci, nil
}

// GeneratePrintImages generates card image for print from card contents
func (g Generate) GeneratePrintImages(styles, cards []byte) ([]domain.Image, error) {
	ci, err := g.GenerateCardImages(styles, cards)
	if err != nil {
		return nil, fmt.Errorf("failed to generate card images: %w", err)
	}
	pi, err := g.imageAdapter.GeneratePrintImages(ci)
	if err != nil {
		return nil, fmt.Errorf("failed to merge card images: %w", err)
	}
	return pi, nil
}
