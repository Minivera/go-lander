package go_lander

import (
	"sync"
	"syscall/js"
)

// This has been borrowed from vugu, make it our own
type EventEnv interface {
	Lock()         // acquire write lock
	UnlockOnly()   // release write lock
	UnlockRender() // release write lock and request re-render

	RLock()   // acquire read lock
	RUnlock() // release read lock
}

type EventListener func(Node, *DOMEvent) error

type DOMEvent struct {
	browserEvent js.Value
	this         js.Value
	environment  *eventEnv
}

func (e *DOMEvent) JSEvent() js.Value {
	return e.browserEvent
}

func (e *DOMEvent) JSEventThis() js.Value {
	return e.this
}

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
