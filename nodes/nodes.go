//go:build js && wasm

package nodes

type Child = Node
type Children = []Child

type NodeType uint8

const (
	NoneType NodeType = iota
	HTMLNodeType
	TextNodeType
	FuncNodeType
	FragmentNodeType
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
