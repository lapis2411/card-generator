package adapter

import (
	"fmt"
	"sort"

	"github.com/lapis2411/card-generator/domain"
)

type (
	cardAdapter struct {
		decoder Decoder
	}

	Decoder interface {
		// DecodeStyle returns single style
		DecodeStyles([]byte) (map[string]*domain.Style, error)
		// DecodeCard returns single card
		DecodeCards([]byte, map[string]*domain.Style) (map[string][]domain.FormattedText, error)
	}
)

func NewCardAdapter(d Decoder) domain.CardAdapter {
	return &cardAdapter{
		decoder: d,
	}
}

func (cr cardAdapter) Fetch(style, card []byte) ([]domain.Card, error) {
	ss, err := cr.decoder.DecodeStyles(style)
	if err != nil {
		return nil, fmt.Errorf("failed to decode styles: %w", err)
	}
	cts, err := cr.decoder.DecodeCards(card, ss)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cards: %w", err)
	}
	keys := make([]string, 0, len(cts))
	for key := range cts {
		keys = append(keys, key)
	}
	sk := sort.StringSlice(keys)

	cds := make([]domain.Card, 0, len(cts))
	for _, key := range sk {
		cds = append(cds, domain.NewCard(cts[key], key))
	}

	return cds, nil
}
