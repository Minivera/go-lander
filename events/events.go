//go:build js && wasm

package events

import (
	"syscall/js"
)

// EventListenerFunc is the type definition for a DOM event in javascript.
type EventListenerFunc func(*DOMEvent) error

type EventListener struct {
	Name    string
	Func    EventListenerFunc
	Wrapper js.Func
}

// DOMEvent is the base struct that contains the data for a DOM event triggered on the client.
type DOMEvent struct {
	browserEvent js.Value
	this         js.Value
}

func NewDOMEvent(browserEvent, this js.Value) *DOMEvent {
	return &DOMEvent{
		browserEvent: browserEvent,
		this:         this,
	}
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
