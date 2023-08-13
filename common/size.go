package common

type (
	Size struct {
		width  int
		height int
	}
)

func (s Size) Width() int {
	return s.width
}

func (s Size) Height() int {
	return s.height
}
