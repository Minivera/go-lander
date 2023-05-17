package nodes

import (
	"reflect"

	"github.com/google/go-cmp/cmp"

	"github.com/minivera/go-lander/context"
)

// Props is a utility type for empty properties to assign to a component. Technically
// interchangeable with struct{}, this type is provided for convenience.
type Props struct{}

// FunctionComponent is the type definition for a function component's factory. This should be the
// definition given to a FuncNode when its created.
type FunctionComponent[T any] func(ctx context.Context, props T, children Children) Child

type noGenericFunctionComponent func(ctx context.Context, props interface{}, children Children) Child

// FuncNode is an implementation of the Node interface which implements the logic to handle
// and render components inside Lander.
type FuncNode struct {
	baseNode

	factory       noGenericFunctionComponent
	givenChildren []Node

	// Properties are the node's properties, which are passed to the factory on render.
	Properties interface{}

	// RenderResult is a reference to the result of the factory render, which is kept to allow
	// diffing later in the algorithm.
	RenderResult Node
}

// NewFuncNode creates a new component node with the provided information.
func NewFuncNode(factory noGenericFunctionComponent, props interface{}, givenChildren []Node) *FuncNode {
	return &FuncNode{
		Properties:    props,
		factory:       factory,
		givenChildren: givenChildren,
	}
}

// Render triggers the component's factory, passing the properties and children of the node.
// It will save the result in the node's memory for later diffs.
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

	if reflect.TypeOf(otherAsFunc.Properties) != reflect.TypeOf(n.Properties) {
		return true
	}

	if !cmp.Equal(otherAsFunc.Properties, n.Properties) {
		return true
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

// Clone clones the component node by creating a new node with all information and state provided,
// except the render result. This is necessary for the diffing algorithm as triggering a render for
// diffing would mutate the existing tree's render result.
func (n *FuncNode) Clone() *FuncNode {
	return &FuncNode{
		baseNode:      n.baseNode,
		factory:       n.factory,
		givenChildren: n.givenChildren,
		Properties:    n.Properties,
		RenderResult:  nil,
	}
}
