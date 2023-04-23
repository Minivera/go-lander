//go:build js && wasm

package lander

import (
	"github.com/minivera/go-lander/nodes"
)

func Html(tag string, attributes nodes.Attributes, children nodes.Children) *nodes.HTMLNode {
	return nodes.NewHTMLNode(tag, attributes, children)
}

func Svg(tag string, attributes nodes.Attributes, children nodes.Children) *nodes.HTMLNode {
	node := Html(tag, attributes, children)
	node.Namespace = "http://www.w3.org/2000/svg"
	return node
}

func Text(text string) *nodes.TextNode {
	return nodes.NewTextNode(text)
}

func Component(factory nodes.FunctionComponent, attributes nodes.Props, children nodes.Children) *nodes.FuncNode {
	return nodes.NewFuncNode(factory, attributes, children)
}
