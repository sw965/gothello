package gothello

import (
	"math/bits"
)

const (
	ROW = 8
	COLUMN = 8
	FLAT_SIZE = ROW * COLUMN
)

type Point struct {
	Row int
	Column int
}

type Points []Point

var ALL_POINTS = func() Points {
	points := make(Points, 0, FLAT_SIZE)
	for row := 0; row < ROW; row++ {
		for col := 0; col < COLUMN; col++ {
			points = append(points, Point{Row:row, Column:col})
		}
	}
	return points
}()

type BitBoard uint64

func (b BitBoard) ToggleBit(p *Point) BitBoard {
	return b ^ (1 << (p.Row * COLUMN + p.Column))
}

func (b BitBoard) ToPoints() Points {
	count := bits.OnesCount64(uint64(b))
	points := make(Points, 0, count)

	for b != 0 {
		// 最下位の1ビットまでを抽出 (例 10101100 ならば 100を抽出)
		lsb := b & -b
		// lsb 1のインデックスを取得 (上記の例であれば、100なので、2を取得する)
		idx := bits.TrailingZeros64(uint64(lsb))
		points = append(points, ALL_POINTS[idx])		
		// 抽出済みのビットをクリア(上記の例であれば、10101000 になる)
		b ^= lsb
	}
	return points
}

/*
	https://blog.qmainconts.dev/articles/yxiplk2_dd
	https://qiita.com/sensuikan1973/items/459b3e11d91f3cb37e43
	上記の参考文献では、最上位ビットを(0, 0)の地点と見なしているが、
	このライブラリでは、(0, 0)の地点を最下位ビットと見なしている。
	なので、左シフトと右シフトの役割が逆になっている。
*/
func (b BitBoard) LegalPointBitBoard(opponent BitBoard) BitBoard {
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
	horizontal := opponent & 0b01111110_01111110_01111110_01111110_01111110_01111110_01111110_01111110

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
	vertical := opponent & 0b00000000_11111111_11111111_11111111_11111111_11111111_00000000

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
	sideCut := opponent & 0b00000000_01111110_01111110_01111110_01111110_01111110_01111110_00000000

	space := ^(b | opponent)

	/*
		右方向への探索に関して説明。他の方向も考え方は同じ。
		1.(b << 1) で左シフトすわなち、自分の石を、全て右方向に移動する(右方向移動した後のボードをAとする)。
		2.Aとhorizontalを用いて、相手の石と重なっている場所をAND演算求める(重なっている場所を1とする)(このボードをrightとする)。
		3.rightは、この時点で、自分の石の1つ右側にある相手の石座標を保持している。
		4.rightを左シフトすわなち、再度右方向に移動させて、horizontalとAND演算をする事で、自分の石の1つ右側にある相手の石の座標と、更に1つ右側にある相手の石の座標を求める。
		5. right |= とする事で、自分の石の1つ右側にある相手の石座標と、更に1つ右側にある相手の石の座標を保持する。
		6. 後は4と5を繰り返す事で、自分の石から右側にあるn個連結した相手の石の座標を求める事が出来る。
	*/
	right := horizontal & ShiftRight(b)
	left := horizontal & ShiftLeft(b)
	up := vertical & ShiftUp(b)
	down := vertical & ShiftDown(b)
	upRight := sideCut & ShiftUpRight(b)
	upLeft := sideCut & ShiftUpLeft(b)
	downRight := sideCut & ShiftDownRight(b)
	downLeft := sideCut & ShiftDownLeft(b)

	for i := 0; i < 5; i++ {
		right |= horizontal & ShiftRight(right)
		left |= horizontal & ShiftLeft(left)
		up |= vertical & ShiftUp(up)
		down |= vertical & ShiftDown(down)
		upRight |= sideCut & ShiftUpRight(upRight)
		upLeft |= sideCut & ShiftUpLeft(upLeft)
		downRight |= sideCut & ShiftDownRight(downRight)
		downLeft |= sideCut & ShiftDownLeft(downLeft)
	}

	//最後に、1番右の相手の石より、更に1つ右側に移動し、その場所が空白であれば1を返す。
	legal := space & ShiftRight(right)
	legal |= space & ShiftLeft(left)
	legal |= space & ShiftUp(up)
	legal |= space & ShiftDown(down)
	legal |= space & ShiftUpRight(upRight)
	legal |= space & ShiftUpLeft(upLeft)
	legal |= space & ShiftDownRight(downRight)
	legal |= space & ShiftDownLeft(downLeft)
	return legal
}

