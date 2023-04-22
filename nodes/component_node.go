package nodes

import "github.com/minivera/go-lander/context"

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

func (n *FuncNode) Render(ctx context.Context) Node {
	return n.factory(ctx, n.Properties, n.givenChildren)
}

func (n *FuncNode) Diff(other Node) bool {
	otherAsFunc, ok := other.(*FuncNode)
	if !ok {
		return true
	}

	if &otherAsFunc.factory != &n.factory {
		return true
	}

	if len(otherAsFunc.Properties) != len(n.Properties) {
		return true
	}

	for key, val := range n.Properties {
		otherVal, ok := otherAsFunc.Properties[key]
		if !ok {
			return true
		}

		if val != otherVal {
			return true
		}
	}

	// We check if any of the given children were dirty in the general diff code
	if len(otherAsFunc.givenChildren) != len(n.givenChildren) {
		return true
	}

	return false
}
