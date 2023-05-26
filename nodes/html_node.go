package nodes

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/minivera/go-lander/events"
)

// Attributes is a map of properties to assign to an element. Technically interchangeable with
// Props or `map[string]interface{}`, this type is provided for convenience.
type Attributes = map[string]interface{}

// HTMLNode is an implementation of the Node interface which implements the logic to handle
// and render HTML elements inside Lander.
type HTMLNode struct {
	baseNode

	// DomNode is the real DOM node associated with this virtual node. If set, this node is
	// mounted.
	DomNode js.Value

	// ActiveClass is an active, random, class given to this element by the styling function.
	ActiveClass string
	// Namespace is th XHTML namespace of this element, if any. Will be used to create the DOM
	// element with a namespace if needed.
	Namespace string
	// DomID is the ID of the element in the DOM, set as "id" on the element itself.
	DomID string
	// Tag is the HMTL tag of this element, such as "div" or "span".
	Tag string
	// Classes is a list of CSS classes to assign to this element.
	Classes []string
	// Attributes is a map of string only attributes to assign using setAttribute on the DOM element
	// properties should be stored in Properties.
	Attributes map[string]string
	// Properties is a map of any type to assign on the DOM element directly as object properties.
	// attributes should be stored in Attributes.
	Properties map[string]interface{}
	// EventListeners is a map of event listeners to their events. Does not expect the "on" prefix.
	// events listeners are added directly on the DOM element.
	EventListeners map[string]*events.EventListener
	// Children is a slice of the children provided to this element.
	Children []Node
	// A slice of the styles assigned to this element as CSS strings. Not minified, they are
	// valid CSS definitions.
	Styles []string
}

// NewHTMLNode creates a new HTML node with the provided information.
func NewHTMLNode(tag string, attributes Attributes, children []Node) *HTMLNode {
	attrs, props, listeners := ExtractAttributes(attributes)

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
		Properties:     props,
		EventListeners: listeners,
		Children:       children,
		Styles:         []string{},
	}
}

// Update updates this HTML node with the provided attributes map. The map will be extracted to
// attributes, props, and event listeners using ExtractAttributes, then applied to the virtual
// DOM node and the underlying real DOM node.
func (n *HTMLNode) Update(newAttributes map[string]interface{}) {
	oldAttributes := n.Attributes
	oldProps := n.Properties
	attrs, props, listeners := ExtractAttributes(newAttributes)

	n.DomID = ""

	if val, ok := attrs["id"]; ok {
		n.DomID = val
		delete(attrs, "id")
	}

	if val, ok := attrs["class"]; ok {
		n.Classes = strings.Split(val, " ")
		delete(attrs, "class")
	}

	n.Attributes = attrs
	n.Properties = props
	n.EventListeners = listeners

	// Remove, then set the new attributes/properties
	for key := range oldAttributes {
		n.DomNode.Call("removeAttribute", key)
	}
	for key := range oldProps {
		n.DomNode.Set(key, nil)
	}

	for key, value := range n.Attributes {
		n.DomNode.Call("setAttribute", key, value)
	}
	for key, value := range n.Properties {
		n.DomNode.Set(key, value)
	}

	// Clear the old class list, then set the new classes
	classList := n.DomNode.Get("classList")
	classesLength := classList.Get("length").Int()
	for i := 0; i < classesLength; i += 1 {
		classList.Call("remove", classList.Call("item", i))
	}

	for _, value := range n.Classes {
		classList.Call("add", value)
	}

	// Set the ID if needed, if not, remove it
	if n.DomID != "" {
		n.DomNode.Set("id", n.DomID)
	} else {
		n.DomNode.Delete("id")
	}
}

