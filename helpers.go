package go_lander

import (
	"math/rand"
	"strings"
	"syscall/js"
	"time"
)

func mergeAttributes(a1, a2 map[string]string) map[string]string {
	for key, _ := range a2 {
		if val, ok := a2[key]; ok {
			a1[key] = val
		} else {
			delete(a1, key)
		}
	}
	return a1
}

func hyperscript(tag string) (string, string, []string) {
	tagParts := strings.Split(tag, ".")
	if len(tagParts) <= 0 {
		// Always create a div by default
		return "div", "", []string{}
	}
	if len(tagParts) == 1 {
		if strings.Index(tagParts[0], "#") >= 0 {
			tagAndID := strings.Split(tagParts[0], "#")
			return tagAndID[0], tagAndID[1], []string{}
		}
		return tagParts[0], "", []string{}
	}

	var tagname, id string
	classes := make([]string, len(tagParts))
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
			classes = append(classes, part)
		}
	}
	return tagname, id, classes
}

func walkTree(currentNode Node, callback func(Node) error) error {
	err := callback(currentNode)
	if err != nil {
		return err
	}

	for _, child := range currentNode.GetChildren() {
		err := walkTree(child, callback)
		if err != nil {
			return err
		}
	}

	return nil
}

func newHTMLElement(document js.Value, currentElement *HtmlNode) js.Value {
	var domElement js.Value
	if currentElement.namespace != "" {
		domElement = document.Call("createElementNS", currentElement.namespace, currentElement.Tag)
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

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()),
)

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
