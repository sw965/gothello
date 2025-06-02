package gothello

const (
	UpLeftCornerIndex = 0
	UpRightCornerIndex = 7
	DownLeftCornerIndex = 56
	DownRightCornerIndex = 63
)

const (
	UpLeftXIndex = 9
	UpRightXIndex = 14
	DownLeftXIndex = 49
    DownRightXIndex = 54
)

var (
	UpSideIndices = func() []int {
		idxs := make([]int, Cols)
		for i := range idxs {
			idxs[i] = i
		}
		return idxs
	}()

	DownSideIndices = func() []int {
		idxs := make([]int, Cols)
		for i := range idxs {
			idxs[i] = DownLeftCornerIndex + i
		}
		return idxs
	}()

	LeftSideIndices = func() []int {
		idxs := make([]int, Rows)
		for i := range idxs {
			idxs[i] = i * Cols
		}
		return idxs
	}()

	RightSideIndices = func() []int {
		idxs := make([]int, Rows)
		for i := range idxs {
			idxs[i] = UpRightCornerIndex + (i * Cols)
		}
		return idxs
	}()

	UpEdgeIndices = UpSideIndices[1:7]
	DownEdgeIndices = DownSideIndices[1:7]
	LeftEdgeIndices = LeftSideIndices[1:7]
	RightEdgeIndices = RightSideIndices[1:7]
)

func IndexToRowAndColumn(idx int) (int, int) {
	return idx/Cols, idx%Cols
}

func RowAndColumnToIndex(row, col int) int {
	return row * Cols + col
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