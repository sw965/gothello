package policy

import (
	"fmt"
	//"github.com/sw965/crow/model/linear"
	"github.com/sw965/gothello"
	"github.com/sw965/omw/funcs"
	omwbits "github.com/sw965/omw/math/bits"
	omwslices "github.com/sw965/omw/slices"
	"slices"
)

func NewSymmetryFeaturesByGroupFromCanonicalIndices(canonicalIdxs []int) ([][]gothello.Feature, error) {
	if !omwslices.IsUnique(canonicalIdxs) {
		return nil, fmt.Errorf("idxs is not unique")
	}

	perspectivesToFeature := func(perspectives []gothello.Perspective) (gothello.Feature, error) {
		feature := gothello.Feature{}
		var err error

		for i, perspective := range perspectives {
			idx := canonicalIdxs[i]

			switch perspective {
			case gothello.Self:
				feature.Self, err = omwbits.ToggleBit64(feature.Self, idx)
			case gothello.Opponent:
				feature.Opponent, err = omwbits.ToggleBit64(feature.Opponent, idx)
			}

			if err != nil {
				return gothello.Feature{}, err
			}
		}
		return feature, nil
	}

	perspectiveSequences := omwslices.Sequences[[]gothello.Perspective](gothello.AllPerspectives, len(canonicalIdxs))
	features, err := funcs.MapErr(perspectiveSequences, perspectivesToFeature)
	if err != nil {
		return nil, err
	}

	bases := make([]gothello.Feature, 0, len(features))
	for _, feature := range features {
		syms := funcs.Juxt(feature, gothello.FeatureSymmetryFuncs)
		if omwslices.AllFunc(syms, func(f gothello.Feature) bool { return !slices.Contains(bases, f) }) {
			bases = append(bases, feature)
		}
	}

	return funcs.Map(bases, func(f gothello.Feature) []gothello.Feature {
		return funcs.Juxt(f, gothello.FeatureSymmetryFuncs)
	}), nil
}

type PatternIndexer []map[gothello.BitBoard]map[gothello.Feature]int

func NewPatternIndexerFromCanonicalIndices(canonicalMaskIdxs, canonicalMoveIdxs []int) (PatternIndexer, error) {
	if !omwslices.IsUnique(canonicalMaskIdxs) {
		return nil, fmt.Errorf("maskIdxs is not unique")
	}
	if !omwslices.IsUnique(canonicalMoveIdxs) {
		return nil, fmt.Errorf("moveIdxs is not unique")
	}
	if !omwslices.IsSubset(canonicalMaskIdxs, canonicalMoveIdxs) {
		return nil, fmt.Errorf("maskIdxs must be a subset of moveIdxs")
	}

	symFeaturesByGroup, err := NewSymmetryFeaturesByGroupFromCanonicalIndices(canonicalMaskIdxs)
	if err != nil {
		return nil, err
	}

	canonicalMask, err := omwbits.New64FromIndices[gothello.BitBoard](canonicalMaskIdxs)
	if err != nil {
		return nil, err
	}
	symMasks := funcs.Juxt(canonicalMask, gothello.BitBoardSymmetryFuncs)

	indexer := make(PatternIndexer, gothello.BoardSize)
	patternIdx := 0

	for _, symFeatures := range symFeaturesByGroup {
		for symI := range gothello.FeatureSymmetryFuncs {
			feature := symFeatures[symI]
			mask := symMasks[symI]
			empties := feature.Empties()
			moves := empties & mask
			featureMoveIdxs := omwbits.OneIndices64(moves)

			for _, canonicalMoveIdx := range canonicalMoveIdxs {
				symMoveIdx := gothello.IndexSymmetryFuncs[symI](canonicalMoveIdx)

				if slices.Contains(featureMoveIdxs, symMoveIdx) {
					if indexer[symMoveIdx] == nil {
						indexer[symMoveIdx] = map[gothello.BitBoard]map[gothello.Feature]int{}
					}

					if indexer[symMoveIdx][mask] == nil {
						indexer[symMoveIdx][mask] = map[gothello.Feature]int{}
					}

					if _, ok := indexer[symMoveIdx][mask][feature]; !ok {
						indexer[symMoveIdx][mask][feature] = patternIdx
					}
				}
			}
		}
		patternIdx++
	}
	return indexer, nil
}

func (pi PatternIndexer) MatchIndices(feature gothello.Feature, moveIdxs []int) []int {
	var matchIdxs []int = nil
	for _, idx := range moveIdxs {
		indexer := pi[idx]
		for mask, featureByIdx := range indexer {
			masked := feature.AndBitBoard(mask)
			patternIdx, ok := featureByIdx[masked]
			if ok {
				matchIdxs = append(matchIdxs, patternIdx)
			}
		}
	}
	return matchIdxs
}
