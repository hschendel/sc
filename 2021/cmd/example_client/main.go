package main

import (
	"github.com/hschendel/sc"
	"github.com/hschendel/sc/2021/blokus"
)

func main() {
	player := new(simplePlayer)
	blokus.ClientMain(player)
}

type simplePlayer struct{}

func (s simplePlayer) NextMove(state blokus.State, color blokus.Color, timeout sc.Timeout) blokus.Move {
	moves := blokus.PossibleNextMoves(state, color)
	return moves[0]
}

func (s simplePlayer) End() {
}
