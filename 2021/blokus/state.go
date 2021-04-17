package blokus

import "testing"

type State interface {
	At(x, y uint8) (c Color, hasPiece bool)
	// NotPlayedPiecesFor returns all pieces that have not yet been played by c, in the order
	// that is preferred by the implementation. The order will be used directly by PossibleNextMoves()
	NotPlayedPiecesFor(c Color) []Piece
	IsPiecePlayed(c Color, p Piece) bool
	IsLastMoveMono(c Color) bool
	HasPlayed(c Color) bool
	// IsPlayerOneFirst is true if the first player is/was the first to move. For initial or reset state, it must be true
	IsPlayerOneFirst() bool
	// IsColorValid is true if the color c is still playing
	IsColorValid(c Color) bool
	// CurrentColor returns the color that is currently playing. For initial or reset State, it must return ColorBlue
	CurrentColor() Color
	// StartPiece returns the game's start piece. For initial or reset State, it must return PieceMono
	StartPiece() Piece
}

type MutableState interface {
	State
	Reset()
	Set(x, y uint8, c Color, hasPiece bool)
	SetNotPlayedPiecesFor(c Color, pieces []Piece)
	SetPiecePlayed(c Color, p Piece, isPlayed bool)
	SetLastMoveMono(c Color, isLastMoveMono bool)
	SetStartPiece(piece Piece)
	SetPlayerOneFirst(isPlayerOneFirst bool)
	SetColorValid(c Color, isValid bool)
	SetCurrentColor(c Color)
}

// CopyState copies from a State instance to a MutableState instance.
// You probably only want to use this when filling your own MutableState implementation
// at the beginning of the move calculation, but not as part of the move calculation
// algorithm. This is because it is quite inefficient, relying only on the interface
// methods of State and MutableState. In your own implementation you probably could use
// more efficient mechanisms for copying, or you never copy, but always work on one
// state variable.
func CopyState(into MutableState, from State) {
	into.SetStartPiece(from.StartPiece())
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

// IsCurrentPlayerOne returns true if the first player is the current player
func IsCurrentPlayerOne(s State) bool {
	return s.CurrentColor() == ColorBlue || s.CurrentColor() == ColorRed
}

// HasGameEnded is true if the game has ended
func HasGameEnded(s State) bool {
	for c := Color(0); c < 4; c++ {
		if s.IsColorValid(c) {
			return false
		}
	}
	return true
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

var benchmarkStateAtResult Color

func BenchmarkStateAt(b *testing.B, s State) {
	b.Helper()
	b.StopTimer()
	positions := []Position{{0, 0}, {19, 19}, {19, 0}, {0, 19}, {9, 10}, {15, 15}}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for _, pos := range positions {
			benchmarkStateAtResult, _ = s.At(pos.X, pos.Y)
		}
	}
}

func BenchmarkMutableStateResetSet(b *testing.B, s MutableState) {
	b.Helper()
	for i := 0; i < b.N; i++ {
		s.Reset()
		applyDummyStateSets(s)
	}
}

func BenchmarkStateIsPiecePlayed(b *testing.B, s MutableState) {
	b.Helper()
	b.StopTimer()
	setDummyPiecesPlayed(s)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		s.IsPiecePlayed(ColorBlue, PieceTrioI)
		s.IsPiecePlayed(ColorYellow, PieceTetroL)
		s.IsPiecePlayed(ColorGreen, PiecePentoT)
		s.IsPiecePlayed(ColorRed, PieceTetroZ)

		s.IsPiecePlayed(ColorBlue, PieceDomino)
		s.IsPiecePlayed(ColorYellow, PieceDomino)
		s.IsPiecePlayed(ColorGreen, PieceDomino)
		s.IsPiecePlayed(ColorRed, PieceDomino)

		s.IsPiecePlayed(ColorBlue, PieceTetroO)
		s.IsPiecePlayed(ColorYellow, PieceTetroL)
		s.IsPiecePlayed(ColorGreen, PieceTetroT)
		s.IsPiecePlayed(ColorRed, PieceTetroI)

		s.IsPiecePlayed(ColorBlue, PieceTetroZ)
		s.IsPiecePlayed(ColorYellow, PiecePentoW)
		s.IsPiecePlayed(ColorGreen, PieceMono)
		s.IsPiecePlayed(ColorRed, PiecePentoU)
	}
}

