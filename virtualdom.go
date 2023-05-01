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
	"github.com/minivera/go-lander/internal"
	"github.com/minivera/go-lander/nodes"
)

var document js.Value

func init() {
	document = js.Global().Get("document")
}

// DomEnvironment is the lander DOM environment that stores the necessary information to allow mounting
// and rendering a lander app. Keep this environment in the main method or in global memory to avoid any
// memory loss.
type DomEnvironment struct {
	sync.RWMutex

	root string

	tree *nodes.FuncNode

	prevContext context.Context
}

// RenderInto renders the provided root component node into the given DOM root. The root selector must
// lead to a valid node, otherwise the mounting will error. The tree is only mounted in this method,
// no diffing will happen. Returns the mounted DOM environment, which can be used to trigger updates.
//
// This function is thread safe and will not allow any updates while the first mount is in progress. Event
// listeners or effects triggered during the mount process will have to wait.
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

// Update triggers the diffing process and updates the virtual and real DOM tree. The app provided to
// RenderInto will rerender fully and be diffed against the previously store tree. The diffing process
// generates a set of patches, which are executed in sequence against both the real and virtual DOM trees
// to update the stored tree with the new changes.
//
// This function is NOT thread safe and many allow other updates while another is in progress. Trigger an
// Update in an event listener to use the thread safe features of Lander.
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
		baseIndex := 0
		patches, renderedStyles, err := diffing.GeneratePatches(
			e.handleDOMEvent,
			nil,
			rootElem,
			&baseIndex,
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

	internal.Debugf("%s Node %p %T (%v)\n", prefix, currentNode, currentNode, currentNode)

	var children nodes.Children
	switch typedNode := currentNode.(type) {
	case *nodes.FuncNode:
		children = []nodes.Node{typedNode.RenderResult}
	case *nodes.FragmentNode:
		children = typedNode.Children
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
