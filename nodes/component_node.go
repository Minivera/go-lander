package nodes

type FuncNode struct {
	baseNode

	factory       FunctionComponent
	givenChildren []Node

	Properties Props
}

func NewFuncNode(factory FunctionComponent, attributes Props, givenChildren []Node) *FuncNode {
	return &FuncNode{
		Properties:    attributes,
		factory:       factory,
		givenChildren: givenChildren,
	}
}

func (n *FuncNode) Update(newAttributes Props, newChildren []Node) {
	n.Properties = newAttributes
	n.givenChildren = newChildren
}

func (n *FuncNode) Render() Node {
	return n.factory(n.Properties, n.givenChildren)
}

func (n *FuncNode) Diff(other Node) bool {
	otherAsFunc, ok := other.(*FuncNode)
	if !ok {
		return false
	}

	if &otherAsFunc.factory != &n.factory {
		return false
	}

	if len(otherAsFunc.Properties) != len(n.Properties) {
		return false
	}

	for key, val := range n.Properties {
		otherVal, ok := otherAsFunc.Properties[key]
		if !ok {
			return false
		}

		if val != otherVal {
			return false
		}
	}

	// We check if any of the given children were dirty in the general diff code
	if len(otherAsFunc.givenChildren) != len(n.givenChildren) {
		return false
	}

	return true
}
