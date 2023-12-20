package decoder

import (
	"github.com/lapis2411/card-generator/adapter"
	"github.com/lapis2411/card-generator/domain"
)

// Converter is an interface for converting csv file
type (
	// csvConverter is a struct for converting csv file
	jsonDecoder struct {
	}
)

func NewJsonDecoder() adapter.Decoder {
	return &jsonDecoder{}
}

// DecodeStyle returns single style by csv format byte array
func (c *jsonDecoder) DecodeStyles(data []byte) (map[string]*domain.Style, error) {
	s := map[string]*domain.Style{}

	return s, nil
}

// DecodeCard returns information for generating cards by csv format byte array
func (c *jsonDecoder) DecodeCards(data []byte, styles map[string]*domain.Style) (map[string][]domain.FormattedText, error) {
	type jsonCard []struct {
		ID   int `json:"id"`
		Data []struct {
			StyleID int    `json:"styleID"`
			Text    string `json:"text"`
		} `json:"data"`
	}
	cds := make(map[string][]domain.FormattedText, 0)

	return cds, nil
}
