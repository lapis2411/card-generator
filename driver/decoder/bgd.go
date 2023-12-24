package decoder

import (
	"encoding/json"
	"fmt"
	"github.com/lapis2411/card-generator/adapter"
	"github.com/lapis2411/card-generator/domain"
	"strconv"
)

// Converter is an interface for converting csv file
type (
	// csvConverter is a struct for converting csv file
	bgdDecoder struct {
		csvDecoder
	}
	cardData []struct {
		ID       int        `json:"id"`
		Name     string     `json:"name"`
		TextData []textData `json:"data"`
	}
	textData struct {
		StyleID int    `json:"styleID"`
		Text    string `json:"text"`
	}
)

func NewBGDDecoder() adapter.Decoder {
	return &bgdDecoder{}
}

// DecodeStyle は既存のものを使う (fileから読み込むことを想定)
//func (c *jsonDecoder) DecodeStyles(data []byte) (map[string]*domain.Style, error) {
//}

// DecodeCard returns information for generating cards by csv format byte array
func (c *bgdDecoder) DecodeCards(data []byte, styles map[string]*domain.Style) (map[string][]domain.FormattedText, error) {

	var jc cardData
	cds := make(map[string][]domain.FormattedText, 0)
	err := json.Unmarshal(data, &jc)
	if err != nil {
		return nil, err
	}

	duplicate := map[string]struct{}{}
	for _, v := range jc {
		ids := strconv.Itoa(v.ID)
		if _, exist := duplicate[ids]; exist {
			return nil, fmt.Errorf("%s is duplicated", ids)
		}
		duplicate[ids] = struct{}{}
		cds[ids], err = toFormattedText(v.TextData, styles)
		if err != nil {
			return nil, err
		}
	}
	return cds, nil
}

func toFormattedText(td []textData, st map[string]*domain.Style) ([]domain.FormattedText, error) {
	ft := make([]domain.FormattedText, 0, len(td))
	for _, v := range td {
		sname := strconv.Itoa(v.StyleID)
		st, ok := st[sname]
		if !ok {
			return nil, fmt.Errorf("%s style is undefined", sname)
		}
		ft = append(ft, domain.NewFormattedText(v.Text, st))
	}
	return ft, nil
}
