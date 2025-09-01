package gothello

import (
	"fmt"
	omwbits "github.com/sw965/omw/math/bits"
)

type Feature struct {
	Self     BitBoard
	Opponent BitBoard
}

func (f Feature) AndBitBoard(bb BitBoard) Feature {
	f.Self &= bb
	f.Opponent &= bb
	return f
}

func (f Feature) Empties() BitBoard {
	return ^(f.Self | f.Opponent)
}

func (f Feature) Legals() BitBoard {
	return f.Self.Legals(f.Opponent)
}

func (f Feature) FlipsByLegal() map[BitBoard]BitBoard {
	m := map[BitBoard]BitBoard{}
	singleLegals := omwbits.ToSingles64(f.Legals())
	for _, legal := range singleLegals {
		m[legal] = f.Self.Flips(f.Opponent, legal)
	}
	return m
}

func (f Feature) Put(move BitBoard) (Feature, error) {
	flips := f.Self.Flips(f.Opponent, move)
	if flips == 0 {
		return Feature{}, fmt.Errorf("非合法な手を打とうとした。")
	}
	f.Self |= move | flips
	f.Opponent &^= flips
	return f, nil
}

func (f Feature) Transpose() Feature {
	f.Self = f.Self.Transpose()
	f.Opponent = f.Opponent.Transpose()
	return f
}

func (f Feature) MirrorHorizontal() Feature {
	f.Self = f.Self.MirrorHorizontal()
	f.Opponent = f.Opponent.MirrorHorizontal()
	return f
}

func (f Feature) MirrorVertical() Feature {
	f.Self = f.Self.MirrorVertical()
	f.Opponent = f.Opponent.MirrorVertical()
	return f
}

func (f Feature) Rotate90() Feature {
	f.Self = f.Self.Rotate90()
	f.Opponent = f.Opponent.Rotate90()
	return f
}

func (f Feature) Rotate180() Feature {
	f.Self = f.Self.Rotate180()
	f.Opponent = f.Opponent.Rotate180()
	return f
}

func (f Feature) Rotate270() Feature {
	f.Self = f.Self.Rotate270()
	f.Opponent = f.Opponent.Rotate270()
	return f
}

func (f Feature) ToArray() ([][]Disc, error) {
	selfArr := f.Self.ToArray()
	oppArr := f.Opponent.ToArray()
	arr := make([][]Disc, Rows)
	for i := range arr {
		arr[i] = make([]Disc, Cols)
	}
	for i := range arr {
		for j := range arr[i] {
			se := selfArr[i][j]
			oe := oppArr[i][j]
			if se == 1 && oe == 1 {
				return nil, fmt.Errorf("(%d, %d) 地点に自分の石と相手の石の両方にビットが入力されている。", i, j)
			}

			if se == 1 {
				arr[i][j] = Black
			}
			
			if oe == 1 {
				arr[i][j] = White
			}
		}
	}
	return arr, nil
}