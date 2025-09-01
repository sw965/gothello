package policy

import (
	"fmt"
	//"github.com/sw965/crow/model/linear"
	"github.com/sw965/gothello"
	omwslices "github.com/sw965/omw/slices"
	omwbits "github.com/sw965/omw/math/bits"
	"github.com/sw965/omw/funcs"
	"slices"
)

func NewSymmetryFeaturesByGroupFromCanonicalIndices(canonicalIdxs []int) ([][]gothello.Feature, error) {
	if !omwslices.IsUnique(canonicalIdxs) {
		return nil, fmt.Errorf("idxs is not unique")
	}

	discsToFeature := func(discs []gothello.Disc) (gothello.Feature, error) {
		feature := gothello.Feature{}
		var err error
		for i, disc := range discs {
			idx := canonicalIdxs[i]
			switch disc {
			case gothello.Black:
				feature.Self, err = omwbits.ToggleBit64(feature.Self, idx)
			case gothello.White:
				feature.Opponent, err = omwbits.ToggleBit64(feature.Opponent, idx)
			}
			if err != nil {
				return gothello.Feature{}, err
			}
		}
		return feature, nil
	}

	discSequences := omwslices.Sequences[[]gothello.Disc](gothello.AllDiscs, len(canonicalIdxs))
	features, err := funcs.MapErr(discSequences, discsToFeature)
	if err != nil {
		return nil, err
	}

	bases := make([]gothello.Feature, 0, len(features))
	for _, feature := range features {
		syms := funcs.Juxt(feature, gothello.FeatureSymmetryFuncs)
		if omwslices.AllFunc(syms, func(f gothello.Feature) bool {
			return !slices.Contains(bases, f)
		}) {
			bases = append(bases, feature)
		}
	}

	return funcs.Map(bases, func(f gothello.Feature) []gothello.Feature {
		return funcs.Juxt(f, gothello.FeatureSymmetryFuncs)
	}), nil
}

type PatternIndexer []map[gothello.BitBoard]map[gothello.Feature]int

func NewPatternIndexerFromCanonicalIndices(canonicalMaskIdxs, canonicalEmptyMaskIdxs []int) (PatternIndexer, error) {
	if !omwslices.IsUnique(canonicalMaskIdxs) {
		return nil, fmt.Errorf("maskIdxs is not unique")
	}

	if !omwslices.IsUnique(canonicalEmptyMaskIdxs) {
		return nil, fmt.Errorf("emptyIdxs is not unique")
	}

	if !omwslices.IsSubset(canonicalMaskIdxs, canonicalEmptyMaskIdxs) {
		return nil, fmt.Errorf("emptyMaskIdxs is not subset")
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

	emptyMaskIdxsByGroup := funcs.Map(canonicalEmptyMaskIdxs, func(idx int) []int {
		return funcs.Juxt(idx, gothello.IndexSymmetryFuncs)
	})
	emptyMaskIdxTable := omwslices.Transpose(emptyMaskIdxsByGroup)

	indexer := make(PatternIndexer, gothello.BoardSize)
	for i := range indexer {
		indexer[i] = map[gothello.BitBoard]map[gothello.Feature]int{}
		for _, mask := range symMasks {
			indexer[i][mask] = map[gothello.Feature]int{}
		}
	}

	patternIdx := 0
	for _, symFeatures := range symFeaturesByGroup {
		for emptyI := 0; emptyI < len(canonicalEmptyMaskIdxs); emptyI++ {
			for symId, feature := range symFeatures {
				mask := symMasks[symId]
				emptyMaskIdx := emptyMaskIdxTable[symId][emptyI]
				empties := feature.Empties()
				emptyIdxs := omwbits.OneIndices64(empties)
				if slices.Contains(emptyIdxs, emptyMaskIdx) {
					indexer[emptyMaskIdx][mask][feature] = patternIdx
				}
			}
			patternIdx++
		}
	}
	return indexer, nil
}