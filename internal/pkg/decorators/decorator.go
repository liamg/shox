package decorators

type Anchor uint8

const (
	AnchorTop Anchor = iota
	AnchorBottom
)

type Decorator interface {
	Draw(rows uint16, cols uint16)
	GetAnchor() Anchor
	GetHeight() (rows uint16)
}
