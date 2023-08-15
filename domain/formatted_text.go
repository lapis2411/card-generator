package domain

import "golang.org/x/image/math/fixed"

// FormattedText is a struct for text and style
// This value is used for generating Card
type FormattedText struct {
	text  string
	style *Style // have pointer to style to save memory
}

func NewFormattedText(sentence string, style *Style) FormattedText {
	return FormattedText{
		text:  sentence,
		style: style,
	}
}

func (s *FormattedText) Point26_6() fixed.Point26_6 {
	return fixed.Point26_6{
		X: fixed.Int26_6(s.style.position.X * 64),
		Y: fixed.Int26_6(s.style.position.Y * 64),
	}
}

func (ft FormattedText) Style() Style {
	return *ft.style
}

func (ft FormattedText) Text() string {
	return ft.text
}
