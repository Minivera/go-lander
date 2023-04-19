//go:build js && wasm

package diffing

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

type Patch interface {
	Execute(js.Value) error
}

type patchText struct {
	parent  nodes.Node
	oldNode *nodes.TextNode
	newText string
}

func newPatchText(parentNode nodes.Node, old *nodes.TextNode, text string) Patch {
	return &patchText{
		parent:  parentNode,
		oldNode: old,
		newText: text,
	}
}

func (p *patchText) Execute(document js.Value) error {
	p.oldNode.Update(p.newText)

	return nil
}

type patchHTML struct {
	listenerFunc     func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{}
	oldNode, newNode *nodes.HTMLNode
}

func newPatchHTML(
	listenerFunc func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{},
	old,
	new *nodes.HTMLNode,
) Patch {
	return &patchHTML{
		listenerFunc: listenerFunc,
		oldNode:      old,
		newNode:      new,
	}
}

func (p *patchHTML) Execute(document js.Value) error {
	newAttributes := make(map[string]interface{}, len(p.newNode.Attributes)+len(p.newNode.EventListeners)+2)
	for key, value := range p.newNode.Attributes {
		newAttributes[key] = value
	}

	for key, value := range p.newNode.EventListeners {
		newAttributes[key] = value
	}

	if p.newNode.DomID != "" {
		newAttributes["id"] = p.newNode.DomID
	}
	if len(p.newNode.Classes) > 0 {
		newAttributes["class"] = strings.Join(p.newNode.Classes, " ")
	}

	// Remove any event listeners using the direct attribute rather than addEventListener
	for event, listener := range p.oldNode.EventListeners {
		fmt.Printf("removing event in diffing\n %s", event)
		p.oldNode.DomNode.Call("removeEventListener", event, listener.Wrapper)
	}

	p.oldNode.Update(newAttributes)

	// Add new event listeners using the attributes
	for event, listener := range p.oldNode.EventListeners {
		fmt.Printf("adding event in diffing\n %s", event)
		listener.Wrapper = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return p.listenerFunc(listener.Func, this, args)
		})
		p.oldNode.DomNode.Call("addEventListener", event, listener.Wrapper)
	}

	return nil
}

type patchInsert struct {
	parent, newNode nodes.Node
}

func newPatchInsert(parent, new nodes.Node) Patch {
	return &patchInsert{
		parent:  parent,
		newNode: new,
	}
}

func (p *patchInsert) Execute(document js.Value) error {
	htmlParent, ok := p.parent.(*nodes.HTMLNode)
	if !ok {
		return nil
	}

	err := htmlParent.InsertChildren(p.newNode, -1)
	if err != nil {
		return err
	}

	var domElement js.Value
	switch typedNode := p.newNode.(type) {
	case *nodes.HTMLNode:
		domElement = nodes.NewHTMLElement(document, typedNode)
	case *nodes.TextNode:
		domElement = document.Call("createTextNode", typedNode.Text)
	default:
		// Ignore anything that's not dom related
		return nil
	}

	htmlParent.DomNode.Call("appendChild", domElement)

	return nil
}

type patchRemove struct {
	parent, oldNode nodes.Node
}

func newPatchRemove(parent, old nodes.Node) Patch {
	return &patchRemove{
		parent:  parent,
		oldNode: old,
	}
}

func (p *patchRemove) Execute(document js.Value) error {
	htmlParent, ok := p.parent.(*nodes.HTMLNode)
	if !ok {
		return nil
	}

	err := htmlParent.RemoveChildren(p.oldNode)
	if err != nil {
		return err
	}

	switch typedNode := p.oldNode.(type) {
	case *nodes.HTMLNode:
		htmlParent.DomNode.Call("removeChild", typedNode.DomNode)
	case *nodes.TextNode:
		htmlParent.DomNode.Call("removeChild", typedNode.DomNode)
	}

	return nil
}

type patchReplace struct {
	parent, newNode, oldNode nodes.Node
}

func newPatchReplace(parent, old, new nodes.Node) Patch {
	return &patchReplace{
		parent:  parent,
		newNode: new,
		oldNode: old,
	}
}

func (p *patchReplace) Execute(document js.Value) error {
	htmlParent, ok := p.parent.(*nodes.HTMLNode)
	if !ok {
		return nil
	}

	err := htmlParent.ReplaceChildren(p.oldNode, p.newNode)
	if err != nil {
		return err
	}

	switch typedNode := p.newNode.(type) {
	case *nodes.HTMLNode:
		domElement := nodes.NewHTMLElement(document, typedNode)

		htmlParent.DomNode.Call("replaceChild", typedNode.DomNode, domElement)
	case *nodes.TextNode:
		domElement := document.Call("createTextNode", typedNode.Text)

		htmlParent.DomNode.Call("replaceChild", typedNode.DomNode, domElement)
	}

	return nil
}
