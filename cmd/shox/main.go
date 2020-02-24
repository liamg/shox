package main

import (
	"github.com/liamg/shox/internal/pkg/decorators"
	"github.com/liamg/shox/internal/pkg/terminal"
)

func main() {
	t := terminal.NewTerminal()

	bar := decorators.NewSimpleBar()

	t.AddDecorator(bar)

	t.Run()
}
