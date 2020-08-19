package proxy

import (
	"fmt"
	"strconv"
	"strings"
)

type csiHandler func(params []string, intermediate string) (output []byte, redraw bool, err error)

func (p *Proxy) handleCSI(pty chan byte) (output []byte, original []byte, redraw bool, err error) {

	var csiSequences = map[byte]csiHandler{
		'd': p.csiLinePositionAbsolute,
		'f': p.csiCursorPositionHandler,
		'n': p.csiDeviceStatusReportHandler,
		'G': p.csiCursorCharacterAbsoluteHandler,
		'H': p.csiCursorPositionHandler,
		'h': p.csiSetModeHandler,
		'l': p.csiResetModeHandler,
		'J': p.csiEraseInDisplayHandler,
		'r': p.csiSetMarginHandler,
		's': p.csiSavePositionHandler,
		'u': p.csiRestorePositionHandler,
	}

	var final byte
	var b byte
	param := ""
	intermediate := ""
CSI:
	for {
		b = <-pty
		original = append(original, b)
		switch true {
		case b >= 0x30 && b <= 0x3F:
			param = param + string(b)
		case b >= 0x20 && b <= 0x2F:
			intermediate += string(b)
		case b >= 0x40 && b <= 0x7e:
			final = b
			break CSI
		}
	}

	params := strings.Split(param, ";")
	if param == "" {
		params = []string{}
	}

	handler, ok := csiSequences[final]
	if !ok {
		return nil, original, false, ErrorCommandNotSupported
	}

	output, redraw, err = handler(params, intermediate)
	if err != nil {
		return output, original, redraw, ErrorCommandNotSupported
	}
	return output, nil, redraw, nil
}

func (p *Proxy) csiDeviceStatusReportHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {

	if !p.canRender {
		return nil, false, ErrorCommandNotSupported
	}

	switch params[0] {
	case "6": // report cursor position
		// TODO "\x1b[%d;%dR", keep track of position and forward? Or proxy command responses from terminal -> shell?
	}

	return nil, false, ErrorCommandNotSupported
}

func (p *Proxy) csiCursorCharacterAbsoluteHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {

	if !p.canRender {
		return nil, false, ErrorCommandNotSupported
	}

	col := 1
	if len(params) > 0 {
		var err error
		col, err = strconv.Atoi(params[0])
		if err != nil || params[0] == "" {
			col = 1
		}
	}

	_, adjustedCol := p.HandleCoordinates(0, uint16(col))
	return []byte(fmt.Sprintf("\033[%dG", adjustedCol)), false, nil
}

func (p *Proxy) csiCursorPositionHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {

	if !p.canRender {
		return nil, false, ErrorCommandNotSupported
	}

	x, y := 1, 1
	if len(params) == 2 {
		var err error
		if params[0] != "" {
			y, err = strconv.Atoi(string(params[0]))
			if err != nil || y < 1 {
				y = 1
			}
		}
		if params[1] != "" {
			x, err = strconv.Atoi(string(params[1]))
			if err != nil || x < 1 {
				x = 1
			}
		}
	}

	row, col := p.HandleCoordinates(uint16(y), uint16(x))
	return []byte(fmt.Sprintf("\033[%d;%dH", row, col)), false, nil
}

func (p *Proxy) csiLinePositionAbsolute(params []string, intermediate string) (output []byte, redraw bool, err error) {

	if !p.canRender {
		return nil, false, ErrorCommandNotSupported
	}

	row := 1
	if len(params) > 0 {
		var err error
		row, err = strconv.Atoi(params[0])
		if err != nil || row < 1 {
			row = 1
		}
	}

	newRow, _ := p.HandleCoordinates(uint16(row), 0)
	return []byte(fmt.Sprintf("\033[%dd", newRow)), false, nil
}

// CSI Pt Pb r
func (p *Proxy) csiSetMarginHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {
	// pass through command, and inject reset position afterwards
	row, col := p.HandleCoordinates(1, 1)
	return []byte(fmt.Sprintf("\033[%d;%dH", row, col)), true, ErrorCommandNotSupported
}

// CSI Ps J
func (p *Proxy) csiEraseInDisplayHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {

	if !p.canRender {
		return nil, false, ErrorCommandNotSupported
	}

	n := "0"
	if len(params) > 0 {
		n = params[0]
	}

	switch n {
	case "2", "3":
		row, col := p.HandleCoordinates(1, 1)
		return []byte(fmt.Sprintf("\033[%d;%dH", row, col)), true, ErrorCommandNotSupported
	}

	return nil, false, ErrorCommandNotSupported
}

func (p *Proxy) csiResetModeHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {
	return p.csiSetModes(params, false)
}

func (p *Proxy) csiSetModeHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {
	return p.csiSetModes(params, true)
}

func (p *Proxy) csiSetModes(modes []string, enabled bool) (output []byte, redraw bool, err error) {
	if len(modes) == 0 {
		return nil, false, ErrorCommandNotSupported
	}
	if len(modes) == 1 {
		return p.csiSetMode(modes[0], enabled)
	}
	// should we propagate DEC prefix?
	const decPrefix = '?'
	isDec := len(modes[0]) > 0 && modes[0][0] == decPrefix

	// iterate through params, propagating DEC prefix to subsequent elements
	for i, v := range modes {
		updatedMode := v
		if i > 0 && isDec {
			updatedMode = string(decPrefix) + v
		}
		_, forceRedraw, _ := p.csiSetMode(updatedMode, enabled)
		redraw = redraw || forceRedraw
	}

	return nil, redraw, ErrorCommandNotSupported
}

func (p *Proxy) csiSetMode(modeStr string, enabled bool) (output []byte, redraw bool, err error) {

	switch modeStr {
	case "?47", "?1047", "?1049":
		// switching to alt buffer
		if enabled {
			// switched to alt buffer - disable rendering for a while
			p.DisableRendering()
		} else {
			p.EnableRendering()
			redraw = enabled
		}
	}

	return nil, redraw, ErrorCommandNotSupported
}

// CSI Ps s
func (p *Proxy) csiSavePositionHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {
	p.pauseDrawing = true
	return nil, false, ErrorCommandNotSupported
}

// CSI Ps u
func (p *Proxy) csiRestorePositionHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {
	p.pauseDrawing = false
	return nil, false, ErrorCommandNotSupported
}
