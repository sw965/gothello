package gothello_test

import (
	"fmt"
	"github.com/sw965/gothello"
	"testing"
)

func Test(t *testing.T) {
	partialBoard1d := gothello.PartialFeature1D{
		0,
		2,
		2,
		1,
		2,
		1,
		2,
		0,
	}

	fmt.Println(partialBoard1d.Move(0))
}
