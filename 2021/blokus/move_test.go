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
