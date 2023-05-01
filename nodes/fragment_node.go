package nodes

// FragmentNode is an implementation of the Node interface which implements the logic to handle
// and render fragments inside Lander. Fragments are slice of nodes, which are only handled in
// the virtual tree. They are reconciled to the closest DOM parent when mounted.
type FragmentNode struct {
	baseNode

	// Children is a slice of the children provided to this fragment.
	Children []Node
}

// NewFragmentNode creates a new fragment node with the provided information.
func NewFragmentNode(children []Node) *FragmentNode {
	return &FragmentNode{
		Children: children,
	}
}

func (n *FragmentNode) ToString() string {
	content := ""
	for _, child := range n.Children {
		content += child.ToString()
	}

	return content
}

func (n *FragmentNode) Diff(other Node) bool {
	otherAsFragment, ok := other.(*FragmentNode)
	if !ok {
		return true
	}

	return len(otherAsFragment.Children) != len(n.Children)
}

func (n *FragmentNode) Type() NodeType {
	return FragmentNodeType
}

// InsertChildren inserts a children at the provided position in the fragment's children.
// Returns an error if the children cannot be inserted. Inserts at the end if provided -1
// as the position.
func (n *FragmentNode) InsertChildren(node Node, position int) error {
	// Insert at the end on a -1
	if position < 0 {
		n.Children = append(n.Children, node)
		return nil
	}

	newChildren := make([]Node, len(n.Children)+1)
	index := 0
	for _, child := range n.Children {
		if index == position {
			newChildren[index] = node
			index++
		}
		newChildren[index] = child
		index++
	}

	n.Children = newChildren

	return nil
}

// ReplaceChildren replaces the provided node with the new node in the fragment's children.
// Returns an error if the children cannot be replaced.
func (n *FragmentNode) ReplaceChildren(old, new Node) error {
	for index, child := range n.Children {
		if child == old {
			n.Children[index] = new
			break
		}
	}
	return nil
}

// RemoveChildren removed the provided children, if found, from the fragment's children.
// Returns an error if the children cannot be inserted.
func (n *FragmentNode) RemoveChildren(node Node) error {
	newChildren := make([]Node, len(n.Children)-1)
	index := 0
	for _, child := range n.Children {
		if child == node {
			continue
		}
		newChildren[index] = child
		index++
	}

	n.Children = newChildren

	return nil
}
