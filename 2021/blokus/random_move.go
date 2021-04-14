package blokus

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func RandomMove(moves []Move) Move {
	if len(moves) == 0 {
		panic("len(moves) == 0")
	}
	i := RandomInt(len(moves))
	return moves[i]
}

func RandomMoveOrEmpty(moves []Move) Move {
	i := RandomInt(len(moves) + 1)
	if i == len(moves) {
		return EmptyMove
	}
	return moves[i]
}

func RandomInt(max int) int {
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
