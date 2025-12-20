package gothello

import (
	"fmt"
	omwbits "github.com/sw965/omw/math/bits"
)

type State struct {
	Blacks BitBoard
	Whites BitBoard
	Turn   Turn
}

func NewInitState() State {
	blacks := BitBoard(0b00000000_00000000_00000000_00001000_00010000_00000000_00000000_00000000)
	whites := BitBoard(0b00000000_00000000_00000000_00010000_00001000_00000000_00000000_00000000)
	return State{Blacks: blacks, Whites: whites, Turn: BlackTurn}
}

func (s State) Empties() BitBoard {
	return ^(s.Blacks | s.Whites)
}

func (s State) Legals() BitBoard {
	if s.Turn == BlackTurn {
		return s.Blacks.Legals(s.Whites)
	} else {
		return s.Whites.Legals(s.Blacks)
	}
}

func (s State) FlipsByLegal() (map[BitBoard]BitBoard, error) {
	var self, opp BitBoard
	if s.Turn == BlackTurn {
		self = s.Blacks
		opp = s.Whites
	} else if s.Turn == WhiteTurn {
		self = s.Whites
		opp = s.Blacks
	} else {
		return nil, fmt.Errorf("State.Turn == Empty")
	}

	m := map[BitBoard]BitBoard{}
	singles := omwbits.ToSingles64(s.Legals())
	for _, legal := range singles {
		m[legal] = self.Flips(opp, legal)
	}
	return m, nil
}

func (s State) Move(move BitBoard) (State, error) {
	var self BitBoard
	var opp BitBoard

	if s.Turn == BlackTurn {
		self = s.Blacks
		opp = s.Whites
	} else if s.Turn == WhiteTurn {
		self = s.Whites
		opp = s.Blacks
	} else {
		return State{}, fmt.Errorf("State.Turn == Empty")
	}

	flips := self.Flips(opp, move)
	if flips == 0 {
		return State{}, fmt.Errorf("非合法の手を打とうとした。")
	}

	//石を置いて、ひっくり返す。
	self |= move | flips
	//ひっくり返される石を消す。
	opp &^= flips

	if s.Turn == BlackTurn {
		s.Blacks = self
		s.Whites = opp
	} else {
		s.Whites = self
		s.Blacks = opp
	}

	if opp.Legals(self) != 0 {
		s.Turn = s.Turn.Opposite()
	}
	return s, nil
}

func (s State) ToFeature() Feature {
	if s.Turn == BlackTurn {
		return Feature{Self: s.Blacks, Opponent: s.Whites}
	} else {
		return Feature{Self: s.Whites, Opponent: s.Blacks}
	}
}

func (s State) Transpose() State {
	s.Blacks = s.Blacks.Transpose()
	s.Whites = s.Whites.Transpose()
	return s
}

func (s State) MirrorHorizontal() State {
	s.Blacks = s.Blacks.MirrorHorizontal()
	s.Whites = s.Whites.MirrorHorizontal()
	return s
}

func (s State) MirrorVertical() State {
	s.Blacks = s.Blacks.MirrorVertical()
	s.Whites = s.Whites.MirrorVertical()
	return s
}

func (s State) Rotate90() State {
	s.Blacks = s.Blacks.Rotate90()
	s.Whites = s.Whites.Rotate90()
	return s
}

func (s State) Rotate180() State {
	s.Blacks = s.Blacks.Rotate180()
	s.Whites = s.Whites.Rotate180()
	return s
}

func (s State) Rotate270() State {
	s.Blacks = s.Blacks.Rotate270()
	s.Whites = s.Whites.Rotate270()
	return s
}

func (s State) ToArray() ([][]Color, error) {
	bArr := s.Blacks.ToArray()
	wArr := s.Whites.ToArray()

	arr := make([][]Color, Rows)
	for i := range arr {
		arr[i] = make([]Color, Cols)
	}

	for i := range arr {
		for j := range arr[i] {
			b := bArr[i][j]
			w := wArr[i][j]

			if b == 1 && w == 1 {
				return nil, fmt.Errorf("(%d, %d) 地点に黒石と白石の両方にビットが入力されている。", i, j)
			}

			if b == 1 {
				arr[i][j] = Black
			}

			if w == 1 {
				arr[i][j] = White
			}
		}
	}
	return arr, nil
}
