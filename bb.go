package gothello

import (
	"fmt"
	omwbits "github.com/sw965/omw/math/bits"
	"math"
	"math/bits"
)

type BitBoard uint64

const MaxBitBoard = ^BitBoard(0)

const (
	UpLeftBitBoard    = BitBoard(0b00000000_00000000_00000000_00000000_00001111_00001111_00001111_00001111)
	UpRightBitBoard   = BitBoard(0b00000000_00000000_00000000_00000000_11110000_11110000_11110000_11110000)
	DownLeftBitBoard  = BitBoard(0b00001111_00001111_00001111_00001111_00000000_00000000_00000000_00000000)
	DownRightBitBoard = BitBoard(0b11110000_11110000_11110000_11110000_00000000_00000000_00000000_00000000)

	BoxBitBoard          = BitBoard(0b00000000_00000000_00111100_00111100_00111100_00111100_00000000_00000000)
	BoxUpSideBitBoard    = BitBoard(0b00000000_00000000_00000000_00000000_00000000_00111100_00000000_00000000)
	BoxDownSideBitBoard  = BitBoard(0b00000000_00000000_00111100_00000000_00000000_00000000_00000000_00000000)
	BoxLeftSideBitBoard  = BitBoard(0b00000000_00000000_00000001_00000001_00000001_00000001_00000000_00000000)
	BoxRightSideBitBoard = BitBoard(0b00000000_00000000_10000000_10000000_10000000_10000000_00000000_00000000)

	WhiteLineBitBoard = BitBoard(0b10000000_01000000_00100000_00010000_00001000_00000100_00000010_00000001)
	BlackLineBitBoard = BitBoard(0b00000001_00000010_00000100_00001000_00010000_00100000_01000000_10000000)

	UpSideBitBoard    = BitBoard(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_11111111)
	DownSideBitBoard  = BitBoard(0b11111111_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
	LeftSideBitBoard  = BitBoard(0b00000001_00000001_00000001_00000001_00000001_00000001_00000001_00000001)
	RightSideBitBoard = BitBoard(0b10000000_10000000_10000000_10000000_10000000_10000000_10000000_10000000)
	SideBitBoard      = UpSideBitBoard | DownSideBitBoard | LeftSideBitBoard | RightSideBitBoard

	UpEdgeBitBoard    = BitBoard(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01111110)
	DownEdgeBitBoard  = BitBoard(0b01111110_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
	LeftEdgeBitBoard  = BitBoard(0b00000000_00000001_00000001_00000001_00000001_00000001_00000001_00000000)
	RightEdgeBitBoard = BitBoard(0b00000000_10000000_10000000_10000000_10000000_10000000_10000000_00000000)
	EdgeBitBoard      = UpEdgeBitBoard | DownEdgeBitBoard | LeftEdgeBitBoard | RightEdgeBitBoard

	CornerBitBoard = BitBoard(0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_10000001)
	CBitBoard      = BitBoard(0b01000010_10000001_00000000_00000000_00000000_00000000_10000001_01000010)
	ABitBoard      = BitBoard(0b00100100_00000000_10000001_00000000_00000000_10000001_00000000_00100100)
	BBitBoard      = BitBoard(0b00011000_00000000_00000000_10000001_10000001_00000000_00000000_00011000)
	XBitBoard      = BitBoard(0b00000000_01000010_00000000_00000000_00000000_00000000_01000010_00000000)

	UpBAABBitBoard    = BitBoard(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_00111100)
	DownBAABBitBoard  = BitBoard(0b00111100_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
	LeftBAABBitBoard  = BitBoard(0b00000000_00000000_00000001_00000001_00000001_00000001_00000000_00000000)
	RightBAABBitBoard = BitBoard(0b00000000_00000000_10000000_10000000_10000000_10000000_00000000_00000000)

	UpMidSideBitBoard    = BitBoard(0b00000000_00000000_00000000_00000000_00000000_00000000_00111100_00000000)
	DownMidSideBitBoard  = BitBoard(0b00000000_00111100_00000000_00000000_00000000_00000000_00000000_00000000)
	LeftMidSideBitBoard  = BitBoard(0b00000000_00000000_00000010_00000010_00000010_00000010_00000000_00000000)
	RightMidSideBitBoard = BitBoard(0b00000000_00000000_01000000_01000000_01000000_01000000_00000000_00000000)
)

var AdjacentBySingle = func() map[BitBoard]BitBoard {
	m := map[BitBoard]BitBoard{}
	for i := 0; i < BoardSize; i++ {
		single := BitBoard(1) << i
		var adj BitBoard = 0
		row, col := IndexToRowColumn(i)

		// 右端でなければ右方向へ
		if col < Cols-1 {
			adj |= ShiftRight(single)
		}
		// 左端でなければ左方向へ
		if col > 0 {
			adj |= ShiftLeft(single)
		}
		// 上端でなければ上方向へ
		if row > 0 {
			adj |= ShiftUp(single)
		}
		// 下端でなければ下方向へ
		if row < Rows-1 {
			adj |= ShiftDown(single)
		}
		// 右上方向
		if row > 0 && col < Cols-1 {
			adj |= ShiftUpRight(single)
		}
		// 左上方向
		if row > 0 && col > 0 {
			adj |= ShiftUpLeft(single)
		}
		// 右下方向
		if row < Rows-1 && col < Cols-1 {
			adj |= ShiftDownRight(single)
		}
		// 左下方向
		if row < Rows-1 && col > 0 {
			adj |= ShiftDownLeft(single)
		}
		m[single] = adj
	}
	return m
}()

func (bb BitBoard) ToggleBit(idx int) (BitBoard, error) {
	max := BoardSize - 1
	if idx < 0 || idx > max {
		return 0, fmt.Errorf("idxは0から%dでなければならない。", max)
	}
	return bb ^ (1 << idx), nil
}

// 転置行列
func (bb BitBoard) Transpose() BitBoard {
	var t BitBoard
	t = (bb ^ (bb >> 7)) & 0b00000000_10101010_00000000_10101010_00000000_10101010_00000000_10101010
	bb = bb ^ t ^ (t << 7)
	t = (bb ^ (bb >> 14)) & 0b00000000_00000000_11001100_11001100_00000000_00000000_11001100_11001100
	bb = bb ^ t ^ (t << 14)
	t = (bb ^ (bb >> 28)) & 0b00000000_00000000_00000000_00000000_11110000_11110000_11110000_11110000
	bb = bb ^ t ^ (t << 28)
	return bb
}

// 横回転
func (bb BitBoard) MirrorHorizontal() BitBoard {
	bb = ((bb & 0b11110000_11110000_11110000_11110000_11110000_11110000_11110000_11110000) >> 4) |
		((bb & 0b00001111_00001111_00001111_00001111_00001111_00001111_00001111_00001111) << 4)
	bb = ((bb & 0b11001100_11001100_11001100_11001100_11001100_11001100_11001100_11001100) >> 2) |
		((bb & 0b00110011_00110011_00110011_00110011_00110011_00110011_00110011_00110011) << 2)
	bb = ((bb & 0b10101010_10101010_10101010_10101010_10101010_10101010_10101010_10101010) >> 1) |
		((bb & 0b01010101_01010101_01010101_01010101_01010101_01010101_01010101_01010101) << 1)
	return bb
}

// 縦回転
func (bb BitBoard) MirrorVertical() BitBoard {
	return ((bb & 0b11111111_00000000_00000000_00000000_00000000_00000000_00000000_00000000) >> 56) |
		((bb & 0b00000000_11111111_00000000_00000000_00000000_00000000_00000000_00000000) >> 40) |
		((bb & 0b00000000_00000000_11111111_00000000_00000000_00000000_00000000_00000000) >> 24) |
		((bb & 0b00000000_00000000_00000000_11111111_00000000_00000000_00000000_00000000) >> 8) |
		((bb & 0b00000000_00000000_00000000_00000000_11111111_00000000_00000000_00000000) << 8) |
		((bb & 0b00000000_00000000_00000000_00000000_00000000_11111111_00000000_00000000) << 24) |
		((bb & 0b00000000_00000000_00000000_00000000_00000000_00000000_11111111_00000000) << 40) |
		((bb & 0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_11111111) << 56)
}

func (bb BitBoard) Rotate90() BitBoard {
	bb = bb.MirrorVertical()
	bb = bb.Transpose()
	return bb
}

func (bb BitBoard) Rotate180() BitBoard {
	return BitBoard(bits.Reverse64(uint64(bb)))
}

func (bb BitBoard) Rotate270() BitBoard {
	bb = bb.Transpose()
	bb = bb.MirrorVertical()
	return bb
}

/*
https://blog.qmainconts.dev/articles/yxiplk2_dd
https://qiita.com/sensuikan1973/items/459b3e11d91f3cb37e43
上記の参考文献では、最上位ビットを(0, 0)の地点と見なしているが、
このライブラリでは、(0, 0)の地点を最下位ビットと見なしている為、
左シフトと右シフトの役割が逆になっている。
*/
func (bb BitBoard) Legals(opp BitBoard) BitBoard {
	/*
		横方向を返す合法座標を探す為のビットボード。
		[[0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]]

		0の座標に置かれた石は、横方向からひっくり返す事は出来ない。
	*/
	horizontal := opp & 0b01111110_01111110_01111110_01111110_01111110_01111110_01111110_01111110

	/*
		縦方向を返す合法座標を探す為のビットボード。
		[[0 0 0 0 0 0 0 0]
		 [1 1 1 1 1 1 1 1]
		 [1 1 1 1 1 1 1 1]
		 [1 1 1 1 1 1 1 1]
		 [1 1 1 1 1 1 1 1]
		 [1 1 1 1 1 1 1 1]
		 [1 1 1 1 1 1 1 1]
		 [0 0 0 0 0 0 0 0]]

		0の座標に置かれた石は、縦方向からひっくり返す事は出来ない。
	*/
	vertical := opp & 0b00000000_11111111_11111111_11111111_11111111_11111111_11111111_00000000

	/*
		斜め方向を返す合法座標を探す為のビットボード。
		[[0 0 0 0 0 0 0 0]
		 [0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]
		 [0 1 1 1 1 1 1 0]
		 [0 0 0 0 0 0 0 0]]

		0の座標に置かれた石は、斜め方向からひっくり返す事は出来ない。
	*/
	sideCut := opp & 0b00000000_01111110_01111110_01111110_01111110_01111110_01111110_00000000

	empty := ^(bb | opp)

	/*
		右方向への探索に関して説明。他の方向も考え方は同じ。
		1.(b << 1) で左シフトすわなち、自分の石を、全て右方向に移動する(右方向移動した後のボードをAとする)。
		2.Aとhorizontalを用いて、相手の石と重なっている場所をAND演算求める(重なっている場所を1とする)(このボードをrightとする)。
		3.rightは、この時点で、自分の石の1つ右側にある相手の石座標を保持している。
		4.rightを左シフトすわなち、再度右方向に移動させて、horizontalとAND演算をする事で、自分の石の1つ右側にある相手の石の座標と、更に1つ右側にある相手の石の座標を求める。
		5. right |= とする事で、自分の石の1つ右側にある相手の石座標と、更に1つ右側にある相手の石の座標を保持する。
		6. 後は4と5を繰り返す事で、自分の石から右側にあるn個連結した相手の石の座標を求める事が出来る。
	*/
	right := horizontal & ShiftRight(bb)
	left := horizontal & ShiftLeft(bb)
	up := vertical & ShiftUp(bb)
	down := vertical & ShiftDown(bb)
	upRight := sideCut & ShiftUpRight(bb)
	upLeft := sideCut & ShiftUpLeft(bb)
	downRight := sideCut & ShiftDownRight(bb)
	downLeft := sideCut & ShiftDownLeft(bb)

	right |= horizontal & ShiftRight(right)
	left |= horizontal & ShiftLeft(left)
	up |= vertical & ShiftUp(up)
	down |= vertical & ShiftDown(down)
	upRight |= sideCut & ShiftUpRight(upRight)
	upLeft |= sideCut & ShiftUpLeft(upLeft)
	downRight |= sideCut & ShiftDownRight(downRight)
	downLeft |= sideCut & ShiftDownLeft(downLeft)

	right |= horizontal & ShiftRight(right)
	left |= horizontal & ShiftLeft(left)
	up |= vertical & ShiftUp(up)
	down |= vertical & ShiftDown(down)
	upRight |= sideCut & ShiftUpRight(upRight)
	upLeft |= sideCut & ShiftUpLeft(upLeft)
	downRight |= sideCut & ShiftDownRight(downRight)
	downLeft |= sideCut & ShiftDownLeft(downLeft)

	right |= horizontal & ShiftRight(right)
	left |= horizontal & ShiftLeft(left)
	up |= vertical & ShiftUp(up)
	down |= vertical & ShiftDown(down)
	upRight |= sideCut & ShiftUpRight(upRight)
	upLeft |= sideCut & ShiftUpLeft(upLeft)
	downRight |= sideCut & ShiftDownRight(downRight)
	downLeft |= sideCut & ShiftDownLeft(downLeft)

	right |= horizontal & ShiftRight(right)
	left |= horizontal & ShiftLeft(left)
	up |= vertical & ShiftUp(up)
	down |= vertical & ShiftDown(down)
	upRight |= sideCut & ShiftUpRight(upRight)
	upLeft |= sideCut & ShiftUpLeft(upLeft)
	downRight |= sideCut & ShiftDownRight(downRight)
	downLeft |= sideCut & ShiftDownLeft(downLeft)

	right |= horizontal & ShiftRight(right)
	left |= horizontal & ShiftLeft(left)
	up |= vertical & ShiftUp(up)
	down |= vertical & ShiftDown(down)
	upRight |= sideCut & ShiftUpRight(upRight)
	upLeft |= sideCut & ShiftUpLeft(upLeft)
	downRight |= sideCut & ShiftDownRight(downRight)
	downLeft |= sideCut & ShiftDownLeft(downLeft)

	//1番右の相手の石より、更に1つ右側に移動し、その場所が空白であれば1を返す。
	legals := empty & ShiftRight(right)
	legals |= empty & ShiftLeft(left)
	legals |= empty & ShiftUp(up)
	legals |= empty & ShiftDown(down)
	legals |= empty & ShiftUpRight(upRight)
	legals |= empty & ShiftUpLeft(upLeft)
	legals |= empty & ShiftDownRight(downRight)
	legals |= empty & ShiftDownLeft(downLeft)
	return legals
}

func (bb BitBoard) Flips(opp, move BitBoard) BitBoard {
	occupied := bb | opp

	//既に石が置かれている場合
	if (occupied & move) != 0 {
		return 0
	}

	shifts := []func(BitBoard) BitBoard{
		ShiftRight, ShiftLeft, ShiftUp, ShiftDown,
		ShiftUpRight, ShiftUpLeft, ShiftDownRight, ShiftDownLeft,
	}

	masks := []BitBoard{
		/*
			右方向のマスク
			[[1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]]
		*/
		0b01111111_01111111_01111111_01111111_01111111_01111111_01111111_01111111,

		/*
			左方向のマスク
			[[0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]]
		*/
		0b11111110_11111110_11111110_11111110_11111110_11111110_11111110_11111110,

		/*
			上方向のマスク
			[[0 0 0 0 0 0 0 0]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]]
		*/
		0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_00000000,

		/*
			下方向のマスク
			[[1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [1 1 1 1 1 1 1 1]
			 [0 0 0 0 0 0 0 0]]
		*/
		0b00000000_11111111_11111111_11111111_11111111_11111111_11111111_11111111,

		/*
			右上方向のマスク
			[[0 0 0 0 0 0 0 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]]
		*/
		0b01111111_01111111_01111111_01111111_01111111_01111111_01111111_00000000,

		/*
			左上方向のマスク
			[[0 0 0 0 0 0 0 0]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]]
		*/
		0b11111110_11111110_11111110_11111110_11111110_11111110_11111110_00000000,

		/*
			右下方向のマスク
			[[1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [1 1 1 1 1 1 1 0]
			 [0 0 0 0 0 0 0 0]]
		*/
		0b00000000_01111111_01111111_01111111_01111111_01111111_01111111_01111111,

		/*
			左下方向のマスク
			[[0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 1 1 1 1 1 1 1]
			 [0 0 0 0 0 0 0 0]]
		*/
		0b00000000_11111110_11111110_11111110_11111110_11111110_11111110_11111110,
	}

	var flips BitBoard
	for i, shift := range shifts {
		mask := masks[i]

		//置く石の座標からスタートして、ある方向に1つ進む
		currentRaw := shift(move)
		currentMasked := currentRaw & mask
		var between BitBoard

		for currentMasked != 0 && (currentMasked&opp) != 0 {
			between |= currentMasked
			currentRaw = shift(currentRaw)
			currentMasked = currentRaw & mask
		}

		if (currentRaw & bb) != 0 {
			flips |= between
		}
	}
	return flips
}

func (bb BitBoard) ToArray() [Rows][Cols]int {
	var arr [Rows][Cols]int
	for i := 0; i < BoardSize; i++ {
		row, col := IndexToRowColumn(i)
		if bb&(1<<i) != 0 {
			arr[row][col] = 1
		} else {
			arr[row][col] = 0
		}
	}
	return arr
}

func ShiftRight(bb BitBoard) BitBoard {
	return bb << 1
}

func ShiftLeft(bb BitBoard) BitBoard {
	return bb >> 1
}

func ShiftUp(bb BitBoard) BitBoard {
	return bb >> 8
}

func ShiftDown(bb BitBoard) BitBoard {
	return bb << 8
}

func ShiftUpRight(bb BitBoard) BitBoard {
	return bb >> 7
}

func ShiftUpLeft(bb BitBoard) BitBoard {
	return bb >> 9
}

func ShiftDownRight(bb BitBoard) BitBoard {
	return bb << 9
}

func ShiftDownLeft(bb BitBoard) BitBoard {
	return bb << 7
}

var GroupBitBoards = []BitBoard{
	CornerBitBoard,
	CBitBoard,
	ABitBoard,
	BBitBoard,
	XBitBoard,
	0b00000000_00100100_01000010_00000000_00000000_01000010_00100100_00000000,
	0b00000000_00011000_00000000_01000010_01000010_00000000_00011000_00000000,
	0b00000000_00000000_00100100_00000000_00000000_00100100_00000000_00000000,
	0b00000000_00000000_00011000_00100100_00100100_00011000_00000000_00000000,
	0b00000000_00000000_00000000_00011000_00011000_00000000_00000000_00000000,
}

var Singles = omwbits.ToSingles64[BitBoard](math.MaxUint64)
