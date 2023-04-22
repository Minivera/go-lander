//go:build js && wasm

package lander

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"

	"github.com/minivera/go-lander/diffing"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

var document js.Value

func init() {
	document = js.Global().Get("document")
}

type DomEnvironment struct {
	root string

	app  nodes.Node
	tree nodes.Node
}

func RenderInto(rootNode *nodes.FuncNode, root string) (*DomEnvironment, error) {
	env := &DomEnvironment{
		root: root,
		app:  rootNode,
	}

	err := env.renderIntoRoot()
	if err != nil {
		return nil, err
	}

	return env, nil
}

func (e *DomEnvironment) Update() error {
	err := e.patchDom()
	if err != nil {
		return err
	}

	return nil
}

func (e *DomEnvironment) renderIntoRoot() error {
	rootElem := document.Call("querySelector", e.root)
	if !rootElem.Truthy() {
		return fmt.Errorf("failed to find mount parent using query selector %q", e.root)
	}

	e.tree = e.generateTree(e.app)
	styles := diffing.RecursivelyMount(e.handleDOMEvent, document, rootElem, e.tree)

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	stylesString, err := m.String("text/css", strings.Join(styles, " "))
	if err != nil {
		return err
	}

	head := document.Call("querySelector", "head")
	if !head.Truthy() {
		return fmt.Errorf("failed to find head using query selector")
	}

	styleTag := document.Call("createElement", "style")
	styleTag.Set("id", "lander-style-tag")
	styleTag.Set("innerHTML", stylesString)
	head.Call("appendChild", styleTag)

	return nil
}

func (e *DomEnvironment) patchDom() error {
	patches, styles, err := diffing.GeneratePatches(e.handleDOMEvent, nil, e.tree, e.generateTree(e.app))
	if err != nil {
		return err
	}

	for _, patch := range patches {
		err := patch.Execute(document, &styles)
		if err != nil {
			return err
		}
	}

	styleTag := document.Call("querySelector", "#lander-style-tag")
	if !styleTag.Truthy() {
		return fmt.Errorf("failed to find the style selector, failing %s", "#lander-style-tag")
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	stylesString, err := m.String("text/css", strings.Join(styles, " "))
	if err != nil {
		return err
	}

	styleTag.Set("innerHTML", stylesString)

	return nil
}

func (e *DomEnvironment) generateTree(currentNode nodes.Node) nodes.Node {
	var toReturn nodes.Node
	var children []nodes.Node

	fmt.Printf("Generating %T, %v\n", currentNode, currentNode)
	// Check the current node's type
	switch typedNode := currentNode.(type) {
	case *nodes.FuncNode:
		// If the current node is a func node, we want to render it and "forget" it exists
		// replacing it with whatever it rendered.
		toReturn = typedNode.Render()
		children = []nodes.Node{toReturn}
	case *nodes.HTMLNode:
		// For all other nodes, use it as the node to return. We should get a tree of only "HTML" nodes.
		toReturn = typedNode
		children = typedNode.Children
	default:
		toReturn = typedNode
	}

	// Render all the children of the current node, if any
	for i, child := range children {
		// Render the current children, get the result
		renderResult := e.generateTree(child)

		if typedNode, ok := renderResult.(*nodes.FuncNode); ok {
			// If the child was another function node, then we should recursively render it until we
			// have a pure HTML node
			child = e.generateTree(typedNode)
		}

		// If the current node is an HTML node, replace the child in its children array with
		// the final child here. For most cases, that should do nothing, but for function nodes
		// it should replace it with the real final result.
		if typedNode, ok := currentNode.(*nodes.HTMLNode); ok {
			typedNode.Children[i] = renderResult
		}
	}

	// Return the final node, we should only have a pure HTML tree here
	return toReturn
}

func (e *DomEnvironment) handleDOMEvent(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		panic(fmt.Errorf("args should be at least 1 element, instead was: %#v", args))
	}

	jsEvent := args[0]

	event := events.NewDOMEvent(jsEvent, this)

	// acquire exclusive lock before we actually process event
	event.Lock()
	defer event.Unlock()
	err := listener(event)
	if err != nil {
		// Return the error message
		return err.Error()
	}

	return true
}
