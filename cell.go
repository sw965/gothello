package gothello

import (
	"slices"
)

type Cell struct {
	Row    int
	Column string
}

func (c *Cell) ToBitBoard() BitBoard {
	cs := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	row := c.Row - 1
	col := slices.Index(cs, c.Column)
	idx := RowColumnToIndex(row, col)
	return 1 << idx
}
