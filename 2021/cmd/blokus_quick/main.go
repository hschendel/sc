package main

import (
	"crypto/rand"
	"fmt"
	"github.com/hschendel/sc/2021/blokus"
	"log"
	"math/big"
	"net"
	"os"
	"sort"
	"strconv"
	"time"
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
	client, err := blokus.OpenClient(addr, &quickPlayer{}, &blokus.BasicState{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot connect to server: %s", err)
		os.Exit(1)
	}
	// client.DebugTo = os.Stderr
	err = client.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while running: %s", err)
		os.Exit(2)
	}
}

type quickPlayer struct{}

func (r *quickPlayer) FirstMove(state blokus.State, color blokus.Color, startPiece blokus.Piece, timeout blokus.Timeout) blokus.Move {
	moves := blokus.PossibleFirstMoves(state, startPiece)
	move := pickBestMove(state, color, moves)
	return move
}

func (r *quickPlayer) NextMove(state blokus.State, color blokus.Color, timeout blokus.Timeout) blokus.Move {
	moves := blokus.PossibleNextMoves(state, color)
	move := pickBestMove(state, color, moves)
	return move
}

func (r *quickPlayer) End() {
}

func pickBestMove(s blokus.State, c blokus.Color, moves []blokus.Move) blokus.Move {
	log.Printf("Pick move for %s", c.String())
	var ms blokus.BasicState
	blokus.CopyState(&ms, s)
	ratedMoves := make(sortRatedMoves, 0, len(moves))
	startRateT := time.Now()
	for _, move := range moves {
		ratedMoves = append(ratedMoves, rateMove(&ms, c, move))
	}
	rateD := time.Since(startRateT)
	startSortT := time.Now()
	sort.Sort(ratedMoves)
	sortD := time.Since(startSortT)
	log.Printf("took %s to rate all %d moves, and %s to sort them.", rateD.String(), len(ratedMoves), sortD.String())
	sameRatingIdx := 0
	for i, ratedMove := range ratedMoves {
		if i == 0 {
			continue
		}
		if ratedMove.move.Transformation.Piece().NumPoints() < ratedMoves[0].move.Transformation.Piece().NumPoints() {
			break
		}
		if ratedMove.stateRating == ratedMoves[0].stateRating {
			sameRatingIdx = i
		} else {
			break
		}
	}
	moveIdx := randomInt(sameRatingIdx + 1)
	move := ratedMoves[moveIdx].move
	log.Printf("Picked move: %s\n  volume diff: %d\n  field coverage diff: %d\n", move.FormatPretty('X', "  "), ratedMoves[moveIdx].volumeDiff, ratedMoves[moveIdx].countDiff)
	return move
}

type sortRatedMoves []ratedMove

func (s sortRatedMoves) Len() int {
	return len(s)
}

func (s sortRatedMoves) Less(i, j int) bool {
	// volume covered = distance to center?
	if s[i].volumeDiff != s[j].volumeDiff {
		return s[i].volumeDiff > s[j].volumeDiff
	}
	// just how many fields are filled
	return s[i].countDiff > s[j].countDiff
}

func (s sortRatedMoves) Swap(i, j int) {
	tmp := s[i]
	s[j] = s[i]
	s[i] = tmp
}

func rateMove(s blokus.MutableState, c blokus.Color, m blokus.Move) (r ratedMove) {
	if err := blokus.ApplyMove(s, c, m); err != nil {
		panic(err)
	}
	r.move = m
	r.stateRating = rateState(s, c)
	blokus.UndoMove(s, c, m)
	return
}

type ratedMove struct {
	stateRating
	move blokus.Move
}

type stateRating struct {
	countDiff  int64
	volumeDiff int64
}

func rateState(s blokus.State, c blokus.Color) (rating stateRating) {
	enemyColors := blokus.EnemyColors(c)
	ownColors := blokus.OwnColors(c)
	volume, count := countExtent(s)
	rating.countDiff = count[ownColors[0]] + count[ownColors[1]] - count[enemyColors[0]] - count[enemyColors[1]]
	rating.volumeDiff = volume[ownColors[0]] + volume[ownColors[1]] - volume[enemyColors[0]] - volume[enemyColors[1]]
	return
}

func countExtent(s blokus.State) (v, n [4]int64) {
	var minX, minY, maxX, maxY [4]uint8
	for x := uint8(0); x < 20; x++ {
		for y := uint8(0); y < 20; y++ {
			c, hasPiece := s.At(x, y)
			if !hasPiece {
				continue
			}
			n[c] += 1
			if x < minX[c] {
				minX[c] = x
			}
			if y < minY[c] {
				minY[c] = y
			}
			if x > maxX[c] {
				maxX[c] = x
			}
			if y > maxY[c] {
				maxY[c] = y
			}
		}
	}
	for c := uint8(0); c < 4; c++ {
		v[c] = int64(maxX[c]-minX[c]) * int64(maxY[c]-minY[c])
	}
	return
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
