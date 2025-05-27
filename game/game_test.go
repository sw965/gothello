package game_test

import (
	"testing"
	"fmt"
	orand "github.com/sw965/omw/math/rand"
	"github.com/sw965/gothello/game"
	"github.com/sw965/gothello"
	"runtime"
	cgame "github.com/sw965/crow/game/sequential"
	"github.com/sw965/omw/fn"
	"slices"
)

func Test(t *testing.T) {
	gameLogic := game.NewLogic()

	p := runtime.NumCPU()
	playersByWorker := make([]cgame.PlayerByAgent[gothello.State, gothello.BitBoards, gothello.BitBoard, gothello.Color], p)
	for i := range playersByWorker {
		rng := orand.NewMt19937()
		players := gameLogic.MakePlayerByAgent()
		players[gothello.BLACK] = gameLogic.NewRandActionPlayer(rng)
		players[gothello.WHITE] = func(state gothello.State, legals gothello.BitBoards) (gothello.BitBoard, error) {
			corners := fn.Filter(legals, func(bb gothello.BitBoard) bool {
				return slices.Contains(gothello.CORNER_BIT_BOARD.ToSingles(), bb)
			})

			if len(corners) != 0 {
				return orand.Choice(corners, rng), nil
			}

			notXs := fn.Filter(legals, func(bb gothello.BitBoard) bool {
				return !slices.Contains(gothello.X_BIT_BOARD.ToSingles(), bb)
			})

			if len(notXs) != 0 {
				return orand.Choice(notXs, rng), nil
			}
			return orand.Choice(legals, rng), nil
		}
		playersByWorker[i] = players
	}

	playoutNum := 196000
	initStates := make([]gothello.State, playoutNum)
	for i := range initStates {
		initStates[i] = gothello.NewInitState()
	}

	finalStates, err := gameLogic.Playouts(initStates, playersByWorker)
	if err != nil {
		panic(err)
	}

	blackWinCount := 0
	drawCount := 0

	for _, final := range finalStates {
		blackCount := final.Black.Count()
		whiteCount := final.White.Count()
		if blackCount > whiteCount {
			blackWinCount++
		} else if blackCount == whiteCount {
			drawCount++
		}
	}
	fmt.Println(finalStates[0].Black.Count(), finalStates[0].White.Count())

	whiteWinCount := playoutNum - blackWinCount - drawCount

	fmt.Println("黒勝率", float64(blackWinCount) / float64(playoutNum))
	fmt.Println("白勝率", float64(whiteWinCount) / float64(playoutNum))
	fmt.Println("引き分け率", float64(drawCount) / float64(playoutNum))
}