package decorators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/liamg/shox/pkg/helpers"

	"github.com/liamg/shox/pkg/ansi"
)

var helperRegex = regexp.MustCompile(`{[^}]+}`)

type StatusBar struct {
	anchor  Anchor
	format  string
	bg      ansi.Colour
	fg      ansi.Colour
	padding uint16
}

func NewStatusBar() *StatusBar {
	return &StatusBar{
		anchor:  AnchorTop,
		format:  "|{time}|",
		bg:      ansi.ColourRed.Bg(),
		fg:      ansi.ColourWhite.Fg(),
		padding: 0,
	}
}

func (b *StatusBar) SetFormat(format string) {
	b.format = format
}

func (b *StatusBar) SetBg(colour ansi.Colour) {
	b.bg = colour.Bg()
}

func (b *StatusBar) SetFg(colour ansi.Colour) {
	b.fg = colour.Fg()
}

func (b *StatusBar) Draw(rows uint16, cols uint16) {

	var row, col uint16
	switch b.anchor {
	case AnchorBottom:
		row = rows - 1
	}
	ansi.SaveCursorPosition()
	ansi.MoveCursorTo(row+1, col+1)
	ansi.ClearLine()

	// set colours
	fmt.Printf("\r\033[%dm\033[%dm", b.bg, b.fg)

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
			output = fmt.Sprintf("%-"+strconv.Itoa(colSize)+"v", output)
		case 1: // centre
			padSize := colSize + midExtra
			leftPad := padSize / 2
			rightPad := padSize
			output = fmt.Sprintf("%"+strconv.Itoa(leftPad)+"v", output)
			output = fmt.Sprintf("%-"+strconv.Itoa(rightPad)+"v", output)
		case 2: // right align
			output = fmt.Sprintf("%"+strconv.Itoa(colSize)+"v", output)
		}

		fmt.Printf("%s", output)

	}

	for i := uint16(0); i < b.padding; i++ {
		fmt.Printf("\n")
	}

	ansi.RestoreCursorPosition()
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

func (b *StatusBar) SetPadding(pad uint16) {
	b.padding = pad
}

func (b *StatusBar) GetAnchor() Anchor {
	return b.anchor
}

func (b *StatusBar) GetHeight() (rows uint16) {
	return b.padding + 1
}