func setDummyPiecesPlayed(s MutableState) {
	s.SetPiecePlayed(ColorBlue, PiecePentoL, true)
	s.SetPiecePlayed(ColorGreen, PieceDomino, true)
	s.SetPiecePlayed(ColorYellow, PiecePentoW, true)
	s.SetPiecePlayed(ColorRed, PieceTrioL, true)

	s.SetPiecePlayed(ColorRed, PieceTetroO, true)
	s.SetPiecePlayed(ColorBlue, PiecePentoR, true)
	s.SetPiecePlayed(ColorYellow, PieceDomino, true)
	s.SetPiecePlayed(ColorGreen, PieceTetroZ, true)

	s.SetPiecePlayed(ColorYellow, PieceTetroL, true)
	s.SetPiecePlayed(ColorGreen, PiecePentoX, true)
	s.SetPiecePlayed(ColorBlue, PieceDomino, true)
	s.SetPiecePlayed(ColorRed, PieceTrioI, true)
}

var dummyPiecesPlayed = []struct {
	c Color
	p Piece
	e bool
}{
	{ColorBlue, PiecePentoL, true},
	{ColorGreen, PieceDomino, true},
	{ColorYellow, PiecePentoW, true},
	{ColorRed, PieceTrioL, true},

	{ColorRed, PieceTetroO, true},
	{ColorBlue, PiecePentoR, true},
	{ColorYellow, PieceDomino, true},
	{ColorGreen, PieceTetroZ, true},

	{ColorYellow, PieceTetroL, true},
	{ColorGreen, PiecePentoX, true},
	{ColorBlue, PieceDomino, true},
	{ColorRed, PieceTrioI, true},
}

func checkDummyPiecesPlayed(t *testing.T, s State) {
	for _, set := range dummyPiecesPlayed {
		g := s.IsPiecePlayed(set.c, set.p)
		if g != set.e {
			t.Errorf("expected IsPiecePlayed(%s, %s) to be %v, but got %v", set.c, set.p, set.e, g)
		}
	}
}

func applyDummyStateSets(s MutableState) {
	s.Set(0, 0, ColorBlue, true)
	s.Set(19, 19, ColorRed, true)
	s.Set(19, 0, ColorGreen, false)
	s.Set(0, 19, ColorYellow, false)
	s.Set(9, 10, ColorBlue, true)
	s.Set(15, 15, ColorRed, false)
}

var dummyStateSets = []struct {
	x        uint8
	y        uint8
	c        Color
	hasPiece bool
}{
	{0, 0, ColorBlue, true},
	{19, 19, ColorRed, true},
	{19, 0, ColorGreen, false},
	{0, 19, ColorYellow, false},
	{9, 10, ColorBlue, true},
	{15, 15, ColorRed, false},
}

func checkDummyStateSets(t *testing.T, s State) {
	for _, set := range dummyStateSets {
		gotColor, gotHasPiece := s.At(set.x, set.y)
		if gotHasPiece != set.hasPiece {
			t.Errorf("expected hasPiece=%v at x=%d, y=%d, but got %v", set.hasPiece, set.x, set.y, gotHasPiece)
		}
		if !set.hasPiece || !gotHasPiece {
			continue
		}
		if gotColor != set.c {
			t.Errorf("expected color=%s at x=%d, y=%d, but got %v", set.c, set.x, set.y, gotColor)
		}
	}
}

func TestMutableStateSet(t *testing.T, s MutableState) {
	t.Helper()
	applyDummyStateSets(s)
	checkDummyStateSets(t, s)
}

func TestMutableStateSetLastMoveMono(t *testing.T, s MutableState) {
	sets := [4]bool{false, true, false, true}
	for c := Color(0); c < 4; c++ {
		s.SetLastMoveMono(c, sets[c])
	}
	for c := Color(0); c < 4; c++ {
		g := s.IsLastMoveMono(c)
		if sets[c] != g {
			t.Errorf("expected s.IsLastMoveMono(%s) to be %v but got %v", c.String(), sets[c], g)
		}
	}
}

func TestMutableStateSetPiecePlayed(t *testing.T, s MutableState) {
	t.Helper()
	setDummyPiecesPlayed(s)
	checkDummyPiecesPlayed(t, s)
}

