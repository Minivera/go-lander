//go:build js && wasm

package utils

import (
	"fmt"
	"math/rand"
	"syscall/js"
	"time"

	lEvents "github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

func ExtractAttributes(attributes map[string]interface{}) (attrs map[string]string, events map[string]lEvents.EventListener) {
	for key, value := range attributes {
		switch casted := value.(type) {
		case string:
			attrs[key] = casted
		case int:
			attrs[key] = fmt.Sprintf("%d", casted)
		case bool:
			// Bool attributes only adds the attribute if true, like required=""
			if casted {
				attrs[key] = ""
			}
		case lEvents.EventListener:
			events[key] = casted
		default:
			// attributes only support vars of type string, bool, int or EventListener
			// Any other attribute is ignored to avoid panicking.
			continue
		}
	}

	return attrs, events
}

func NewHTMLElement(document js.Value, currentElement *nodes.HTMLNode) js.Value {
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

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
