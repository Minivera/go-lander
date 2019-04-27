package go_lander

import (
	"fmt"
	"strings"
	js "syscall/js"
)

type HashID string

type Node interface {
	Create() error
	Mount(js.Value) error
	Position(parent, next, prev Node) error
	Update(map[string]string) error
	Remove() error
	Render() error
	ToString() string
	GetChildren() []Node
	Clone() Node
}

type FunctionComponent func(attributes map[string]string, children []Node) []Node

type HtmlNode struct {
	id          HashID
	domNode     js.Value
	namespace   string
	activeClass string

	DomID                                string
	Tag                                  string
	Classes                              []string
	Attributes                           map[string]string
	Children                             []Node
	Parent, NextSibling, PreviousSibling Node
	Styles                               []string
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

func (n *HtmlNode) Create() error {
	return nil
}

func (n *HtmlNode) Mount(newNode js.Value) error {
	n.domNode = newNode
	return nil
}

func (n *HtmlNode) Position(parent, next, prev Node) error {
	n.Parent = parent
	n.NextSibling = next
	n.PreviousSibling = prev
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

func (n *HtmlNode) Render() error {
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

func (n *HtmlNode) Clone() Node {
	return &HtmlNode{
		id:         n.id,
		domNode:    n.domNode,
		namespace:  n.namespace,
		Tag:        n.Tag,
		DomID:      n.DomID,
		Attributes: n.Attributes,
		Classes:    n.Classes,
		Children:   n.Children,
		Styles:     n.Styles,
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
	id      HashID
	domNode js.Value

	Text                                 string
	Parent, NextSibling, PreviousSibling Node
}

func newTextNode(text string) *TextNode {
	return &TextNode{
		Text: text,
	}
}

func (n *TextNode) Create() error {
	return nil
}

func (n *TextNode) Mount(newNode js.Value) error {
	n.domNode = newNode
	return nil
}

func (n *TextNode) Position(parent, next, prev Node) error {
	n.Parent = parent
	n.NextSibling = next
	n.PreviousSibling = prev
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

func (n *TextNode) Render() error {
	return nil
}

func (n *TextNode) ToString() string {
	return n.Text
}

func (n *TextNode) GetChildren() []Node {
	return []Node{}
}

func (n *TextNode) Clone() Node {
	return &TextNode{
		id:      n.id,
		domNode: n.domNode,
		Text:    n.Text,
	}
}

type FragmentNode struct {
	id HashID

	Children                             []Node
	Parent, NextSibling, PreviousSibling Node
}

func newFragmentNode(children []Node) *FragmentNode {
	return &FragmentNode{
		Children: children,
	}
}

func (n *FragmentNode) Create() error {
	return nil
}

func (n *FragmentNode) Mount(newNode js.Value) error {
	return nil
}

func (n *FragmentNode) Position(parent, next, prev Node) error {
	n.Parent = parent
	n.NextSibling = next
	n.PreviousSibling = prev
	return nil
}

func (n *FragmentNode) Update(newAttributes map[string]string) error {
	return nil
}

func (n *FragmentNode) Remove() error {
	return nil
}

func (n *FragmentNode) Render() error {
	return nil
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

func (n *FragmentNode) Clone() Node {
	return &FragmentNode{
		id:       n.id,
		Children: n.Children,
	}
}

type FuncNode struct {
	id            HashID
	factory       FunctionComponent
	givenChildren []Node

	Attributes                           map[string]string
	Children                             []Node
	Parent, NextSibling, PreviousSibling Node
}

func newFuncNode(factory FunctionComponent, attributes map[string]string, givenChildren []Node) *FuncNode {
	return &FuncNode{
		Attributes:    attributes,
		factory:       factory,
		givenChildren: givenChildren,
	}
}

func (n *FuncNode) Create() error {
	return nil
}

func (n *FuncNode) Mount(newNode js.Value) error {
	return nil
}

func (n *FuncNode) Position(parent, next, prev Node) error {
	n.Parent = parent
	n.NextSibling = next
	n.PreviousSibling = prev
	return nil
}

func (n *FuncNode) Update(newAttributes map[string]string) error {
	n.Attributes = mergeAttributes(n.Attributes, newAttributes)
	return nil
}

func (n *FuncNode) Remove() error {
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
	return &FuncNode{
		id:         n.id,
		Attributes: n.Attributes,
		Children:   n.Children,
	}
}
