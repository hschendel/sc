package blokus

import "testing"

func Test_uniquePieceTransformations(t *testing.T) {
	for p, tpl := range uniquePieceTransformations {
		if len(tpl) == 0 {
			t.Errorf("uniquePieceTransformations empty for %s", Piece(p).String())
		}
		for i, tpi := range tpl {
			pi := tpi.Positions()
			for dj, tpj := range tpl[i+1:] {
				pj := tpj.Positions()
				if PositionsEqual(pi, pj) {
					t.Errorf("Piece %s: equal transformations %d (%s) and %d (%s)\n\n%s\n\n%s", Piece(p).String(), i, tpi.String(), i+dj+1, tpj.String(), tpi.PrettyFormat('I', "  "), tpj.PrettyFormat('J', "  "))
				}
			}
		}
	}
}

func Test_rotateLeft(t *testing.T) {
	cases := []struct {
		p Piece
		e []Position
	}{
		{
			p: PieceDomino,
			e: []Position{{0, 0}, {0, 1}},
		},
		{
			p: PieceTrioL,
			e: []Position{{1, 0}, {0, 1}, {1, 1}},
		},
	}
	for i, tc := range cases {
		o := rotateLeft(tc.p.Points())
		if !PositionsEqual(tc.e, o) {
			t.Errorf("case %d failed.\nexpected:\n%#v\ngot:\n%#v", i, tc.e, o)
		}
	}
}

func Test_rotateRight(t *testing.T) {
	cases := []struct {
		p Piece
		e []Position
	}{
		{
			p: PieceDomino,
			e: []Position{{0, 0}, {0, 1}},
		},
		{
			p: PieceTrioL,
			e: []Position{{0, 0}, {1, 0}, {0, 1}},
		},
	}
	for i, tc := range cases {
		o := rotateRight(tc.p.Points())
		if !PositionsEqual(tc.e, o) {
			t.Errorf("case %d failed.\nexpected:\n%#v\ngot:\n%#v", i, tc.e, o)
		}
	}
}

func Test_mirror(t *testing.T) {
	cases := []struct {
		p Piece
		e []Position
	}{
		{
			p: PieceDomino,
			e: []Position{{0, 0}, {1, 0}},
		},
		{
			p: PieceTrioL,
			e: []Position{{0, 0}, {1, 0}, {1, 1}},
		},
	}
	for i, tc := range cases {
		o := mirror(tc.p.Points())
		if !PositionsEqual(tc.e, o) {
			t.Errorf("case %d failed.\nexpected:\n%#v\ngot:\n%#v", i, tc.e, o)
		}
	}
}

func Test_flipPositions(t *testing.T) {
	cases := []struct {
		i []Position
		e []Position
	}{
		{
			i: []Position{{0, 0}, {1, 0}},
			e: []Position{{0, 0}, {1, 0}},
		},
		{
			i: []Position{{0, 0}, {1, 0}, {1, 1}},
			e: []Position{{0, 0}, {1, 0}, {0, 1}},
		},
	}
	for i, tc := range cases {
		o := flipPositions(tc.i)
		if !PositionsEqual(tc.e, o) {
			t.Errorf("case %d failed.\nexpected:\n%#v\ngot:\n%#v", i, tc.e, o)
		}
	}
}
