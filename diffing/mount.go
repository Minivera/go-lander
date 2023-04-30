package diffing

import (
	"fmt"
	"syscall/js"

	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

func RecursivelyMount(listenerFunc func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{},
	document js.Value, lastElement js.Value, currentNode nodes.Node) []string {

	if currentNode == nil {
		return []string{}
	}

	add := false
	domElement := lastElement
	var styles []string
	var children []nodes.Node

	fmt.Printf("Mounting %T node, %v\n", currentNode, currentNode)
	switch typedNode := currentNode.(type) {
	case *nodes.FuncNode:
		// If the current node is a func node, we want to render it and keep going
		// so we eventually hit a normal HTML node.
		context.RegisterComponent(typedNode)
		context.RegisterComponentContext("mount", typedNode)
		context.RegisterComponentContext("render", typedNode)
		children = []nodes.Node{typedNode.Render(context.CurrentContext)}
	case *nodes.FragmentNode:
		children = typedNode.Children
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
		return []string{}
	}

	for _, child := range children {
		if child == nil {
			continue
		}

		childStyles := RecursivelyMount(listenerFunc, document, domElement, child)
		for _, style := range childStyles {
			styles = append(styles, style)
		}
	}

	if add {
		lastElement.Call("appendChild", domElement)
	}

	return styles
}
