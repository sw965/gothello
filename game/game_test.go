package game_test

import (
	"testing"
	"fmt"
	orand "github.com/sw965/omw/math/rand"
	"github.com/sw965/gothello/game"
	"github.com/sw965/gothello"
	"runtime"
	"github.com/sw965/omw/fn"
	"slices"
	"math/rand"
	cgame "github.com/sw965/crow/game/sequential"
)

func Test(t *testing.T) {
	gameLogic := game.NewLogic()
	p := runtime.NumCPU()
	rngs := make([]*rand.Rand, p)
	for i := range rngs {
		rngs[i] = orand.NewMt19937()
	}

	policyProvider := func(state gothello.State, legals []gothello.BitBoard) cgame.Policy[gothello.BitBoard] {
		if state.Hand == gothello.Black {
			n := len(legals)
			p := 1.0 / float32(n)
			m := cgame.Policy[gothello.BitBoard]{}
			for _, a := range legals {
				m[a] = p
			}
			return m
		}

		corners := fn.Filter(legals, func(bb gothello.BitBoard) bool {
			return slices.Contains(gothello.CornerBitBoard.ToSingles(), bb)
		})

		if len(corners) != 0 {
			m := cgame.Policy[gothello.BitBoard]{}
			cornerBitBoards := gothello.CornerBitBoard.ToSingles()
			for _, a := range legals {
				if slices.Contains(cornerBitBoards, a) {
					m[a] = 1.0
				} else {
					m[a] = 0.0
				}
			}
			return m
		}

		notXs := fn.Filter(legals, func(bb gothello.BitBoard) bool {
			return !slices.Contains(gothello.XBitBoard.ToSingles(), bb)
		})

		if len(notXs) != 0 {
			xBitBoards := gothello.XBitBoard.ToSingles()
			m := cgame.Policy[gothello.BitBoard]{}
			for _, a := range legals {
				if slices.Contains(xBitBoards, a) {
					m[a] = 0.0
				} else {
					m[a] = 1.0
				}
			}
			return m
		}

		n := len(legals)
		p := 1.0 / float32(n)
		m := cgame.Policy[gothello.BitBoard]{}
		for _, a := range legals {
			m[a] = p
		}
		return m
	}

	playoutNum := 19600
	initStates := make(gothello.States, playoutNum)
	for i := range initStates {
		initStates[i] = gothello.NewInitState()
	}

	finalStates, err := gameLogic.Playouts([]gothello.State(initStates), policyProvider, rngs)
	if err != nil {
		panic(err)
	}

	counter, err := gothello.States(finalStates).CountResult()
	if err != nil {
		panic(err)
	}

	fmt.Println("黒勝率", counter.BlackWinRate(true))
	fmt.Println("白勝率", counter.WhiteWinRate(true))
	fmt.Println("引き分け率", counter.DrawRate())
}