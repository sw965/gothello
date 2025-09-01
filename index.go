package gothello

import (
	"slices"
	omwbits "github.com/sw965/omw/math/bits"
)

const (
	UpLeftCornerIndex    = 0
	UpRightCornerIndex   = 7
	DownLeftCornerIndex  = 56
	DownRightCornerIndex = 63
)

var (
    UpLeft16MassIndices    = omwbits.OneIndices64(UpLeftBitBoard)
    UpRight16MassIndices   = omwbits.OneIndices64(UpRightBitBoard)
    DownLeft16MassIndices  = omwbits.OneIndices64(DownLeftBitBoard)
    DownRight16MassIndices = omwbits.OneIndices64(DownRightBitBoard)

    WhiteLineIndices = omwbits.OneIndices64(WhiteLineBitBoard)
    BlackLineIndices = omwbits.OneIndices64(BlackLineBitBoard)

    UpSideIndices   = omwbits.OneIndices64(UpSideBitBoard)
    DownSideIndices = omwbits.OneIndices64(DownSideBitBoard)
    LeftSideIndices  = omwbits.OneIndices64(LeftSideBitBoard)
    RightSideIndices = omwbits.OneIndices64(RightSideBitBoard)

    UpEdgeIndices    = omwbits.OneIndices64(UpEdgeBitBoard)
    DownEdgeIndices  = omwbits.OneIndices64(DownEdgeBitBoard)
    LeftEdgeIndices  = omwbits.OneIndices64(LeftEdgeBitBoard)
    RightEdgeIndices = omwbits.OneIndices64(RightEdgeBitBoard)
    EdgeIndices      = omwbits.OneIndices64(EdgeBitBoard)

    CornerIndices = omwbits.OneIndices64(CornerBitBoard)
    CIndices      = omwbits.OneIndices64(CBitBoard)
    AIndices      = omwbits.OneIndices64(ABitBoard)
    BIndices      = omwbits.OneIndices64(BBitBoard)
    XIndices      = omwbits.OneIndices64(XBitBoard)
)

func IsCornerIndex(idx int) bool {
	return slices.Contains(CornerIndices, idx)
}

func IndexToRowAndColumn(idx int) (int, int) {
	return idx/Cols, idx%Cols
}

func RowAndColumnToIndex(row, col int) int {
	return row * Cols + col
}

func TransposeIndex(idx int) int {
    row, col := idx/Cols, idx%Cols
    return col*Cols + row
}

func MirrorHorizontalIndex(idx int) int {
	row, col := IndexToRowAndColumn(idx)
	newCol := (Cols - 1) - col
	return row*Cols + newCol
}

func MirrorVerticalIndex(idx int) int {
	row, col := IndexToRowAndColumn(idx)
	newRow := (Rows - 1) - row
	return newRow*Cols + col
}

func Rotate90Index(idx int) int {
	row, col := IndexToRowAndColumn(idx)
	newRow := col
	newCol := (Rows - 1) - row
	return newRow*Cols + newCol
}

func Rotate180Index(idx int) int {
	row, col := IndexToRowAndColumn(idx)
	newRow := (Rows - 1) - row
	newCol := (Cols - 1) - col
	return newRow*Cols + newCol
}

func Rotate270Index(idx int) int {
	row, col := IndexToRowAndColumn(idx)
	newRow := (Cols - 1) - col
	newCol := row
	return newRow*Cols + newCol
}

var GroupIndexTable = [][]int{
	[]int{0, 7, 56, 63},
	[]int{1, 6, 8, 15, 48, 55, 57, 62},
	[]int{2, 5, 16, 23, 40, 47, 58, 61},
	[]int{3, 4, 24, 31, 32, 39, 59, 60},
	[]int{9, 14, 49, 54},
	[]int{10, 13, 17, 22, 41, 46, 50, 53},
	[]int{11, 12, 25, 30, 33, 38, 51, 52},
	[]int{18, 21, 42, 45},
	[]int{19, 20, 26, 29, 34, 37, 43, 44},
	[]int{27, 28, 35, 36},
}

var GroupIdByIndex = func() []int {
	groupIds := make([]int, BoardSize)
	for id, gropuIdxs := range GroupIndexTable {
		for _, idx := range gropuIdxs {
			groupIds[idx] = id
		}
	}
	return groupIds
}()