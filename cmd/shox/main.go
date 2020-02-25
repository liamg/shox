package main

import (
	"fmt"
	"os"

	"github.com/liamg/shox/internal/app/shox"
)

func main() {
	if err := shox.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to start shox: %s\n", err)
		os.Exit(1)
	}
}
