package policy_test

import (
	"testing"
	"github.com/sw965/gothello/bunny/policy"
	omwbits "github.com/sw965/omw/math/bits"
	"fmt"
)

func Test(t *testing.T) {
	patternIndexer, err := policy.NewPatternIndexerFromCanonicalIndices([]int{0, 1, 8, 9}, []int{0, 1, 8, 9})
	if err != nil {
		panic(err)
	}

	for emptyIdx := range patternIndexer {
		for mask, v := range patternIndexer[emptyIdx] {
			for feature, idx := range v {
				maskIdx := omwbits.OneIndices64(mask)
				arr, err := feature.ToArray()
				if err != nil {
					panic(err)
				}
				fmt.Println(emptyIdx, maskIdx, idx)
				fmt.Println(arr)
				fmt.Println("")
			}
		}
	}
}

