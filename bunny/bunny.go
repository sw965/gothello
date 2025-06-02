package bunny

import (
	"fmt"
	"github.com/sw965/gothello"
	"slices"
	"github.com/sw965/crow/model/linear"
	oslices "github.com/sw965/omw/slices"
)

func MakePatterns() {
	patterns := oslices.Sequence[[][]int, []int]([]int{0, 1, 2}, 8)
	for _, pattern := range patterns {
		fmt.Println(pattern)
	}
}

func NewInput(state gothello.State) linear.Input {
	legal := state.LegalBitBoard()
	legals := legal.ToSingles()
	handPairBitBoard := state.NewHandPairBitBoard()

	input := make(linear.Input, gothello.BoardSize)

	for i, legalIdx := range legal.OneIndices() {
		flip := handPairBitBoard.Self.FlipBitBoard(handPairBitBoard.Opponent, legals[i])
		flipIdxs := flip.OneIndices()
		for weightIdx, t := range gothello.GroupIndexTable {
			if slices.Contains(t, legalIdx) {
				input[legalIdx] = append(input[legalIdx], linear.Entry{
					X:float32(len(flipIdxs)),
					WeightIndex:weightIdx,
				})
			}
		}
	}

	return input
}