package gothello_test

import (
	"testing"
	"fmt"
	"math/bits"
	"github.com/sw965/gothello"
	omwrand "github.com/sw965/omw/math/rand"
)

func TestSideIndices(t *testing.T) {
	fmt.Println("upSideIdxs =", gothello.UP_SIDE_INDICES)
	fmt.Println("downSideIdxs =", gothello.DOWN_SIDE_INDICES)
	fmt.Println("leftSideIdxs =", gothello.LEFT_SIDE_INDICES)
	fmt.Println("rightSideIdxs =", gothello.RIGHT_SIDE_INDICES)
}

func TestEdgeIndices(t *testing.T) {
	fmt.Println("upEdgeIdxs =", gothello.UP_EDGE_INDICES)
	fmt.Println("downEdgeIdxs =", gothello.DOWN_EDGE_INDICES)
	fmt.Println("leftEdgeIdxs =", gothello.LEFT_EDGE_INDICES)
	fmt.Println("rightEdgeIdxs =", gothello.RIGHT_EDGE_INDICES)
}

func TestAdjacentBySingleBitBoard(t *testing.T) {
	for k, v := range gothello.ADJACENT_BY_SINGLE_BIT_BOARD {
		fmt.Println(k.ToArray())
		fmt.Println(v.ToArray())
		fmt.Println("")
	}
}

func TestNewRandState(t *testing.T) {
	testNum := 12800
	r := omwrand.NewMt19937()
	blackSum := 0
	whiteSum := 0
	spaceSum := 0
	for i := 0; i < testNum; i++ {
		state := gothello.NewRandState(0.50, 0.25, r)
		if bits.OnesCount64(uint64(state.Black&gothello.UP_SIDE_BIT_BOARD)) == 8 {
			fmt.Println("終わり", i, state.Black.ToArray())
			fmt.Println(state.White.ToArray())
			break
		}
		blackCount := bits.OnesCount64(uint64(state.Black))
		blackSum += blackCount
		whiteCount := bits.OnesCount64(uint64(state.White))
		whiteSum += whiteCount
		spaceCount := bits.OnesCount64(uint64(state.SpaceBitBoard()))
		spaceSum += spaceCount
	}
	fmt.Println(float64(blackSum) / float64(testNum))
	fmt.Println(float64(whiteSum) / float64(testNum))
	fmt.Println(float64(spaceSum) / float64(testNum))
}

func TestPut(t *testing.T) {
	init := gothello.NewInitState()

	//1手目
	move1 := gothello.Cell{Row:3, Column:"e"}
	state1 := init.Put(move1.ToBitBoard())
	legal1 := state1.LegalBitBoard()
	expectedState1 := gothello.State{
		Black:0b00000000_00000000_00000000_00010000_00011000_00010000_00000000_00000000,
		White:0b00000000_00000000_00000000_00001000_00000000_00000000_00000000_00000000,
		Hand:gothello.WHITE,
	}
	expectedLegal1 := gothello.BitBoard(0b00000000_00000000_00000000_00100000_00000000_00101000_00000000_00000000)

	if state1 != expectedState1 {
		fmt.Println(state1.ToString())
		t.Errorf("テスト失敗")
	}

	if expectedLegal1 != legal1 {
		fmt.Println(legal1.ToArray())
		t.Errorf("テスト失敗")
	}

	//2手目
	move2 := gothello.Cell{Row:3, Column:"f"}
	state2 := state1.Put(move2.ToBitBoard())
	legal2 := state2.LegalBitBoard()
	expectedState2 := gothello.State{
		Black:0b00000000_00000000_00000000_00010000_00001000_00010000_00000000_00000000,
		White:0b00000000_00000000_00000000_00001000_00010000_00100000_00000000_00000000,
		Hand:gothello.BLACK,
	}
	expectedLegal2 := gothello.BitBoard(0b00000000_00000000_00001000_00000100_00100000_01000000_00000000_00000000)

	if expectedState2 != state2 {
		fmt.Println(state2.ToString())
		t.Errorf("テスト失敗")
	}

	if expectedLegal2 != legal2 {
		fmt.Println(legal2.ToArray())
		t.Errorf("テスト失敗")
	}

	//3手目
	move3 := gothello.Cell{Row:4, Column:"f"}
	state3 := state2.Put(move3.ToBitBoard())
	legal3 := state3.LegalBitBoard()
	expectedState3 := gothello.State{
		Black:0b00000000_00000000_00000000_00010000_00111000_00010000_00000000_00000000,
		White:0b00000000_00000000_00000000_00001000_00000000_00100000_00000000_00000000,
		Hand:gothello.WHITE,
	}
	expectedLegal3 := gothello.BitBoard(0b00000000_00000000_00000000_00100000_00000000_00001000_00000000_00000000)

	if expectedState3 != state3 {
		fmt.Println(state3.ToString())
		t.Errorf("テスト失敗")
	}

	if expectedLegal3 != legal3 {
		fmt.Println(legal3.ToArray())
		t.Errorf("テスト失敗")
	}
}

func TestFlipBitBoard(t *testing.T) {
	r := omwrand.NewMt19937()
	testGameNum := 6400
	for i := 0; i < testGameNum; i++ {
		state := gothello.NewInitState()
		for {
			legalBitBoards := state.LegalBitBoard().ToSingles()

			/*
				合法手に石を置けば、ひっくり返せる石があるという事。
				FlipPointBitBoardが0でないかテストする。
				0ならば合法手であるのに、ひっくり返せる石がないという事にどこかにロジックミスがある。
			*/
			for _, bb := range legalBitBoards {
				hand := state.NewHandPairBitBoard()
				if hand.Self.FlipBitBoard(hand.Opponent, bb) == 0 {
					t.Errorf("テスト失敗")
					return
				}
			}

			//ランダムに手を選択する。
			bb := omwrand.Choice(legalBitBoards, r)
			state = state.Put(bb)

			//両プレイヤーの合法手がなくなった場合、ゲームが終了する。
			black := state.Black.LegalBitBoard(state.White)
			white := state.White.LegalBitBoard(state.Black)
			if black == 0 && white == 0 {
				break
			}
		}
	}
}

func TestMirrorHorizontalIndex(t *testing.T) {
	//左上の隅
	result := gothello.MirrorHorizontalIndex(0)
	expected := 7

	if result != expected {
		t.Errorf("テスト失敗")
	}

	result = gothello.MirrorHorizontalIndex(19)
	expected = 20

	if result != expected {
		t.Errorf("テスト失敗")
	}

	//左下のX
	result = gothello.MirrorHorizontalIndex(49)
	expected = 54
	if result != expected {
		fmt.Println(result)
		t.Errorf("テスト失敗")
	}
}

func TestMirrorVerticalIndex(t *testing.T) {

}