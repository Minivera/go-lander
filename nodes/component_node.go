package nodes

import (
	"reflect"

	"github.com/minivera/go-lander/context"
)

type Props = map[string]interface{}
type FunctionComponent func(ctx context.Context, attributes Props, children Children) Child

type FuncNode struct {
	baseNode

	factory       FunctionComponent
	givenChildren []Node

	Properties   Props
	RenderResult Node
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
	n.RenderResult = n.factory(ctx, n.Properties, n.givenChildren)
	return n.RenderResult
}

func (n *FuncNode) Diff(other Node) bool {
	otherAsFunc, ok := other.(*FuncNode)
	if !ok {
		return true
	}

	if context.CurrentContext.IsDirty() {
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

		if reflect.TypeOf(val).Comparable() && reflect.TypeOf(otherVal).Comparable() {
			if val != otherVal {
				return true
			}
		}
	}

	// We check if any of the given children were dirty in the general diff code
	if len(otherAsFunc.givenChildren) != len(n.givenChildren) {
		return true
	}

	return false
}

func (n *FuncNode) Type() NodeType {
	return FuncNodeType
}

func (n *FuncNode) Clone() *FuncNode {
	return &FuncNode{
		baseNode:      n.baseNode,
		factory:       n.factory,
		givenChildren: n.givenChildren,
		Properties:    n.Properties,
		RenderResult:  nil,
	}
}
