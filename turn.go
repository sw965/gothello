package gothello

type Turn int

const (
	BlackTurn Turn = 1
	WhiteTurn Turn = 2
)

func (t Turn) Opposite() Turn {
	if t == BlackTurn {
		return WhiteTurn
	}
	return BlackTurn
}