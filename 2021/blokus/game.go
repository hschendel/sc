package blokus

import (
	"fmt"
	"io"
)

func RunRepeatedGames(player1, player2 Player, player1Name, player2Name string, repetitions uint, logTo io.Writer) (result RepeatGameResult) {
	for ri := uint(0); ri < repetitions; ri++ {
		one, two := player1, player2
		if (ri % 2) == 1 {
			one, two = player2, player1
		}
		gameResult, score1, score2, err1, err2 := RunGame(one, two)
		if (ri % 2) == 1 {
			if gameResult == GameResultPlayer1Won {
				gameResult = GameResultPlayer2Won
			} else if gameResult == GameResultPlayer2Won {
				gameResult = GameResultPlayer1Won
			}
			score1, score2 = score2, score1
			err1, err2 = err2, err1
		}

		var resultText string
		switch gameResult {
		case GameResultPlayer1Won:
			result.Player1Wins++
			resultText = fmt.Sprintf("%s wins", player1Name)
		case GameResultPlayer2Won:
			result.Player2Wins++
			resultText = fmt.Sprintf("%s wins", player2Name)
		case GameResultDraw:
			result.Draws++
			resultText = "draw"
		}

		if err1 != nil {
			result.Player1Errors = append(result.Player1Errors, err1)
		}
		if err2 != nil {
			result.Player2Errors = append(result.Player2Errors, err2)
		}

		result.Player1TotalScore += score1
		result.Player2TotalScore += score2

		if logTo != nil {
			fmt.Fprintf(logTo, "game %3d: %s, score %d:%d\n", ri+1, resultText, score1, score2)
			if err1 != nil {
				fmt.Fprintf(logTo, "  error reported from %s: %s\n", player1Name, err1)
			}
			if err2 != nil {
				fmt.Fprintf(logTo, "  error reported from %s: %s\n", player2Name, err2)
			}
		}
	}

	result.Player1AvgScore = float64(result.Player1TotalScore) / float64(repetitions)
	result.Player2AvgScore = float64(result.Player2TotalScore) / float64(repetitions)

	return
}

type RepeatGameResult struct {
	Draws uint
	Player1Wins uint
	Player2Wins uint
	Player1Errors []error
	Player2Errors []error
	Player1TotalScore uint
	Player2TotalScore uint
	Player1AvgScore float64
	Player2AvgScore float64
}

func (r *RepeatGameResult) Print(w io.Writer) {
	fmt.Fprintf(w, "Draws:                %3d\n", r.Draws)
	fmt.Fprintf(w, "Wins for Player 1:    %3d\n", r.Player1Wins)
	fmt.Fprintf(w, "Wins for Player 2:    %3d\n\n", r.Player2Wins)
	fmt.Fprintf(w, "Total score Player 1: %12d\n", r.Player1TotalScore)
	fmt.Fprintf(w, "Avg. score Player 1:  %5.1f\n\n", r.Player1AvgScore)
	fmt.Fprintf(w, "Total score Player 2: %12d\n", r.Player2TotalScore)
	fmt.Fprintf(w, "Avg. score Player 2:  %5.1f\n\n", r.Player2AvgScore)

	printErrors(w, "Player 1", r.Player1Errors)
	printErrors(w, "Player 2", r.Player2Errors)
}

func printErrors(w io.Writer, playerName string, errors []error) {
	if len(errors) != 0 {
		fmt.Fprintf(w,"\nReported errors for %s:\n", playerName)
		for _, err := range errors {
			fmt.Fprintf(w, "  %s\n", err)
		}
	}
}

type GameResult uint8

const (
	GameResultDraw = GameResult(iota)
	GameResultPlayer1Won
	GameResultPlayer2Won
)

