package ansi

import "fmt"

func ClearLine() {
	fmt.Printf("\033[K")
}

func SaveCursorPosition() {
	fmt.Printf("\033[s")
}

func RestoreCursorPosition() {
	fmt.Printf("\033[u")
}

// MoveCursorTo 1-indexed position
func MoveCursorTo(row, col uint16) {
	fmt.Printf("\033[%d;%dH", row, col)
}

func IsOSCTerminator(char byte) bool {
	return char == 0x5c || char == 0x07
}

func MoveToCol(col uint16) {
	fmt.Printf("\033[%dG", col)
}
