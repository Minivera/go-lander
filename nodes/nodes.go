//go:build js && wasm

package nodes

type Child = Node
type Children = []Child

type NodeType uint8

const (
	NoneType     NodeType = 0
	HTMLNodeType NodeType = 1
	TextNodeType NodeType = 2
	FuncNodeType NodeType = 3
)

type Node interface {
	Position(parent Node)
	ToString() string
	Diff(other Node) bool
	Type() NodeType
}

type baseNode struct {
	Parent Node
}

func (n *baseNode) Position(parent Node) {
	n.Parent = parent
}

func (n *baseNode) ToString() string {
	return ""
}

func (n *baseNode) Diff(other Node) bool {
	otherAsBase, ok := other.(*baseNode)
	if !ok {
		return false
	}

	// Don't check for siblings since those may change without impacting this
	// node.
	return otherAsBase.Parent == n.Parent
}

func (n *baseNode) Type() NodeType {
	return NoneType
}
