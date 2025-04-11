package gothello

import (
	"math/bits"
	omwslices "github.com/sw965/omw/slices"
)

const (
	ROW = 8
	COLUMN = 8
	SIDE_SIZE = 8
	FLAT_SIZE = ROW * COLUMN
)

type BitBoard uint64

const (
	CORNER_POINT_BIT_BOARD = BitBoard(0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_10000001)
	C_POINT_BIT_BOARD = BitBoard(0b01000010_10000001_00000000_00000000_00000000_00000000_10000001_01000010)
	A_POINT_BIT_BOARD = BitBoard(0b00100100_00000000_10000001_00000000_00000000_10000001_00000000_00100100)
	B_POINT_BIT_BOARD = BitBoard(0b00011000_00000000_00000000_10000001_10000001_00000000_00000000_00011000)
	X_POINT_BIT_BOARD = BitBoard(0b00000000_01000010_00000000_00000000_00000000_00000000_01000010_00000000)
)

func (bb BitBoard) ToggleBit(p *Point) BitBoard {
	return bb ^ (1 << (p.Row * COLUMN + p.Column))
}

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

func (bb BitBoard) MirrorHorizontal() BitBoard {
	bb = ((bb & 0b11110000_11110000_11110000_11110000_11110000_11110000_11110000_11110000) >> 4) |
		((bb & 0b00001111_00001111_00001111_00001111_00001111_00001111_00001111_00001111) << 4)
	bb = ((bb & 0b11001100_11001100_11001100_11001100_11001100_11001100_11001100_11001100) >> 2) |
		((bb & 0b00110011_00110011_00110011_00110011_00110011_00110011_00110011_00110011) << 2)
	bb = ((bb & 0b10101010_10101010_10101010_10101010_10101010_10101010_10101010_10101010) >> 1) |
		((bb & 0b01010101_01010101_01010101_01010101_01010101_01010101_01010101_01010101) << 1)
	return bb
}

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

func (bb BitBoard) ToOneHots() BitBoards {
	count := bits.OnesCount64(uint64(bb))
	oneHots := make(BitBoards, 0, count)
	for bb != 0 {
		// 最下位の1ビットを抽出
		lsb := bb & -bb
		oneHots = append(oneHots, lsb)
		// 抽出したビットをクリア
		bb ^= lsb
	}
	return oneHots
}

