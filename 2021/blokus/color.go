package blokus

import (
	"fmt"
	"strings"
)

type Color uint8

const (
	ColorBlue = Color(iota)
	ColorYellow
	ColorRed
	ColorGreen
)

func (c Color) String() string {
	switch c {
	case ColorBlue:
		return "BLUE"
	case ColorYellow:
		return "YELLOW"
	case ColorRed:
		return "RED"
	case ColorGreen:
		return "GREEN"
	default:
		panic(fmt.Sprintf("unknown color value: %d", c))
	}
}

func ParseColor(s string) (c Color, err error) {
	switch strings.TrimSpace(s) {
	case "BLUE": c = ColorBlue
	case "YELLOW": c = ColorYellow
	case "RED": c = ColorRed
	case "GREEN": c = ColorGreen
	default:
		err = fmt.Errorf("unknown color value: %q", s)
	}
	return
}