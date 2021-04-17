package blokus

import (
	"github.com/hschendel/sc"
)

type Player interface {
	FirstMove(state State, color Color, startPiece Piece, timeout sc.Timeout) Move
	NextMove(state State, color Color, timeout sc.Timeout) Move
	End()
}
