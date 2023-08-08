package domain

import "golang.org/x/image/math/fixed"

// StyledText is a struct for text and style
// This value is used for generating Card
type FormattedText struct {
	sentence string
	style    *Style // have pointer to style to save memory
}

func (s *FormattedText) Point26_6() fixed.Point26_6 {
	return fixed.Point26_6{
		X: fixed.Int26_6(s.style.position.X * 64),
		Y: fixed.Int26_6(s.style.position.Y * 64),
	}
}
