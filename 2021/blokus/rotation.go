package blokus

import (
	"fmt"
	"strings"
)

type Rotation uint8

const (
	RotationNone = Rotation(iota)
	RotationRight
	RotationMirror
	RotationLeft
)

func (r Rotation) String() string {
	switch r {
	case RotationNone:
		return "NONE"
	case RotationRight:
		return "RIGHT"
	case RotationMirror:
		return "MIRROR"
	case RotationLeft:
		return "LEFT"
	default:
		panic(fmt.Sprintf("unkown Rotation value: %d", r))
	}
}

func ParseRotation(s string) (rotation Rotation, err error) {
	switch strings.TrimSpace(s) {
	case "NONE":
		rotation = RotationNone
	case "RIGHT":
		rotation = RotationRight
	case "MIRROR":
		rotation = RotationMirror
	case "LEFT":
		rotation = RotationLeft
	default:
		err = fmt.Errorf("unknown Rotation value: %q", s)
	}
	return
}
