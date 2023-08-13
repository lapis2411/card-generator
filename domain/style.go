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

	StyleAdapter interface {
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

func (s Style) FontSize() float64 {
	return s.fontsize
}

func (s Style) RGBA() color.RGBA {
	return s.rgba
}
