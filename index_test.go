package gothello_test

import (
	"testing"
	"fmt"
	"github.com/sw965/gothello"
)

func TestIndices(t *testing.T) {
	fmt.Println("upSideRightFlowIdxs =", gothello.UpSideRightFlowIndices)
	fmt.Println("upSideLeftFlowIdxs = ", gothello.UpSideLeftFlowIndices)

	fmt.Println("downSideRightFlowIdxs = ", gothello.DownSideRightFlowIndices)
	fmt.Println("downSideLeftFlowIdxs = ", gothello.DownSideLeftFlowIndices)

	fmt.Println("leftSideDownFlowIdxs = ", gothello.LeftSideDownFlowIndices)
	fmt.Println("leftSideUpFlowIdxs = ", gothello.LeftSideUpFlowIndices)

	fmt.Println("rightSideDownFlowIdxs = ", gothello.RightSideDownFlowIndices)
	fmt.Println("rightSideUpFlowIdxs = ", gothello.RightSideUpFlowIndices)
}