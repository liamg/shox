package ansi

import "fmt"

// ClearLine clears the current terminal line at the cursor position
func ClearLine() {
	fmt.Printf("\033[K")
}

// Clear clears all content from the terminal
func Clear() {
	fmt.Printf("\033[2J")
}

// Reset performs a full reset on the terminal
func Reset() {
	fmt.Printf("\033c")
}

// SaveCursorPosition pushes the cursor position to the stack
func SaveCursorPosition() {
	fmt.Printf("\033[s")
}

// RestoreCursorPosition pops the cursor position from the stack
func RestoreCursorPosition() {
	fmt.Printf("\033[u")
}

// MoveCursorTo a 1-indexed position
func MoveCursorTo(row, col uint16) {
	fmt.Printf("\033[%d;%dH", row, col)
}
