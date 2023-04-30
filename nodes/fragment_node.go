package nodes

type FragmentNode struct {
	baseNode

	Children []Node
}

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

func (n *FragmentNode) ReplaceChildren(old, new Node) error {
	for index, child := range n.Children {
		if child == old {
			n.Children[index] = new
			break
		}
	}
	return nil
}

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
