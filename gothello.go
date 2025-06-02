package gothello

import (
	"fmt"
	"math/bits"
	"slices"
)

const (
	Rows = 8
	Cols = 8
	BoardSize = Rows * Cols
)

type BitBoard uint64

const (
	UpEdgeBitBoard    = BitBoard(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_01111110)
	DownEdgeBitBoard  = BitBoard(0b01111110_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
	LeftEdgeBitBoard  = BitBoard(0b00000000_00000001_00000001_00000001_00000001_00000001_00000001_00000000)
	RightEdgeBitBoard = BitBoard(0b00000000_10000000_10000000_10000000_10000000_10000000_10000000_00000000)

	UpSideBitBoard    = BitBoard(0b00000000_00000000_00000000_00000000_00000000_00000000_00000000_11111111)
	DownSideBitBoard  = BitBoard(0b11111111_00000000_00000000_00000000_00000000_00000000_00000000_00000000)
	LeftSideBitBoard  = BitBoard(0b00000001_00000001_00000001_00000001_00000001_00000001_00000001_00000001)
	RightSideBitBoard = BitBoard(0b10000000_10000000_10000000_10000000_10000000_10000000_10000000_10000000)
	SideBitBoard      = BitBoard(0b11111111_10000001_10000001_10000001_10000001_10000001_10000001_11111111)

	CornerBitBoard = BitBoard(0b10000001_00000000_00000000_00000000_00000000_00000000_00000000_10000001)
	CBitBoard      = BitBoard(0b01000010_10000001_00000000_00000000_00000000_00000000_10000001_01000010)
	ABitBoard      = BitBoard(0b00100100_00000000_10000001_00000000_00000000_10000001_00000000_00100100)
	BBitBoard      = BitBoard(0b00011000_00000000_00000000_10000001_10000001_00000000_00000000_00011000)
	XBitBoard      = BitBoard(0b00000000_01000010_00000000_00000000_00000000_00000000_01000010_00000000)

	WhiteLineBitBoard = BitBoard(0b10000000_01000000_00100000_00010000_00001000_00000100_00000010_00000001)
	BlackLineBitBoard = BitBoard(0b00000001_00000010_00000100_00001000_00010000_00100000_01000000_10000000)
)

var AdjacentBySingle = func() map[BitBoard]BitBoard {
	m := map[BitBoard]BitBoard{}
	for i := 0; i < BoardSize; i++ {
		single := BitBoard(1) << i
		var adj BitBoard = 0
		row, col := IndexToRowAndColumn(i)

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

//転置行列と同じ
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

//横回転
func (bb BitBoard) MirrorHorizontal() BitBoard {
	bb = ((bb & 0b11110000_11110000_11110000_11110000_11110000_11110000_11110000_11110000) >> 4) |
		((bb & 0b00001111_00001111_00001111_00001111_00001111_00001111_00001111_00001111) << 4)
	bb = ((bb & 0b11001100_11001100_11001100_11001100_11001100_11001100_11001100_11001100) >> 2) |
		((bb & 0b00110011_00110011_00110011_00110011_00110011_00110011_00110011_00110011) << 2)
	bb = ((bb & 0b10101010_10101010_10101010_10101010_10101010_10101010_10101010_10101010) >> 1) |
		((bb & 0b01010101_01010101_01010101_01010101_01010101_01010101_01010101_01010101) << 1)
	return bb
}

//縦回転
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

func (bb BitBoard) Count() int {
	return bits.OnesCount64(uint64(bb))
}

func (bb BitBoard) OneIndices() []int {
	b := uint64(bb)
	idxs := make([]int, 0, bits.OnesCount64(b))
	for b != 0 {
		// 最下位の1ビットの位置を求める
		idx := bits.TrailingZeros64(b)
		idxs = append(idxs, idx)
		// 最下位の1ビットをクリア
		b &= b - 1
	}
	return idxs
}

func (bb BitBoard) ToSingles() []BitBoard {
	count := bb.Count()
	singles := make([]BitBoard, 0, count)
	for bb != 0 {
		// 最下位の1ビットを抽出
		lsb := bb & -bb
		singles = append(singles, lsb)
		// 抽出したビットをクリア
		bb ^= lsb
	}
	return singles
}

/*
	https://blog.qmainconts.dev/articles/yxiplk2_dd
	https://qiita.com/sensuikan1973/items/459b3e11d91f3cb37e43
	上記の参考文献では、最上位ビットを(0, 0)の地点と見なしているが、
	このライブラリでは、(0, 0)の地点を最下位ビットと見なしている為、
	左シフトと右シフトの役割が逆になっている。
*/
func (bb BitBoard) LegalBitBoard(opponent BitBoard) BitBoard {
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

func (bb BitBoard) FlipBitBoard(opponent, move BitBoard) BitBoard {
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

func (bb BitBoard) ToArray() [Rows][Cols]int {
	var arr [Rows][Cols]int
	for i := 0; i < BoardSize; i++ {
		row, col := IndexToRowAndColumn(i)
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

type State struct {
	Black BitBoard
	White BitBoard
	Hand  Color
}

func NewInitState() State {
	black := BitBoard(0b00000000_00000000_00000000_00010000_00001000_00000000_00000000_00000000)
	white := BitBoard(0b00000000_00000000_00000000_00001000_00010000_00000000_00000000_00000000)
	return State{Black:black, White:white, Hand:Black}
}

func (s State) SpaceBitBoard() BitBoard {
	return ^(s.Black | s.White)
}

func (s State) LegalBitBoard() BitBoard {
	if s.Hand == Black {
		return s.Black.LegalBitBoard(s.White)
	} else {
		return s.White.LegalBitBoard(s.Black)
	}
}

func (s State) NewHandPairBitBoard() HandPairBitBoard {
	if s.Hand == Black {
		return HandPairBitBoard{Self:s.Black, Opponent:s.White}
	} else {
		return HandPairBitBoard{Self:s.White, Opponent:s.Black}
	}
}

func (s State) Put(move BitBoard) (State, error) {
	var self BitBoard
	var opponent BitBoard

	if s.Hand == Black {
		self = s.Black
		opponent = s.White
	} else {
		self = s.White
		opponent = s.Black
	}

	flips := self.FlipBitBoard(opponent, move)
	if flips == 0 {
		return State{}, fmt.Errorf("非合法の手を打とうとした。")
	}

	//石を置いて、ひっくり返す。
	self |= move | flips
	//ひっくり返される石を消す。
	opponent &^= flips

	if s.Hand == Black {
		s.Black = self
		s.White = opponent
	} else {
		s.White = self
		s.Black = opponent
	}

	if opponent.LegalBitBoard(self) != 0 {
		s.Hand = s.Hand.Opposite()
	}
	return s, nil
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

func (s State) ToArray() ([][]int, error) {
	bArr := s.Black.ToArray()
	wArr := s.White.ToArray()

	arr := make([][]int, Rows)
	for i := range arr {
		arr[i] = make([]int, Cols)
	}

	for i := range arr {
		for j := range arr[i] {
			be := bArr[i][j]
			we := wArr[i][j]
			if be == 1 && we == 1 {
				return nil, fmt.Errorf("(%d, %d) 地点に黒石と白石の両方にビットが入力されている。", i, j)
			}

			if be == 1 {
				arr[i][j] = 1
			}
			
			if we == 1 {
				arr[i][j] = 2
			}
		}
	}
	return arr, nil
}

type States []State

func (ss States) CountResult() (ResultCounter, error) {
	counter := ResultCounter{}
	for _, s := range ss {
		blackLegal := s.Black.LegalBitBoard(s.White)
		whiteLegal := s.White.LegalBitBoard(s.Black)

		if blackLegal != 0 || whiteLegal != 0 {
			return ResultCounter{}, fmt.Errorf("ゲームが終了していない為、カウント出来ない。")
		}

		blackCount := s.Black.Count()
		whiteCount := s.White.Count()
		if blackCount > whiteCount {
			counter.BlackWin++
		} else if blackCount < whiteCount {
			counter.WhiteWin++
		} else {
			counter.Draw++
		}
	}
	return counter, nil
}

type ColorPairBitBoard struct {
	Black BitBoard
	White BitBoard
}

func (cp ColorPairBitBoard) SpaceBitBoard() BitBoard {
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

type HandPairBitBoard struct {
	Self     BitBoard
	Opponent BitBoard
}

func (hp HandPairBitBoard) SpaceBitBoard() BitBoard {
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

func (hp HandPairBitBoard) ToArray() ([][]int, error) {
	selfArr := hp.Self.ToArray()
	oppArr := hp.Opponent.ToArray()
	arr := make([][]int, Rows)
	for i := range arr {
		arr[i] = make([]int, Cols)
	}
	for i := range arr {
		for j := range arr[i] {
			se := selfArr[i][j]
			oe := oppArr[i][j]
			if se == 1 && oe == 1 {
				return nil, fmt.Errorf("(%d, %d) 地点に自分の石と相手の石の両方にビットが入力されている。", i, j)
			}

			if se == 1 {
				arr[i][j] = 1
			}
			
			if oe == 1 {
				arr[i][j] = 2
			}
		}
	}
	return arr, nil
}

type Cell struct {
	Row int
	Column string
}

func (c *Cell) ToBitBoard() BitBoard {
	cs := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	row := c.Row-1
	col := slices.Index(cs, c.Column)
	idx := RowAndColumnToIndex(row, col)
	return 1 << idx
}

type ResultCounter struct {
	BlackWin int
	WhiteWin int
	Draw     int
}

func (r ResultCounter) Total() int {
	return r.BlackWin + r.WhiteWin + r.Draw
}

func (r ResultCounter) TotalWin() int {
	return r.BlackWin + r.WhiteWin
}

func (r ResultCounter) BlackWinRate(includeDraws bool) float64 {
	w := float64(r.BlackWin)
	if includeDraws {
		n := float64(r.Total())
		d := float64(r.Draw)
		return (w + (0.5 * d)) / n
	}
	n := float64(r.TotalWin())
	return w / float64(n)
}

func (r ResultCounter) WhiteWinRate(includeDraws bool) float64 {
	w := float64(r.WhiteWin)
	if includeDraws {
		n := float64(r.Total())
		d := float64(r.Draw)
		return (w + (0.5 * d)) / n
	}

	n := float64(r.TotalWin())
	return w / float64(n)
}

func (r ResultCounter) DrawRate() float64 {
	n := float64(r.Total())
	return float64(r.Draw) / n
}