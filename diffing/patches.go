//go:build js && wasm

package diffing

import (
	"fmt"
	"strings"
	"syscall/js"
)

type Patch interface {
	Execute(js.Value) error
}

type patchText struct {
	parentDomNode            jsValue
	parent, oldNode, newNode Node
	newText                  string
}

func newPatchText(parent jsValue, parentNode, old Node, text string) patch {
	return &patchText{
		parentDomNode: parent,
		parent:        parentNode,
		oldNode:       old,
		newText:       text,
	}
}

func (p *patchText) execute(document jsValue) error {
	err := p.oldNode.Update(map[string]interface{}{
		"text": p.newText,
	})
	if err != nil {
		return err
	}

	index := 0
	for _, node := range p.parent.GetChildren() {
		if node == p.oldNode {
			break
		}
		index++
	}

	if index >= len(p.parent.GetChildren()) {
		return fmt.Errorf("could not find the child in the parent")
	}

	domNode := p.parentDomNode.Get("childNodes").Index(index)
	if !domNode.Truthy() {
		return fmt.Errorf("could not find node at index %d for id [data-lander-id=%d]", index, p.parent.ID())
	}

	domNode.Set("nodeValue", p.newText)

	return nil
}

type patchHTML struct {
	oldNode, newNode Node
}

func newPatchHTML(old, new Node) patch {
	return &patchHTML{
		oldNode: old,
		newNode: new,
	}
}

func (p *patchHTML) execute(document jsValue) error {
	newHtml, ok := p.newNode.(*HTMLNode)
	if !ok {
		return fmt.Errorf("new node was not of type HTMLNode, %T given instead", p.newNode)
	}

	// TODO: Find a way to bind new event listeners
	newAttributes := make(map[string]interface{}, len(newHtml.Attributes)+len(newHtml.EventListeners)+2)
	for key, value := range newHtml.Attributes {
		newAttributes[key] = value
	}

	// TODO: Fix the memory leak here when a node is removed, but not its event listeners
	for key, value := range newHtml.EventListeners {
		newAttributes[key] = value
	}

	if newHtml.DomID != "" {
		newAttributes["id"] = newHtml.DomID
	}
	if len(newHtml.Classes) > 0 {
		newAttributes["class"] = strings.Join(newHtml.Classes, " ")
	}

	err := p.oldNode.Update(newAttributes)
	if err != nil {
		return err
	}

	oldHtml, ok := p.oldNode.(*HTMLNode)
	if !ok {
		return fmt.Errorf("new node was not of type HTMLNode, %T given instead", p.oldNode)
	}

	oldNode := document.Call("querySelector", fmt.Sprintf(`[data-lander-id="%d"]`, p.oldNode.ID()))
	if !oldNode.Truthy() {
		return fmt.Errorf("could not find node for id [data-lander-id=%d]", p.oldNode.ID())
	}

	// TODO: Find a way to clear the attributes on the dom node before setting the new ones
	for key, value := range oldHtml.Attributes {
		oldNode.Call("setAttribute", key, value)
	}

	// TODO: Find a way to clear the class list before adding to it
	classList := oldNode.Get("classList")
	for _, value := range oldHtml.Classes {
		classList.Call("add", value)
	}

	if oldHtml.DomID != "" {
		oldNode.Set("id", oldHtml.DomID)
	}

	return nil
}

type patchInsert struct {
	parentDomNode   jsValue
	parent, newNode Node
}

func newPatchInsert(parentElem jsValue, parent, new Node) patch {
	return &patchInsert{
		parentDomNode: parentElem,
		parent:        parent,
		newNode:       new,
	}
}

func (p *patchInsert) execute(document jsValue) error {
	err := p.parent.InsertChildren(p.newNode, -1)
	if err != nil {
		println(err)
		return err
	}

	var domElement jsValue
	switch typedNode := p.newNode.(type) {
	case *HTMLNode:
		domElement = newHTMLElement(document, typedNode)
	case *TextNode:
		domElement = document.Call("createTextNode", typedNode.Text)
	default:
		// Ignore anything that's not dom related
		return nil
	}

	p.parentDomNode.Call("appendChild", domElement)

	return nil
}

type patchRemove struct {
	parentDomNode   jsValue
	parent, oldNode Node
}

func newPatchRemove(parentElem jsValue, parent, old Node) patch {
	return &patchRemove{
		parentDomNode: parentElem,
		parent:        parent,
		oldNode:       old,
	}
}

func (p *patchRemove) execute(document jsValue) error {
	index := 0
	for _, node := range p.parent.GetChildren() {
		if node == p.oldNode {
			break
		}
		index++
	}

	if index >= len(p.parent.GetChildren()) {
		return fmt.Errorf("could not find the child in the parent")
	}

	domNode := p.parentDomNode.Get("childNodes").Index(index)
	if !domNode.Truthy() {
		return fmt.Errorf("could not find node at index %d for id [data-lander-id=%d]", index, p.parent.ID())
	}

	err := p.parent.RemoveChildren(p.oldNode)
	if err != nil {
		return err
	}

	switch p.oldNode.(type) {
	case *HTMLNode:
		p.parentDomNode.Call("removeChild", domNode)
	case *TextNode:
		p.parentDomNode.Call("removeChild", domNode)
	}

	return nil
}

type patchReplace struct {
	parentDomNode            jsValue
	parent, newNode, oldNode Node
}

func newPatchReplace(parentElem jsValue, parent, old, new Node) patch {
	return &patchReplace{
		parentDomNode: parentElem,
		parent:        parent,
		newNode:       new,
		oldNode:       old,
	}
}

func (p *patchReplace) execute(document jsValue) error {
	index := 0
	for _, node := range p.parent.GetChildren() {
		if node == p.oldNode {
			break
		}
		index++
	}

	if index >= len(p.parent.GetChildren()) {
		return fmt.Errorf("could not find the child in the parent")
	}

	domNode := p.parentDomNode.Get("childNodes").Index(index)
	if !domNode.Truthy() {
		return fmt.Errorf("could not find node at index %d for id [data-lander-id=%d]", index, p.parent.ID())
	}

	err := p.parent.ReplaceChildren(p.oldNode, p.newNode)
	if err != nil {
		return err
	}

	switch typedNode := p.newNode.(type) {
	case *HTMLNode:
		domElement := newHTMLElement(document, typedNode)

		p.parentDomNode.Call("replaceChild", domNode, domElement)
	case *TextNode:
		domElement := document.Call("createTextNode", typedNode.Text)

		p.parentDomNode.Call("replaceChild", domNode, domElement)
	}

	return nil
}
