package decorators

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/liamg/shox/pkg/helpers"
	"github.com/mattn/go-runewidth"

	"github.com/liamg/shox/pkg/ansi"
)

var helperRegex = regexp.MustCompile(`{[^}]+}`)

// StatusBar is a full width bar containing useful info which can be added to the terminal using a proxy
type StatusBar struct {
	anchor  Anchor
	format  string
	bg      ansi.Colour
	fg      ansi.Colour
	padding uint16
	visible bool
}

// NewStatusBar creates a new status bar instance
func NewStatusBar() *StatusBar {
	return &StatusBar{
		anchor:  AnchorTop,
		format:  "|{time}|",
		bg:      ansi.ColourRed.Bg(),
		fg:      ansi.ColourWhite.Fg(),
		padding: 0,
		visible: true,
	}
}

func (b *StatusBar) IsVisible() bool {
	return b.visible
}

func (b *StatusBar) SetVisible(visible bool) {
	b.visible = visible
}

// SetFormat controls the output format of the status bar
func (b *StatusBar) SetFormat(format string) {
	b.format = format
}

// SetBg sets the background colour of the status bar
func (b *StatusBar) SetBg(colour ansi.Colour) {
	b.bg = colour.Bg()
}

// SetFg sets the background colour of the status bar
func (b *StatusBar) SetFg(colour ansi.Colour) {
	b.fg = colour.Fg()
}

// Draw renders the decorator to StdOut
func (b *StatusBar) Draw(rows uint16, cols uint16, writeFunc func(data []byte)) {

	var row, col uint16
	switch b.anchor {
	case AnchorBottom:
		row = rows - 1
	}

	// move cursor to status bar location
	writeFunc([]byte(fmt.Sprintf("\033[%d;%dH", row+1, col+1)))

	// clear line
	writeFunc([]byte("\x1b[K"))

	// set colours
	writeFunc([]byte(fmt.Sprintf("\r\033[%dm\033[%dm", b.bg, b.fg)))

	segments := strings.SplitN(b.format, "|", 3)
	colSize := int(cols) / len(segments)
	midExtra := int(cols) - (colSize * len(segments))
	for i, segment := range segments {
		output := b.applyHelpers(segment)
		if len(output) > colSize {
			output = output[:colSize]
		}
		switch i {
		case 0: // left align
			output = padRight(output, colSize)
		case 1: // centre
			padSize := colSize + midExtra
			leftPad := padSize / 2
			output = padLeft(output, leftPad)
			output = padRight(output, padSize)
		case 2: // right align
			output = padLeft(output, colSize)
		}

		writeFunc([]byte(output))

	}

	for i := uint16(0); i < b.padding; i++ {
		writeFunc([]byte{0x0a})
	}

}

// SetPadding sets a vertical padding on the status bar
func (b *StatusBar) SetPadding(pad uint16) {
	b.padding = pad
}

// GetAnchor returns the anchor e.g. Top/Bottom
func (b *StatusBar) GetAnchor() Anchor {
	return b.anchor
}

// GetHeight returns the height of the decorator in terminal character rows
func (b *StatusBar) GetHeight() (rows uint16) {
	return b.padding + 1
}

func padLeft(input string, totalLen int) string {
	pad := totalLen - runewidth.StringWidth(input) // utf8.RuneCountInString(input)
	if pad > 0 {
		input = strings.Repeat(" ", pad) + input
	}
	return input
}

func padRight(input string, totalLen int) string {
	pad := totalLen - runewidth.StringWidth(input) // utf8.RuneCountInString(input)
	if pad > 0 {
		input += strings.Repeat(" ", pad)
	}
	return input
}

func (b *StatusBar) applyHelpers(segment string) string {
	formatted := segment
	// run helpers
	helperPatterns := helperRegex.FindAllString(segment, -1)
	for _, pattern := range helperPatterns {
		args := strings.SplitN(pattern[1:len(pattern)-1], ":", 2)
		helper := args[0]
		var config string
		if len(args) > 1 {
			config = args[1]
		}
		output, err := helpers.Run(helper, config)
		if err != nil {
			continue
		}
		formatted = strings.Replace(formatted, pattern, output, 1)
	}
	return formatted
}
