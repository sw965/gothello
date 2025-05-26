package game

import (
	"github.com/sw965/gothello"
	game "github.com/sw965/crow/game/sequential"
)

func NewLogic() game.Logic[gothello.State, gothello.BitBoards, gothello.BitBoard, gothello.Color] {
	legalActionsProvider := func(state gothello.State) gothello.BitBoards {
		return state.LegalBitBoard().ToSingles()
	}

	transitioner := func(state gothello.State, move gothello.BitBoard) (gothello.State, error) {
		return state.Put(move)
	}

	comparator := func(s1, s2 gothello.State) bool {
		return s1 == s2
	}

	currentTurnAgentGetter := func(state gothello.State) gothello.Color {
		return state.Hand
	}

	placementsJudger := func(state gothello.State) (game.PlacementByAgent[gothello.Color], error) {
		blackLegal := state.Black.LegalBitBoard(state.White)
		whiteLegal := state.White.LegalBitBoard(state.Black)

		//ゲームが終了している場合
		if blackLegal == 0 && whiteLegal == 0 {
			blackCount := state.Black.Count()
			whiteCount := state.White.Count()
			placements := game.PlacementByAgent[gothello.Color]{}

			if blackCount > whiteCount {
				placements[gothello.BLACK] = 1
				placements[gothello.WHITE] = 2
			} else if blackCount < whiteCount {
				placements[gothello.BLACK] = 2
				placements[gothello.WHITE] = 1
			} else {
				placements[gothello.BLACK] = 1
				placements[gothello.WHITE] = 1
			}
			return placements, nil
		}

		//まだゲームが終了していない場合
		return nil, nil
	}

	logic := game.Logic[gothello.State, gothello.BitBoards, gothello.BitBoard, gothello.Color]{
		LegalActionsProvider:legalActionsProvider,
		Transitioner:transitioner,
		Comparator:comparator,
		CurrentTurnAgentGetter:currentTurnAgentGetter,
		PlacementsJudger:placementsJudger,
	}

	logic.SetStandardResultScoresEvaluator()
	return logic
}