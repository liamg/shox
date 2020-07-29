package proxy

import (
	"bytes"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/liamg/shox/pkg/ansi"

	"github.com/liamg/shox/pkg/decorators"
)

// Proxy sits between the terminal and the shell and adds decorators such as a status bar to the output
type Proxy struct {
	output                []byte
	mutex                 sync.Mutex
	workChan              chan byte
	started               bool
	closeChan             chan struct{}
	closeOnce             sync.Once
	processCompletionChan chan struct{}
	decorators            []decorators.Decorator
	decMutex              sync.Mutex
	realRows              uint16
	realCols              uint16
	canRender             bool
	redrawChan            chan struct{}
}

// NewProxy creates a new proxy instance
func NewProxy() *Proxy {
	return &Proxy{
		workChan:              make(chan byte, 0xffff),
		closeChan:             make(chan struct{}),
		processCompletionChan: make(chan struct{}),
		canRender:             true,
		redrawChan:            make(chan struct{}, 1),
	}
}

// Start runs the proxy
func (p *Proxy) Start() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.started == true {
		return
	}
	go p.process()
	p.started = true
}

// DisableRendering prevents the proxy from rendering it's decorators
func (p *Proxy) DisableRendering() {
	p.canRender = false
}

// EnableRendering allows the proxy to render it's decorators
func (p *Proxy) EnableRendering() {
	p.canRender = true
}

// Close shuts down the proxy
func (p *Proxy) Close() {
	p.closeOnce.Do(func() {
		close(p.closeChan)
		<-p.processCompletionChan
		ansi.Reset()
	})
}

// Write writes data to the proxy, to be filtered
func (p *Proxy) Write(data []byte) (n int, err error) {
	if !p.started {
		return 0, fmt.Errorf("proxy not started")
	}
	for _, d := range data {
		p.workChan <- d
	}
	return len(data), nil
}

// Read reads data from the proxy, which has been filtered
func (p *Proxy) Read(data []byte) (n int, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.started {
		return 0, fmt.Errorf("proxy not started")
	}
	n, err = bytes.NewBuffer(p.output).Read(data)
	if err != nil && err != io.EOF {
		return n, err
	}
	p.output = p.output[n:]
	return n, nil
}

// HandleCoordinates converts coordinates received from the shell to those with an area for decorators reserved
func (p *Proxy) HandleCoordinates(row, col uint16) (outRow uint16, outCol uint16) {

	p.decMutex.Lock()
	defer p.decMutex.Unlock()

	for _, dec := range p.decorators {
		if !dec.IsVisible() {
			continue
		}
		rows := dec.GetHeight()
		switch dec.GetAnchor() {
		case decorators.AnchorTop:
			row += rows
		}
	}

	return row, col
}

// HandleResize takes the new dimensions and proxies them, returning a new size taking into account any decorators
func (p *Proxy) HandleResize(rows, cols uint16) (outRows uint16, outCols uint16) {

	p.decMutex.Lock()
	defer p.decMutex.Unlock()

	p.realRows = rows
	p.realCols = cols

	for _, dec := range p.decorators {
		if !dec.IsVisible() {
			continue
		}
		h := dec.GetHeight()

		if h >= rows {
			rows = 0
		} else {
			rows -= h
		}

	}

	return rows, cols
}

// AddDecorator adds a decorator, such as a status bar, to the proxy
func (p *Proxy) AddDecorator(d decorators.Decorator) {
	p.decMutex.Lock()
	defer p.decMutex.Unlock()
	p.decorators = append(p.decorators, d)
}

func (p *Proxy) process() {

	p.requestRedraw()

	tickDuration := time.Millisecond * 10
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	renderInterval := time.Second

	lastRender := time.Now()

	for {
		select {
		case b := <-p.workChan:

			if b == 0x1b {
				output, original, redraw := p.proxyANSICommand(p.workChan)
				if original != nil {
					p.writeOutput(append([]byte{b}, original...))
				}
				if output != nil {
					p.writeOutput(output)
				}

				// can requestRedraw even if command errored, as we may not support a command but still want to requestRedraw
				// if we know it affects our rendering
				if redraw {
					p.requestRedraw()
				}
			} else if b == '\n' {
				p.writeOutput([]byte{b})
				p.requestRedraw()
			} else {
				p.writeOutput([]byte{b})
			}
		case <-ticker.C:
			select {
			case <-p.redrawChan:
				p.redraw()
				lastRender = time.Now()
			default:
				if time.Since(lastRender) >= renderInterval {
					p.redraw()
					lastRender = time.Now()
				}
			}
		case <-p.closeChan:
			close(p.processCompletionChan)
			return

		}

	}

}

func (p *Proxy) ForceRedraw() {
	p.redraw()
}

func (p *Proxy) requestRedraw() {
	select {
	case p.redrawChan <- struct{}{}:
	default:
	}
}

func (p *Proxy) redraw() {
	p.decMutex.Lock()
	defer p.decMutex.Unlock()
	if !p.canRender {
		return
	}
	//save cursor pos
	p.writeOutput([]byte("\x1b[s"))
	for _, decorator := range p.decorators {
		if !decorator.IsVisible() {
			continue
		}
		decorator.Draw(p.realRows, p.realCols, p.writeOutput)
	}
	// restore cursor position
	p.writeOutput([]byte(fmt.Sprintf("\x1b[u")))
}

func (p *Proxy) writeOutput(data []byte) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.output = append(p.output, data...)
}
