package policy_test

import (
	"testing"
	"github.com/sw965/gothello/bunny/policy"
)

func Test(t *testing.T) {
	board := []int{0, 1, 1, 2, 0, 2, 2, 1}
	policy.Legal(board)
}