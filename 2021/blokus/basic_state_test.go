package blokus

import "testing"

func BenchmarkBasicState_At(b *testing.B) {
	b.StopTimer()
	s := earlyTestState()
	b.StartTimer()
	BenchmarkStateAt(b, s)
}

func BenchmarkBasicState_Set(b *testing.B) {
	b.StopTimer()
	s := earlyTestState()
	b.StartTimer()
	BenchmarkMutableStateResetSet(b, s)
}

func BenchmarkBasicState_IsPiecePlayed(b *testing.B) {
	b.StopTimer()
	var s BasicState
	b.StartTimer()
	BenchmarkStateIsPiecePlayed(b, &s)
}

func TestBasicState_Set(t *testing.T) {
	var s BasicState
	TestMutableStateSet(t, &s)
}

func TestBasicState_Reset(t *testing.T) {
	var s BasicState
	TestMutableStateReset(t, &s)
}

func TestBasicState_SetNotPlayedPiecesFor(t *testing.T) {
	var s BasicState
	TestMutableStateSetNotPlayedPiecesFor(t, &s)
}

func TestBasicState_SetPiecePlayed(t *testing.T) {
	var s BasicState
	TestMutableStateSetPiecePlayed(t, &s)
}

func TestBasicState_SetLastMoveMono(t *testing.T) {
	var s BasicState
	TestMutableStateSetLastMoveMono(t, &s)
}