//go:build js && wasm

package events

import (
	"sync"
	"syscall/js"

	"github.com/minivera/go-lander/nodes"
)

// EventListener is the type definition for a DOM event in javascript.
type EventListener func(nodes.Node, *DOMEvent) error

// DOMEvent is the base struct that contains the data for a DOM event triggered on the client.
type DOMEvent struct {
	sync.RWMutex

	browserEvent js.Value
	this         js.Value
}

// JSEvent returns the browser event value which contains what would usually be the first argument of an event listener.
func (e *DOMEvent) JSEvent() js.Value {
	return e.browserEvent
}

// JSEventThis returns the value of the "this" variable for the Javascript event listener.
func (e *DOMEvent) JSEventThis() js.Value {
	return e.this
}

// PreventDefault calls preventDefault() on the underlying DOM event.
// May only be used within event handler in same goroutine.
func (e *DOMEvent) PreventDefault() {
	e.browserEvent.Call("preventDefault")
}
