package example_players

import (
	"crypto/rand"
	"fmt"
	"github.com/hschendel/sc/2021/blokus"
	"math/big"
)

func randomMove(moves []blokus.Move) blokus.Move {
	if len(moves) == 0 {
		panic("len(moves) == 0")
	}
	i := randomInt(len(moves))
	return moves[i]
}

func randomMoveOrEmpty(moves []blokus.Move) blokus.Move {
	i := randomInt(len(moves) + 1)
	if i == len(moves) {
		return blokus.EmptyMove
	}
	return moves[i]
}

func randomInt(max int) int {
	if max <= 1 {
		return 0
	}
	var bmax big.Int
	bmax.SetInt64(int64(max))
	bn, err := rand.Int(rand.Reader, &bmax)
	if err != nil {
		panic(fmt.Errorf("cannot read from random device: %s", err))
	}
	return int(bn.Int64())
}
