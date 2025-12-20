package gothello

import (
	"fmt"
	omwbits "github.com/sw965/omw/math/bits"
	omwslices "github.com/sw965/omw/slices"
	"slices"
)

type Feature struct {
	Self     BitBoard
	Opponent BitBoard
}

func NewFeatureFromIndices(selfIdxs, oppIdxs []int) (Feature, error) {
	if !omwslices.IsMutuallyExclusive(selfIdxs, oppIdxs) {
		return Feature{}, fmt.Errorf("self と opp は 排反でなければならない")
	}

	self, err := omwbits.New64FromIndices[BitBoard](selfIdxs)
	if err != nil {
		return Feature{}, err
	}

	opp, err := omwbits.New64FromIndices[BitBoard](oppIdxs)
	if err != nil {
		return Feature{}, err
	}

	return Feature{
		Self:     self,
		Opponent: opp,
	}, nil
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

func (f Feature) Move(move BitBoard) (Feature, error) {
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

func (f Feature) ToArray() ([][]Perspective, error) {
	selfArr := f.Self.ToArray()
	oppArr := f.Opponent.ToArray()
	arr := make([][]Perspective, Rows)
	for i := range arr {
		arr[i] = make([]Perspective, Cols)
	}
	for i := range arr {
		for j := range arr[i] {
			se := selfArr[i][j]
			oe := oppArr[i][j]
			if se == 1 && oe == 1 {
				return nil, fmt.Errorf("(%d, %d) 地点に自分の石と相手の石の両方にビットが入力されている。", i, j)
			}

			if se == 1 {
				arr[i][j] = Self
			}

			if oe == 1 {
				arr[i][j] = Opponent
			}
		}
	}
	return arr, nil
}

type PartialFeature1D []Perspective

func (f PartialFeature1D) Move(idx int) (PartialFeature1D, error) {
	fn := len(f)
	if fn <= idx {
		return nil, fmt.Errorf("len(f) <= idx")
	}

	if f[idx] != None {
		return nil, fmt.Errorf("空白ではない場所に置こうとした")
	}

	flipIdxs := make([]int, 0, fn-2)

	if idx > 1 {
		betweenIdxs := make([]int, 0, idx-2)

	LeftScan:
		for i := idx - 1; i > 0; i-- {
			switch f[i] {
			case None:
				betweenIdxs = nil
				break LeftScan
			case Self:
				break LeftScan
			case Opponent:
				betweenIdxs = append(betweenIdxs, i)
			}
		}

		for _, betweenIdx := range betweenIdxs {
			flipIdxs = append(flipIdxs, betweenIdx)
		}
	}

	if idx <= fn-2 {
		betweenIdxs := make([]int, 0, fn-idx-1)
	RightScan:
		for i := idx + 1; i < fn; i++ {
			switch f[i] {
			case None:
				betweenIdxs = nil
				break RightScan
			case Self:
				break RightScan
			case Opponent:
				betweenIdxs = append(betweenIdxs, i)
			}
		}

		for _, betweenIdx := range betweenIdxs {
			flipIdxs = append(flipIdxs, betweenIdx)
		}
	}

	f = slices.Clone(f)
	f[idx] = Self
	for _, flipIdx := range flipIdxs {
		f[flipIdx] = Self
	}
	return f, nil
}

func (f PartialFeature1D) CountLeading(p Perspective) int {
	c := 0
	for _, e := range f {
		if e != p {
			return c
		}
		c++
	}
	return c
}

func (f PartialFeature1D) ToFeature(idxs []int) (Feature, error) {
	if len(f) != len(idxs) {
		return Feature{}, fmt.Errorf("len(f) != len(idxs)")
	}

	if !omwslices.IsUnique(idxs) {
		return Feature{}, fmt.Errorf("idxs is not unique")
	}

	feature := Feature{}
	for i, idx := range idxs {
		var err error
		switch f[i] {
		case Self:
			feature.Self, err = omwbits.ToggleBit64(feature.Self, idx)
		case Opponent:
			feature.Opponent, err = omwbits.ToggleBit64(feature.Opponent, idx)
		}
		if err != nil {
			return Feature{}, err
		}
	}
	return feature, nil
}