package go_lander

type DomEnvironment struct {
	root string
	tree Node
}

func NewLander(root string, rootNode Node) *DomEnvironment {
	return &DomEnvironment{
		root: root,
		tree: rootNode,
	}
}

func (e *DomEnvironment) Mount() error {
	err := walkTree(e.tree, buildNode)
	if err != nil {
		return err
	}

	err = mountToDom(e.root, e.tree)
	if err != nil {
		return err
	}

	return nil
}

func buildNode(currentNode Node) error {
	err := currentNode.Create()
	if err != nil {
		return err
	}

	err = currentNode.Render()
	if err != nil {
		return err
	}

	return nil
}
