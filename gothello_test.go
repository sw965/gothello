package gothello_test

import (
	"testing"
	"fmt"
	"github.com/sw965/gothello"
	"golang.org/x/exp/slices"
	omwrand "github.com/sw965/omw/math/rand"
)

func TestOneIndices(t *testing.T) {
	bb := gothello.BitBoard(0b10000000_01000000_00100000_00010000_00001000_00000100_00000010_00000001)
	fmt.Println(bb.OneIndices())
}

func TestSidePoints(t *testing.T) {
	fmt.Println(gothello.RIGHT_FLOW_UP_SIDE_POINTS)
	fmt.Println(gothello.LEFT_FLOW_UP_SIDE_POINTS)
	
	fmt.Println(gothello.RIGHT_FLOW_DOWN_SIDE_POINTS)
	fmt.Println(gothello.LEFT_FLOW_DOWN_SIDE_POINTS)

	fmt.Println(gothello.DOWN_FLOW_LEFT_SIDE_POINTS)
	fmt.Println(gothello.UP_FLOW_LEFT_SIDE_POINTS)

	fmt.Println(gothello.DOWN_FLOW_RIGHT_SIDE_POINTS)
	fmt.Println(gothello.UP_FLOW_RIGHT_SIDE_POINTS)
}

func TestOneHotBitBoards(t *testing.T) {
	for _, oneHot := range gothello.ONE_HOT_BIT_BOARDS {
		fmt.Println(oneHot.ToArray())
	}
}

func TestPointToIndex(t *testing.T) {
	for _, point := range gothello.ALL_POINTS {
		fmt.Println(point, point.ToIndex())
	}
}

func TestToOneHots(t *testing.T) {
	b := gothello.BitBoard(0b10000000_01000000_00100000_00010000_00001000_00000100_00000010_00000001)
	oneHots := b.ToOneHots()
	for _, oneHot := range oneHots {
		fmt.Println(oneHot.ToArray())
	} 
}

func TestToPoints(t *testing.T) {
	b := gothello.BitBoard(0b10000000_01000000_00100000_00010000_00001000_00000100_00000010_00000001)
	points := b.ToPoints()
	fmt.Println(points)
}

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

	result1 := state.Put(gothello.BIT_BOARD_BY_POINT[gothello.Point{Row:4, Column:7}])

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

func TestRotate90(t *testing.T) {
	state := gothello.NewInitState()
	state.Black = state.Black.ToggleBit(&gothello.Point{Row:0, Column:0})
	state.Black = state.Black.ToggleBit(&gothello.Point{Row:0, Column:1})
	state.Black = state.Black.ToggleBit(&gothello.Point{Row:0, Column:2})
	state.White = state.White.ToggleBit(&gothello.Point{Row:0, Column:4})
	state.White = state.White.ToggleBit(&gothello.Point{Row:0, Column:6})

	fmt.Println(state.ToString())

	rotated90 := state.Rotate90()
	rotated180 := state.Rotate180()
	rotated270 := state.Rotate270()
	mirrorH := state.MirrorHorizontal()
	mirrorV := state.MirrorVertical()

	fmt.Println(rotated90.ToString())
	fmt.Println(rotated180.ToString())
	fmt.Println(rotated270.ToString())
	fmt.Println(mirrorH.ToString())
	fmt.Println(mirrorV.ToString())
}

func TestFlipPointBitBoard(t *testing.T) {
	r := omwrand.NewMt19937()
	for i := 0; i < 6400; i++ {
		state := gothello.NewInitState()
		for {
			legalBitBoardPoints := state.LegalPointBitBoard().ToOneHots()
			for _, bbPoint := range legalBitBoardPoints {
				hand := state.NewHandPairBitBoard()
				if hand.Self.FlipPointBitBoard(hand.Opponent, bbPoint) == 0 {
					t.Errorf("テスト失敗")
					return
				}
			}
			bbPoint := omwrand.Choice(legalBitBoardPoints, r)
			state = state.Put(bbPoint)

			black := state.Black.LegalPointBitBoard(state.White)
			white := state.White.LegalPointBitBoard(state.Black)

			if black == 0 && white == 0 {
				break
			}
		}
	}
}