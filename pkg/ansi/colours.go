package ansi

import "fmt"

// Colour represents an ANSI colour
type Colour uint8

const (
	ColourBlack Colour = iota + 30
	ColourRed
	ColourGreen
	ColourYellow
	ColourBlue
	ColourMagenta
	ColourCyan
	ColourLightGrey
)

const (
	ColourDarkGrey Colour = iota + 90
	ColourLightRed
	ColourLightGreen
	ColourLightYellow
	ColourLightBlue
	ColourLightMagenta
	ColourLightCyan
	ColourWhite
)

// Fg converts the colour to an ANSI foreground SGR code
func (c Colour) Fg() Colour {
	return c
}

// Bg converts the colour to an ANSI background SGR code
func (c Colour) Bg() Colour {
	return c + 10
}

// ColourFromString creates an ANSI colour code from a name string e.g. "red"
func ColourFromString(c string) (Colour, error) {
	switch c {
	case "black":
		return ColourBlack, nil
	case "red":
		return ColourRed, nil
	case "green":
		return ColourGreen, nil
	case "yellow":
		return ColourYellow, nil
	case "blue":
		return ColourBlue, nil
	case "magenta":
		return ColourMagenta, nil
	case "cyan":
		return ColourCyan, nil
	case "lightgrey":
		return ColourLightGrey, nil
	case "darkgrey":
		return ColourDarkGrey, nil
	case "lightred":
		return ColourLightRed, nil
	case "lightgreen":
		return ColourLightGreen, nil
	case "lightyellow":
		return ColourLightYellow, nil
	case "lightblue":
		return ColourLightBlue, nil
	case "lightmagenta":
		return ColourLightMagenta, nil
	case "lightcyan":
		return ColourLightCyan, nil
	case "white":
		return ColourWhite, nil
	default:
		return 0, fmt.Errorf("colour not found")
	}
}
