package gothello

const (
	Rows = 8
	Cols = 8
	BoardSize = Rows * Cols
)

type Disc int

const (
	Empty Disc = iota
	Black
	White
)

func (d Disc) Opposite() Disc {
	switch d {
	case Black:
		return White
	case White:
		return Black
	}
	return Empty
}

var AllDiscs = []Disc{Empty, Black, White}