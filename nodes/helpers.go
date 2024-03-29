//go:build js && wasm

package nodes

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"

	lEvents "github.com/minivera/go-lander/events"
)

// ExtractAttributes extracts the relevant attributes, props, and listeners for an HTML node given the
// attributes map. This allows extracting based on types, which can then be reconciled with the DOM nodes
// attributes and properties.
func ExtractAttributes(attributes map[string]interface{}) (
	map[string]string, map[string]interface{}, map[string]*lEvents.EventListener) {

	attrs := map[string]string{}
	props := map[string]interface{}{}
	events := map[string]*lEvents.EventListener{}

	for key, value := range attributes {
		switch casted := value.(type) {
		case string:
			attrs[key] = casted
			props[key] = casted
		case int:
			attrs[key] = fmt.Sprintf("%d", casted)
			props[key] = casted
		case bool:
			// Bool attributes only adds the attribute if true, like required=""
			props[key] = casted
			if casted {
				attrs[key] = ""
			}
		case func(*lEvents.DOMEvent) error:
			events[key] = &lEvents.EventListener{
				Name: key,
				Func: casted,
			}
		case lEvents.EventListenerFunc:
			events[key] = &lEvents.EventListener{
				Name: key,
				Func: casted,
			}
		default:
			// attributes only support vars of type string, bool, int or EventListener
			// Any other attribute is ignored to avoid panicking.
			continue
		}
	}

	return attrs, props, events
}

// NewHTMLElement creates a new HTML node and sets all its attributes, properties, and event listeners
// on creation.
func NewHTMLElement(document js.Value, currentElement *HTMLNode) js.Value {
	var domElement js.Value
	if currentElement.Namespace != "" {
		domElement = document.Call("createElementNS", currentElement.Namespace, currentElement.Tag)
	} else {
		domElement = document.Call("createElement", currentElement.Tag)
	}

	for key, value := range currentElement.Attributes {
		domElement.Call("setAttribute", key, value)
	}

	classList := domElement.Get("classList")
	for _, value := range currentElement.Classes {
		classList.Call("add", value)
	}

	if currentElement.DomID != "" {
		domElement.Set("id", currentElement.DomID)
	}

	return domElement
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()),
)

// RandomString generates a random string of the provided length from a specific charset. This is useful
// for generating semi-random class names for HTML elements.
func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