func TestMutableStateSetNotPlayedPiecesFor(t *testing.T, s MutableState) {
	notPlayedPieces := []Piece{PieceDomino, PieceMono, PieceTetroL, PieceTrioI, PiecePentoR, PiecePentoW}
	var expectedIsPlayed [NumPieces]bool
	for _, p := range AllPieces {
		isPlayed := true
		for _, np := range notPlayedPieces {
			if p == np { // lol
				isPlayed = false
				break
			}
		}
		if isPlayed {
			expectedIsPlayed[p] = true
		}
	}

	for c := Color(0); c < 4; c++ { // exhaust all colors in case of compact storage / bitshifting errors
		s.SetNotPlayedPiecesFor(c, nil)
		colorFailed := false
		if s.IsPiecePlayed(c, PiecePentoU) != true {
			t.Errorf("expected IsPiecePlayed(%s, %s) to be false after SetNotPlayedPiecesFor(%s, nil), but got true", c.String(), PiecePentoU, c.String())
			colorFailed = true
		}
		s.SetNotPlayedPiecesFor(c, notPlayedPieces)
		for _, p := range AllPieces {
			g := s.IsPiecePlayed(c, p)
			e := expectedIsPlayed[p]
			if e != g {
				t.Errorf("expected IsPiecePlayed(%s, %s) to be %v but got %v", c.String(), p.String(), e, g)
				colorFailed = true
			}
		}
		if colorFailed {
			// ignore repeat errors
			break
		}
	}
}

func TestMutableStateReset(t *testing.T, s MutableState) {
	applyDummyStateSets(s)
	setDummyPiecesPlayed(s)
	s.Reset()
	if s.StartPiece() != PieceMono {
		t.Errorf("expected StartPiece() to be %s after Reset(), but got %s", PieceMono.String(), s.StartPiece().String())
	}
	if s.CurrentColor() != ColorBlue {
		t.Errorf("expected CurrentColor() to be %s after Reset(), but got %s", ColorBlue.String(), s.CurrentColor().String())
	}
	if s.IsPlayerOneFirst() != true {
		t.Errorf("expected IsPlayerOneFirst() to be true after Reset(), but got false")
	}
	for x := uint8(0); x < 20; x++ {
		for y := uint8(0); y < 20; y++ {
			_, hasPiece := s.At(x, y)
			if hasPiece {
				t.Errorf("expected At(%d, %d) to yield hasPiece=false, but got true", x, y)
			}
		}
	}
	for c := Color(0); c < 4; c++ {
		if !s.IsColorValid(c) {
			t.Errorf("expected IsColorValid(%s) to be true, but got false", c.String())
		}
		if s.HasPlayed(c) {
			t.Errorf("expected HasPlayed(%s) to be false, but got true", c.String())
		}
		gotNotPlayedPieces := s.NotPlayedPiecesFor(c)
		if len(gotNotPlayedPieces) != NumPieces {
			t.Errorf("expected len(NotPlayedPiecesFor(%s)) to be %d, but got %d", c.String(), NumPieces, len(gotNotPlayedPieces))
		}
		for _, p := range AllPieces {
			if s.IsPiecePlayed(c, p) {
				t.Errorf("expected IsPiecePlayed(%s, %s) to be false, but got true", c.String(), p.String())
				break
			}
		}
		if s.IsLastMoveMono(c) {
			t.Errorf("expected IsLastMoveMono(%s) to be false, but got true", c.String())
		}

	}
}

func TestMutableStateSetStartPiece(t *testing.T, s MutableState) {
	for _, p := range AllPieces {
		s.SetStartPiece(p)
		g := s.StartPiece()
		if p != g {
			t.Errorf("expected StartPiece() after SetStartPiece(%s) to be %s, but got %s", p.String(), p.String(), g.String())
		}
	}
}

func TestMutableStateSetCurrentColor(t *testing.T, s MutableState) {
	for c := Color(0); c < 4; c++ {
		s.SetCurrentColor(c)
		g := s.CurrentColor()
		if c != g {
			t.Errorf("expected CurrentColor() after SetCurrentColor(%s) to be %s, but got %s", c.String(), c.String(), g.String())
		}
	}
}

func TestMutableStateSetColorValid(t *testing.T, s MutableState) {
	for c := Color(0); c < 4; c++ {
		s.SetColorValid(c, true)
		if !s.IsColorValid(c) {
			t.Errorf("expected IsColorValid(%s) to be true after setting it, but got false", c.String())
		}
		s.SetColorValid(c, false)
		if s.IsColorValid(c) {
			t.Errorf("expected IsColorValid(%s) to be false after unsetting it, but got true", c.String())
		}
	}
}

func TestMutableStateSetPlayerOneFirst(t *testing.T, s MutableState) {
	s.SetPlayerOneFirst(true)
	if !s.IsPlayerOneFirst() {
		t.Errorf("expected IsPlayerOneFirst() to be true after setting it, but got false")
	}
	s.SetPlayerOneFirst(false)
	if s.IsPlayerOneFirst() {
		t.Errorf("expected IsPlayerOneFirst() to be false after setting it, but got true")
	}
}