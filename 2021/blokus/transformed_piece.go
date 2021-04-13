package blokus

import (
	"fmt"
	"io"
	"strings"
)

type TransformedPiece uint8

const maskPiece = 0b00011111
const shiftPiece = 0
const maskRotation = 0b01100000
const shiftRotation = 5
const maskFlipped = 0b10000000
const shiftFlipped = 7

func NewTransformedPiece(piece Piece, rotation Rotation, flipped bool) TransformedPiece {
	var u8Flipped uint8
	if flipped {
		u8Flipped = 1
	}
	return TransformedPiece(((uint8(piece) << shiftPiece) & maskPiece) |
		((uint8(rotation) << shiftRotation) & maskRotation) |
		((u8Flipped << shiftFlipped) & maskFlipped))
}

func (p *TransformedPiece) String() string {
	flippedS := ""
	if p.Flipped() {
		flippedS = " flipped"
	}
	return fmt.Sprintf("%s %s%s", p.Piece().String(), p.Rotation().String(), flippedS)
}

func (p *TransformedPiece) Piece() Piece {
	return Piece(((*p) & maskPiece) >> shiftPiece)
}

func (p *TransformedPiece) Rotation() Rotation {
	return Rotation(((*p) & maskRotation) >> shiftRotation)
}

func (p *TransformedPiece) Flipped() bool {
	return (((*p) & maskFlipped) >> shiftFlipped) > 0
}

func (p *TransformedPiece) Width() uint8 {
	switch p.Rotation() {
	case RotationNone, RotationMirror:
		return p.Piece().Width()
	default:
		return p.Piece().Height()
	}
}

func (p *TransformedPiece) Height() uint8 {
	switch p.Rotation() {
	case RotationNone, RotationMirror:
		return p.Piece().Height()
	default:
		return p.Piece().Width()
	}
}

func (p *TransformedPiece) Positions() []Position {
	u8Flipped := 0
	if p.Flipped() {
		u8Flipped = 1
	}
	return transformedPieces[p.Piece()][p.Rotation()][u8Flipped]
}

func (p *TransformedPiece) Fills(x, y uint8) bool {
	positions := p.Positions()
	for _, pos := range positions {
		if pos.X == x && pos.Y == y {
			return true
		}
	}
	return false
}

func (p *TransformedPiece) PrettyFormat(r rune, indent string) string {
	var sb strings.Builder
	err := p.PrettyPrint(r, indent, &sb)
	if err != nil {
		panic(err)
	}
	return sb.String()
}

func (p *TransformedPiece) PrettyPrint(r rune, indent string, w io.Writer) (err error) {
	for y := uint8(0); y < p.Height(); y++ {
		if _, err = fmt.Fprint(w, indent); err != nil {
			return
		}
		for x := uint8(0); x < p.Width(); x++ {
			c := ' '
			if p.Fills(x, y) {
				c = r
			}
			if _, err = fmt.Fprintf(w, "%c", c); err != nil {
				return
			}
		}
		if _, err = fmt.Fprintln(w); err != nil {
			return
		}
	}
	return
}

func (p Piece) Transformations() []TransformedPiece {
	return uniquePieceTransformations[p]
}
