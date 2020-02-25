package proxy

import (
	"fmt"
)

type escapeSequenceHandler func(pty chan byte) (output []byte, discard []byte, redraw bool, err error)

// ErrorCommandNotSupported means the command is not supported and should be parsed to the underlying terminal
var ErrorCommandNotSupported = fmt.Errorf("command not supported")

func (p *Proxy) proxyANSICommand(input chan byte) (output []byte, original []byte, requiredRedraw bool) {

	b := <-input
	original = append(original, b)

	switch b {
	case 'c': //RIS
		row, col := p.HandleCoordinates(0, 0)
		output := []byte(fmt.Sprintf("\033[%d;%dH", row, col))
		return output, original, true
	case '[': // CSI
		output, original2, redraw, err := p.handleCSI(input)
		if err != nil {
			return output, append(original, original2...), redraw
		}
		return output, nil, redraw
	}

	return nil, original, false
}
