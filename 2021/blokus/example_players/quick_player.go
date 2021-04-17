package example_players

import (
	"github.com/hschendel/sc"
	"github.com/hschendel/sc/2021/blokus"
	"log"
	"sort"
	"time"
)

type QuickPlayer struct{}

func (q *QuickPlayer) NextMove(state blokus.State, color blokus.Color, timeout sc.Timeout) blokus.Move {
	moves := blokus.PossibleNextMoves(state, color)
	move := q.pickBestMove(state, color, moves)
	return move
}

func (q *QuickPlayer) End() {
}

func (q *QuickPlayer) pickBestMove(s blokus.State, c blokus.Color, moves []blokus.Move) blokus.Move {
	log.Printf("Pick move for %s", c.String())
	var ms blokus.BasicState
	blokus.CopyState(&ms, s)
	ratedMoves := make(sortRatedMovesQuick, 0, len(moves))
	startRateT := time.Now()
	for _, move := range moves {
		ratedMoves = append(ratedMoves, rateMoveQuick(&ms, c, move))
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
		if ratedMove.stateRatingQuick == ratedMoves[0].stateRatingQuick {
			sameRatingIdx = i
		} else {
			break
		}
	}
	moveIdx := blokus.RandomInt(sameRatingIdx + 1)
	move := ratedMoves[moveIdx].move
	log.Printf("Picked move: %s\n  volume diff: %d\n  field coverage diff: %d\n", move.FormatPretty('X', "  "), ratedMoves[moveIdx].volumeDiff, ratedMoves[moveIdx].countDiff)
	return move
}

type sortRatedMovesQuick []ratedMoveQuick

func (s sortRatedMovesQuick) Len() int {
	return len(s)
}

func (s sortRatedMovesQuick) Less(i, j int) bool {
	// volume covered = distance to center?
	if s[i].volumeDiff != s[j].volumeDiff {
		return s[i].volumeDiff > s[j].volumeDiff
	}
	// just how many fields are filled
	return s[i].countDiff > s[j].countDiff
}

func (s sortRatedMovesQuick) Swap(i, j int) {
	tmp := s[i]
	s[j] = s[i]
	s[i] = tmp
}

func rateMoveQuick(s blokus.MutableState, c blokus.Color, m blokus.Move) (r ratedMoveQuick) {
	if err := blokus.ApplyMove(s, c, m); err != nil {
		panic(err)
	}
	r.move = m
	r.stateRatingQuick = rateStateQuick(s, c)
	blokus.UndoMove(s, c, m)
	return
}

type ratedMoveQuick struct {
	stateRatingQuick
	move blokus.Move
}

type stateRatingQuick struct {
	countDiff  int64
	volumeDiff int64
}

func rateStateQuick(s blokus.State, c blokus.Color) (rating stateRatingQuick) {
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
