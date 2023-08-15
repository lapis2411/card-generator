package common

type (
	Size struct {
		width  int
		height int
	}
)

func NewSize(w, h int) Size {
	return Size{
		width:  w,
		height: h,
	}
}

func (s Size) Width() int {
	return s.width
}

func (s Size) Height() int {
	return s.height
}
