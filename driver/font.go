package driver

import (
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
)

func DrawStringWithWrap(d *font.Drawer, s string, w int) {
	originX := d.Dot.X
	progress := fixed.Int26_6(0)
	width := fixed.I(w)

	prevC := rune(-1)
	met := d.Face.Metrics()
	for _, c := range s {
		if prevC >= 0 {
			d.Dot.X += d.Face.Kern(prevC, c)
		}
		advance, ok := d.Face.GlyphAdvance(c)
		if !ok {
			c = '?'
			advance, _ = d.Face.GlyphAdvance(c)
		}
		progress += advance
		if progress > width {
			d.Dot.X = originX
			d.Dot.Y += met.Height
			progress = 0
		}
		dr, mask, maskp, advance, _ := d.Face.Glyph(d.Dot, c)
		if !dr.Empty() {
			draw.DrawMask(d.Dst, dr, d.Src, image.Point{}, mask, maskp, draw.Over)
		}
		d.Dot.X += advance
		prevC = c
	}
}
