package blokus

func OfficialRatingForColor(s State, c Color) (points uint) {
	if len(s.NotPlayedPiecesFor(c)) == 0 {
		points += PointsForAllPieces
		if s.IsLastMoveMono(c) {
			points += AdditionalPointsIfLastMoveMono
		}
		return
	}
	for _, p := range AllPieces {
		if s.IsPiecePlayed(c, p) {
			points += p.NumPoints()
		}
	}
	return
}

func OfficialRatingForPlayer(s State, isFirstPlayer bool) uint {
	if isFirstPlayer {
		return OfficialRatingForColor(s, ColorBlue) + OfficialRatingForColor(s, ColorGreen)
	}
	return OfficialRatingForColor(s, ColorYellow) + OfficialRatingForColor(s, ColorRed)
}
