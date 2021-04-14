package blokus

import "testing"

func TestPossibleFirstMoves(t *testing.T) {
	cases := []struct {
		fp Piece
		e  []Move
	}{
		{
			fp: PieceMono,
			e: []Move{
				NewMove(NewTransformedPiece(PieceMono, RotationNone, false), 0, 0),
				NewMove(NewTransformedPiece(PieceMono, RotationNone, false), 19, 0),
				NewMove(NewTransformedPiece(PieceMono, RotationNone, false), 19, 19),
				NewMove(NewTransformedPiece(PieceMono, RotationNone, false), 0, 19),
			},
		},
		{
			fp: PieceDomino,
			e: []Move{
				NewMove(NewTransformedPiece(PieceDomino, RotationNone, false), 0, 0),
				NewMove(NewTransformedPiece(PieceDomino, RotationRight, false), 0, 0),
				NewMove(NewTransformedPiece(PieceDomino, RotationNone, false), 18, 0),
				NewMove(NewTransformedPiece(PieceDomino, RotationRight, false), 19, 0),
				NewMove(NewTransformedPiece(PieceDomino, RotationNone, false), 18, 19),
				NewMove(NewTransformedPiece(PieceDomino, RotationRight, false), 19, 18),
				NewMove(NewTransformedPiece(PieceDomino, RotationNone, false), 0, 19),
				NewMove(NewTransformedPiece(PieceDomino, RotationRight, false), 0, 18),
			},
		},
		{
			fp: PieceTrioL,
			e: []Move{
				NewMove(NewTransformedPiece(PieceTrioL, RotationNone, false), 0, 0),
				NewMove(NewTransformedPiece(PieceTrioL, RotationNone, false), 18, 18),
				NewMove(NewTransformedPiece(PieceTrioL, RotationNone, false), 0, 18),
				NewMove(NewTransformedPiece(PieceTrioL, RotationLeft, false), 18, 0),
				NewMove(NewTransformedPiece(PieceTrioL, RotationLeft, false), 18, 18),
				NewMove(NewTransformedPiece(PieceTrioL, RotationLeft, false), 0, 18),
				NewMove(NewTransformedPiece(PieceTrioL, RotationRight, false), 0, 0),
				NewMove(NewTransformedPiece(PieceTrioL, RotationRight, false), 18, 0),
				NewMove(NewTransformedPiece(PieceTrioL, RotationRight, false), 0, 18),
				NewMove(NewTransformedPiece(PieceTrioL, RotationMirror, false), 0, 0),
				NewMove(NewTransformedPiece(PieceTrioL, RotationMirror, false), 18, 0),
				NewMove(NewTransformedPiece(PieceTrioL, RotationMirror, false), 18, 18),
			},
		},
	}

	var s BasicState
	for i, tc := range cases {
		o := PossibleFirstMoves(&s, tc.fp)
		if !MovesEqual(tc.e, o) {
			t.Errorf("case %d for start piece %s: failed\nExpected moves:\n%s\nGot moves:\n%s", i, tc.fp.String(), FormatPrettyMoves(tc.e, 'E', "  "), FormatPrettyMoves(o, 'G', "  "))
		}
	}
}

func TestPossibleNextMoves(t *testing.T) {
	// color is always ColorRed
	var s0 BasicState
	s0.Set(0, 0, ColorRed, true)
	s0.SetNotPlayedPiecesFor(ColorRed, []Piece{PieceTetroO})
	var s1 BasicState
	s1.Set(0, 0, ColorRed, true)
	s1.Set(1,1, ColorRed, true)
	s1.Set(0, 2, ColorGreen, true)
	s1.SetNotPlayedPiecesFor(ColorRed, []Piece{PieceMono})
	tpTetroO := NewTransformedPiece(PieceTetroO, RotationNone, false)
	tpMono := NewTransformedPiece(PieceMono, RotationNone, false)
	var cases = []struct {
		s State
		e []Move
	}{
		{ &s0, []Move {NewMove(tpTetroO, 1, 1)}},
		{ &s1, []Move {NewMove(tpMono, 2, 2), NewMove(tpMono, 2, 0)}},
	}

	for i, tc := range cases {
		o := PossibleNextMoves(tc.s, ColorRed)
		if !MovesEqual(tc.e, o) {
			t.Errorf("case %d: failed.\nexpected:\n%s\ngot:\n%s", i, FormatPrettyMoves(tc.e, 'X', "  "), FormatPrettyMoves(tc.e, 'X', "  "))
		}
	}
}
