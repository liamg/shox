package proxy

import (
	"fmt"
	"strconv"
	"strings"
)

type csiHandler func(params []string, intermediate string) (output []byte, redraw bool, err error)

func (proxy *Proxy) handleCSI(pty chan byte) (output []byte, discard []byte, redraw bool, err error) {

	var csiSequences = map[byte]csiHandler{
		'd': proxy.csiLinePositionAbsolute,
		'f': proxy.csiCursorPositionHandler,
		'n': proxy.csiDeviceStatusReportHandler,
		'G': proxy.csiCursorCharacterAbsoluteHandler,
		'H': proxy.csiCursorPositionHandler,
		'h': proxy.csiSetModeHandler,
		'l': proxy.csiResetModeHandler,
		'J': proxy.csiEraseInDisplayHandler,
		//'K': proxy.csiEraseInLineHandler,
		//'r': proxy.csiSetMarginsHandler,

		//'X': proxy.csiEraseCharactersHandler,
	}

	var final byte
	var b byte
	param := ""
	intermediate := ""
CSI:
	for {
		b = <-pty
		discard = append(discard, b)
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
		return nil, discard, false, ErrorCommandNotSupported
	}

	output, redraw, err = handler(params, intermediate)
	if err != nil {
		return nil, discard, redraw, ErrorCommandNotSupported
	}
	return output, nil, redraw, nil
}

// TODO - below this line

func (proxy *Proxy) csiDeviceStatusReportHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {

	if len(params) == 0 {
		return nil, false, fmt.Errorf("Missing Device Status Report identifier")
	}

	switch params[0] {
	case "5":
		return []byte("\x1b[0n"), false, nil // everything is cool
	case "6": // report cursor position
		/*
			_ = terminal.Write([]byte(fmt.Sprintf(
					"\x1b[%d;%dR",
					terminal.ActiveBuffer().CursorLine()+1,
					terminal.ActiveBuffer().CursorColumn()+1,
				)))
		*/

		// TODO keep track of position and forward? Or proxy command responses from terminal -> shell?
		return nil, false, fmt.Errorf("Not supported yet")

	default:
		return nil, false, ErrorCommandNotSupported
	}
}

func (proxy *Proxy) csiCursorCharacterAbsoluteHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {
	col := 1
	if len(params) > 0 {
		var err error
		col, err = strconv.Atoi(params[0])
		if err != nil || params[0] == "" {
			col = 1
		}
	}

	_, adjustedCol := proxy.HandleCoordinates(0, uint16(col))
	return []byte(fmt.Sprintf("\033[%dG", adjustedCol)), false, nil
}

func (proxy *Proxy) csiCursorPositionHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {
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

	row, col := proxy.HandleCoordinates(uint16(y), uint16(x))
	return []byte(fmt.Sprintf("\033[%d;%dH", row, col)), false, nil
}

func (proxy *Proxy) csiLinePositionAbsolute(params []string, intermediate string) (output []byte, redraw bool, err error) {
	row := 1
	if len(params) > 0 {
		var err error
		row, err = strconv.Atoi(params[0])
		if err != nil || row < 1 {
			row = 1
		}
	}

	newRow, _ := proxy.HandleCoordinates(uint16(row), 0)
	return []byte(fmt.Sprintf("\033[%dd", newRow)), false, nil
}

func (proxy *Proxy) csiResetModeHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {
	return proxy.csiSetModes(params, false)
}

func (proxy *Proxy) csiSetModeHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {
	return proxy.csiSetModes(params, true)
}

// CSI Ps J
func (proxy *Proxy) csiEraseInDisplayHandler(params []string, intermediate string) (output []byte, redraw bool, err error) {

	n := "0"
	if len(params) > 0 {
		n = params[0]
	}

	switch n {
	case "2", "3":
		return nil, true, ErrorCommandNotSupported
	}

	return nil, false, ErrorCommandNotSupported
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

	/*
	   Mouse support
	   		#define SET_X10_MOUSE               9
	        #define SET_VT200_MOUSE             1000
	        #define SET_VT200_HIGHLIGHT_MOUSE   1001
	        #define SET_BTN_EVENT_MOUSE         1002
	        #define SET_ANY_EVENT_MOUSE         1003
	        #define SET_FOCUS_EVENT_MOUSE       1004
	        #define SET_EXT_MODE_MOUSE          1005
	        #define SET_SGR_EXT_MODE_MOUSE      1006
	        #define SET_URXVT_EXT_MODE_MOUSE    1015
	        #define SET_ALTERNATE_SCROLL        1007
	*/

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

/*

// DECSTBM
func (proxy *Proxy) csiSetMarginsHandler(params []string, intermediate string) (output []byte, requestRedraw bool, err error) {
	top := 1
	bottom := int(terminal.ActiveBuffer().ViewHeight())

	if len(params) > 2 {
		return fmt.Errorf("Not set margins")
	}

	if len(params) > 0 {
		var err error
		top, err = strconv.Atoi(params[0])
		if err != nil || top < 1 {
			top = 1
		}

		if len(params) > 1 {
			var err error
			bottom, err = strconv.Atoi(params[1])
			if err != nil || bottom > int(terminal.ActiveBuffer().ViewHeight()) || bottom < 1 {
				bottom = int(terminal.ActiveBuffer().ViewHeight())
			}
		}
	}
	top--
	bottom--

	terminal.ActiveBuffer().SetVerticalMargins(uint(top), uint(bottom))
	terminal.ActiveBuffer().SetPosition(0, 0)

	return nil
}



// CSI Ps K
func (proxy *Proxy) csiEraseInLineHandler(params []string, intermediate string) (output []byte, requestRedraw bool, err error) {

	n := "0"
	if len(params) > 0 {
		n = params[0]
	}

	switch n {
	case "0", "": //erase adter cursor
		terminal.ActiveBuffer().EraseLineFromCursor()
	case "1": // erase to cursor inclusive
		terminal.ActiveBuffer().EraseLineToCursor()
	case "2": // erase entire
		terminal.ActiveBuffer().EraseLine()
	default:
		return fmt.Errorf("Unsupported EL: CSI %s K", n)
	}
	return nil
}
*/
