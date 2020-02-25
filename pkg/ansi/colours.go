package ansi

import "fmt"

// Colour represents an ANSI colour
type Colour uint8

const (
	// ColourBlack is black
	ColourBlack Colour = iota + 30
	// ColourRed is the colour red
	ColourRed
	// ColourGreen is the colour green
	ColourGreen
	// ColourYellow is the colour yellow
	ColourYellow
	// ColourBlue is the colour blue
	ColourBlue
	// ColourMagenta is the colour magenta
	ColourMagenta
	// ColourCyan is the colour cyan
	ColourCyan
	// ColourLightGrey is the colour grey
	ColourLightGrey
)

const (
	// ColourDarkGrey is the colour dark grey
	ColourDarkGrey Colour = iota + 90
	// ColourLightRed is the colour light red
	ColourLightRed
	// ColourLightGreen is the colour light green
	ColourLightGreen
	// ColourLightYellow is the colour light yellow
	ColourLightYellow
	// ColourLightBlue is the colour light blue
	ColourLightBlue
	// ColourLightMagenta is the colour light magenta
	ColourLightMagenta
	// ColourLightCyan is the colour light cyan
	ColourLightCyan
	// ColourWhite is the colour white
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