func (bb BitBoard) ToPoints() Points {
	oneHots := bb.ToOneHots()
	points := make(Points, len(oneHots))
	for i, oneHot := range oneHots {
		points[i] = POINT_BY_BIT_BOARD[oneHot]
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
func (bb BitBoard) LegalPointBitBoard(opponent BitBoard) BitBoard {
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

	space := ^(bb | opponent)

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

func (bb BitBoard) FlipPointBitBoard(opponent, move BitBoard) BitBoard {
	occupied := bb | opponent

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

		//置く石の座標からスタートして、ある方向に1つ進む
		currentRaw := shift(move)
		currentMasked := currentRaw & mask
		var candidate BitBoard
		
		for currentMasked != 0 && (currentMasked & opponent) != 0 {
			candidate |= currentMasked
			currentRaw = shift(currentRaw)
			currentMasked = currentRaw & mask
		}

		if (currentRaw & bb) != 0 {
			flips |= candidate
		}
	}
	return flips
}

func (bb BitBoard) ToArray() [ROW][COLUMN]int {
	var arr [ROW][COLUMN]int
	for i, p := range ALL_POINTS {
		r, c := p.Row, p.Column
		if bb&(1<<i) != 0 {
			arr[r][c] = 1
		} else {
			arr[r][c] = 0
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

type BitBoards []BitBoard

var ONE_HOT_BIT_BOARDS = BitBoard(0b11111111_11111111_11111111_11111111_11111111_11111111_11111111_11111111).ToOneHots()

var GROUP_BIT_BOARDS = BitBoards{
	CORNER_POINT_BIT_BOARD,
	C_POINT_BIT_BOARD,
	A_POINT_BIT_BOARD,
	B_POINT_BIT_BOARD,
	X_POINT_BIT_BOARD,
	0b00000000_00100100_01000010_00000000_00000000_01000010_00100100_00000000,
	0b00000000_00011000_00000000_01000010_01000010_00000000_00011000_00000000,
	0b00000000_00000000_00100100_00000000_00000000_00100100_00000000_00000000,
	0b00000000_00000000_00011000_00100100_00100100_00011000_00000000_00000000,
	0b00000000_00000000_00000000_00011000_00011000_00000000_00000000_00000000,
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

func (s State) Put(move BitBoard) State {
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
	self |= move | flips
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

func (s *State) ToString() string {
	blackArr := s.Black.ToArray()
	whiteArr := s.White.ToArray()
	legalArr := s.LegalPointBitBoard().ToArray()
	str := ""

	for row := 0; row < ROW; row++ {
		for col := 0; col < COLUMN; col++ {
			var mark string
			if blackArr[row][col] == 1 {
				mark = "b"
			} else if whiteArr[row][col] == 1 {
				mark = "w"
			} else if legalArr[row][col] == 1 {
				mark = "p"
			} else {
				mark = "-"
			}

			if col != COLUMN-1 {
				mark += " "
			}
			str += mark
		}
		str += "\n"
	}

	if s.Hand == BLACK {
		str += "hand = black\n"
	} else {
		str += "hand = white\n"
	}
	return str
}

func (s State) MirrorHorizontal() State {
	s.Black = s.Black.MirrorHorizontal()
	s.White = s.White.MirrorHorizontal()
	return s
}

func (s State) MirrorVertical() State {
	s.Black = s.Black.MirrorVertical()
	s.White = s.White.MirrorVertical()
	return s
}

func (s State) Rotate90() State {
	s.Black = s.Black.Rotate90()
	s.White = s.White.Rotate90()
	return s
}

func (s State) Rotate180() State {
	s.Black = s.Black.Rotate180()
	s.White = s.White.Rotate180()
	return s
}

func (s State) Rotate270() State {
	s.Black = s.Black.Rotate270()
	s.White = s.White.Rotate270()
	return s
}

type ColorPairBitBoard struct {
	Black BitBoard
	White BitBoard
}

func (cp *ColorPairBitBoard) SpaceBitBoard() BitBoard {
	return ^(cp.Black | cp.White)
}

func (cp ColorPairBitBoard) MirrorHorizontal() ColorPairBitBoard {
	cp.Black = cp.Black.MirrorHorizontal()
	cp.White = cp.White.MirrorHorizontal()
	return cp
}

func (cp ColorPairBitBoard) MirrorVertical() ColorPairBitBoard {
	cp.Black = cp.Black.MirrorVertical()
	cp.White = cp.White.MirrorVertical()
	return cp
}

func (cp ColorPairBitBoard) Rotate90() ColorPairBitBoard {
	cp.Black = cp.Black.Rotate90()
	cp.White = cp.White.Rotate90()
	return cp
}

func (cp ColorPairBitBoard) Rotate180() ColorPairBitBoard {
	cp.Black = cp.Black.Rotate180()
	cp.White = cp.White.Rotate180()
	return cp
}

func (cp ColorPairBitBoard) Rotate270() ColorPairBitBoard {
	cp.Black = cp.Black.Rotate270()
	cp.White = cp.White.Rotate270()
	return cp
}

func (cp *ColorPairBitBoard) ToString() string {
	blackArr := cp.Black.ToArray()
	whiteArr := cp.White.ToArray()
	str := ""

	for row := 0; row < ROW; row++ {
		for col := 0; col < COLUMN; col++ {
			var mark string
			if blackArr[row][col] == 1 {
				mark = "b"
			} else if whiteArr[row][col] == 1 {
				mark = "w"
			} else {
				mark = "-"
			}
			if col != COLUMN-1 {
				mark += " "
			}
			str += mark
		}
		str += "\n"
	}
	return str
}

type HandPairBitBoard struct {
	Self     BitBoard
	Opponent BitBoard
}

func (hp *HandPairBitBoard) SpaceBitBoard() BitBoard {
	return ^(hp.Self | hp.Opponent)
}

func (hp HandPairBitBoard) MirrorHorizontal() HandPairBitBoard {
	hp.Self = hp.Self.MirrorHorizontal()
	hp.Opponent = hp.Opponent.MirrorHorizontal()
	return hp
}

func (hp HandPairBitBoard) MirrorVertical() HandPairBitBoard {
	hp.Self = hp.Self.MirrorVertical()
	hp.Opponent = hp.Opponent.MirrorVertical()
	return hp
}

func (hp HandPairBitBoard) Rotate90() HandPairBitBoard {
	hp.Self = hp.Self.Rotate90()
	hp.Opponent = hp.Opponent.Rotate90()
	return hp
}

func (hp HandPairBitBoard) Rotate180() HandPairBitBoard {
	hp.Self = hp.Self.Rotate180()
	hp.Opponent = hp.Opponent.Rotate180()
	return hp
}

func (hp HandPairBitBoard) Rotate270() HandPairBitBoard {
	hp.Self = hp.Self.Rotate270()
	hp.Opponent = hp.Opponent.Rotate270()
	return hp
}

func (hp *HandPairBitBoard) ToArray() [][]int {
	selfArr := hp.Self.ToArray()
	oppArr := hp.Opponent.ToArray()
	arr := make([][]int, ROW)
	for i := range arr {
		arr[i] = make([]int, COLUMN)
	}
	for i := range arr {
		for j := range arr[i] {
			if selfArr[i][j] == 1 {
				arr[i][j] = 1
			} else if oppArr[i][j] == 1 {
				arr[i][j] = 2
			}
		}
	}
	return arr
}

func (hp *HandPairBitBoard) ToString() string {
	selfArr := hp.Self.ToArray()
	oppArr := hp.Opponent.ToArray()
	str := ""

	for row := 0; row < ROW; row++ {
		for col := 0; col < COLUMN; col++ {
			var mark string
			if selfArr[row][col] == 1 {
				mark = "s"
			} else if oppArr[row][col] == 1 {
				mark = "o"
			} else {
				mark = "-"
			}
			if col != COLUMN-1 {
				mark += " "
			}
			str += mark
		}
		str += "\n"
	}
	return str
}

type Point struct {
	Row int
	Column int
}

func (p *Point) ToIndex() int {
	return p.Row * COLUMN + p.Column
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

var (
	RIGHT_FLOW_UP_SIDE_POINTS = func() Points {
		points := make(Points, COLUMN)
			for col := 0; col < COLUMN; col++ {
			points[col] = Point{Row:0, Column:col}
		}
		return points
	}()

	LEFT_FLOW_UP_SIDE_POINTS = omwslices.Reverse(RIGHT_FLOW_UP_SIDE_POINTS)

	RIGHT_FLOW_DOWN_SIDE_POINTS = func() Points {
		points := make(Points, COLUMN)
		for col := 0; col < COLUMN; col++ {
			points[col] = Point{Row:ROW-1, Column:col}
		}
		return points
	}()

	LEFT_FLOW_DOWN_SIDE_POINTS = omwslices.Reverse(RIGHT_FLOW_DOWN_SIDE_POINTS)

	DOWN_FLOW_LEFT_SIDE_POINTS = func() Points {
		points := make(Points, ROW)
		for row := 0; row < ROW; row++ {
			points[row] = Point{Row:row, Column:0}
		}
		return points
	}()

	UP_FLOW_LEFT_SIDE_POINTS = omwslices.Reverse(DOWN_FLOW_LEFT_SIDE_POINTS)

	DOWN_FLOW_RIGHT_SIDE_POINTS = func() Points {
		points := make(Points, ROW)
		for row := 0; row < ROW; row++ {
			points[row] = Point{Row:row, Column:COLUMN-1}
		}
		return points
	}()

	UP_FLOW_RIGHT_SIDE_POINTS = omwslices.Reverse(DOWN_FLOW_RIGHT_SIDE_POINTS)
)

var POINT_BY_BIT_BOARD = func() map[BitBoard]Point {
	m := map[BitBoard]Point{}
	for i, point := range ALL_POINTS {
		bb := BitBoard(1)
		bb = bb << (i)
		m[bb] = point
	}
	return m
}()

var BIT_BOARD_BY_POINT = func() map[Point]BitBoard {
	m := map[Point]BitBoard{}
	for k, v := range POINT_BY_BIT_BOARD {
		m[v] = k
	}
	return m
}()