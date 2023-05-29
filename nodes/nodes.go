//go:build js && wasm

package nodes

// Child is a utility type that is interchangeable with Node. It defines a single child
// a node can receive.
type Child = Node

// Children is a utility type that is interchangeable with []Node. Is defines a slice of children that
// a node accepting children can receive.
type Children = []Child

// NodeType is an enum of node types assigned to the node implementations. This is not used, but
// may be in the future when type casting cannot be used.
type NodeType uint8

const (
	NoneType NodeType = iota
	HTMLNodeType
	TextNodeType
	FuncNodeType
	FragmentNodeType
)

// Node is a generic interface for a Node in the virtual DOM tree. All nodes should implement this
// interface through the baseNode concrete struct.
type Node interface {
	// ToString returns the node's content as valid HTML for rendering on the server side. Not used
	// at the moment.
	ToString() string

	// Diff checks if the current node is different to the other node. Will return true of the nodes are
	// different and false if they are the same. Each node is tasked with implementing their own version
	// of Diff.
	Diff(other Node) bool

	// Type returns the node's type as a NodeType enum value.
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
