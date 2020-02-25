package shox

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/liamg/shox/pkg/ansi"

	"github.com/liamg/shox/pkg/decorators"
	"github.com/liamg/shox/pkg/terminal"
)

func Run() error {

	term := terminal.NewTerminal()
	bar := decorators.NewStatusBar()
	bar.SetFormat("{time}||CPU: {cpu} MEM: {memory}")
	term.AddDecorator(bar)

	shell := os.Getenv("SHELL")
	if filepath.Base(shell) != "shox" {
		term.SetShell(shell)
	}

	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config file: %s", err)
	}

	if config != nil {
		if config.Shell != "" {
			term.SetShell(config.Shell)
		}
		if config.Bar.Format != "" {
			bar.SetFormat(config.Bar.Format)
		}
		if bg, err := ansi.ColourFromString(config.Bar.Colours.Bg); err == nil {
			bar.SetBg(bg)
		}
		if fg, err := ansi.ColourFromString(config.Bar.Colours.Fg); err == nil {
			bar.SetFg(fg)
		}
		bar.SetPadding(config.Bar.Padding)
	}

	return term.Run()
}
