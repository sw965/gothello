package gothello

import (
	"math/bits"
)

type BitBoard uint64

func (bb BitBoard) OnesCount() BitBoard {
	return BitBoard(bits.OnesCount64(uint64(bb)))
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
