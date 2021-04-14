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
	s0.Set(19, 19, ColorRed, true)
	s0.SetNotPlayedPiecesFor(ColorRed, []Piece{PieceTetroO})
	var s1 BasicState
	s1.Set(0, 19, ColorRed, true)
	s1.Set(1,18, ColorRed, true)
	s1.Set(0, 17, ColorGreen, true)
	s1.SetNotPlayedPiecesFor(ColorRed, []Piece{PieceMono})
	s2 := earlyTestState()
	tpTetroO := NewTransformedPiece(PieceTetroO, RotationNone, false)
	tpMono := NewTransformedPiece(PieceMono, RotationNone, false)
	// tpDomino := NewTransformedPiece(PieceDomino, RotationNone, false)
	// tpDominoI := NewTransformedPiece(PieceDomino, RotationRight, false)
	var cases = []struct {
		s State
		e []Move
		l int // expectend length to avoid listing all moves
	}{
		{ s: &s0, e: []Move {NewMove(tpTetroO, 17, 17)}},
		{ s: &s1, e: []Move {NewMove(tpMono, 2, 17), NewMove(tpMono, 2, 19)}},
		{ s: s2, l: 311},
	}

	for i, tc := range cases {
		o := PossibleNextMoves(tc.s, ColorRed)
		if tc.e == nil {
			if tc.l != len(o) {
				t.Errorf("case %d failed. Expected %d moves but got %d", i, tc.l, len(o))
			}
		} else if !MovesEqual(tc.e, o) {
			t.Errorf("case %d failed.\nexpected:\n%s\ngot:\n%s", i, FormatPrettyMoves(tc.e, 'X', "  "), FormatPrettyMoves(o, 'X', "  "))
		}
	}
}

func BenchmarkPossibleNextMovesEarly(b *testing.B) {
	b.StopTimer()
	s := earlyTestState()
	b.StartTimer()
	benchmarkNextMovesImpl(b, s, ColorBlue, PossibleNextMoves)
}

func BenchmarkPossibleNextMovesSimpleEarly(b *testing.B) {
	b.StopTimer()
	s := earlyTestState()
	b.StartTimer()
	benchmarkNextMovesImpl(b, s, ColorBlue, possibleNextMovesSimple)
}

func benchmarkNextMovesImpl(b *testing.B, s State, c Color, f func(s State, c Color) []Move) {
	b.Helper()
	for n := 0; n < b.N; n++ {
		_ = f(s, c)
	}
}

func earlyTestState() State {
	s := new(BasicState)

	MustApplyMove(s, ColorBlue, NewMove(NewTransformedPiece(PieceTetroO, RotationNone, false), 0, 0))
	MustApplyMove(s, ColorYellow, NewMove(NewTransformedPiece(PieceTetroO, RotationNone, false), 18, 0))
	MustApplyMove(s, ColorRed, NewMove(NewTransformedPiece(PieceTetroO, RotationNone, false), 0, 18))
	MustApplyMove(s, ColorGreen, NewMove(NewTransformedPiece(PieceTetroO, RotationNone, false), 18, 18))

	MustApplyMove(s, ColorBlue, NewMove(NewTransformedPiece(PiecePentoW, RotationNone, false), 2, 2))
	MustApplyMove(s, ColorYellow, NewMove(NewTransformedPiece(PiecePentoI, RotationRight, false), 13, 2))
	MustApplyMove(s, ColorRed, NewMove(NewTransformedPiece(PiecePentoL, RotationLeft, false), 2, 16))
	MustApplyMove(s, ColorGreen, NewMove(NewTransformedPiece(PiecePentoY, RotationRight, true), 14, 16))
	return s
}