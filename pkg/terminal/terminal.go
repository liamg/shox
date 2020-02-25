package terminal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/liamg/shox/pkg/decorators"

	"github.com/liamg/shox/pkg/proxy"

	"github.com/creack/pty"
	"golang.org/x/crypto/ssh/terminal"
)

type Terminal struct {
	shell string
	proxy *proxy.Proxy
}

func NewTerminal() *Terminal {

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	return &Terminal{
		shell: shell,
		proxy: proxy.NewProxy(),
	}
}

func (t *Terminal) SetShell(shell string) {
	t.shell = shell
}

func (t *Terminal) AddDecorator(d decorators.Decorator) {
	t.proxy.AddDecorator(d)
}

func (t *Terminal) Run() error {

	if os.Getenv("SHOX") != "" {
		return fmt.Errorf("shox is already running in this terminal")
	}

	_ = os.Setenv("SHOX", "1")

	t.proxy.Start()
	defer t.proxy.Close()
	t.proxy.Write([]byte("\033c")) // reset term

	// Create arbitrary command.
	c := exec.Command(t.shell)

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		panic(err)
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {

			size, err := pty.GetsizeFull(os.Stdin)
			if err != nil {
				continue
			}

			rows, cols := t.proxy.HandleResize(size.Rows, size.Cols)
			size.Rows = rows
			size.Cols = cols

			if err := pty.Setsize(ptmx, size); err != nil {
				continue
			}

			// successful resize
		}
	}()
	ch <- syscall.SIGWINCH // Initial resize.

	// Set stdin in raw mode.
	oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	// Copy stdin to the pty and the pty to stdout.
	go func() { _, _ = io.Copy(ptmx, os.Stdin); fmt.Println("Finished read from stdin") }()
	go func() {
		_, _ = io.Copy(os.Stdout, t.proxy)
		fmt.Println("Finished read from proxy")
	}()
	_, _ = io.Copy(t.proxy, ptmx)
	fmt.Printf("\r\n")
}
