package main

import (
	"fmt"
	"os"

	"github.com/liamg/shox/internal/app/shox"
)

func main() {
	if err := shox.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to parse config file: %s", err)
		os.Exit(1)
	}
}