// Mount sets the real DOM node on this HTML node, the applies the attributes, props, and event listeners
// on the underlying real DOM node.
func (n *HTMLNode) Mount(domNode js.Value) {
	n.DomNode = domNode

	// Attributes
	for name, value := range n.Attributes {
		n.DomNode.Call("setAttribute", name, value)
	}

	// Properties
	for name, value := range n.Properties {
		n.DomNode.Set(name, value)
	}

	// Classes
	classList := n.DomNode.Get("classList")
	for _, value := range n.Classes {
		classList.Call("add", value)
	}

	// Add the active class
	if n.ActiveClass != "" {
		classList.Call("add", n.ActiveClass)
	}

	// ID if set
	if n.DomID != "" {
		n.DomNode.Set("id", n.DomID)
	}
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

	html := fmt.Sprintf("<%s ", n.Tag)
	if n.DomID != "" {
		html += fmt.Sprintf("id=\"%s\" ", n.DomID)
	}

	if len(n.Classes) > 0 || n.ActiveClass != "" {
		html += fmt.Sprintf("class=\"%s\" ", strings.Join(append(n.Classes, n.ActiveClass), " "))
	}

	return fmt.Sprintf(
		"%s%s>%s</%s>",
		html,
		strings.Join(attributesString, " "),
		content,
		n.Tag,
	)
}

func (n *HTMLNode) Diff(other Node) bool {
	otherAsHtml, ok := other.(*HTMLNode)
	if !ok {
		return true
	}

	if otherAsHtml.Tag != n.Tag || otherAsHtml.DomID != n.DomID {
		return true
	}

	if len(otherAsHtml.Classes) != len(n.Classes) {
		return true
	}

	for i, class := range n.Classes {
		if class != otherAsHtml.Classes[i] {
			return true
		}
	}

	if len(otherAsHtml.Attributes) != len(n.Attributes) {
		return true
	}

	for key, val := range n.Attributes {
		otherVal, ok := otherAsHtml.Attributes[key]
		if !ok {
			return true
		}

		if val != otherVal {
			return true
		}
	}

	// We don't check event listeners here, they should always be updated

	if len(otherAsHtml.Styles) != len(n.Styles) {
		return true
	}

	// Check the styles, but remove the randomized active style class so we don't constantly
	// update the node due to randomness.
	for i, style := range n.Styles {
		if strings.Replace(style, n.ActiveClass, "", 1) !=
			strings.Replace(otherAsHtml.Styles[i], otherAsHtml.ActiveClass, "", 1) {
			return true
		}
	}

	return false
}

func (n *HTMLNode) Type() NodeType {
	return HTMLNodeType
}

// InsertChildren inserts a children at the provided position in the element's children.
// Returns an error if the children cannot be inserted. Inserts at the end if provided -1
// as the position.
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

// ReplaceChildren replaces the provided node with the new node in the element's children.
// Returns an error if the children cannot be replaced.
func (n *HTMLNode) ReplaceChildren(old, new Node) error {
	for index, child := range n.Children {
		if child == old {
			n.Children[index] = new
			break
		}
	}
	return nil
}

// RemoveChildren removed the provided children, if found, from the elements's children.
// Returns an error if the children cannot be inserted.
func (n *HTMLNode) RemoveChildren(node Node) error {
	newChildren := make([]Node, len(n.Children)-1)
	index := 0
	for _, child := range n.Children {
		if child == node {
			continue
		}
		newChildren[index] = child
		index++
	}

	n.Children = newChildren

	return nil
}

// Style will assign a random CSS class name to this node and assign the passed CSS styles to it on
// render and mount. Calling Style multiple time will override the previous styles.
func (n *HTMLNode) Style(styling string) *HTMLNode {
	// Generate a random CSS class name of length 10
	className := RandomString(10)

	n.Styles = append(n.Styles, fmt.Sprintf(".%s{%s}", className, styling))
	n.ActiveClass = className
	return n
}

// SelectorStyle uses the provided selector and creates a CSS definition using the passed CSS styles,
// which will be added to the head on render and mounts. SelectorStyle must be called after Style as it
// uses the active class name generated from Style to create the selector.
func (n *HTMLNode) SelectorStyle(selector, styling string) *HTMLNode {
	n.Styles = append(n.Styles, fmt.Sprintf(".%s %s{%s}", n.ActiveClass, selector, styling))
	return n
}
