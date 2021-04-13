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
	case "BLUE":
		c = ColorBlue
	case "YELLOW":
		c = ColorYellow
	case "RED":
		c = ColorRed
	case "GREEN":
		c = ColorGreen
	default:
		err = fmt.Errorf("unknown color value: %q", s)
	}
	return
}

func OtherOwnColor(c Color) Color {
	switch c {
	case ColorBlue:
		return ColorRed
	case ColorYellow:
		return ColorGreen
	case ColorRed:
		return ColorBlue
	case ColorGreen:
		return ColorYellow
	default:
		panic(fmt.Sprintf("unknown color value: %d", c))
	}
}

func EnemyColors(c Color) [2]Color {
	switch c {
	case ColorBlue:
		return [2]Color{ColorYellow, ColorGreen}
	case ColorYellow:
		return [2]Color{ColorBlue, ColorRed}
	case ColorRed:
		return [2]Color{ColorYellow, ColorGreen}
	case ColorGreen:
		return [2]Color{ColorBlue, ColorRed}
	default:
		panic(fmt.Sprintf("unknown color value: %d", c))
	}
}

func OwnColors(c Color) [2]Color {
	switch c {
	case ColorBlue:
		return [2]Color{ColorBlue, ColorRed}
	case ColorYellow:
		return [2]Color{ColorYellow, ColorGreen}
	case ColorRed:
		return [2]Color{ColorBlue, ColorRed}
	case ColorGreen:
		return [2]Color{ColorYellow, ColorGreen}
	default:
		panic(fmt.Sprintf("unknown color value: %d", c))
	}
}
