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

type BasicState struct {
	board           [20][20]boardValue
	notPlayedPieces [4][]Piece
	lastMoveMono    [4]bool
}

func (b *BasicState) At(x, y uint8) (c Color, hasPiece bool) {
	if x > 19 || y > 19 {
		return
	}
	bv := b.board[x][y]
	c = bv.color
	hasPiece = bv.hasPiece
	return
}

func (b *BasicState) ensureNotPlayedPieces() {
	if b.notPlayedPieces[0] == nil {
		for c := Color(0); c < 4; c++ {
			b.notPlayedPieces[c] = make([]Piece, NumPieces)
			copy(b.notPlayedPieces[c], AllPieces[:])
		}
	}
}
func (b *BasicState) NotPlayedPiecesFor(c Color) []Piece {
	b.ensureNotPlayedPieces()
	return b.notPlayedPieces[c]
}

func (b *BasicState) IsPiecePlayed(c Color, p Piece) bool {
	b.ensureNotPlayedPieces()
	for _, pp := range b.notPlayedPieces[c] {
		if pp == p {
			return false
		}
	}
	return true
}

func (b *BasicState) IsLastMoveMono(c Color) bool {
	return b.lastMoveMono[c]
}

func (b *BasicState) HasPlayed(c Color) bool {
	return len(b.notPlayedPieces[c]) < NumPieces
}

func (b *BasicState) Reset() {
	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			b.board[x][y] = boardValue{}
		}
	}
	for c := ColorBlue; c < 4; c++ {
		b.notPlayedPieces[c] = make([]Piece, NumPieces)
		copy(b.notPlayedPieces[c], AllPieces[:])
		b.lastMoveMono[c] = false
	}
}

func (b *BasicState) Set(x, y uint8, c Color, hasPiece bool) {
	b.board[x][y] = boardValue{
		color:    c,
		hasPiece: hasPiece,
	}
}

func (b *BasicState) SetNotPlayedPiecesFor(c Color, pieces []Piece) {
	b.ensureNotPlayedPieces()
	b.notPlayedPieces[c] = b.notPlayedPieces[c][0:0]
	for _, p := range pieces {
		b.notPlayedPieces[c] = append(b.notPlayedPieces[c], p)
	}
}

func (b *BasicState) SetPiecePlayed(c Color, p Piece, isPlayed bool) {
	b.ensureNotPlayedPieces()
	l := len(b.notPlayedPieces[c])
	insertI := 0
	for i, pi := range b.notPlayedPieces[c] {
		if pi == p {
			if isPlayed {
				if (i + 1) < l {
					copy(b.notPlayedPieces[c][i:l-1], b.notPlayedPieces[c][i+1:l])
				}
				b.notPlayedPieces[c] = b.notPlayedPieces[c][0 : l-1]
			}
			return
		}
		if p.NumPoints() > pi.NumPoints() {
			insertI = i
		}
	}
	if isPlayed {
		return
	}
	b.notPlayedPieces[c] = append(b.notPlayedPieces[c], p)
	copy(b.notPlayedPieces[c][insertI+1:], b.notPlayedPieces[c][insertI:l])
	b.notPlayedPieces[c][insertI] = p
}

func (b *BasicState) SetLastMoveMono(c Color, isLastMoveMono bool) {
	b.lastMoveMono[c] = isLastMoveMono
}

type boardValue struct {
	color    Color
	hasPiece bool
}
