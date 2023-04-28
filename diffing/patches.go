//go:build js && wasm

package diffing

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/minivera/go-lander/context"

	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

type Patch interface {
	Execute(js.Value, *[]string) error
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

func (p *patchText) Execute(document js.Value, styles *[]string) error {
	fmt.Printf("Executing patch Text on %T, %v\n", p.oldNode, p.oldNode)
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

func (p *patchHTML) Execute(document js.Value, styles *[]string) error {
	fmt.Printf("Executing patch HTML on %T, %v\n", p.oldNode, p.oldNode)
	newAttributes := make(map[string]interface{}, len(p.newNode.Attributes)+len(p.newNode.EventListeners)+2)

	// Run only on the properties since they retain their original type prior to extraction
	for key, value := range p.newNode.Properties {
		newAttributes[key] = value
	}

	for key, value := range p.newNode.EventListeners {
		newAttributes[key] = value.Func
	}

	if p.newNode.DomID != "" {
		newAttributes["id"] = p.newNode.DomID
	}
	if len(p.newNode.Classes) > 0 {
		newAttributes["class"] = strings.Join(p.newNode.Classes, " ")
	}

	// Remove any event listeners to avoid old closures leaking into them
	for event, listener := range p.oldNode.EventListeners {
		p.oldNode.DomNode.Call("removeEventListener", event, listener.Wrapper)
	}

	p.oldNode.Update(newAttributes)

	// Add new event listeners using the attributes
	for event, listener := range p.oldNode.EventListeners {
		listener.Wrapper = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return p.listenerFunc(listener.Func, this, args)
		})
		p.oldNode.DomNode.Call("addEventListener", event, listener.Wrapper)
	}

	// Update the active class with the new value, replace the styles
	p.oldNode.ActiveClass = p.newNode.ActiveClass
	p.oldNode.Styles = p.newNode.Styles

	classList := p.oldNode.DomNode.Get("classList")
	if p.oldNode.ActiveClass != "" {
		classList.Call("add", p.oldNode.ActiveClass)
	}

	return nil
}

type patchListeners struct {
	listenerFunc func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{}
	oldNode      *nodes.HTMLNode
}

func newPatchListeners(
	listenerFunc func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{},
	old *nodes.HTMLNode,
) Patch {
	return &patchListeners{
		listenerFunc: listenerFunc,
		oldNode:      old,
	}
}

func (p *patchListeners) Execute(_ js.Value, _ *[]string) error {
	fmt.Printf("Executing patch listeners on %T, %v\n", p.oldNode, p.oldNode)
	// Remove any event listeners to avoid old closures leaking into them
	for event, listener := range p.oldNode.EventListeners {
		p.oldNode.DomNode.Call("removeEventListener", event, listener.Wrapper)
	}

	// Add new event listeners using the attributes
	for event, listener := range p.oldNode.EventListeners {
		listener.Wrapper = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return p.listenerFunc(listener.Func, this, args)
		})
		p.oldNode.DomNode.Call("addEventListener", event, listener.Wrapper)
	}

	return nil
}

type patchInsert struct {
	listenerFunc     func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{}
	closestDOMParent js.Value
	parent, newNode  nodes.Node
}

func newPatchInsert(
	listenerFunc func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{},
	closestDOMParent js.Value,
	parent,
	new nodes.Node,
) Patch {
	return &patchInsert{
		listenerFunc:     listenerFunc,
		closestDOMParent: closestDOMParent,
		parent:           parent,
		newNode:          new,
	}
}

func (p *patchInsert) Execute(document js.Value, styles *[]string) error {
	fmt.Printf("Executing patch insert on %T, %v\n", p.newNode, p.newNode)
	fmt.Printf("Parent is %T, %v\n", p.parent, p.parent)
	switch parent := p.parent.(type) {
	case *nodes.FuncNode:
		parent.RenderResult = p.newNode

		// Trigger a recursive mount for its render result
		childStyles := RecursivelyMount(p.listenerFunc, document, p.closestDOMParent, parent.RenderResult)

		for _, style := range childStyles {
			*styles = append(*styles, style)
		}
	case *nodes.HTMLNode:
		err := parent.InsertChildren(p.newNode, -1)
		if err != nil {
			return err
		}

		var domElement js.Value
		toAdd := false
		switch typedNode := p.newNode.(type) {
		case *nodes.HTMLNode:
			toAdd = true
			domElement = nodes.NewHTMLElement(document, typedNode)
			typedNode.Mount(domElement)

			// Trigger a recursive mount for all its children
			for _, child := range typedNode.Children {
				if child == nil {
					continue
				}

				childStyles := RecursivelyMount(p.listenerFunc, document, domElement, child)

				for _, style := range childStyles {
					*styles = append(*styles, style)
				}
			}
		case *nodes.TextNode:
			toAdd = true
			domElement = document.Call("createTextNode", typedNode.Text)
			typedNode.Mount(domElement)
		case *nodes.FuncNode:
			context.RegisterComponent(typedNode)
			context.RegisterComponentContext("render", typedNode)
			context.RegisterComponentContext("mount", typedNode)

			return newPatchInsert(p.listenerFunc, parent.DomNode, typedNode, typedNode.Clone().Render(context.CurrentContext)).
				Execute(document, styles)
		default:
			// Ignore anything that's not dom related
			return nil
		}

		if toAdd {
			parent.DomNode.Call("appendChild", domElement)
		}
	default:
		// Ignore anything that's not dom related
		return nil
	}

	return nil
}

