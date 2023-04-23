package nodes

import "github.com/minivera/go-lander/context"

type Props = map[string]interface{}
type FunctionComponent func(ctx context.Context, attributes Props, children Children) Child

type FuncNode struct {
	baseNode

	factory       FunctionComponent
	givenChildren []Node

	Properties Props
}

func NewFuncNode(factory FunctionComponent, props Props, givenChildren []Node) *FuncNode {
	return &FuncNode{
		Properties:    props,
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

func (n *FuncNode) Type() NodeType {
	return FuncNodeType
}
