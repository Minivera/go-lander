// +build js,wasm

package lander

import (
	"sync"
	"syscall/js"
)

// EventEnv is the interface for an environement for a DOM event triggered in Javascript.
// This has been borrowed from vugu, make it our own
type EventEnv interface {
	Lock()         // acquire write lock
	UnlockOnly()   // release write lock
	UnlockRender() // release write lock and request re-render

	RLock()   // acquire read lock
	RUnlock() // release read lock
}

// EventListener is the type definition for a DOM event in javascript.
type EventListener func(Node, *DOMEvent) error

// DOMEvent is the base struct that contains the data for a DOM event triggered on the client.
type DOMEvent struct {
	browserEvent js.Value
	this         js.Value
	environment  *eventEnv
}

// JSEvent returns the browser event value which contains what would usually be the first argument of an event listener.
func (e *DOMEvent) JSEvent() js.Value {
	return e.browserEvent
}

// JSEventThis returns the value of the "this" variable for the Javascript event listener.
func (e *DOMEvent) JSEventThis() js.Value {
	return e.this
}

// EventEnv returns the Environement for the current event, allowing to block or unblock rendering if necessary.
func (e *DOMEvent) EventEnv() EventEnv {
	return e.environment
}

// PreventDefault calls preventDefault() on the underlying DOM event.
// May only be used within event handler in same goroutine.
func (e *DOMEvent) PreventDefault() {
	e.browserEvent.Call("preventDefault")
}

type wasmEvent struct {
	listener EventListener
	nodeHash uint64
}

// eventEnv implements EventEnv
type eventEnv struct {
	rwmu            *sync.RWMutex
	requestRenderCH chan bool
}

// Lock will acquire write lock
func (ee *eventEnv) Lock() {
	ee.rwmu.Lock()
}

// UnlockOnly will release the write lock
func (ee *eventEnv) UnlockOnly() {
	ee.rwmu.Unlock()
}

// UnlockRender will release write lock and request re-render
func (ee *eventEnv) UnlockRender() {
	ee.rwmu.Unlock()
	if ee.requestRenderCH != nil {
		// send non-blocking
		select {
		case ee.requestRenderCH <- true:
		default:
		}
	}
}

// RLock will acquire a read lock
func (ee *eventEnv) RLock() {
	ee.rwmu.RLock()
}

// RUnlock will release the read lock
func (ee *eventEnv) RUnlock() {
	ee.rwmu.RUnlock()
}
