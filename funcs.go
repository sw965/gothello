package gothello

import (
	"github.com/sw965/omw/funcs"
)

var BitBoardSymmetryFuncs = []func(BitBoard) BitBoard{
	funcs.Identity[BitBoard],
	BitBoard.Rotate90,
	BitBoard.Rotate180,
	BitBoard.Rotate270,
	BitBoard.MirrorHorizontal,
	BitBoard.MirrorVertical,
	BitBoard.Transpose,
	func(bb BitBoard) BitBoard {
		t := bb.Transpose()
		return t.Rotate180()
	},
}

var IndexSymmetryFuncs = []func(int) int{
	funcs.Identity[int],
	Rotate90Index,
	Rotate180Index,
	Rotate270Index,
	MirrorHorizontalIndex,
	MirrorVerticalIndex,
	TransposeIndex,
	func(idx int) int {
		t := TransposeIndex(idx)
		return Rotate180Index(t)
	},
}

var StateSymmetryFuncs = []func(State) State{
	funcs.Identity[State],
	State.Rotate90,
	State.Rotate180,
	State.Rotate270,
	State.MirrorHorizontal,
	State.MirrorVertical,
	State.Transpose,
	func(state State) State {
		t := state.Transpose()
		return t.Rotate180()
	},
}

var FeatureSymmetryFuncs = []func(Feature) Feature{
	funcs.Identity[Feature],
	Feature.Rotate90,
	Feature.Rotate180,
	Feature.Rotate270,
	Feature.MirrorHorizontal,
	Feature.MirrorVertical,
	Feature.Transpose,
	func(f Feature) Feature {
		t := f.Transpose()
		return t.Rotate180()
	},
}
