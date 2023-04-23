package diffing

import (
	"syscall/js"

	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

func RecursivelyMount(listenerFunc func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{},
	document js.Value, lastElement js.Value, currentNode nodes.Node) (nodes.Node, []string) {

	if currentNode == nil {
		return currentNode, []string{}
	}

	add := false
	domElement := lastElement
	var styles []string
	var children []nodes.Node
	toReturn := currentNode

	switch typedNode := currentNode.(type) {
	case *nodes.FuncNode:
		// If the current node is a func node, we want to render it and "forget" it exists
		// replacing it with whatever it rendered.
		context.RegisterComponentContext("mount", typedNode)
		context.RegisterComponentContext("render", typedNode)
		toReturn = typedNode.Render(context.CurrentContext)
		children = []nodes.Node{toReturn}
	case *nodes.HTMLNode:
		add = true
		domElement = nodes.NewHTMLElement(document, typedNode)
		typedNode.Mount(domElement)

		for event, listener := range typedNode.EventListeners {
			listener.Wrapper = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				return listenerFunc(listener.Func, this, args)
			})
			domElement.Call("addEventListener", event, listener.Wrapper)
		}

		children = typedNode.Children

		for _, style := range typedNode.Styles {
			styles = append(styles, style)
		}
	case *nodes.TextNode:
		add = true
		domElement = document.Call("createTextNode", typedNode.Text)
		typedNode.Mount(domElement)
	default:
		return toReturn, []string{}
	}

	for i, child := range children {
		if child == nil {
			continue
		}

		child.Position(currentNode)

		renderResult, childStyles := RecursivelyMount(listenerFunc, document, domElement, child)

		if renderResult.Type() == nodes.FuncNodeType {
			typedNode := any(renderResult).(*nodes.FuncNode)
			// If the child was another function node, then we should recursively render it until we
			// have a pure HTML node
			child, _ = RecursivelyMount(listenerFunc, document, domElement, typedNode)
		}

		// If the current node is an HTML node, replace the child in its children array with
		// the final child here. For most cases, that should do nothing, but for function nodes
		// it should replace it with the real final result.
		if typedNode, ok := currentNode.(*nodes.HTMLNode); ok {
			typedNode.Children[i] = renderResult
		}

		for _, style := range childStyles {
			styles = append(styles, style)
		}
	}

	if add {
		lastElement.Call("appendChild", domElement)
	}

	return toReturn, styles
}
