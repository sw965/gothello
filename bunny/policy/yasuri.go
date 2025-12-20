package policy

import (
	"github.com/sw965/gothello"
	omwbits "github.com/sw965/omw/math/bits"
	omwslices "github.com/sw965/omw/slices"
)

type YasuriScorer map[int]map[gothello.BitBoard]map[gothello.Feature]int

func NewYasuriScorer() (YasuriScorer, error) {
	perspectiveSequences := omwslices.Sequences(gothello.AllPerspectives, gothello.Rows)
	sideIdxTable := [][]int{
		gothello.UpSideIndices,
		gothello.DownSideIndices,
		gothello.LeftSideIndices,
		gothello.RightSideIndices,
	}

	scorer := YasuriScorer{}
	for _, idxs := range sideIdxTable {
		for _, idx := range idxs {
			scorer[idx] = map[gothello.BitBoard]map[gothello.Feature]int{}
		}
	}

	for _, perspectives := range perspectiveSequences {
		partialFeature1D := perspectives.ToPartialFeature1D()
		for _, idxs := range sideIdxTable {
			feature, err := partialFeature1D.ToFeature(idxs)
			if err != nil {
				return nil, err
			}

			mask, err := omwbits.New64FromIndices[gothello.BitBoard](idxs)
			if err != nil {
				return nil, err
			}

			for i, idx := range idxs {
				if partialFeature1D[i] == gothello.None {
					moved, err := partialFeature1D.Move(i)
					if err != nil {
						return nil, err
					}

					var gain int
					if omwslices.Count(moved, gothello.Self) == gothello.Rows {
						leftStable := partialFeature1D.CountLeading(gothello.Self)
						rightStable := omwslices.Reversed(partialFeature1D).CountLeading(gothello.Self)
						gain = gothello.Rows - leftStable + rightStable
					} else {
						leftGain := moved.CountLeading(gothello.Self) - partialFeature1D.CountLeading(gothello.Self)
						rightGain := omwslices.Reversed(moved).CountLeading(gothello.Self) - omwslices.Reversed(partialFeature1D).CountLeading(gothello.Self)
						gain = leftGain + rightGain
					}

					if _, ok := scorer[idx][mask]; !ok {
						scorer[idx][mask] = map[gothello.Feature]int{}
					}
					scorer[idx][mask][feature] = gain
				}
			}
		}
	}
	return scorer, nil
}

func (ys YasuriScorer) ScoreByMoveIndex(f gothello.Feature) map[int]int {
	scores := map[int]int{}
	for moveIdx, scorerByIdx := range ys {
		for mask, gainByFeature := range scorerByIdx {
			for feature, gain := range gainByFeature {
				masked := f.AndBitBoard(mask)
				if masked == feature {
					if _, ok := scores[moveIdx]; ok {
						scores[moveIdx] -= 1
					} 
					scores[moveIdx] += gain
				}
			}
		}
	}
	return scores
}
