package go_lander

import (
	"fmt"
	"strings"
	"syscall/js"
)

type patch interface {
	execute(js.Value) error
}

type patchText struct {
	parentDomNode    js.Value
	oldNode, newNode Node
	newText          string
}

func newPatchText(parent js.Value, old Node, text string) patch {
	return &patchText{
		parentDomNode: parent,
		oldNode:       old,
		newText:       text,
	}
}

func (p *patchText) execute(document js.Value) error {
	err := p.oldNode.Update(map[string]interface{}{
		"text": p.newText,
	})
	if err != nil {
		return err
	}

	oldText, ok := p.oldNode.(*TextNode)
	if !ok {
		return fmt.Errorf("old node was not of type TextNode, %T given instead", p.oldNode)
	}

	newNode := document.Call("createTextNode", p.newText)
	p.parentDomNode.Call("replaceChild", newNode, oldText.domNode)

	err = p.newNode.Mount(newNode)
	if err != nil {
		return err
	}

	return nil
}

type patchHtml struct {
	oldNode, newNode Node
}

func newPatchHtml(old, new Node) patch {
	return &patchHtml{
		oldNode: old,
		newNode: new,
	}
}

func (p *patchHtml) execute(_ js.Value) error {
	newHtml, ok := p.newNode.(*HtmlNode)
	if !ok {
		return fmt.Errorf("old node was not of type HtmlNode, %T given instead", p.oldNode)
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

	newAttributes["id"] = newHtml.DomID
	newAttributes["class"] = strings.Join(newHtml.Classes, " ")

	err := p.oldNode.Update(newAttributes)
	if err != nil {
		return err
	}

	return nil
}

type patchInsert struct {
	parentDomNode   js.Value
	parent, newNode Node
}

func newPatchInsert(parentElem js.Value, parent, new Node) patch {
	return &patchInsert{
		parentDomNode: parentElem,
		parent:        parent,
		newNode:       new,
	}
}

func (p *patchInsert) execute(document js.Value) error {
	println("inserting children")
	err := p.parent.InsertChildren(p.newNode, -1)
	if err != nil {
		println(err)
		return err
	}
	println("children inserted")

	println(fmt.Sprintf("parent node: %# v", p.parent))
	println(fmt.Sprintf("new node: %# v", p.newNode))

	var domElement js.Value
	switch typedNode := p.newNode.(type) {
	case *HtmlNode:
		println("insert html node")
		domElement = newHTMLElement(document, typedNode)
	case *TextNode:
		println("insert text node")
		domElement = document.Call("createTextNode", typedNode.Text)
	default:
		// Ignore anything that's not dom related
		return nil
	}

	p.parentDomNode.Call("appendChild", domElement)

	err = p.newNode.Mount(domElement)
	if err != nil {
		return err
	}

	return nil
}

type patchRemove struct {
	parentDomNode   js.Value
	parent, oldNode Node
}

func newPatchRemove(parentElem js.Value, parent, old Node) patch {
	return &patchRemove{
		parentDomNode: parentElem,
		parent:        parent,
		oldNode:       old,
	}
}

func (p *patchRemove) execute(document js.Value) error {
	err := p.parent.RemoveChildren(p.oldNode)
	if err != nil {
		return err
	}

	switch typedNode := p.oldNode.(type) {
	case *HtmlNode:
		p.parentDomNode.Call("removeChild", typedNode.domNode)
	case *TextNode:
		p.parentDomNode.Call("removeChild", typedNode.domNode)
	}

	err = p.oldNode.Remove()
	if err != nil {
		return err
	}

	return nil
}

type patchReplace struct {
	parentDomNode            js.Value
	parent, newNode, oldNode Node
}

func newPatchReplace(parentElem js.Value, parent, old, new Node) patch {
	return &patchReplace{
		parentDomNode: parentElem,
		parent:        parent,
		newNode:       new,
		oldNode:       old,
	}
}

func (p *patchReplace) execute(document js.Value) error {
	err := p.parent.ReplaceChildren(p.oldNode, p.newNode)
	if err != nil {
		return err
	}

	switch typedNode := p.newNode.(type) {
	case *HtmlNode:
		domElement := newHTMLElement(document, typedNode)

		oldHtml, ok := p.oldNode.(*HtmlNode)
		if !ok {
			return fmt.Errorf("old node was not of type HtmlNode, %T given instead", p.oldNode)
		}

		p.parentDomNode.Call("replaceChild", oldHtml.domNode, domElement)

		err := p.oldNode.Mount(domElement)
		if err != nil {
			return err
		}
	case *TextNode:
		domElement := document.Call("createTextNode", typedNode.Text)

		p.parentDomNode.Call("replaceChild", typedNode.domNode, domElement)

		err := p.oldNode.Mount(domElement)
		if err != nil {
			return err
		}
	}

	err = p.oldNode.Remove()
	if err != nil {
		return err
	}

	return nil
}
