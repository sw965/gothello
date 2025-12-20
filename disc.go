package gothello

const (
	Rows      = 8
	Cols      = 8
	BoardSize = Rows * Cols
)

type Color int

const (
	Empty Color = iota
	Black
	White
)

func (c Color) Opposite() Color {
	switch c {
	case Black:
		return White
	case White:
		return Black
	}
	return Empty
}

var AllColors = []Color{Empty, Black, White}

type Perspective int

const (
	None Perspective = iota
	Self
	Opponent
)

type Perspectives []Perspective

var AllPerspectives = Perspectives{None, Self, Opponent}

func (ps Perspectives) ToPartialFeature1D() PartialFeature1D {
	f := make(PartialFeature1D, len(ps))
	for i, p := range ps {
		f[i] = p
	}
	return f
}