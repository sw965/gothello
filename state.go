package gothello

import (
	"fmt"
)

type State struct {
	Blacks BitBoard
	Whites BitBoard
	turn   Turn
}

func NewInitState() State {
	blacks := BitBoard(0b00000000_00000000_00000000_00001000_00010000_00000000_00000000_00000000)
	whites := BitBoard(0b00000000_00000000_00000000_00010000_00001000_00000000_00000000_00000000)
	return State{Blacks: blacks, Whites: whites, turn: BlackTurn}
}

func (s State) Turn() Turn {
	return s.turn
}

func (s State) Empties() BitBoard {
	return ^(s.Blacks | s.Whites)
}

func (s State) Legals() BitBoard {
	if s.turn == BlackTurn {
		return s.Blacks.Legals(s.Whites)
	} else {
		return s.Whites.Legals(s.Blacks)
	}
}

func (s State) Move(move BitBoard) (State, error) {
	var self BitBoard
	var opp BitBoard

	if s.turn == BlackTurn {
		self = s.Blacks
		opp = s.Whites
	} else if s.turn == WhiteTurn {
		self = s.Whites
		opp = s.Blacks
	} else {
		return State{}, fmt.Errorf("State.Turn == Empty")
	}

	flips := self.Flips(opp, move)
	if flips == 0 {
		return State{}, fmt.Errorf("非合法の手を打とうとした。")
	}

	//石を置いて、ひっくり返す。
	self |= move | flips
	//ひっくり返される石を消す。
	opp &^= flips

	if s.turn == BlackTurn {
		s.Blacks = self
		s.Whites = opp
	} else {
		s.Whites = self
		s.Blacks = opp
	}

	if opp.Legals(self) != 0 {
		s.turn = s.turn.Opposite()
	}
	return s, nil
}