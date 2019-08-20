// +build js,wasm

package lander

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"syscall/js"

	"github.com/cespare/xxhash"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

type jsValue interface {
	Call(string, ...interface{}) js.Value
	Get(string) js.Value
	Index(int) js.Value
	Set(string, interface{})
	Truthy() bool
}

var document js.Value

func init() {
	document = js.Global().Get("document")
}

type DomEnvironment struct {
	root string
	tree Node

	instancesToHash map[uint64]Node
	eventsToHash    map[uint64]*wasmEvent
	eventsCallback  js.Func

	// TODO: This has been borrowed from vugu, find a way to make it our own
	eventWaitCh chan bool    // events send to this and EventWait receives from it
	eventRWMU   sync.RWMutex // make sure Render and event handling are not attempted at the same time (not totally sure if this is necessary in terms of the wasm threading model but enforce it with a rwmutex all the same)
	eventEnv    *eventEnv    // our EventEnv implementation that exposes eventRWMU and eventWaitCh to events in a clean way

	patches []patch
}

func NewLander(root string, rootNode Node) *DomEnvironment {
	env := &DomEnvironment{
		root: root,
		tree: rootNode,

		instancesToHash: make(map[uint64]Node, 1024),
		eventsToHash:    make(map[uint64]*wasmEvent, 32),
		eventWaitCh:     make(chan bool, 64),

		patches: []patch{},
	}

	env.eventsCallback = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		return env.handleDOMEvent(this, args)
	})

	env.eventEnv = &eventEnv{
		rwmu:            &env.eventRWMU,
		requestRenderCH: env.eventWaitCh,
	}

	return env
}

// EventWait will block until an event occurs and will return after the event is completed.
// This is our first attempt at making a "render loop".  Will return false if the JSEnv
// becomes invalid and should exit.
// TODO: This has been borrowed from vugu, try to make it our own
func (e *DomEnvironment) EventWait() (ok bool) {
	ok = js.Global().Get("document").Truthy()
	if !ok {
		return
	}

	// FIXME: this should probably have some sort of "debouncing" on it to handle the case of
	// several events in rapid succession causing multiple renders - maybe we read from eventWaitCH
	// continuously until it's empty, with a max of like 20ms pause between each or something, and then
	// only return after we don't see anything for that time frame.

	ok = <-e.eventWaitCh
	return
}

func (e *DomEnvironment) Mount() error {
	err := e.buildNodeRecursively(e.tree, "[root]")
	if err != nil {
		return err
	}

	err = e.mountToDom(e.root, e.tree)
	if err != nil {
		return err
	}

	return nil
}

func (e *DomEnvironment) Update() error {
	newTree := e.tree.Clone()

	err := e.buildNodeRecursively(newTree, "[root]")
	if err != nil {
		return err
	}

	err = e.patchDom(newTree)
	if err != nil {
		return err
	}

	e.patches = []patch{}

	return nil
}

func (e *DomEnvironment) buildNodeRecursively(currentNode Node, positionString string) error {
	if currentNode == nil {
		return nil
	}

	id := hashPosition(positionString)

	if _, ok := e.instancesToHash[id]; !ok {
		err := currentNode.Create(id)
		if err != nil {
			return err
		}
	} else {
		// Reconciliate the ids
		currentNode.SetID(id)
	}

	err := currentNode.Render()
	if err != nil {
		return err
	}

	for index, child := range currentNode.GetChildren() {
		if child == nil {
			continue
		}

		err := e.buildNodeRecursively(child, fmt.Sprintf("%s[%d]", positionString, index))
		if err != nil {
			return err
		}

		if len(currentNode.GetChildren()) <= 1 {
			err := child.Position(currentNode, nil, nil)
			if err != nil {
				return err
			}

			continue
		}

		if index <= 0 {
			err := child.Position(currentNode, currentNode.GetChildren()[index+1], nil)
			if err != nil {
				return err
			}
		} else if index < len(currentNode.GetChildren())-1 {
			err := child.Position(currentNode, currentNode.GetChildren()[index+1], currentNode.GetChildren()[index-1])
			if err != nil {
				return err
			}
		} else {
			err := child.Position(currentNode, nil, currentNode.GetChildren()[index-1])
			if err != nil {
				return err
			}
		}
	}

	e.instancesToHash[currentNode.ID()] = currentNode

	return nil
}

