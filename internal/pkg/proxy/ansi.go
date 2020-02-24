package proxy

import (
	"fmt"
)

/*
var ansiSequenceMap = map[byte]escapeSequenceHandler{
	'[': proxy.handleCSI,
		'D': indexHandler,
		'E': nextLineHandler, // NEL
		'H': tabSetHandler,   // HTS
		'M': reverseIndexHandler,
		'P': sixelHandler,
		'c': risHandler, //RIS
		'#': screenStateHandler,
		'(': scs0Handler, // select character set into G0
		')': scs1Handler, // select character set into G1
}
*/

type escapeSequenceHandler func(pty chan byte) (output []byte, discard []byte, redraw bool, err error)

var ErrorCommandNotSupported = fmt.Errorf("command not supported")

func (p *Proxy) proxyANSICommand(input chan byte) (discard []byte, requiredRedraw bool, err error) {

	b := <-input
	discard = append(discard, b)

	if b == '[' {
		output, discard2, redraw, err := p.handleCSI(input)
		if err != nil {
			return append(discard, discard2...), redraw, ErrorCommandNotSupported
		}

		p.writeOutput(output)
		return nil, redraw, nil
	}

	return discard, false, ErrorCommandNotSupported
}
