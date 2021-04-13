package blokus

import (
	"fmt"
	"strings"
)

const NumPieces = 21

type Piece uint8

const (
	PieceMono = Piece(iota)
	PieceDomino
	PieceTrioL
	PieceTrioI
	PieceTetroO
	PieceTetroT
	PieceTetroI
	PieceTetroL
	PieceTetroZ
	PiecePentoL
	PiecePentoT
	PiecePentoV
	PiecePentoS
	PiecePentoZ
	PiecePentoI
	PiecePentoP
	PiecePentoW
	PiecePentoU
	PiecePentoR
	PiecePentoX
	PiecePentoY
)

func (p Piece) String() string {
	switch p {
	case PieceMono:
		return "MONO"
	case PieceDomino:
		return "DOMINO"
	case PieceTrioL:
		return "TRIO_L"
	case PieceTrioI:
		return "TRIO_I"
	case PieceTetroO:
		return "TETRO_O"
	case PieceTetroT:
		return "TETRO_T"
	case PieceTetroI:
		return "TETRO_I"
	case PieceTetroL:
		return "TETRO_L"
	case PieceTetroZ:
		return "TETRO_Z"
	case PiecePentoL:
		return "PENTO_L"
	case PiecePentoT:
		return "PENTO_T"
	case PiecePentoV:
		return "PENTO_V"
	case PiecePentoS:
		return "PENTO_S"
	case PiecePentoZ:
		return "PENTO_Z"
	case PiecePentoI:
		return "PENTO_I"
	case PiecePentoP:
		return "PENTO_P"
	case PiecePentoW:
		return "PENTO_W"
	case PiecePentoU:
		return "PENTO_U"
	case PiecePentoR:
		return "PENTO_R"
	case PiecePentoX:
		return "PENTO_X"
	case PiecePentoY:
		return "PENTO_Y"
	default:
		panic(fmt.Errorf("unknown Piece value: %d", p))
	}
}

func ParsePiece(s string) (p Piece, err error) {
	switch strings.TrimSpace(s) {
	case "MONO":
		p = PieceMono
	case "DOMINO":
		p = PieceDomino
	case "TRIO_L":
		p = PieceTrioL
	case "TRIO_I":
		p = PieceTrioI
	case "TETRO_O":
		p = PieceTetroO
	case "TETRO_T":
		p = PieceTetroT
	case "TETRO_I":
		p = PieceTetroI
	case "TETRO_L":
		p = PieceTetroL
	case "TETRO_Z":
		p = PieceTetroZ
	case "PENTO_L":
		p = PiecePentoL
	case "PENTO_T":
		p = PiecePentoT
	case "PENTO_V":
		p = PiecePentoV
	case "PENTO_S":
		p = PiecePentoS
	case "PENTO_Z":
		p = PiecePentoZ
	case "PENTO_I":
		p = PiecePentoI
	case "PENTO_P":
		p = PiecePentoP
	case "PENTO_W":
		p = PiecePentoW
	case "PENTO_U":
		p = PiecePentoU
	case "PENTO_R":
		p = PiecePentoR
	case "PENTO_X":
		p = PiecePentoX
	case "PENTO_Y":
		p = PiecePentoY
	default:
		err = fmt.Errorf("unknown Piece value: %q", s)
	}
	return
}

func (p Piece) Name() string {
	return p.String()
}

func (p Piece) Width() uint8 {
	switch p {
	case PieceMono:
		return 1
	case PieceDomino:
		return 2
	case PieceTrioL:
		return 2
	case PieceTrioI:
		return 1
	case PieceTetroO:
		return 2
	case PieceTetroT:
		return 3
	case PieceTetroI:
		return 1
	case PieceTetroL:
		return 2
	case PieceTetroZ:
		return 3
	case PiecePentoL:
		return 2
	case PiecePentoT:
		return 3
	case PiecePentoV:
		return 3
	case PiecePentoS:
		return 4
	case PiecePentoZ:
		return 3
	case PiecePentoI:
		return 1
	case PiecePentoP:
		return 2
	case PiecePentoW:
		return 3
	case PiecePentoU:
		return 3
	case PiecePentoR:
		return 3
	case PiecePentoX:
		return 3
	case PiecePentoY:
		return 2
	default:
		panic(fmt.Errorf("unknown Piece value: %d", p))
	}
}

func (p Piece) Height() uint8 {
	switch p {
	case PieceMono:
		return 1
	case PieceDomino:
		return 1
	case PieceTrioL:
		return 2
	case PieceTrioI:
		return 3
	case PieceTetroO:
		return 2
	case PieceTetroT:
		return 2
	case PieceTetroI:
		return 4
	case PieceTetroL:
		return 3
	case PieceTetroZ:
		return 2
	case PiecePentoL:
		return 4
	case PiecePentoT:
		return 3
	case PiecePentoV:
		return 3
	case PiecePentoS:
		return 2
	case PiecePentoZ:
		return 3
	case PiecePentoI:
		return 5
	case PiecePentoP:
		return 3
	case PiecePentoW:
		return 3
	case PiecePentoU:
		return 2
	case PiecePentoR:
		return 3
	case PiecePentoX:
		return 3
	case PiecePentoY:
		return 4
	default:
		panic(fmt.Errorf("unknown Piece value: %d", p))
	}
}

