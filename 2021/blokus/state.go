package blokus

type State interface {
	At(x, y uint8) (c Color, hasPiece bool)
	// NotPlayedPiecesFor returns all pieces that have not yet been played by c, in the order
	// that is preferred by the implementation. The order will be used directly by PossibleNextMoves()
	NotPlayedPiecesFor(c Color) []Piece
	IsPiecePlayed(c Color, p Piece) bool
	IsLastMoveMono(c Color) bool
	HasPlayed(c Color) bool
}

type MutableState interface {
	State
	Reset()
	Set(x, y uint8, c Color, hasPiece bool)
	SetNotPlayedPiecesFor(c Color, pieces []Piece)
	SetPiecePlayed(c Color, p Piece, isPlayed bool)
	SetLastMoveMono(c Color, isLastMoveMono bool)
}

// CopyState copies from a State instance to a MutableState instance.
// You probably only want to use this when filling your own MutableState implementation
// at the beginning of the move calculation, but not as part of the move calculation
// algorithm. This is because it is quite inefficient, relying only on the interface
// methods of State and MutableState. In your own implementation you probably could use
// more efficient mechanisms for copying, or you never copy, but always work on one
// state variable.
func CopyState(into MutableState, from State) {
	for x := uint8(0); x < 20; x++ {
		for y := uint8(0); y < 20; y++ {
			color, hasPiece := from.At(x, y)
			into.Set(x, y, color, hasPiece)
		}
	}
	for c := 0; c < 4; c++ {
		color := Color(c)
		into.SetNotPlayedPiecesFor(color, from.NotPlayedPiecesFor(color))
		into.SetLastMoveMono(color, from.IsLastMoveMono(color))
	}
}

// ColorDirection gives the "direction" the color is playing, depending on its
// start corner
func ColorDirection(s State, c Color) (goDown, goRight bool) {
	for _, pos := range StartCorners {
		if cc, found := s.At(pos.X, pos.Y); found && cc == c {
			if pos.X == 0 {
				goRight = true
			}
			if pos.Y == 0 {
				goDown = true
			}
			return
		}
	}
	return
}

// StartCorner determines the start corner of a color
func StartCorner(s State, c Color) (x, y uint8, found bool) {
	for _, pos := range StartCorners {
		if cc, cFound := s.At(pos.X, pos.Y); cFound && cc == c {
			x = pos.X
			y = pos.Y
			found = true
			return
		}
	}
	return
}
