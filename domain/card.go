package domain

import "errors"

type (
	// CardContent is a struct for generating card image
	// @TODO CARDの数が多すぎる場合はPoolを検討する
	Card struct {
		formattedTexts []FormattedText
		styles         map[string]struct{}
	}

	CardRepository interface {
		// Cards returns all cards by map name to FormattedText and styles
		Fetch() ([]Card, error)
	}
)

// AddStyleText adds styled text to the card of cards
// name is card name, text is description or some value, style is style for text pointer
// func (c *CardContents) AddStyleText(card string, text string, style *Style) error {
// 	if style == nil {
// 		return errors.New("style is undefined")
// 	}
// 	if _, ok := (*c)[card]; !ok {
// 		(*c)[card] = &CardContent{}
// 	}
// 	return (*c)[card].addStyleText(text, style)
// }

// AddStyleText adds styled text to single card
// text is description or some value, style is style for text pointer
func (c *Card) AddStyleText(sentence string, style *Style) error {
	if _, ok := (*c).styles[style.name]; ok {
		return errors.New("style is duplicated")
	}
	if (*c).styles == nil {
		(*c).styles = make(map[string]struct{})
	}
	(*c).styles[style.name] = struct{}{}
	if (*c).formattedTexts == nil {
		(*c).formattedTexts = make([]FormattedText, 0)
	}
	(*c).formattedTexts = append((*c).formattedTexts, FormattedText{
		sentence: sentence,
		style:    style,
	})
	return nil
}
