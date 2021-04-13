package main

import (
	"crypto/rand"
	"fmt"
	"github.com/hschendel/sc/2021/blokus"
	"math/big"
	"net"
	"os"
	"strconv"
)

func main() {
	addr := blokus.DefaultServerAddress
	if len(os.Args) > 1 {
		port, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid value for port: %q", os.Args[1])
		}
		addr = &net.TCPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: port,
		}
	}
	client, err := blokus.OpenClient(addr, &randomPlayer{}, &blokus.BasicState{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot connect to server: %s", err)
		os.Exit(1)
	}
	client.DebugTo = os.Stderr
	err = client.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while running: %s", err)
		os.Exit(2)
	}
}

type randomPlayer struct {}

func (r *randomPlayer) FirstMove(state blokus.State, color blokus.Color, startPiece blokus.Piece, timeout blokus.Timeout) blokus.Move {
	moves := blokus.PossibleFirstMoves(state, startPiece)
	move := randomMove(moves)
	return move
}

func (r *randomPlayer) NextMove(state blokus.State, color blokus.Color, timeout blokus.Timeout) blokus.Move {
	moves := blokus.PossibleNextMoves(state, color)
	move := randomMoveOrEmpty(moves)
	return move
}

func (r *randomPlayer) End() {
}

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