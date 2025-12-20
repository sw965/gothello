package policy_test

import (
	"fmt"
	"github.com/sw965/gothello"
	"github.com/sw965/gothello/bunny/policy"
	"testing"
)

func Test(t *testing.T) {
	scorer, err := policy.NewYasuriScorer()
	if err != nil {
		panic(err)
	}

	feature, err := gothello.NewFeatureFromIndices([]int{4, 32}, []int{1, 2, 3, 8, 16, 24})
	if err != nil {
		panic(err)
	}
	arr, err := feature.ToArray()
	if err != nil {
		panic(err)
	}
	fmt.Println(arr)

	gainByMoveIdx := scorer.ScoreByMoveIndex(feature)
	fmt.Println(gainByMoveIdx)
}
