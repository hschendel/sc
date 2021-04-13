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

func (r *quickPlayer) FirstMove(state blokus.State, color blokus.Color, startPiece blokus.Piece, timeout blokus.Timeout) (move blokus.Move) {
	moves := blokus.PossibleFirstMoves(state, startPiece)
	bestVolume := uint8(0)
	for _, m := range moves {
		moveVol := m.Transformation.Width() * m.Transformation.Height()
		if moveVol > bestVolume {
			move = m
			bestVolume = moveVol
		}
	}
	return
}

func (r *quickPlayer) NextMove(state blokus.State, color blokus.Color, timeout blokus.Timeout) blokus.Move {
	moves := blokus.PossibleNextMoves(state, color)
	presortMoves(state, color, moves)
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
	const rateTimeout = time.Millisecond * 1500
	for _, move := range moves {
		ratedMoves = append(ratedMoves, rateMove(&ms, c, move))
		if time.Since(startRateT) > rateTimeout {
			log.Println("reached rating timeout")
			break
		}
	}
	rateD := time.Since(startRateT)
	startSortT := time.Now()
	sort.Sort(ratedMoves)
	sortD := time.Since(startSortT)
	log.Printf("took %s to rate %d/%d moves, and %s to sort them.", rateD.String(), len(ratedMoves), len(moves), sortD.String())
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
	log.Printf("Picked move: %s\n  enemy moves possible: %d\nown extend: %d\n", move.FormatPretty('X', "  "), ratedMoves[moveIdx].enemyMoves, ratedMoves[moveIdx].volumeDiff)
	return move
}

func presortMoves(s blokus.State, c blokus.Color, moves []blokus.Move) {
	var sm sortMoves
	sm.m = moves
	for _, pos := range blokus.StartCorners {
		if cc, found := s.At(pos.X, pos.Y); found && cc == c {
			if pos.X == 0 {
				sm.goRight = true
			}
			if pos.Y == 0 {
				sm.goDown = true
			}
			break
		}
	}
	sort.Sort(&sm)
}

type sortMoves struct {
	m       []blokus.Move
	goDown  bool
	goRight bool
}

func (s *sortMoves) Len() int {
	return len(s.m)
}

func (s *sortMoves) Less(i, j int) bool {
	if s.m[i].Transformation.Piece().NumPoints() != s.m[j].Transformation.Piece().NumPoints() {
		return s.m[i].Transformation.Piece().NumPoints() > s.m[j].Transformation.Piece().NumPoints()
	}
	if s.goDown {
		iy := s.m[i].Y + s.m[i].Transformation.Height() - 1
		jy := s.m[j].Y + s.m[j].Transformation.Height() - 1
		if iy != jy {
			return iy > jy
		}
	} else {
		if s.m[i].Y != s.m[j].Y {
			return s.m[i].Y < s.m[j].Y
		}
	}
	if s.goRight {
		ix := s.m[i].X + s.m[i].Transformation.Width() - 1
		jx := s.m[j].X + s.m[j].Transformation.Width() - 1
		if ix != jx {
			return ix > jx
		}
	} else {
		if s.m[i].X != s.m[j].X {
			return s.m[i].X < s.m[j].X
		}
	}
	return false
}

func (s sortMoves) Swap(i, j int) {
	tmp := s.m[i]
	s.m[j] = s.m[i]
	s.m[i] = tmp
}

type sortRatedMoves []ratedMove

func (s sortRatedMoves) Len() int {
	return len(s)
}

func (s sortRatedMoves) Less(i, j int) bool {
	// 1) least possible enemy moves
	if s[i].enemyMoves != s[j].enemyMoves {
		return s[i].enemyMoves < s[j].enemyMoves
	}
	// 2) piece size (larger first)
	if s[i].move.Transformation.Piece().NumPoints() != s[j].move.Transformation.Piece().NumPoints() {
		return s[i].move.Transformation.Piece().NumPoints() > s[j].move.Transformation.Piece().NumPoints()
	}
	// 3) own volume extend
	return s[i].volumeDiff > s[j].volumeDiff
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
	enemyMoves uint64
	volumeDiff uint64
}

func rateState(s blokus.State, c blokus.Color) (rating stateRating) {
	enemyColors := blokus.EnemyColors(c)
	ownColors := blokus.OwnColors(c)
	rating.enemyMoves = uint64(len(blokus.PossibleNextMoves(s, enemyColors[0]))) + uint64(len(blokus.PossibleNextMoves(s, enemyColors[1])))
	volume := volumeExtend(s)
	rating.volumeDiff = volume[ownColors[0]] + volume[ownColors[1]] - volume[enemyColors[0]] - volume[enemyColors[1]]
	return
}

func volumeExtend(s blokus.State) (v [4]uint64) {
	var minX, minY, maxX, maxY [4]uint8
	for x := uint8(0); x < 20; x++ {
		for y := uint8(0); y < 20; y++ {
			c, hasPiece := s.At(x, y)
			if !hasPiece {
				continue
			}
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
		v[c] = uint64(maxX[c]-minX[c]) * uint64(maxY[c]-minY[c])
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
