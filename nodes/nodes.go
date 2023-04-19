//go:build js && wasm

package nodes

type Child = Node
type Children = []Child

type Props = map[string]interface{}

type FunctionComponent func(attributes Props, children Children) Child

type Node interface {
	Position(parent Node)
	ToString() string
	Diff(other Node) bool
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
