package go_lander

import (
	"fmt"
	"regexp"
	"strings"
	"syscall/js"
)

var document js.Value

func init() {
	document = js.Global().Get("document")
}

func validElement(value js.Value) bool {
	return value != (js.Value{}) || value != js.Null() || value != js.Undefined()
}

func mountToDom(rootElement string, vTree Node) error {
	root := document.Call("querySelector", rootElement)
	if !validElement(root) {
		return fmt.Errorf("failed to find mount parent using query selector %q", rootElement)
	}

	styles, err := recursivelyMount(root, vTree)
	if err != nil {
		return err
	}

	replaceRegex := regexp.MustCompile(`'.+'|\".+\"|\\S+`)
	stylesString := replaceRegex.ReplaceAllString(strings.Join(styles, ""), "")

	head := document.Call("querySelector", "head")
	if !validElement(head) {
		return fmt.Errorf("failed to find heads using query selector")
	}

	styleTag := document.Call("createElement", "style")
	styleTag.Set("id", "lander-style-tag")
	styleTag.Set("innerHTML", stylesString)
	head.Call("appendChild", styleTag)

	return nil
}

func recursivelyMount(lastElement js.Value, currentNode Node) ([]string, error) {
	add := false
	domElement := lastElement
	styles := make([]string, 250)

	switch typedNode := currentNode.(type) {
	case *HtmlNode:
		add = true
		if typedNode.namespace != "" {
			domElement = document.Call("createElementNS", typedNode.namespace, typedNode.Tag)
		} else {
			domElement = document.Call("createElement", typedNode.Tag)
		}

		for key, value := range typedNode.Attributes {
			domElement.Call("setAttribute", key, value)
		}

		classList := domElement.Get("classList")
		for _, value := range typedNode.Classes {
			classList.Call("add", value)
		}

		if typedNode.DomID != "" {
			domElement.Set("id", typedNode.DomID)
		}

		for _, style := range typedNode.Styles {
			styles = append(styles, style)
		}
	case *TextNode:
		add = true
		domElement = document.Call("createTextNode", typedNode.Text)
	}

	err := currentNode.Mount(domElement)
	if err != nil {
		return styles, err
	}

	for _, child := range currentNode.GetChildren() {
		childStyles, err := recursivelyMount(domElement, child)
		if err != nil {
			return styles, err
		}

		for _, style := range childStyles {
			styles = append(styles, style)
		}
	}

	if add {
		lastElement.Call("appendChild", domElement)
	}

	return styles, nil
}
