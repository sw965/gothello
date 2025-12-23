package gothello

import (
	"slices"
)

var CellCols = []string{"a", "b", "c", "d", "e", "f", "g", "h"}

type Cell struct {
	Row    int
	Col    string
}

func (c Cell) ToIndex() int {
	row := c.Row - 1
	col := slices.Index(CellCols, c.Col)
	return (row*8) + col
}

func (c Cell) ToBitBoard() BitBoard {
	idx := c.ToIndex() 
	return 1 << idx
}
