package blokus

import (
	"github.com/hschendel/sc"
)

// Player is the interface a Blokus player must implement
type Player interface {
	// NextMove must return the player's next move before timeout.Reached() == true
	NextMove(state State, color Color, timeout sc.Timeout) Move
	// End is called when the game has ended, so the player can stop any ongoing calculations.
	End()
}