type patchRemove struct {
	closestDOMParent js.Value
	parent, oldNode  nodes.Node
}

func newPatchRemove(parent nodes.Node, closestDOMParent js.Value, old nodes.Node) Patch {
	return &patchRemove{
		parent:           parent,
		closestDOMParent: closestDOMParent,
		oldNode:          old,
	}
}

func (p *patchRemove) Execute(document js.Value, styles *[]string) error {
	fmt.Printf("Executing patch remove on %T, %v\n", p.oldNode, p.oldNode)
	switch typedNode := p.parent.(type) {
	case *nodes.FuncNode:
		typedNode.RenderResult = nil
	case *nodes.HTMLNode:
		err := typedNode.RemoveChildren(p.oldNode)
		if err != nil {
			return err
		}
	default:
		// Ignore anything that's not dom related
		return nil
	}

	switch typedNode := p.oldNode.(type) {
	case *nodes.HTMLNode:
		p.closestDOMParent.Call("removeChild", typedNode.DomNode)
	case *nodes.TextNode:
		p.closestDOMParent.Call("removeChild", typedNode.DomNode)
	case *nodes.FuncNode:
		return newPatchRemove(typedNode, p.closestDOMParent, typedNode.RenderResult).Execute(document, styles)
	}

	return nil
}

type patchReplace struct {
	listenerFunc             func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{}
	closestDOMParent         js.Value
	parent, newNode, oldNode nodes.Node
}

func newPatchReplace(
	listenerFunc func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{},
	closestDOMParent js.Value,
	parent,
	old,
	new nodes.Node,
) Patch {
	return &patchReplace{
		listenerFunc:     listenerFunc,
		closestDOMParent: closestDOMParent,
		parent:           parent,
		newNode:          new,
		oldNode:          old,
	}
}

func (p *patchReplace) Execute(document js.Value, styles *[]string) error {
	fmt.Printf("Executing patch replace on %T, %v\n", p.oldNode, p.oldNode)
	switch parent := p.parent.(type) {
	case *nodes.FuncNode:
		parent.RenderResult = p.newNode

		// Trigger a recursive mount for its render result
		childStyles := RecursivelyMount(p.listenerFunc, document, p.closestDOMParent, parent.RenderResult)

		for _, style := range childStyles {
			*styles = append(*styles, style)
		}
	case *nodes.HTMLNode:
		err := parent.ReplaceChildren(p.oldNode, p.newNode)
		if err != nil {
			return err
		}

		var oldDomNode js.Value
		if converted, ok := p.oldNode.(*nodes.HTMLNode); ok {
			oldDomNode = converted.DomNode
		} else if converted, ok := p.oldNode.(*nodes.TextNode); ok {
			oldDomNode = converted.DomNode
		}

		switch typedNode := p.newNode.(type) {
		case *nodes.HTMLNode:
			domElement := nodes.NewHTMLElement(document, typedNode)
			typedNode.Mount(domElement)

			// Trigger a recursive mount for all its children
			for _, child := range typedNode.Children {
				if child == nil {
					continue
				}

				childStyles := RecursivelyMount(p.listenerFunc, document, domElement, child)

				for _, style := range childStyles {
					*styles = append(*styles, style)
				}
			}

			parent.DomNode.Call("replaceChild", domElement, oldDomNode)
		case *nodes.TextNode:
			domElement := document.Call("createTextNode", typedNode.Text)
			typedNode.Mount(domElement)

			parent.DomNode.Call("replaceChild", domElement, oldDomNode)
		case *nodes.FuncNode:
			context.RegisterComponent(typedNode)
			context.RegisterComponentContext("render", typedNode)
			context.RegisterComponentContext("mount", typedNode)

			return newPatchReplace(p.listenerFunc, parent.DomNode, typedNode, p.oldNode, typedNode.Clone().Render(context.CurrentContext)).
				Execute(document, styles)
		}
	default:
		// Ignore anything that's not dom related
		return nil
	}

	return nil
}
