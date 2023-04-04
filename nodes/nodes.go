//go:build js && wasm

package nodes

type Node interface {
	Position(parent, next, prev Node)
	Render() Node
	ToString() string
	Diff(other Node) bool
}

type baseNode struct {
	Parent, NextSibling, PreviousSibling Node
}

func (n *baseNode) Position(parent, next, prev Node) {
	n.Parent = parent
	n.NextSibling = next
	n.PreviousSibling = prev
}

func (n *baseNode) Render() Node {
	return n
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
