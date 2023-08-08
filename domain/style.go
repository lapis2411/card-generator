package domain

import (
	"image"
	"image/color"
)

type (

	// Styles is a map for text style definition by style name
	// Styles map[string]*Style

	// Style is a struct for text style definition
	// can define position, fontsize and color
	Style struct {
		name     string
		position image.Point
		fontsize float64
		rgba     color.RGBA
	}

	StyleRepository interface {
		// Styles returns all styles
		StylesMap() (map[string]Style, error)
	}
)

func NewStyle(name string, position image.Point, fontsize float64, rgba color.RGBA) Style {
	return Style{
		name:     name,
		position: position,
		fontsize: fontsize,
		rgba:     rgba,
	}
}

// StylePointer returns pointer to style by style name
// func (s Styles) StylePointer(name string) (*Style, error) {
// 	if _, e := s[name]; !e {
// 		return nil, errors.New("style is undefined")
// 	}
// 	return s[name], nil
// }
