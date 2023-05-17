//go:build js && wasm

package lander

import (
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

// Html creates an HTML node in the virtual tree using the provided tag, attributes, and children.
func Html(tag string, attributes nodes.Attributes, children nodes.Children) *nodes.HTMLNode {
	return nodes.NewHTMLNode(tag, attributes, children)
}

// Svg is a helper function to generate an HTML node with the SVG namespace set. This is necessary
// when handling SVGs, otherwise they will not render as expected.
func Svg(tag string, attributes nodes.Attributes, children nodes.Children) *nodes.HTMLNode {
	node := Html(tag, attributes, children)
	node.Namespace = "http://www.w3.org/2000/svg"
	return node
}

// Text creates a text node using the provided text. It will render as a text element in the tree and
// can only be used as a children.
func Text(text string) *nodes.TextNode {
	return nodes.NewTextNode(text)
}

// Component creates a function component node using the provided factory, props, and children. The factory
// will be executed on every render cycle with the most up-to-date props and children, and is expected to
// return a single child. Components also take a context, which provides hooks to lister to mount, render
// and unmount events. See Context.
//
// Lander expects a component as its first node, this component could then render an HTML element or text
// node, but only a component can be given to RenderInto.
func Component[T any](factory nodes.FunctionComponent[T], props T, children nodes.Children) *nodes.FuncNode {
	// Create an intermediary function so we hide the generic away. The generic is here only for
	// developer convenience.
	return nodes.NewFuncNode(func(ctx context.Context, props interface{}, children nodes.Children) nodes.Child {
		return factory(ctx, props.(T), children)
	}, props, children)
}

// Fragment creates a fragment node, which is a utility node that allows returning multiple children from
// function components. A fragment node does not appear in the DOM, its children are assigned to the closest
// DOM parent.
func Fragment(children nodes.Children) *nodes.FragmentNode {
	return nodes.NewFragmentNode(children)
}