func (e *DomEnvironment) mountToDom(rootElement string, vTree Node) error {
	rootElem := document.Call("querySelector", rootElement)
	if !rootElem.Truthy() {
		return fmt.Errorf("failed to find mount parent using query selector %q", rootElement)
	}

	styles, err := e.recursivelyMount(rootElem, vTree)
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

func (e *DomEnvironment) recursivelyMount(lastElement jsValue, currentNode Node) ([]string, error) {
	if currentNode == nil {
		return []string{}, nil
	}

	add := false
	domElement := lastElement
	styles := make([]string, 250)

	switch typedNode := currentNode.(type) {
	case *HTMLNode:
		add = true
		domElement = newHTMLElement(document, typedNode)

		for key, listener := range typedNode.EventListeners {
			hash := xxhash.Sum64String(string(typedNode.id) + "_" + key)

			e.eventsToHash[hash] = &wasmEvent{
				listener: listener,
				nodeHash: typedNode.ID(),
			}

			keyName := "data-lander-event-" + key + "-id"
			keyVal := fmt.Sprint(hash)

			oldKeyJSVal := domElement.Get(keyName)
			// If we couldn't find the previous key for that even
			if !oldKeyJSVal.Truthy() {
				// Set the event listener to the global listener
				domElement.Call("addEventListener", key, e.eventsCallback)
			}

			// Set the data key on the node
			domElement.Set(keyName, keyVal)
		}

		for _, style := range typedNode.Styles {
			styles = append(styles, style)
		}
	case *TextNode:
		add = true
		domElement = document.Call("createTextNode", typedNode.Text)
	}

	for _, child := range currentNode.GetChildren() {
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

func (e *DomEnvironment) patchDom(tree Node) error {
	rootElem := document.Call("querySelector", e.root)
	if !rootElem.Truthy() {
		return fmt.Errorf("failed to find mount parent using query selector %q", rootElem)
	}

	err := e.generatePatches(nil, e.tree, tree, rootElem)
	if err != nil {
		return err
	}

	for _, patch := range e.patches {
		err := patch.execute(document)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *DomEnvironment) generatePatches(prev, old, new Node, lastDomElement jsValue) error {
	domElement := lastDomElement
	if node, ok := old.(*HTMLNode); ok {
		domElement = document.Call("querySelector", fmt.Sprintf(`[data-lander-id="%d"]`, node.id))
		if !domElement.Truthy() {
			return fmt.Errorf(
				"failed to find mount parent using query selector %s",
				fmt.Sprintf(`[data-lander-id="%d"]`, node.id),
			)
		}
	}

	// If the old is missing, we need to insert the new node
	if old == nil {
		e.patches = append(e.patches, newPatchInsert(lastDomElement, prev, new))
		return nil
	}

	// If the new is missing, we need to remove the old node
	if new == nil {
		e.patches = append(e.patches, newPatchRemove(lastDomElement, prev, old))
		return nil
	}

	// If both nodes are identical, run on children
	if old.ID() == new.ID() && hashNode(old) == hashNode(new) {
		for index, child := range old.GetChildren() {
			var newChild Node
			if index <= len(new.GetChildren()) {
				newChild = new.GetChildren()[index]
			}

			err := e.generatePatches(old, child, newChild, domElement)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// If both nodes are similar
	if old.ID() == new.ID() && reflect.TypeOf(old) == reflect.TypeOf(new) {
		if val, ok := new.(*TextNode); ok {
			e.patches = append(e.patches, newPatchText(lastDomElement, prev, old, val.Text))
			return nil
		}
		if _, ok := new.(*HTMLNode); ok {
			e.patches = append(e.patches, newPatchHTML(old, new))
		}

		for index, child := range old.GetChildren() {
			var newChild Node
			if index <= len(new.GetChildren()) {
				newChild = new.GetChildren()[index]
			}

			err := e.generatePatches(old, child, newChild, domElement)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// if both nodes are not of the same type
	if reflect.TypeOf(old) != reflect.TypeOf(new) {
		// A replace will be needed
		e.patches = append(e.patches, newPatchReplace(lastDomElement, prev, old, new))
	}

	return nil
}

// TODO: This code has been borrowed from vugu, find a way to make it our own
func (e *DomEnvironment) handleDOMEvent(this js.Value, args []js.Value) interface{} {
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
