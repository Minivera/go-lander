// +build js,wasm

package lander

import (
	"fmt"
	"sort"
	"strings"
)

type Node interface {
	ID() uint64
	SetID(uint64)
	Create(uint64) error
	Position(parent, next, prev Node) error
	Update(map[string]interface{}) error
	Render() error
	ToString() string
	GetChildren() []Node
	InsertChildren(Node, int) error
	ReplaceChildren(old, new Node) error
	RemoveChildren(Node) error
	Clone() Node
}

type baseNode struct {
	id uint64

	Parent, NextSibling, PreviousSibling Node
}

func (n *baseNode) ID() uint64 {
	return n.id
}

func (n *baseNode) SetID(id uint64) {
	n.id = id
}

func (n *baseNode) Create(id uint64) error {
	n.id = id
	return nil
}

func (n *baseNode) Position(parent, next, prev Node) error {
	n.Parent = parent
	n.NextSibling = next
	n.PreviousSibling = prev
	return nil
}

func (n *baseNode) Update(newAttributes map[string]interface{}) error {
	return nil
}

func (n *baseNode) Render() error {
	return nil
}

func (n *baseNode) ToString() string {
	return ""
}

func (n *baseNode) GetChildren() []Node {
	return []Node{}
}

func (n *baseNode) InsertChildren(node Node, position int) error {
	return nil
}

func (n *baseNode) ReplaceChildren(old, new Node) error {
	return nil
}

func (n *baseNode) RemoveChildren(node Node) error {
	return nil
}

type FunctionComponent func(attributes map[string]interface{}, children []Node) []Node

type HTMLNode struct {
	baseNode

	namespace   string
	activeClass string

	DomID          string
	Tag            string
	Classes        []string
	Attributes     map[string]string
	EventListeners map[string]EventListener
	Children       []Node
	Styles         []string
}

func newHTMLNode(tag, id string, classes []string, attributes map[string]interface{}, children []Node) (*HTMLNode, error) {
	attrs, events, err := extractAttributes(attributes)
	if err != nil {
		return nil, err
	}

	if val, ok := attrs["id"]; ok {
		id = val
	}

	if val, ok := attrs["class"]; ok {
		classes = strings.Split(val, " ")
	}

	return &HTMLNode{
		DomID:          id,
		Tag:            tag,
		Classes:        classes,
		Attributes:     attrs,
		EventListeners: events,
		Children:       children,
		Styles:         []string{},
	}, nil
}

func (n *HTMLNode) Update(newAttributes map[string]interface{}) error {
	attrs, events, err := extractAttributes(newAttributes)
	if err != nil {
		return err
	}

	if val, ok := attrs["id"]; ok {
		n.DomID = val
		delete(attrs, "id")
	}

	if val, ok := attrs["class"]; ok {
		n.Classes = strings.Split(val, " ")
		delete(attrs, "class")
	}

	n.Attributes = attrs
	n.EventListeners = events

	return nil
}

func (n *HTMLNode) ToString() string {
	content := ""
	for _, child := range n.Children {
		content += child.ToString()
	}

	// Arrange keys in sorted order to make them more predictable
	attrsKeys := make([]string, 0, len(n.Attributes))
	for key, _ := range n.Attributes {
		attrsKeys = append(attrsKeys, key)
	}
	sort.Strings(attrsKeys)

	attributes := ""
	for _, key := range attrsKeys {
		if key == "class" || key == "id" {
			continue
		}
		attributes += fmt.Sprintf(` %s="%s"`, key, n.Attributes[key])
	}

	return fmt.Sprintf(
		`<%s id="%s" class="%s"%s>%s</%s>`,
		n.Tag,
		n.DomID,
		strings.Join(n.Classes, " "),
		attributes,
		content,
		n.Tag,
	)
}