func (p Piece) Fills(x, y uint8) bool {
	switch p {
	case PieceMono:
		return x == 0 && y == 0
	case PieceDomino:
		return x <= 1 && y == 0
	case PieceTrioL:
		return x == 0 && y == 0 || x <= 1 && y == 1
	case PieceTrioI:
		return x == 0 && y <= 2
	case PieceTetroO:
		return x <= 1 && y <= 1
	case PieceTetroT:
		return x == 0 && y <= 2 || x == 1 && y == 1
	case PieceTetroI:
		return x == 0 && y <= 3
	case PieceTetroL:
		return x == 0 && y <= 2 || x == 1 && y == 2
	case PieceTetroZ:
		return y == 0 && x <= 1 || x >= 1 && x <= 2 && y == 1
	case PiecePentoL:
		return x == 0 && y <= 4 || x == 1 && y == 4
	case PiecePentoT:
		return y == 0 && x <= 2 || y <= 2 && x == 1
	case PiecePentoV:
		return x == 0 && y <= 2 || x <= 2 && y == 2
	case PiecePentoS:
		return y == 1 && x <= 1 || y == 0 && x >= 1 && x <= 3
	case PiecePentoZ:
		return x == 0 && y == 0 || x == 1 && y <= 2 || x == 2 && y == 2
	case PiecePentoI:
		return x == 0 && y <= 4
	case PiecePentoP:
		return x == 0 && y <= 2 || x == 1 && y <= 1
	case PiecePentoW:
		return x == 0 && y <= 1 || x == 1 && y >= 1 && y <= 2 || x == 2 && y == 2
	case PiecePentoU:
		return x == 0 && y <= 1 || x == 1 && y == 1 || x == 2 && y <= 1
	case PiecePentoR:
		return y == 0 && x == 2 || y == 1 && x <= 2 || y == 2 && x == 1
	case PiecePentoX:
		return x == 0 && y == 1 || x == 1 && y <= 2 || x == 2 && y == 1
	case PiecePentoY:
		return x == 0 && y == 1 || x == 1 && y <= 3
	default:
		panic(fmt.Errorf("unknown Piece value: %d", p))
	}
}

var piecePoints = [NumPieces][]Position{
	// Mono
	{{0, 0}},
	// Domino
	{{0, 0}, {1, 0}},
	// TrioL
	{{0, 0}, {0, 1}, {1, 1}},
	// TrioI
	{{0, 0}, {0, 1}, {0, 2}},
	// TetroO
	{{0, 0}, {0, 1}, {1, 0}, {1, 1}},
	// TetroT
	{{0, 0}, {1, 0}, {2, 0}, {1, 1}},
	// TetroI
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}},
	// TetroL
	{{0, 0}, {0, 1}, {0, 2}, {1, 2}},
	// TetroZ
	{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
	// PentoL
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {1, 3}},
	// PentoT
	{{0, 0}, {1, 0}, {1, 1}, {1, 2}, {2, 0}},
	// PentoV
	{{0, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}},
	// PentoS
	{{0, 1}, {1, 0}, {1, 1}, {2, 0}, {3, 0}},
	// PentoZ
	{{0, 0}, {1, 0}, {1, 1}, {1, 2}, {2, 2}},
	// PentoI
	{{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}},
	// PentoP
	{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 1}},
	// PentoW
	{{0, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 2}},
	// PentoU
	{{0, 0}, {0, 1}, {1, 1}, {2, 0}, {2, 1}},
	// PentoR
	{{0, 1}, {1, 1}, {1, 2}, {2, 0}, {2, 1}},
	// PentoX
	{{0, 1}, {1, 0}, {1, 1}, {1, 2}, {2, 1}},
	// PentoY
	{{0, 1}, {1, 0}, {1, 1}, {1, 2}, {1, 3}},
}

func (p Piece) Points() []Position {
	return piecePoints[p]
}

// IsFlipInvariant returns true if flipping the piece does not yield new forms
func (p Piece) IsFlipInvariant() bool {
	switch p {
	case PieceMono, PieceDomino, PieceTrioI, PieceTetroO, PieceTetroI, PieceTetroT, PiecePentoI,
		PiecePentoU, PiecePentoX, PiecePentoT:
		return true
	default:
		return false
	}
}

// IsRotationSymmetric returns true if the only true new form is generated by rotating right
func (p Piece) IsRotationSymmetric() bool {
	if p.IsInvariant() {
		return true
	}
	switch p {
	case PieceDomino, PieceTrioI, PieceTetroI, PiecePentoI:
		return true
	default:
		return false
	}
}

