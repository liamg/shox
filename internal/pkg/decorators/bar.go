package decorators

import (
	"fmt"
	"strconv"

	"github.com/liamg/shox/internal/pkg/ansi"
)

type SimpleBar struct {
	anchor Anchor
	text   string
}

func NewSimpleBar() *SimpleBar {
	return &SimpleBar{
		anchor: AnchorTop,
		text:   "I am a simple bar.",
	}
}

func (b *SimpleBar) Draw(rows uint16, cols uint16) {
	// TODO move cursor to top/bottom row, col zero
	var row, col uint16
	if b.anchor == AnchorBottom {
		row = rows - 1
	}
	ansi.SaveCursorPosition()
	ansi.MoveCursorTo(row+1, col+1)
	ansi.ClearLine()
	fmt.Printf("\r\033[44m\033[97m%-"+strconv.Itoa(int(cols))+"v", b.text)
	ansi.RestoreCursorPosition()
}

func (b *SimpleBar) GetAnchor() Anchor {
	return b.anchor
}

func (b *SimpleBar) GetHeight() (rows uint16) {
	return 1
}
