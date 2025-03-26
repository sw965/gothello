package gothello_test

import (
	"testing"
	//"fmt"
	"github.com/sw965/gothello"
	"golang.org/x/exp/slices"
)

func Test(t *testing.T) {
	/*
		引用した局面
		https://youtube.com/shorts/oTp0EQzgy0o?si=ozqqu9Fy3zhOq4QZ
	*/

	state := gothello.State{
		Black:0b01111110_00111100_01111100_00101100_11011100_11100100_10111100_01111110,
		White:0b00000000_00000000_10000000_01010000_00100000_00011000_00000000_00000000,
		Hand:gothello.BLACK,
	}

	legalPoints := state.LegalPointBitBoard().ToPoints()
	if len(legalPoints) != 1 {
		t.Errorf("テスト失敗")
	}

	if legalPoints[0] != (gothello.Point{Row:4, Column:7}) {
		t.Errorf("テスト失敗")
	}

	//4, 7に石を置く。
	expected1 := gothello.State{
		Black:0b01111110_00111100_01111100_11101100_11011100_11100100_10111100_01111110,
		White:0b00000000_00000000_10000000_00010000_00100000_00011000_00000000_00000000,
		Hand:gothello.WHITE,
	}

	result1 := state.Put(&gothello.Point{Row:4, Column:7})

	if result1 != expected1 {
		t.Errorf("テスト失敗")
	}

	legalPoints = result1.LegalPointBitBoard().ToPoints()

	expectedLegalPoints := gothello.Points{
		gothello.Point{Row:0, Column:7},
		gothello.Point{Row:1, Column:1},
		gothello.Point{Row:2, Column:1},
		gothello.Point{Row:3, Column:1},
		gothello.Point{Row:4, Column:1},
		gothello.Point{Row:5, Column:1},
		gothello.Point{Row:6, Column:6},
		gothello.Point{Row:6, Column:7},
	}

	if !slices.Equal(legalPoints, expectedLegalPoints) {
		t.Errorf("テスト失敗")
	}
}