func (b BitBoard) FlipPointBitBoard(opponent BitBoard, movePoint *Point) BitBoard {
	move := BitBoard(0).ToggleBit(movePoint)
	occupied := b | opponent

	//既に石が置かれている場合
	if (occupied&move) != 0 {
		return 0
	}

	shifts := []func(BitBoard) BitBoard {
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
		var candidate BitBoard

		//置く石の座標からスタートして、ある方向に1つ進む
		current := shift(move) & mask

		/*
			マスクの範囲外に飛び出た場合、current == 0 になる。
			currentは現在探索している座標のみに1が立っているので、current&opponentでその座標に相手の石があるかを調べる。
			1ならば、currentの座標に相手の石が存在する。0ならば相手の石は存在しない。
			よって、このfor文は、currentがマスクの範囲内かつ相手の石が存在する場合に、ループする。
			1回ループする度に、currentを1つ先に進める。
		*/

		for current != 0 && (current&opponent) != 0 {
			candidate |= current
			current = shift(current) & mask
		}

		/*
			上記のfor文で、マスクの範囲外に飛び出た場合は、current == 0であり、ひっくり返す事は出来ない。
			相手の石が壁にぶつかる(範囲外に飛び出る)まで続いており、自分の石がおけない状態。

			current&opponent == 0 で止まった場合は、進んだ先に、相手の石が見つからなかった状態。
			もしも、進んだ先が自分の石ならば、candidateで蓄積した座標をひっくり返す事が出来る。

			よって(current & b) != 0 とする事で、ひっくり返せるかどうかをチェック出来る。
		*/
		if (current & b) != 0 {
			flips |= candidate
		}
	}
	return flips
}

func (b BitBoard) ToArray() [ROW][COLUMN]int {
	var arr [ROW][COLUMN]int
	for i, p := range ALL_POINTS {
		r, c := p.Row, p.Column
		if b&(1<<i) != 0 {
			arr[r][c] = 1
		} else {
			arr[r][c] = 0
		}
	}
	return arr
}

func ShiftRight(b BitBoard) BitBoard {
	return b << 1
}

func ShiftLeft(b BitBoard) BitBoard {
	return b >> 1
}

func ShiftUp(b BitBoard) BitBoard {
	return b >> 8
}

func ShiftDown(b BitBoard) BitBoard {
	return b << 8
}

func ShiftUpRight(b BitBoard) BitBoard {
	return b >> 7
}

func ShiftUpLeft(b BitBoard) BitBoard {
	return b >> 9
}

func ShiftDownRight(b BitBoard) BitBoard {
	return b << 9
}

func ShiftDownLeft(b BitBoard) BitBoard {
	return b << 7
}

const (
	BLACK = 0
	WHITE = 1
)

type State struct {
	Black BitBoard
	White BitBoard
	//0なら黒, 1なら白
	Hand int
}

func NewInitState() State {
	black := BitBoard(0b00000000_00000000_00000000_00001000_00010000_00000000_00000000_00000000)
	white := BitBoard(0b00000000_00000000_00000000_00010000_00001000_00000000_00000000_00000000)
	return State{Black:black, White:white, Hand:0}
}

func (s *State) SpaceBitBoard() BitBoard {
	return ^(s.Black | s.White)
}

func (s *State) LegalPointBitBoard() BitBoard {
	if s.Hand == BLACK {
		return s.Black.LegalPointBitBoard(s.White)
	} else {
		return s.White.LegalPointBitBoard(s.Black)
	}
}

func (s *State) NewHandPairBitBoard() HandPairBitBoard {
	if s.Hand == BLACK {
		return HandPairBitBoard{Self:s.Black, Opponent:s.White}
	} else {
		return HandPairBitBoard{Self:s.White, Opponent:s.Black}
	}
}

func (s State) Put(move *Point) State {
	var self BitBoard
	var opponent BitBoard

	if s.Hand == BLACK {
		self = s.Black
		opponent = s.White
	} else {
		self = s.White
		opponent = s.Black
	}

	flips := self.FlipPointBitBoard(opponent, move)
	//石を置いて、ひっくり返す。
	self |= self.ToggleBit(move) | flips
	//ひっくり返される石を消す。
	opponent &^= flips

	if s.Hand == BLACK {
		s.Black = self
		s.White = opponent
	} else {
		s.White = self
		s.Black = opponent
	}

	if opponent.LegalPointBitBoard(self) != 0 {
		s.Hand = map[int]int{BLACK:WHITE, WHITE:BLACK}[s.Hand]
	}
	return s
}

type ColorPairBitBoard struct {
	Black BitBoard
	White BitBoard
}

type HandPairBitBoard struct {
	Self     BitBoard
	Opponent BitBoard
}