// IsInvariant is true, if all transformations lead to the same form
func (p Piece) IsInvariant() bool {
	switch p {
	case PieceMono, PieceTetroO, PiecePentoX:
		return true
	default:
		return false
	}
}

func (p Piece) NumPoints() uint {
	switch p {
	case PieceMono:
		return 1
	case PieceDomino:
		return 2
	case PieceTrioL:
		return 3
	case PieceTrioI:
		return 3
	case PieceTetroO:
		return 4
	case PieceTetroT:
		return 4
	case PieceTetroI:
		return 4
	case PieceTetroL:
		return 4
	case PieceTetroZ:
		return 4
	case PiecePentoL:
		return 5
	case PiecePentoT:
		return 5
	case PiecePentoV:
		return 5
	case PiecePentoS:
		return 5
	case PiecePentoZ:
		return 5
	case PiecePentoI:
		return 5
	case PiecePentoP:
		return 5
	case PiecePentoW:
		return 5
	case PiecePentoU:
		return 5
	case PiecePentoR:
		return 5
	case PiecePentoX:
		return 5
	case PiecePentoY:
		return 5
	default:
		panic(fmt.Errorf("unknown Piece value: %d", p))
	}
}

const AdditionalPointsIfLastMoveMono = 5
const AdditionalPointsIfAllPiecesPlayed = 15
const PointsForAllPieces = 1 + 2 + 2*3 + 5*4 + 12*5 + AdditionalPointsIfAllPiecesPlayed

var AllPieces = [NumPieces]Piece{
	PiecePentoL,
	PiecePentoT,
	PiecePentoV,
	PiecePentoS,
	PiecePentoZ,
	PiecePentoI,
	PiecePentoP,
	PiecePentoW,
	PiecePentoU,
	PiecePentoR,
	PiecePentoX,
	PiecePentoY,
	PieceTetroO,
	PieceTetroT,
	PieceTetroI,
	PieceTetroL,
	PieceTetroZ,
	PieceTrioL,
	PieceTrioI,
	PieceDomino,
	PieceMono,
}

func applyTransformation(positions []Position, rotation Rotation, flipped bool) []Position {
	switch rotation {
	case RotationRight:
		positions = rotateRight(positions)
	case RotationMirror:
		positions = mirror(positions)
	case RotationLeft:
		positions = rotateLeft(positions)
	}
	if flipped {
		positions = flipPositions(positions)
	}
	return positions
}

func flipPositions(positions []Position) (newPositions []Position) {
	newPositions = make([]Position, 0, len(positions))
	var maxX uint8
	for _, pos := range positions {
		if pos.X > maxX {
			maxX = pos.X
		}
	}
	for _, pos := range positions {
		newPositions = append(newPositions, Position{X: maxX - pos.X, Y: pos.Y})
	}
	return
}

func rotateRight(positions []Position) (newPositions []Position) {
	newPositions = make([]Position, 0, len(positions))
	var maxY uint8
	for _, pos := range positions {
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}
	for _, pos := range positions {
		newPositions = append(newPositions, Position{X: maxY - pos.Y, Y: pos.X})
	}
	return
}

func rotateLeft(positions []Position) (newPositions []Position) {
	newPositions = make([]Position, 0, len(positions))
	var maxX uint8
	for _, pos := range positions {
		if pos.X > maxX {
			maxX = pos.X
		}
	}
	for _, pos := range positions {
		newPositions = append(newPositions, Position{X: pos.Y, Y: maxX - pos.X})
	}
	return
}

func mirror(positions []Position) (newPositions []Position) {
	newPositions = make([]Position, 0, len(positions))
	var maxX, maxY uint8
	for _, pos := range positions {
		if pos.X > maxX {
			maxX = pos.X
		}
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}
	for _, pos := range positions {
		newPositions = append(newPositions, Position{X: maxX - pos.X, Y: maxY - pos.Y})
	}
	return
}

var transformedPieces [NumPieces][4][2][]Position
var uniquePieceTransformations [NumPieces][]TransformedPiece

func init() {
	for p := PieceMono; p < NumPieces; p++ {
		for rotation := RotationNone; rotation < 4; rotation++ {
			for flippedI := 0; flippedI < 2; flippedI++ {
				flipped := flippedI == 1
				transformedPositions := applyTransformation(p.Points(), rotation, flipped)
				found := false
				for _, previousTransformedPiece := range uniquePieceTransformations[p] {
					previousFlippedI := 0
					if previousTransformedPiece.Flipped() {
						previousFlippedI = 1
					}
					if PositionsEqual(transformedPositions, previousTransformedPiece.Positions()) {
						transformedPieces[p][rotation][flippedI] = transformedPieces[previousTransformedPiece.Piece()][previousTransformedPiece.Rotation()][previousFlippedI]
						found = true
						break
					}
				}
				if !found {
					transformedPiece := NewTransformedPiece(p, rotation, flipped)
					uniquePieceTransformations[p] = append(uniquePieceTransformations[p], transformedPiece)
					transformedPieces[p][rotation][flippedI] = transformedPositions
				}
			}
		}
	}
}
