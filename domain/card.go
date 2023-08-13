package domain

type (
	// CardContent is a struct for generating card image
	// @TODO CARDの数が多すぎる場合はPoolを検討する
	Card struct {
		formattedTexts []FormattedText
		name           string
	}

	CardAdapter interface {
		// Cards returns all cards by map name to FormattedText and styles
		Fetch(style, card []byte) ([]Card, error)
	}
)

// NewCard returns new card
func NewCard(f []FormattedText, n string) Card {
	return Card{
		formattedTexts: f,
		name:           n,
	}
}

// FormattedTexts returns formatted texts
func (c Card) FormattedTexts() []FormattedText {
	return c.formattedTexts
}

// Name returns name
func (c Card) Name() string {
	return c.name
}
