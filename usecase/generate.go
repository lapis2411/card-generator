package usecase

import (
	"fmt"

	"github.com/lapis2411/card-generator/domain"
)

type (
	Generate struct {
		cardRepository  domain.CardRepository
		imageRepository domain.ImageRepository
	}
)

func NewGenerate(cr domain.CardRepository, ir domain.ImageRepository) *Generate {
	return &Generate{
		cardRepository:  cr,
		imageRepository: ir,
	}
}

func (g Generate) GenerateCardImages() ([]domain.Image, error) {
	c, err := g.cardRepository.Cards()
	if err != nil {
		return nil, fmt.Errorf("failed to get cards: %w", err)
	}
	ci, err := g.imageRepository.GenerateCardImages(c)
	if err != nil {
		return nil, fmt.Errorf("failed to generate card images: %w", err)
	}
	return ci, nil
}

// GeneratePrintImages generates card image for print from card contents
func (g Generate) GeneratePrintImages() ([]domain.Image, error) {
	ci, err := g.GenerateCardImages()
	if err != nil {
		return nil, fmt.Errorf("failed to generate card images: %w", err)
	}
	pi, err := g.imageRepository.GeneratePrintImages(ci)
	if err != nil {
		return nil, fmt.Errorf("failed to merge card images: %w", err)
	}
	return pi, nil
}
