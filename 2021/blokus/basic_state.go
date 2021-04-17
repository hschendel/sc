package blokus

import "fmt"

// BasicState is a quite inefficient way of storing the game state
type BasicState struct {
	board           [20][20]boardValue
	notPlayedPieces [4][]Piece
	lastMoveMono    [4]bool
	startPiece      Piece
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
	b.ensureNotPlayedPieces()
	return len(b.notPlayedPieces[c]) < NumPieces
}

func (b *BasicState) Reset() {
	b.startPiece = PieceMono
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
	if x > 19 || y > 19 {
		panic(fmt.Errorf("trying to set (color=%s, hasPiece=%v) at invalid coordinates x=%d y%d", c.String(), hasPiece, x, y))
	}
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

func (b *BasicState) SetStartPiece(piece Piece) {
	b.startPiece = piece
}

func (b *BasicState) StartPiece() Piece {
	return b.startPiece
}