//go:build js && wasm

package lander

import (
	"fmt"
	"strings"
	"sync"
	"syscall/js"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"

	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/diffing"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

var document js.Value

func init() {
	document = js.Global().Get("document")
}

type DomEnvironment struct {
	sync.RWMutex

	root string

	tree *nodes.FuncNode

	prevContext context.Context
}

func RenderInto(rootNode *nodes.FuncNode, root string) (*DomEnvironment, error) {
	env := &DomEnvironment{
		root: root,
		tree: rootNode,
	}

	env.Lock()
	defer env.Unlock()

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

	var styles []string
	err := context.WithNewContext(e.Update, nil, func() error {
		styles = diffing.RecursivelyMount(e.handleDOMEvent, document, rootElem, e.tree)
		e.prevContext = context.CurrentContext
		return nil
	})
	if err != nil {
		return err
	}

	e.printTree(e.tree, 0)

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	stylesString, err := m.String("text/css", strings.Join(styles, " "))
	if err != nil {
		return fmt.Errorf("could not minify CSS styles from HTML nodes. %w", err)
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
	rootElem := document.Call("querySelector", e.root)
	if !rootElem.Truthy() {
		return fmt.Errorf("failed to find mount parent using query selector %q", e.root)
	}

	var styles []string
	err := context.WithNewContext(e.Update, e.prevContext, func() error {
		patches, renderedStyles, err := diffing.GeneratePatches(
			e.handleDOMEvent,
			nil,
			rootElem,
			e.tree,
			e.tree.Clone(),
		)
		if err != nil {
			return err
		}

		for _, patch := range patches {
			err := patch.Execute(document, &renderedStyles)
			if err != nil {
				return err
			}
		}

		styles = renderedStyles
		e.prevContext = context.CurrentContext
		return nil
	})
	if err != nil {
		return err
	}

	e.printTree(e.tree, 0)

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

func (e *DomEnvironment) handleDOMEvent(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		panic(fmt.Errorf("args should be at least 1 element, instead was: %#v", args))
	}

	jsEvent := args[0]

	event := events.NewDOMEvent(jsEvent, this)

	// acquire exclusive lock before we actually process event
	e.Lock()
	defer e.Unlock()
	err := listener(event)
	if err != nil {
		// Return the error message
		return err.Error()
	}

	return true
}

func (e *DomEnvironment) printTree(currentNode nodes.Node, layers int) {
	prefix := ""

	for i := 0; i < layers; i++ {
		prefix += "|--"
	}

	fmt.Printf("%s Node %p %T (%v)\n", prefix, currentNode, currentNode, currentNode)

	var children nodes.Children
	switch typedNode := currentNode.(type) {
	case *nodes.FuncNode:
		children = []nodes.Node{typedNode.RenderResult}
	case *nodes.HTMLNode:
		children = typedNode.Children
	default:
		return
	}

	for _, child := range children {
		if child == nil {
			continue
		}

		e.printTree(child, layers+1)
	}
}
