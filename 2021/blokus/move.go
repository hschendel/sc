package blokus

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

type Move struct {
	// Transformation is the rotated and flipped piece that should be set at X,Y. If it is nil, this move is the empty move.
	Transformation TransformedPiece
	X              uint8
	Y              uint8
	IsMove         bool
}

func NewMove(tp TransformedPiece, x, y uint8) Move {
	return Move{
		Transformation: tp,
		X:              x,
		Y:              y,
		IsMove:         true,
	}
}

var EmptyMove = Move{}

func (m Move) IsEmpty() bool {
	return !m.IsMove
}

func (m Move) Equal(o Move) bool {
	if m.IsMove != o.IsMove {
		return false
	}
	if !m.IsMove {
		return true
	}
	if m.X != o.X || m.Y != o.Y {
		return false
	}
	return PositionsEqual(m.Transformation.Positions(), o.Transformation.Positions())
}

// MovesEqual returns true if two slices of distinct moves contain the same moves
func MovesEqual(a, b []Move) bool {
	if len(a) != len(b) {
		return false
	}
	for _, am := range a {
		equals := 0
		for _, bm := range b {
			if am.Equal(bm) {
				equals++
			}
		}
		if equals != 1 {
			return false
		}
	}
	return true
}

func (m Move) FormatPretty(r rune, indent string) string {
	var sb strings.Builder
	sb.WriteString(indent)
	fmt.Fprintf(&sb, "Set %s at x=%d y=%d:\n", m.Transformation.Piece().String(), m.X, m.Y)
	m.Transformation.PrettyPrint(r, indent, &sb)
	return sb.String()
}

func FormatPrettyMoves(moves []Move, r rune, indent string) string {
	var sb strings.Builder
	for _, move := range moves {
		sb.WriteString(move.FormatPretty(r, indent))
	}
	return sb.String()
}

var ErrForbiddenMove = errors.New("move not allowed")

func ApplyMove(s MutableState, c Color, m Move) error {
	if !CanApplyMove(s, c, m) {
		return ErrForbiddenMove
	}
	for dx := uint8(0); dx < m.Transformation.Width(); dx++ {
		for dy := uint8(0); dy < m.Transformation.Height(); dy++ {
			if m.Transformation.Fills(dx, dy) {
				s.Set(m.X+dx, m.Y+dy, c, true)
			}
		}
	}
	s.SetPiecePlayed(c, m.Transformation.Piece(), true)
	if m.Transformation.Piece() == PieceMono {
		s.SetLastMoveMono(c, true)
	}
	return nil
}

func MustApplyMove(s MutableState, c Color, m Move) {
	err := ApplyMove(s, c, m)
	if err != nil {
		log.Printf("invalid move:\n%s", m.FormatPretty('X', "  "))
		panic(err)
	}
}

// UndoMove reverts a move.
// Warning: It does not check whether m has actually been played before and thus can leave s in an invalid state.
func UndoMove(s MutableState, c Color, m Move) {
	for dx := uint8(0); dx < m.Transformation.Width(); dx++ {
		for dy := uint8(0); dy < m.Transformation.Height(); dy++ {
			if m.Transformation.Fills(dx, dy) {
				s.Set(m.X+dx, m.Y+dy, c, false)
			}
		}
	}
	s.SetPiecePlayed(c, m.Transformation.Piece(), false)
	if m.Transformation.Piece() == PieceMono {
		s.SetLastMoveMono(c, false)
	}
}

func CanApplyMove(s State, c Color, m Move) bool {
	if !m.IsMove {
		// skip move only allowed after first move
		return s.HasPlayed(c)
	}
	if len(s.NotPlayedPiecesFor(c)) == 0 {
		return false
	}
	if s.IsPiecePlayed(c, m.Transformation.Piece()) {
		return false
	}
	return CanPlayNextPiece(s, c, m.Transformation, m.X, m.Y)
}

func canPlayFirstPiece(s State, p TransformedPiece, x, y uint8) bool {
	if s.StartPiece() != p.Piece() {
		return false
	}
	fillsStartCorner := false
	for dx := uint8(0); dx < p.Width(); dx++ {
		for dy := uint8(0); dy < p.Height(); dy++ {
			if p.Fills(dx, dy) {
				px, py := x+dx, y+dy
				if _, hasPiece := s.At(px, py); hasPiece {
					// already filled
					return false
				}
				if IsStartCorner(px, py) {
					fillsStartCorner = true
				}
			}
		}
	}
	return fillsStartCorner
}

var StartCorners = []Position{
	{0, 0},
	{19, 0},
	{19, 19},
	{0, 19},
}

func IsStartCorner(x, y uint8) bool {
	return (x == 0 || x == 19) && (y == 0 || y == 19)
}

func CanPlayNextPiece(s State, c Color, p TransformedPiece, x, y uint8) bool {
	if !IsOnBoard(p, x, y) {
		return false
	}
	if !s.HasPlayed(c) {
		return canPlayFirstPiece(s, p, x, y)
	}
	hasAdjacent := false
	for _, piecePos := range p.Positions() {
		px, py := x+piecePos.X, y+piecePos.Y
		if _, hasPiece := s.At(px, py); hasPiece {
			// already filled
			return false
		}
		// must not touch same color
		if HasDirectNeighborWithColor(s, c, px, py) {
			return false
		}
		// at least one pixel must be adjacent to an existing pixel of the same color
		if !hasAdjacent {
			hasAdjacent = HasAdjacentWithColor(s, c, px, py)
		}
	}
	return hasAdjacent
}

