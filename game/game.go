package game

import (
	"github.com/sw965/omw/mathx/bitsx"
	"github.com/sw965/gothello"
	game "github.com/sw965/crow/game/sequential"
)

func NewGameLogic() game.Logic[gothello.State, gothello.BitBoard, gothello.Turn] {
	return game.Logic[gothello.State, gothello.BitBoard, gothello.Turn]{
		LegalMovesFunc:func(state gothello.State) []gothello.BitBoard {
			legals := state.Legals()
			return bitsx.Singles(legals)
		},
		MoveFunc:gothello.State.Move,
		EqualFunc:func(s1, s2 gothello.State) bool {
			return s1 == s2
		},
		CurrentAgentFunc:gothello.State.Turn,
	}
}

func NewGameEngine() game.Engine[gothello.State, gothello.BitBoard, gothello.Turn]{
	e := game.Engine[gothello.State, gothello.BitBoard, gothello.Turn]{
		Logic:NewGameLogic(),
		RankByAgentFunc:func(state gothello.State) (game.RankByAgent[gothello.Turn], error) {
			blackLegalCount := state.Blacks.Legals(state.Whites)
			whiteLegalCount := state.Whites.Legals(state.Blacks)
			if blackLegalCount != 0 || whiteLegalCount != 0 {
				return nil, nil
			}

			blackCount := state.Blacks.OnesCount()
			whiteCount := state.Whites.OnesCount()
			var blackRank int
			var whiteRank int

			if blackCount > whiteCount {
				blackRank = 1
				whiteRank = 2
			} else if blackCount < whiteCount {
				blackRank = 2
				whiteRank = 1
			} else {
				blackRank = 1
				whiteRank = 1
			}
			return game.RankByAgent[gothello.Turn]{
				gothello.BlackTurn:blackRank,
				gothello.WhiteTurn:whiteRank,
			}, nil
		},
		Agents:[]gothello.Turn{gothello.BlackTurn, gothello.WhiteTurn},
	}
	e.SetStandardResultScoreByAgentFunc()
	return e
}