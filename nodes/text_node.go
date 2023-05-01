package nodes

import (
	"syscall/js"
)

// TextNode is an implementation of the Node interface which implements the logic to handle
// and render HTML text elements inside Lander.
type TextNode struct {
	baseNode

	// DomNode is the real DOM node associated with this virtual node. If set, this node is
	// mounted.
	DomNode js.Value

	// Text is the stored text of this node, is assigned directly as the text of the DomNode.
	Text string
}

// NewTextNode creates a new HTML node with the provided information.
func NewTextNode(text string) *TextNode {
	return &TextNode{
		Text: text,
	}
}

// Update updates this HTML node with the provided text, then applies the changes to
// the underlying real DOM node.
func (n *TextNode) Update(newText string) {
	n.Text = newText

	n.DomNode.Set("nodeValue", n.Text)
}

// Mount sets the real DOM node on this text node, the applies the text on the underlying real
// DOM node.
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
		return true
	}

	return otherAsText.Text != n.Text
}

func (n *TextNode) Type() NodeType {
	return TextNodeType
}
