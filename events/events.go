//go:build js && wasm

package events

import (
	"syscall/js"
)

// EventListenerFunc is the type definition for a DOM event listener in javascript. Use this type
// to validate the event listeners passed as props to DOM elements.
type EventListenerFunc func(*DOMEvent) error

// EventListener is the concrete definition of a DOM event listener. It is used to keep track of the
// js function associated to the DOM element, and its underlying listener function. This allows the
// patches to remove the listeners safely and revoke the function to avoid any memory leaks.
type EventListener struct {
	Name    string
	Func    EventListenerFunc
	Wrapper js.Func
}

// DOMEvent is the base struct that contains the data for a DOM event triggered on the client.
// it contains the reference to the `this` object referencing the DOM node and the definition for the
// DOM event as a js.Value.
type DOMEvent struct {
	browserEvent js.Value
	this         js.Value
}

// NewDOMEvent generates a new DOM event to be passed to an event listener.
func NewDOMEvent(browserEvent, this js.Value) *DOMEvent {
	return &DOMEvent{
		browserEvent: browserEvent,
		this:         this,
	}
}

// JSEvent returns the browser event value which contains what would usually be the first argument
// of an event listener.
func (e *DOMEvent) JSEvent() js.Value {
	return e.browserEvent
}

// JSEventThis returns the value of the "this" variable for the Javascript event listener.
func (e *DOMEvent) JSEventThis() js.Value {
	return e.this
}

// PreventDefault calls preventDefault() on the underlying DOM event. Is thread safe, but may only be used
// in the same goroutine to avoid memory leaks.
func (e *DOMEvent) PreventDefault() {
	e.browserEvent.Call("preventDefault")
}
