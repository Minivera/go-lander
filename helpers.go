// +build js,wasm

package lander

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func extractAttributes(attributes map[string]interface{}) (map[string]string, map[string]EventListener, error) {
	attrs := make(map[string]string, 32)
	events := make(map[string]EventListener, 32)

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
		case EventListener:
			events[key] = casted
		default:
			return nil, nil, fmt.Errorf("attributes only support vars of type string, bool, int or EventListener, %T received", value)
		}
	}

	return attrs, events, nil
}

func hyperscript(tag string) (string, string, []string) {
	if tag == "" {
		// Always create a div by default
		return "div", "", []string{}
	}

	tagParts := strings.Split(tag, ".")
	if len(tagParts) == 1 {
		if strings.Index(tagParts[0], "#") >= 0 {
			tagAndID := strings.Split(tagParts[0], "#")
			return tagAndID[0], tagAndID[1], []string{}
		}
		return tagParts[0], "", []string{}
	}

	var tagname, id string
	classes := []string{}
	for i, part := range tagParts {
		if strings.Index(part, "#") >= 0 {
			tagAndID := strings.Split(part, "#")
			id = tagAndID[1]
			if i == 0 {
				tagname = tagAndID[0]
			} else {
				classes = append(classes, tagAndID[0])
			}
		} else {
			if i == 0 {
				tagname = part
			} else {
				classes = append(classes, part)
			}
		}
	}
	return tagname, id, classes
}

func newHTMLElement(document jsValue, currentElement *HTMLNode) jsValue {
	var domElement jsValue
	if currentElement.namespace != "" {
		domElement = document.Call("createElementNS", currentElement.namespace, currentElement.Tag)
	} else {
		domElement = document.Call("createElement", currentElement.Tag)
	}

	domElement.Call("setAttribute", "data-lander-id", strconv.FormatUint(currentElement.ID(), 10))

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

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
