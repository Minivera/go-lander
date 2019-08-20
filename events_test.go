// +build js,wasm

package lander

import (
	"syscall/js"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDOMEvent_JSEvent(t *testing.T) {
	t.Run("Returns the JSEvent from a DOMEvent", func(t *testing.T) {
		bevent := js.Value{}

		event := DOMEvent{
			browserEvent: bevent,
		}

		assert.Equal(t, bevent, event.JSEvent(), "The event passed to the struct should be the same as the return value")
	})
}

func TestDOMEvent_JSEventThis(t *testing.T) {
	t.Run("Returns the value of this from a DOMEvent", func(t *testing.T) {
		this := js.Value{}

		event := DOMEvent{
			this: this,
		}

		assert.Equal(t, this, event.JSEventThis(), "The this passed to the struct should be the same as the return value")
	})
}

func TestDOMEvent_EventEnv(t *testing.T) {
	t.Run("Returns the value of the event environment", func(t *testing.T) {
		env := &eventEnv{}

		event := DOMEvent{
			environment: env,
		}

		assert.Equal(t, env, event.EventEnv(), "The event passed to the struct should be the same as the return value")
	})
}