func IsOnBoard(p TransformedPiece, x, y uint8) bool {
	return x < 20 && y < 20 && (x+p.Width()-1) < 20 && (y+p.Height()-1) < 20
}

func HasDirectNeighborWithColor(s State, c Color, x, y uint8) bool {
	return HasColorAt(s, c, x-1, y) ||
		HasColorAt(s, c, x, y-1) ||
		HasColorAt(s, c, x+1, y) ||
		HasColorAt(s, c, x, y+1)
}

func HasAdjacentWithColor(s State, c Color, x, y uint8) bool {
	return HasColorAt(s, c, x-1, y-1) ||
		HasColorAt(s, c, x+1, y-1) ||
		HasColorAt(s, c, x-1, y+1) ||
		HasColorAt(s, c, x+1, y+1)
}

func HasColorAt(s State, c Color, x, y uint8) bool {
	if x >= 20 || y >= 20 {
		return false
	}
	filledColor, filled := s.At(x, y)
	if !filled {
		return false
	}
	return filledColor == c
}

func possibleFirstMoves(s State) (moves []Move) {
	for _, tp := range uniquePieceTransformations[s.StartPiece()] {
		for _, cornerPos := range StartCorners {
			if cornerPos.X > 0 {
				cornerPos.X -= tp.Width() - 1
			}
			if cornerPos.Y > 0 {
				cornerPos.Y -= tp.Height() - 1
			}
			if canPlayFirstPiece(s, tp, cornerPos.X, cornerPos.Y) {
				move := NewMove(tp, cornerPos.X, cornerPos.Y)
				moves = append(moves, move)
			}
		}
	}
	return
}

func HasPossibleNextMoves(s State, c Color) bool {
	return len(PossibleNextMoves(s, c)) > 0
}

func PossibleNextMoves(s State, c Color) (moves []Move) {
	if !s.HasPlayed(c) {
		return possibleFirstMoves(s)
	}
	pieces := s.NotPlayedPiecesFor(c)
	scX, scY, started := StartCorner(s, c)
	if !started {
		panic(fmt.Errorf("color %s has not yet started, use PossibleFirstMoves()", c.String()))
	}
	applyRadiusX := applyRadiusInc
	shiftLeft := false
	applyRadiusY := applyRadiusInc
	shiftUp := false
	if scX == 19 {
		applyRadiusX = applyRadiusDec
		shiftLeft = true
	}
	if scY == 19 {
		applyRadiusY = applyRadiusDec
		shiftUp = true
	}

	for radius := uint8(1); radius < 20; radius++ {
		startX, endX, stepX, rX := applyRadiusX(radius)
		startY, _, stepY, rY := applyRadiusY(radius)
		colorFound := false
		for x := startX; x != endX; x += stepX {
			if cc, cFound := s.At(x, rY); cFound && cc == c {
				colorFound = true
			}
			moves = addPlayableMoves(s, pieces, c, shiftLeft, shiftUp, x, rY, moves)
		}
		for y := startY; y != rY; y += stepY {
			if cc, cFound := s.At(rX, y); cFound && cc == c {
				colorFound = true
			}
			moves = addPlayableMoves(s, pieces, c, shiftLeft, shiftUp, rX, y, moves)
		}
		if !colorFound {
			// cut-off: there will be no more moves with a higher radius
			break
		}
	}
	return
}

func addPlayableMoves(s State, pieces []Piece, c Color, shiftLeft, shiftUp bool, x, y uint8, moves []Move) (movesResult []Move) {
	movesResult = moves
	for _, p := range pieces {
		for _, tp := range uniquePieceTransformations[p] {
			tx, ty := x, y
			if shiftLeft {
				tx -= tp.Width() - 1
			}
			if shiftUp {
				ty -= tp.Height() - 1
			}
			if CanPlayNextPiece(s, c, tp, tx, ty) {
				move := NewMove(tp, tx, ty)
				movesResult = append(movesResult, move)
			}
		}
	}
	return
}

func applyRadiusInc(radius uint8) (start, end, step, fixed uint8) {
	start = 0
	end = radius + 1
	step = 1
	fixed = radius
	return
}

func applyRadiusDec(radius uint8) (start, end, step, fixed uint8) {
	start = 19
	end = 19 - radius - 1
	step = 255
	fixed = 19 - radius
	return
}

func possibleNextMovesSimple(s State, c Color) (moves []Move) {
	pieces := s.NotPlayedPiecesFor(c)
	for _, p := range pieces {
		for _, tp := range uniquePieceTransformations[p] {
			for x := uint8(0); x < 20; x++ {
				for y := uint8(0); y < 20; y++ {
					if CanPlayNextPiece(s, c, tp, x, y) {
						move := NewMove(tp, x, y)
						moves = append(moves, move)
					}
				}
			}
		}
	}
	return
}
