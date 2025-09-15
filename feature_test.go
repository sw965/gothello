package gothello_test

import (
	"testing"
	"github.com/sw965/gothello"
)

func Test(t *testing.T) {
	partialBoard1d := gothello.PartialBoard1D{
		gothello.Black,
		gothello.Black,
		gothello.White,
		gothello.Empty,
		gothello.White,
		gothello.Black,
		gothello.White,
		gothello.Empty,
	}

	partialBoard1d.Put()
}