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
	ToString() string
	Diff(other Node) bool
	Type() NodeType
}

type baseNode struct{}

func (n *baseNode) ToString() string {
	return ""
}

func (n *baseNode) Diff(_ Node) bool {
	return false
}

func (n *baseNode) Type() NodeType {
	return NoneType
}
