package example_players

import (
	"github.com/hschendel/sc/2021/blokus"
)

type RandomPlayer struct{}

func (r *RandomPlayer) FirstMove(state blokus.State, color blokus.Color, startPiece blokus.Piece, timeout blokus.Timeout) blokus.Move {
	moves := blokus.PossibleFirstMoves(state, startPiece)
	move := randomMove(moves)
	return move
}

func (r *RandomPlayer) NextMove(state blokus.State, color blokus.Color, timeout blokus.Timeout) blokus.Move {
	moves := blokus.PossibleNextMoves(state, color)
	move := randomMoveOrEmpty(moves)
	return move
}

func (r *RandomPlayer) End() {
}

