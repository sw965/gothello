package policy_test

import (
	"testing"
	"github.com/sw965/gothello"
	"github.com/sw965/gothello/bunny/policy"
	omwbits "github.com/sw965/omw/math/bits"
	"fmt"
)

func TestPatternIndexer(t *testing.T) {
	patternIndexer, err := policy.NewPatternIndexerFromCanonicalIndices([]int{0, 1, 2, 3, 4, 5, 6, 7}, []int{0, 1, 2, 3, 4, 5, 6, 7})
	if err != nil {
		panic(err)
	}

	mask, err := omwbits.New64FromIndices[gothello.BitBoard]([]int{0, 1, 2, 3, 4, 5, 6, 7})
	if err != nil {
		panic(err)
	}

	feature, err := gothello.NewFeatureFromIndices([]int{1, 4}, []int{2})
	if err != nil {
		panic(err)
	}

	idx, ok := patternIndexer[gothello.UpLeftCornerIndex][mask][feature]
	fmt.Println(idx, ok)

	feature, err = gothello.NewFeatureFromIndices([]int{8, 32}, []int{16})
	if err != nil {
		panic(err)
	}

	mask, err = omwbits.New64FromIndices[gothello.BitBoard]([]int{0, 8, 16, 24, 32, 40, 48, 56})
	if err != nil {
		panic(err)
	}

	idx, ok = patternIndexer[gothello.UpLeftCornerIndex][mask][feature]
	fmt.Println(idx, ok)

	testFeature, err := gothello.NewFeatureFromIndices([]int{1, 3, 8, 24}, []int{2, 16})
	if err != nil {
		panic(err)
	}

	matchIdxs := patternIndexer.MatchIndices(testFeature, []int{0})
	fmt.Println(matchIdxs)
}