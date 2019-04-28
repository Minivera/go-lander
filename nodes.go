package go_lander

import (
	"fmt"
	"strings"
	"syscall/js"
)

type Node interface {
	ID() uint64
	Create(uint64) error
	Mount(js.Value) error
	Position(parent, next, prev Node) error
	Update(map[string]string) error
	Remove() error
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

func (n *baseNode) Create(id uint64) error {
	n.id = id
	return nil
}

func (n *baseNode) Mount(newNode js.Value) error {
	return nil
}

func (n *baseNode) Position(parent, next, prev Node) error {
	n.Parent = parent
	n.NextSibling = next
	n.PreviousSibling = prev
	return nil
}

func (n *baseNode) Update(newAttributes map[string]string) error {
	return nil
}

func (n *baseNode) Remove() error {
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

type FunctionComponent func(attributes map[string]string, children []Node) []Node

type HtmlNode struct {
	baseNode

	domNode     js.Value
	namespace   string
	activeClass string

	DomID      string
	Tag        string
	Classes    []string
	Attributes map[string]string
	Children   []Node
	Styles     []string
}

func newHtmlNode(tag, id string, classes []string, attributes map[string]string, children []Node) *HtmlNode {
	return &HtmlNode{
		DomID:      id,
		Tag:        tag,
		Classes:    classes,
		Attributes: attributes,
		Children:   children,
	}
}

func (n *HtmlNode) Mount(newNode js.Value) error {
	n.domNode = newNode
	return nil
}

func (n *HtmlNode) Update(newAttributes map[string]string) error {
	n.Attributes = mergeAttributes(n.Attributes, newAttributes)

	if val, ok := newAttributes["id"]; ok {
		n.DomID = val
		delete(newAttributes, "id")
	}

	if val, ok := newAttributes["class"]; ok {
		n.Classes = strings.Split(val, " ")
		delete(newAttributes, "class")
	}

	return nil
}

func (n *HtmlNode) Remove() error {
	n.domNode = js.Value{}
	return nil
}

func (n *HtmlNode) ToString() string {
	content := ""
	for _, child := range n.Children {
		content += child.ToString()
	}

	attributes := ""
	for key, value := range n.Attributes {
		attributes += fmt.Sprintf(` %s="%s"`, key, value)
	}

	return fmt.Sprintf(
		`<%s id="%s" class="%s" %s>%s</%s>`,
		n.Tag,
		n.id,
		strings.Join(n.Classes, " "),
		attributes,
		content,
		n.Tag,
	)
}

func (n *HtmlNode) GetChildren() []Node {
	return n.Children
}

func (n *HtmlNode) InsertChildren(node Node, position int) error {
	// Insert at the end on a -1
	if position < 0 {
		n.Children = append(n.Children, node)
	}

	newChildren := make([]Node, len(n.Children)+1)
	index := 0
	for _, child := range n.Children {
		if index == position {
			newChildren[index] = node
			index++
		} else {
			newChildren[index] = child
		}
		index++
	}

	n.Children = newChildren

	return nil
}

func (n *HtmlNode) ReplaceChildren(old, new Node) error {
	for index, child := range n.Children {
		if child == old {
			n.Children[index] = new
			break
		}
	}
	return nil
}

func (n *HtmlNode) RemoveChildren(node Node) error {
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

func (n *HtmlNode) Clone() Node {
	clonedAttrs := make(map[string]string, len(n.Attributes))
	for key, value := range n.Attributes {
		clonedAttrs[key] = value
	}

	clonedChildren := make([]Node, len(n.Children))
	for index, child := range n.Children {
		clonedChildren[index] = child
	}

	clonedClasses := make([]string, len(n.Classes))
	for index, val := range n.Classes {
		clonedClasses[index] = val
	}

	clonedStyles := make([]string, len(n.Styles))
	for index, val := range n.Styles {
		clonedStyles[index] = val
	}

	return &HtmlNode{
		baseNode: baseNode{
			id: n.id,
		},
		domNode:    n.domNode,
		namespace:  n.namespace,
		Tag:        n.Tag,
		DomID:      n.DomID,
		Attributes: clonedAttrs,
		Classes:    clonedClasses,
		Children:   clonedChildren,
		Styles:     clonedStyles,
	}
}

func (n *HtmlNode) Style(styling string) *HtmlNode {
	// Generate a random CSS class name of length 10
	className := randomString(10)

	n.Styles = append(n.Styles, fmt.Sprintf(".%s{%s}", className, styling))
	n.Classes = append(n.Classes, className)
	n.activeClass = className
	return n
}

func (n *HtmlNode) SelectorStyle(selector, styling string) *HtmlNode {
	n.Styles = append(n.Styles, fmt.Sprintf(".%s%s{%s}", n.activeClass, selector, styling))
	return n
}

type TextNode struct {
	baseNode

	domNode js.Value

	Text string
}

func newTextNode(text string) *TextNode {
	return &TextNode{
		Text: text,
	}
}

func (n *TextNode) Mount(newNode js.Value) error {
	n.domNode = newNode
	return nil
}

func (n *TextNode) Update(newAttributes map[string]string) error {
	if val, ok := newAttributes["text"]; ok {
		n.Text = val
	}
	return nil
}

func (n *TextNode) Remove() error {
	n.domNode = js.Value{}
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
		domNode: n.domNode,
		Text:    n.Text,
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
	}

	newChildren := make([]Node, len(n.Children)+1)
	index := 0
	for _, child := range n.Children {
		if index == position {
			newChildren[index] = node
			index++
		} else {
			newChildren[index] = child
		}
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

	Attributes map[string]string
	Children   []Node
}

func newFuncNode(factory FunctionComponent, attributes map[string]string, givenChildren []Node) *FuncNode {
	return &FuncNode{
		Attributes:    attributes,
		factory:       factory,
		givenChildren: givenChildren,
	}
}

func (n *FuncNode) Update(newAttributes map[string]string) error {
	n.Attributes = mergeAttributes(n.Attributes, newAttributes)
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
	clonedAttrs := make(map[string]string, len(n.Attributes))
	for key, value := range n.Attributes {
		clonedAttrs[key] = value
	}

	clonedChildren := make([]Node, len(n.Children))
	for index, child := range n.Children {
		clonedChildren[index] = child
	}

	return &FuncNode{
		baseNode: baseNode{
			id: n.id,
		},
		Attributes: clonedAttrs,
		Children:   clonedChildren,
	}
}
