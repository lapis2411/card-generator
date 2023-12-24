package decoder

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gocarina/gocsv"
	"github.com/lapis2411/card-generator/adapter"
	"github.com/lapis2411/card-generator/domain"
)

// Converter is an interface for converting csv file
type (
	// csvConverter is a struct for converting csv file
	csvDecoder struct {
	}
)

func NewCsvDecoder() adapter.Decoder {
	return &csvDecoder{}
}

// DecodeStyle returns single style by csv format byte array
func (c *csvDecoder) DecodeStyles(data []byte) (map[string]*domain.Style, error) {
	type styleCSV struct {
		Name     string  `csv:"name"`
		FontSize float64 `csv:"font_size"`
		X        int     `csv:"x"`
		Y        int     `csv:"y"`
		ColorR   int     `csv:"color_r"`
		ColorG   int     `csv:"color_g"`
		ColorB   int     `csv:"color_b"`
		ColorA   int     `csv:"color_a"`
	}
	s := map[string]*domain.Style{}
	if err := gocsv.UnmarshalBytesToCallback(data, func(sc styleCSV) error {
		p := image.Point{sc.X, sc.Y}
		rgba := color.RGBA{
			uint8(sc.ColorR),
			uint8(sc.ColorG),
			uint8(sc.ColorB),
			uint8(sc.ColorA),
		}
		if _, ok := s[sc.Name]; ok {
			fmt.Println("duplicated")
			return fmt.Errorf("style name %s is duplicated", sc.Name)
		}
		st := domain.NewStyle(sc.Name, p, sc.FontSize, rgba)
		s[sc.Name] = &st
		return nil
	}); err != nil {
		return nil, err
	}
	return s, nil
}

// DecodeCard returns information for generating cards by csv format byte array
func (c *csvDecoder) DecodeCards(data []byte, styles map[string]*domain.Style) (map[string][]domain.FormattedText, error) {
	type cardCSV struct {
		Name  string `csv:"name"`
		Style string `csv:"style"`
		Text  string `csv:"text"`
	}
	cds := make(map[string][]domain.FormattedText, 0)
	// cardname -> style
	dupStyle := make(map[string]map[string]struct{})
	if err := gocsv.UnmarshalBytesToCallback(data, func(cc cardCSV) error {
		sp, ok := styles[cc.Style]
		if !ok {
			return fmt.Errorf("style %s is not found", cc.Style)
		}
		if _, ok := dupStyle[cc.Name]; !ok {
			dupStyle[cc.Name] = make(map[string]struct{})
		}
		if _, ok := dupStyle[cc.Name][cc.Style]; ok {
			return fmt.Errorf("style %s is duplicated on card %s", cc.Style, cc.Name)
		}
		dupStyle[cc.Name][cc.Style] = struct{}{}
		if _, ok := cds[cc.Name]; !ok {
			cds[cc.Name] = make([]domain.FormattedText, 0)
		}
		cds[cc.Name] = append(cds[cc.Name], domain.NewFormattedText(cc.Text, sp))
		return nil
	}); err != nil {
		return nil, err
	}
	return cds, nil
}