func RunGame(player1, player2 Player) (result GameResult, score1, score2 uint, err1, err2 error) {
	const timeout = DefaultMoveTimeout
	var state, copyState BasicState
	startPiece := AllPieces[RandomInt(len(AllPieces))]
	var tr turnTracker
	tr.players[0] = player1
	tr.players[1] = player2

	for ; !tr.gameEnded; tr.nextColor() {
		player, color := tr.current()
		if !tr.firstRound() && len(PossibleNextMoves(&state, color)) == 0 {
			tr.endCurrent()
			continue
		}
		CopyState(&copyState, &state)
		var move Move
		t := NewTimeout(timeout)
		var timeoutReached bool
		if tr.firstRound() {
			move = player.FirstMove(&copyState, color, startPiece, NewTimeout(timeout))
			timeoutReached = t.Reached()
			if move.IsEmpty() {
				result, err1, err2 = setErrorResult(color, "first move must not be empty")
				return
			}
			if move.Transformation.Piece() != startPiece {
				result, err1, err2 = setErrorResult(color, "start piece was %s but got %s", startPiece.String(), move.Transformation.Piece().String())
				return
			}
			if !CanPlayFirstPiece(&state, move.Transformation, move.X, move.Y) {
				result, err1, err2 = setErrorResult(color, "start move is invalid: %s", move.FormatPretty('X', ""))
				return
			}
		} else {
			move = player.NextMove(&copyState, color, NewTimeout(timeout))
			timeoutReached = t.Reached()
			if !move.IsEmpty() {
				if state.IsPiecePlayed(color, move.Transformation.Piece()) {
					result, err1, err2 = setErrorResult(color, "invalid move: piece %s was already played", move.Transformation.Piece().String())
					return
				}
				if !CanPlayNextPiece(&state, color, move.Transformation, move.X, move.Y) {
					result, err1, err2 = setErrorResult(color, "invalid move: %s", move.FormatPretty('X', "  "))
					return
				}
			}
		}
		if timeoutReached {
			result, err1, err2 = setErrorResult(color, "player hit timeout")
			return
		}
		if !move.IsEmpty() {
			moveErr := ApplyMove(&state, color, move)
			if moveErr != nil {
				result, err1, err2 = setErrorResult(color, "move error: %s", moveErr)
				return
			}
			score1, score2 = updateScore(color, move, score1, score2)
		}
	}
	result, score1, score2 = finalizeScore(score1, score2, &state)
	return
}

type turnTracker struct {
	color           Color
	players         [2]Player
	colorEnded      [4]bool
	roundsCompleted uint
	gameEnded       bool
}

func (t *turnTracker) current() (p Player, c Color) {
	p = t.players[t.color % 2]
	c = t.color
	return
}

func (t *turnTracker) endCurrent() {
	t.colorEnded[t.color] = true
	oldColor := t.color
	t.color = (t.color + 2) % 4
	if t.colorEnded[t.color] {
		t.color = oldColor
		t.gameEnded = true
	}
	return
}

func (t *turnTracker) nextColor() {
	oldColor := t.color
	t.color = (t.color + 1) % 4
	if t.colorEnded[t.color] {
		t.color = (t.color +2) % 4
		t.gameEnded = t.colorEnded[t.color]
		if t.gameEnded {
			t.color = oldColor
		}
	}
	if oldColor == 3 && !t.gameEnded {
		t.roundsCompleted++
	}
	return
}

func (t *turnTracker) firstRound() bool {
	return t.roundsCompleted == 0
}

func setErrorResult(c Color, format string, params ...interface{}) (result GameResult, err1, err2 error) {
	playerIdx := uint8(c) % 2
	msg := fmt.Sprintf(format, params...)
	err := fmt.Errorf("%s: %s", c.String(), msg)
	if playerIdx == 0 {
		result = GameResultPlayer2Won
		err1 = err
	} else {
		result = GameResultPlayer1Won
		err2 = err
	}
	return
}

func updateScore(c Color, m Move, score1, score2 uint) (newScore1, newScore2 uint) {
	newScore1, newScore2 = score1, score2
	playerIdx := uint8(c) % 2
	points := m.Transformation.Piece().NumPoints()
	if playerIdx == 0 {
		newScore1 += points
	} else {
		newScore2 += points
	}
	return
}

func finalizeScore(score1, score2 uint, s State) (result GameResult, finalScore1, finalScore2 uint) {
	finalScore1, finalScore2 = score1, score2
	for c := Color(0); c < 4; c++ {
		addScore := uint(0)
		if len(s.NotPlayedPiecesFor(c)) == 0 {
			addScore += 15
			if s.IsLastMoveMono(c) {
				addScore += 5
			}
		}
		if (c % 2) == 0 {
			finalScore1 += addScore
		} else {
			finalScore2 += addScore
		}
	}
	if finalScore1 > finalScore2 {
		result = GameResultPlayer1Won
	} else if finalScore2 > finalScore1 {
		result = GameResultPlayer2Won
	} else {
		result = GameResultDraw
	}
	return
}