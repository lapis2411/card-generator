package domain

type (
	// CardContent is a struct for generating card image
	// @TODO CARDの数が多すぎる場合はPoolを検討する
	Card struct {
		name           string
		formattedTexts []FormattedText
	}

	CardAdapter interface {
		// Cards returns all cards by map name to FormattedText and styles
		Fetch(style, card []byte) ([]Card, error)
	}
)

// NewCard returns new card
func NewCard(n string, f []FormattedText) Card {
	return Card{
		name:           n,
		formattedTexts: f,
	}
}

// Name returns card name
func (c Card) Name() string {
	return c.name
}

// FormattedTexts returns formatted texts
func (c Card) FormattedTexts() []FormattedText {
	return c.formattedTexts
}
