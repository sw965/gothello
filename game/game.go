package game

import (
	"github.com/sw965/gothello"
	game "github.com/sw965/crow/game/sequential"
)

func NewLogic() game.Logic[gothello.State, gothello.BitBoard, gothello.Disc] {
	legalActionsProvider := func(state gothello.State) []gothello.BitBoard {
		return state.LegalBitBoard().ToSingles()
	}

	transitioner := func(state gothello.State, move gothello.BitBoard) (gothello.State, error) {
		return state.Put(move)
	}

	comparator := func(s1, s2 gothello.State) bool {
		return s1 == s2
	}

	currentAgentGetter := func(state gothello.State) gothello.Disc {
		return state.Turn
	}

	placementsJudger := func(state gothello.State) (game.PlacementByAgent[gothello.Disc], error) {
		blackLegal := state.Black.LegalBitBoard(state.White)
		whiteLegal := state.White.LegalBitBoard(state.Black)

		//ゲームが終了している場合
		if blackLegal == 0 && whiteLegal == 0 {
			blackCount := state.Black.OnesCount()
			whiteCount := state.White.OnesCount()
			placements := game.PlacementByAgent[gothello.Disc]{}

			if blackCount > whiteCount {
				placements[gothello.Black] = 1
				placements[gothello.White] = 2
			} else if blackCount < whiteCount {
				placements[gothello.Black] = 2
				placements[gothello.White] = 1
			} else {
				placements[gothello.Black] = 1
				placements[gothello.White] = 1
			}
			return placements, nil
		}

		//まだゲームが終了していない場合
		return nil, nil
	}

	logic := game.Logic[gothello.State, gothello.BitBoard, gothello.Disc]{
		LegalActionsProvider:legalActionsProvider,
		Transitioner:transitioner,
		Comparator:comparator,
		CurrentAgentGetter:currentAgentGetter,
		PlacementsJudger:placementsJudger,
		Agents:[]gothello.Disc{gothello.Black, gothello.White},
	}

	logic.SetStandardResultScoresEvaluator()
	return logic
}