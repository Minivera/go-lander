package nodes

import "syscall/js"

type TextNode struct {
	baseNode

	DomNode js.Value

	Text string
}

func NewTextNode(text string) *TextNode {
	return &TextNode{
		Text: text,
	}
}

func (n *TextNode) Update(newText string) {
	n.Text = newText

	n.DomNode.Set("nodeValue", n.Text)
}

func (n *TextNode) Mount(domNode js.Value) {
	n.DomNode = domNode
	n.DomNode.Set("nodeValue", n.Text)
}

func (n *TextNode) ToString() string {
	return n.Text
}

func (n *TextNode) Diff(other Node) bool {
	otherAsText, ok := other.(*TextNode)
	if !ok {
		return false
	}

	return otherAsText.Text == n.Text
}
