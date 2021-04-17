package example_players

import (
	"github.com/hschendel/sc"
	"github.com/hschendel/sc/2021/blokus"
	"log"
	"sort"
	"time"
)

// RestrictingPlayer tries to pick the move most restricting the enemy's moves.
type RestrictingPlayer struct{}

func (rp *RestrictingPlayer) NextMove(state blokus.State, color blokus.Color, timeout sc.Timeout) blokus.Move {
	moves := blokus.PossibleNextMoves(state, color)
	presortMoves(state, color, moves)
	move := rp.pickBestMove(state, color, moves)
	return move
}

func (rp *RestrictingPlayer) End() {
}

func (rp *RestrictingPlayer) pickBestMove(s blokus.State, c blokus.Color, moves []blokus.Move) blokus.Move {
	log.Printf("Pick move for %s", c.String())
	var ms blokus.BasicState
	blokus.CopyState(&ms, s)
	var ratedMoves sortRatedMovesRestricting
	startRateT := time.Now()
	const rateTimeout = time.Millisecond * 1500
	for _, move := range moves {
		ratedMoves.m = append(ratedMoves.m, rp.rateMove(&ms, c, move))
		if time.Since(startRateT) > rateTimeout {
			log.Println("reached rating timeout")
			break
		}
	}
	rateD := time.Since(startRateT)
	startSortT := time.Now()
	sort.Sort(&ratedMoves)
	sortD := time.Since(startSortT)
	log.Printf("took %s to rate %d/%d moves, and %s to sort them.", rateD.String(), len(ratedMoves.m), len(moves), sortD.String())
	sameRatingIdx := 0
	log.Printf("best moves:\n%s", ratedMoves.m[0].move.FormatPretty('X', "  "))
	for i, ratedMove := range ratedMoves.m {
		if i == 0 {
			continue
		}
		if ratedMove.move.Transformation.Piece().NumPoints() < ratedMoves.m[0].move.Transformation.Piece().NumPoints() {
			break
		}
		if ratedMove.stateRatingRestricting == ratedMoves.m[0].stateRatingRestricting {
			log.Print(ratedMove.move.FormatPretty('X', "  "))
			sameRatingIdx = i
		} else {
			break
		}
	}
	if sameRatingIdx > 0 {
		log.Printf("choosing randomly between %d moves", sameRatingIdx+1)
	}
	moveIdx := blokus.RandomInt(sameRatingIdx + 1)
	move := ratedMoves.m[moveIdx].move
	log.Printf("Picked move: %s\n  enemy moves possible: %d\nown extend: %d\n", move.FormatPretty('X', "  "), ratedMoves.m[moveIdx].enemyMoves, ratedMoves.m[moveIdx].volumeDiff)
	return move
}

func presortMoves(s blokus.State, c blokus.Color, moves []blokus.Move) {
	var sm sortMovesRestricting
	sm.m = moves
	sm.goDown, sm.goRight = blokus.ColorDirection(s, c)
	sort.Sort(&sm)
}

type sortMovesRestricting struct {
	m       []blokus.Move
	goDown  bool
	goRight bool
}

func (s *sortMovesRestricting) Len() int {
	return len(s.m)
}

func (s *sortMovesRestricting) Less(i, j int) bool {
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

func (s sortMovesRestricting) Swap(i, j int) {
	tmp := s.m[i]
	s.m[j] = s.m[i]
	s.m[i] = tmp
}

type sortRatedMovesRestricting struct {
	m []ratedMoveRestricting
}

func (s *sortRatedMovesRestricting) Len() int {
	return len(s.m)
}

func (s *sortRatedMovesRestricting) Less(i, j int) bool {
	// 1) least possible enemy moves
	if s.m[i].enemyMoves != s.m[j].enemyMoves {
		return s.m[i].enemyMoves < s.m[j].enemyMoves
	}
	// 2) piece size (larger first)
	if s.m[i].move.Transformation.Piece().NumPoints() != s.m[j].move.Transformation.Piece().NumPoints() {
		return s.m[i].move.Transformation.Piece().NumPoints() > s.m[j].move.Transformation.Piece().NumPoints()
	}
	// 3) height * width of own color
	areaI := s.m[i].area()
	areaJ := s.m[j].area()
	if areaI != areaJ {
		return areaI > areaJ
	}
	// 4) own volume extend
	return s.m[i].volumeDiff > s.m[j].volumeDiff
}

func (s sortRatedMovesRestricting) Swap(i, j int) {
	tmp := s.m[i]
	s.m[j] = s.m[i]
	s.m[i] = tmp
}

func (rp *RestrictingPlayer) rateMove(s blokus.MutableState, c blokus.Color, m blokus.Move) (r ratedMoveRestricting) {
	if err := blokus.ApplyMove(s, c, m); err != nil {
		panic(err)
	}
	r.move = m
	r.stateRatingRestricting = rateState(s, c)
	blokus.UndoMove(s, c, m)
	return
}

type ratedMoveRestricting struct {
	stateRatingRestricting
	move blokus.Move
}

type stateRatingRestricting struct {
	enemyMoves uint64
	volumeDiff uint64
	height     uint8
	width      uint8
}

func (r *stateRatingRestricting) area() uint64 {
	return uint64(r.height) * uint64(r.width)
}

func rateState(s blokus.State, c blokus.Color) (rating stateRatingRestricting) {
	enemyColors := blokus.EnemyColors(c)
	ownColors := blokus.OwnColors(c)
	rating.enemyMoves = uint64(len(blokus.PossibleNextMoves(s, enemyColors[0]))) + uint64(len(blokus.PossibleNextMoves(s, enemyColors[1])))
	volume, height, width := volumeExtend(s)
	rating.volumeDiff = volume[ownColors[0]] + volume[ownColors[1]] - volume[enemyColors[0]] - volume[enemyColors[1]]
	rating.height = height[c]
	rating.width = width[c]
	return
}

func volumeExtend(s blokus.State) (v [4]uint64, h, w [4]uint8) {
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
		h[c] = maxY[c] - minY[c]
		w[c] = maxX[c] - minX[c]
		v[c] = uint64(w[c]) * uint64(h[c])
	}
	return
}
