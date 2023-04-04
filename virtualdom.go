//go:build js && wasm

package lander

import (
	"fmt"
	"runtime/debug"
	"strings"
	"syscall/js"

	"github.com/minivera/go-lander/events"

	"github.com/minivera/go-lander/diffing"

	"github.com/minivera/go-lander/utils"

	"github.com/minivera/go-lander/nodes"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

var document js.Value

func init() {
	document = js.Global().Get("document")
}

type DomEnvironment struct {
	root string
	app  nodes.FuncNode[map[string]interface{}]

	tree nodes.Node
}

func New(root string, rootNode nodes.FuncNode[map[string]interface{}]) *DomEnvironment {
	env := &DomEnvironment{
		root: root,
		app:  rootNode,
	}

	return env
}

func (e *DomEnvironment) Render() error {
	e.app.Render()
	e.tree = e.app.RenderResult

	e.recursivelyPosition(e.tree)
	err := e.mountToDom()
	if err != nil {
		return err
	}

	return nil
}

func (e *DomEnvironment) Update() error {
	e.app.Render()
	newTree := e.app.RenderResult

	err := e.patchDom(newTree)
	if err != nil {
		return err
	}

	return nil
}

func (e *DomEnvironment) recursivelyPosition(currentNode nodes.Node) {
	if currentNode == nil {
		panic("no component provided to lander environment")
	}

	var children []nodes.Node
	switch typedNode := currentNode.(type) {
	case *nodes.TextNode:
		return
	case *nodes.HTMLNode:
		children = typedNode.Children
	}

	for index, child := range children {
		if child == nil {
			continue
		}

		if len(children) <= 1 {
			child.Position(currentNode, nil, nil)
			continue
		}

		if index <= 0 {
			child.Position(currentNode, children[index+1], nil)
		} else if index < len(children)-1 {
			child.Position(currentNode, children[index+1], children[index-1])
		} else {
			child.Position(currentNode, nil, children[index-1])
		}

		e.recursivelyPosition(child)
	}
}

func (e *DomEnvironment) mountToDom() error {
	rootElem := document.Call("querySelector", e.root)
	if !rootElem.Truthy() {
		return fmt.Errorf("failed to find mount parent using query selector %q", e.root)
	}

	styles, err := e.recursivelyMount(rootElem, e.tree)
	if err != nil {
		return err
	}

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

func (e *DomEnvironment) recursivelyMount(lastElement js.Value, currentNode nodes.Node) ([]string, error) {
	if currentNode == nil {
		return []string{}, nil
	}

	add := false
	domElement := lastElement
	var styles []string
	var children []nodes.Node

	switch typedNode := currentNode.(type) {
	case *nodes.HTMLNode:
		add = true
		domElement = utils.NewHTMLElement(document, typedNode)
		typedNode.Mount(domElement)

		for key, listener := range typedNode.EventListeners {
			domElement.Call("addEventListener", key, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				return e.handleDOMEvent(listener, this, args)
			}))
		}

		children = typedNode.Children

		for _, style := range typedNode.Styles {
			styles = append(styles, style)
		}
	case *nodes.TextNode:
		add = true
		domElement = document.Call("createTextNode", typedNode.Text)
		typedNode.Mount(domElement)
	}

	for _, child := range children {
		if child == nil {
			continue
		}

		childStyles, err := e.recursivelyMount(domElement, child)
		if err != nil {
			return styles, err
		}

		for _, style := range childStyles {
			styles = append(styles, style)
		}
	}

	if add {
		lastElement.Call("appendChild", domElement)
	}

	return styles, nil
}

func (e *DomEnvironment) patchDom(newTree nodes.Node) error {
	rootElem := document.Call("querySelector", e.root)
	if !rootElem.Truthy() {
		return fmt.Errorf("failed to find mount parent using query selector %q", rootElem)
	}

	patches, err := diffing.GeneratePatches(nil, e.tree, newTree)
	if err != nil {
		return err
	}

	for _, patch := range patches {
		err := patch.Execute(document)
		if err != nil {
			return err
		}
	}

	return nil
}

// TODO: This code has been borrowed from vugu, find a way to make it our own
func (e *DomEnvironment) handleDOMEvent(listener events.EventListener, this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		panic(fmt.Errorf("args should be at least 1 element, instead was: %#v", args))
	}

	// TODO: give this more thought - but for now we just do a non-blocking push to the
	// eventWaitCh, telling the render loop that a render is required, but if a bunch
	// of them stack up we don't wait
	defer func() {

		if r := recover(); r != nil {
			fmt.Println("handleDOMEvent caught panic", r)
			debug.PrintStack()

			// in error case send false to tell event loop to exit
			select {
			case e.eventWaitCh <- false:
			default:
			}
			return

		}

		// in normal case send true to the channel to tell the event loop it should render
		select {
		case e.eventWaitCh <- true:
		default:
		}
	}()

	jsEvent := args[0]

	typeName := jsEvent.Get("type").String()

	key := "data-lander-event-" + typeName + "-id"
	funcHash := this.Get(key).String()
	var funcID uint64
	_, err := fmt.Sscanf(funcHash, "%d", &funcID)
	if err != nil {
		panic(fmt.Errorf("lander could not retreive and convert the dom event key for some reason, %v", err))
	}

	if funcID == 0 {
		panic(fmt.Errorf("looking for %q on 'this' found %q which parsed into value 0 - cannot find the appropriate function to route to", key, funcHash))
	}

	eventDef, ok := e.eventsToHash[funcID]
	if !ok {
		panic(fmt.Errorf("nothing found in eventsToHash for %d", funcID))
	}

	var node Node
	if val, ok := e.instancesToHash[eventDef.nodeHash]; ok {
		node = val
	} else {
		fmt.Printf("could not find the virtual node associated with the dom event, maybe it was removed?")
	}

	event := &DOMEvent{
		browserEvent: jsEvent,
		this:         this,
		//environment: e.
	}

	// acquire exclusive lock before we actually process event
	e.eventRWMU.Lock()
	defer e.eventRWMU.Unlock()
	err = eventDef.listener(node, event)
	if err != nil {
		// Return the error message
		return err.Error()
	}

	err = e.Update()
	if err != nil {
		// Return the error message if we couldn't update
		return err.Error()
	}

	return true
}
