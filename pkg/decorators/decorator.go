package decorators

// Anchor describes the layout of a decorator - i.e. whether it should render at the top/bottom of the terminal
type Anchor uint8

const (
	AnchorTop    Anchor = iota //AnchorTop renders at the top of the terminal
	AnchorBottom               //AnchorBottom renders at the bottom of the terminal
)

// Decorator is an entity which modifies the terminal output in a desirable way
type Decorator interface {
	Draw(rows uint16, cols uint16) // Draw renders the decorator to StdOut
	GetAnchor() Anchor             // GetAnchor returns the anchor e.g. Top/Bottom
	GetHeight() (rows uint16)      // GetHeight returns the height of the decorator in terminal character rows
}
