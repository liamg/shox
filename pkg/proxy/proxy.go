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

func NewProxy() *Proxy {
	return &Proxy{
		workChan:              make(chan byte, 0xffff),
		closeChan:             make(chan struct{}),
		processCompletionChan: make(chan struct{}),
		canRender:             true,
		redrawChan:            make(chan struct{}, 1),
	}
}

func (p *Proxy) Start() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.started == true {
		return
	}
	go p.process()
	p.started = true
}

func (p *Proxy) DisableRendering() {
	p.canRender = false
}

func (p *Proxy) EnableRendering() {
	p.canRender = true
}

func (p *Proxy) Close() {
	p.closeOnce.Do(func() {
		close(p.closeChan)
		<-p.processCompletionChan
		ansi.Reset()
	})
}

func (p *Proxy) Write(data []byte) (n int, err error) {
	if !p.started {
		return 0, fmt.Errorf("proxy not started")
	}
	for _, d := range data {
		p.workChan <- d
	}
	return len(data), nil
}

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

func (p *Proxy) HandleCoordinates(row, col uint16) (outRow uint16, outCol uint16) {

	p.decMutex.Lock()
	defer p.decMutex.Unlock()

	for _, dec := range p.decorators {
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
		h := dec.GetHeight()

		if h >= rows {
			rows = 0
		} else {
			rows -= h
		}

	}

	return rows, cols
}

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
					p.writeOutput(append([]byte{b}, output...))
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
	for _, decorator := range p.decorators {
		decorator.Draw(p.realRows, p.realCols)
	}
}

func (p *Proxy) writeOutput(data []byte) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.output = append(p.output, data...)
}