func (n *HTMLNode) GetChildren() []Node {
	return n.Children
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

func (n *HTMLNode) Clone() Node {
	clonedAttrs := make(map[string]string, len(n.Attributes))
	for key, value := range n.Attributes {
		clonedAttrs[key] = value
	}

	clonedEvents := make(map[string]EventListener, len(n.EventListeners))
	for key, value := range n.EventListeners {
		clonedEvents[key] = value
	}

	clonedChildren := make([]Node, len(n.Children))
	for index, child := range n.Children {
		clonedChildren[index] = child.Clone()
	}

	clonedClasses := make([]string, len(n.Classes))
	for index, val := range n.Classes {
		clonedClasses[index] = val
	}

	clonedStyles := make([]string, len(n.Styles))
	for index, val := range n.Styles {
		clonedStyles[index] = val
	}

	return &HTMLNode{
		baseNode: baseNode{
			id: n.id,
		},
		namespace:      n.namespace,
		Tag:            n.Tag,
		DomID:          n.DomID,
		Attributes:     clonedAttrs,
		EventListeners: clonedEvents,
		Classes:        clonedClasses,
		Children:       clonedChildren,
		Styles:         clonedStyles,
	}
}

func (n *HTMLNode) Style(styling string) *HTMLNode {
	// Generate a random CSS class name of length 10
	className := randomString(10)

	n.Styles = append(n.Styles, fmt.Sprintf(".%s{%s}", className, styling))
	n.Classes = append(n.Classes, className)
	n.activeClass = className
	return n
}

func (n *HTMLNode) SelectorStyle(selector, styling string) *HTMLNode {
	n.Styles = append(n.Styles, fmt.Sprintf(".%s%s{%s}", n.activeClass, selector, styling))
	return n
}

type TextNode struct {
	baseNode

	Text string
}

func newTextNode(text string) *TextNode {
	return &TextNode{
		Text: text,
	}
}

func (n *TextNode) Update(newAttributes map[string]interface{}) error {
	attrs, _, err := extractAttributes(newAttributes)
	if err != nil {
		return err
	}

	if val, ok := attrs["text"]; ok {
		n.Text = val
	}
	return nil
}

func (n *TextNode) ToString() string {
	return n.Text
}

func (n *TextNode) Clone() Node {
	return &TextNode{
		baseNode: baseNode{
			id: n.id,
		},
		Text: n.Text,
	}
}

type FragmentNode struct {
	baseNode

	Children []Node
}

func newFragmentNode(children []Node) *FragmentNode {
	return &FragmentNode{
		Children: children,
	}
}

func (n *FragmentNode) ToString() string {
	content := ""
	for _, child := range n.Children {
		content += child.ToString()
	}
	return content
}

func (n *FragmentNode) GetChildren() []Node {
	return n.Children
}

func (n *FragmentNode) InsertChildren(node Node, position int) error {
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

func (n *FragmentNode) ReplaceChildren(old, new Node) error {
	for index, child := range n.Children {
		if child == old {
			n.Children[index] = new
			break
		}
	}
	return nil
}

func (n *FragmentNode) RemoveChildren(node Node) error {
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

func (n *FragmentNode) Clone() Node {
	clonedChildren := make([]Node, len(n.Children))
	for index, child := range n.Children {
		clonedChildren[index] = child
	}

	return &FragmentNode{
		baseNode: baseNode{
			id: n.id,
		},
		Children: clonedChildren,
	}
}

type FuncNode struct {
	baseNode

	factory       FunctionComponent
	givenChildren []Node

	Attributes map[string]interface{}
	Children   []Node
}

func newFuncNode(factory FunctionComponent, attributes map[string]interface{}, givenChildren []Node) *FuncNode {
	return &FuncNode{
		Attributes:    attributes,
		factory:       factory,
		givenChildren: givenChildren,
	}
}

func (n *FuncNode) Update(newAttributes map[string]interface{}) error {
	n.Attributes = newAttributes
	return nil
}

func (n *FuncNode) Render() error {
	n.Children = n.factory(n.Attributes, n.givenChildren)
	return nil
}

func (n *FuncNode) ToString() string {
	content := ""
	for _, child := range n.Children {
		content += child.ToString()
	}
	return content
}

func (n *FuncNode) GetChildren() []Node {
	return n.Children
}

func (n *FuncNode) Clone() Node {
	clonedAttrs := make(map[string]interface{}, len(n.Attributes))
	for key, value := range n.Attributes {
		clonedAttrs[key] = value
	}

	clonedChildren := make([]Node, len(n.givenChildren))
	for index, child := range n.givenChildren {
		clonedChildren[index] = child
	}

	return &FuncNode{
		baseNode: baseNode{
			id: n.id,
		},
		factory:       n.factory,
		givenChildren: clonedChildren,
		Attributes:    clonedAttrs,
		Children:      []Node{},
	}
}
