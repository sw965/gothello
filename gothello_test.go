package gothello_test

import (
	"testing"
	"fmt"
	"github.com/sw965/gothello"
	omwrand "github.com/sw965/omw/math/rand"
)

func TestSideIndices(t *testing.T) {
	fmt.Println("upSideIdxs =", gothello.UpSideIndices)
	fmt.Println("downSideIdxs =", gothello.DownSideIndices)
	fmt.Println("leftSideIdxs =", gothello.LeftSideIndices)
	fmt.Println("rightSideIdxs =", gothello.RightSideIndices)
}

func TestEdgeIndices(t *testing.T) {
	fmt.Println("upEdgeIdxs =", gothello.UpEdgeIndices)
	fmt.Println("downEdgeIdxs =", gothello.DownEdgeIndices)
	fmt.Println("leftEdgeIdxs =", gothello.LeftEdgeIndices)
	fmt.Println("rightEdgeIdxs =", gothello.RightEdgeIndices)
}

func TestAdjacentBySingleBitBoard(t *testing.T) {
	for k, v := range gothello.AdjacentBySingle {
		fmt.Println(k.ToArray())
		fmt.Println(v.ToArray())
		fmt.Println("")
	}
}

func TestPut(t *testing.T) {
	init := gothello.NewInitState()

	//1手目
	move1 := gothello.Cell{Row:3, Column:"e"}
	state1, err := init.Put(move1.ToBitBoard())
	if err != nil {
		panic(err)
	}

	legal1 := state1.LegalBitBoard()
	expectedState1 := gothello.State{
		Black:0b00000000_00000000_00000000_00010000_00011000_00010000_00000000_00000000,
		White:0b00000000_00000000_00000000_00001000_00000000_00000000_00000000_00000000,
		Hand:gothello.White,
	}
	expectedLegal1 := gothello.BitBoard(0b00000000_00000000_00000000_00100000_00000000_00101000_00000000_00000000)

	if state1 != expectedState1 {
		fmt.Println(state1.ToArray())
		t.Errorf("テスト失敗")
	}

	if expectedLegal1 != legal1 {
		fmt.Println(legal1.ToArray())
		t.Errorf("テスト失敗")
	}

	//2手目
	move2 := gothello.Cell{Row:3, Column:"f"}
	state2, err := state1.Put(move2.ToBitBoard())
	if err != nil {
		panic(err)
	}

	legal2 := state2.LegalBitBoard()
	expectedState2 := gothello.State{
		Black:0b00000000_00000000_00000000_00010000_00001000_00010000_00000000_00000000,
		White:0b00000000_00000000_00000000_00001000_00010000_00100000_00000000_00000000,
		Hand:gothello.Black,
	}
	expectedLegal2 := gothello.BitBoard(0b00000000_00000000_00001000_00000100_00100000_01000000_00000000_00000000)

	if expectedState2 != state2 {
		fmt.Println(state2.ToArray())
		t.Errorf("テスト失敗")
	}

	if expectedLegal2 != legal2 {
		fmt.Println(legal2.ToArray())
		t.Errorf("テスト失敗")
	}

	//3手目
	move3 := gothello.Cell{Row:4, Column:"f"}
	state3, err := state2.Put(move3.ToBitBoard())
	if err != nil {
		panic(err)
	}

	legal3 := state3.LegalBitBoard()
	expectedState3 := gothello.State{
		Black:0b00000000_00000000_00000000_00010000_00111000_00010000_00000000_00000000,
		White:0b00000000_00000000_00000000_00001000_00000000_00100000_00000000_00000000,
		Hand:gothello.White,
	}
	expectedLegal3 := gothello.BitBoard(0b00000000_00000000_00000000_00100000_00000000_00001000_00000000_00000000)

	if expectedState3 != state3 {
		fmt.Println(state3.ToArray())
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
			bb, err := omwrand.Choice(legalBitBoards, r)
			if err != nil {
				panic(err)
			}

			state, err = state.Put(bb)
			if err != nil {
				panic(err)
			}

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