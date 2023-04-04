package nodes

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/utils"
)

type HTMLNode struct {
	baseNode

	DomNode js.Value

	ActiveClass    string
	Namespace      string
	DomID          string
	Tag            string
	Classes        []string
	Attributes     map[string]string
	EventListeners map[string]events.EventListener
	Children       []Node
	Styles         []string
}

func NewHTMLNode(tag string, attributes map[string]interface{}, children []Node) *HTMLNode {
	attrs, listeners := utils.ExtractAttributes(attributes)

	var id string
	if val, ok := attrs["id"]; ok {
		id = val
	}

	var classes []string
	if val, ok := attrs["class"]; ok {
		classes = strings.Split(val, " ")
	}

	return &HTMLNode{
		DomID:          id,
		Tag:            tag,
		Classes:        classes,
		Attributes:     attrs,
		EventListeners: listeners,
		Children:       children,
		Styles:         []string{},
	}
}

func (n *HTMLNode) Render() Node {
	for _, child := range n.Children {
		child.Render()
	}

	return n
}

func (n *HTMLNode) Update(newAttributes map[string]interface{}) {
	attrs, listeners := utils.ExtractAttributes(newAttributes)

	if val, ok := attrs["id"]; ok {
		n.DomID = val
		delete(attrs, "id")
	}

	if val, ok := attrs["class"]; ok {
		n.Classes = strings.Split(val, " ")
		delete(attrs, "class")
	}

	n.Attributes = attrs
	n.EventListeners = listeners
}

func (n *HTMLNode) Mount(domNode js.Value) {
	n.DomNode = domNode
	// TODO: Apply everything
}

func (n *HTMLNode) ToString() string {
	content := ""
	for _, child := range n.Children {
		content += child.ToString()
	}

	attributesString := make([]string, len(n.Attributes))
	count := 0
	for key, val := range n.Attributes {
		attributesString[count] = fmt.Sprintf("%s=\"%s\"", key, val)
		count += 1
	}

	return fmt.Sprintf(
		`<%s id="%s" class="%s"%s>%s</%s>`,
		n.Tag,
		n.DomID,
		strings.Join(n.Classes, " "),
		strings.Join(attributesString, " "),
		content,
		n.Tag,
	)
}

func (n *HTMLNode) Diff(other Node) bool {
	otherAsHtml, ok := other.(*HTMLNode)
	if !ok {
		return false
	}

	if otherAsHtml.Tag != n.Tag || otherAsHtml.DomID != n.DomID {
		return false
	}

	if len(otherAsHtml.Classes) != len(n.Classes) {
		return false
	}

	for i, class := range n.Classes {
		if class != otherAsHtml.Classes[i] {
			return false
		}
	}

	if len(otherAsHtml.Attributes) != len(n.Attributes) {
		return false
	}

	for key, val := range n.Attributes {
		otherVal, ok := otherAsHtml.Attributes[key]
		if !ok {
			return false
		}

		if val != otherVal {
			return false
		}
	}

	// We don't check event listeners here, they should always be updated

	if len(otherAsHtml.Styles) != len(n.Styles) {
		return false
	}

	for i, style := range n.Styles {
		if style != otherAsHtml.Styles[i] {
			return false
		}
	}

	if len(otherAsHtml.Children) != len(n.Children) {
		return false
	}

	return true
}

func (n *HTMLNode) InsertChildren(node Node, position int) error {
	// Insert at the end on a -1
	if position < 0 {
		n.Children = append(n.Children, node)
		return nil
	}

	newChildren := make([]Node, len(n.Children)+1)
	index := 0
	for _, child := range n.Children {
		if index == position {
			newChildren[index] = node
			index++
		}
		newChildren[index] = child
		index++
	}

	n.Children = newChildren

	return nil
}

func (n *HTMLNode) ReplaceChildren(old, new Node) error {
	for index, child := range n.Children {
		if child == old {
			n.Children[index] = new
			break
		}
	}
	return nil
}

func (n *HTMLNode) RemoveChildren(node Node) error {
	newChildren := make([]Node, len(n.Children)-1)
	index := 0
	for _, child := range n.Children {
		if index >= len(newChildren) {
			// Could not find the child
			return nil
		}
		if child == node {
			continue
		}
		newChildren[index] = child
		index++
	}

	n.Children = newChildren

	return nil
}

func (n *HTMLNode) Style(styling string) *HTMLNode {
	// Generate a random CSS class name of length 10
	className := utils.RandomString(10)

	n.Styles = append(n.Styles, fmt.Sprintf(".%s{%s}", className, styling))
	n.Classes = append(n.Classes, className)
	n.ActiveClass = className
	return n
}

func (n *HTMLNode) SelectorStyle(selector, styling string) *HTMLNode {
	n.Styles = append(n.Styles, fmt.Sprintf(".%s%s{%s}", n.ActiveClass, selector, styling))
	return n
}
