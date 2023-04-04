//go:build js && wasm

package lander

import "github.com/minivera/go-lander/nodes"

func Html(tag string, attributes map[string]interface{}, children []nodes.Node) *nodes.HTMLNode {
	return nodes.NewHTMLNode(tag, attributes, children)
}

func Svg(tag string, attributes map[string]interface{}, children []nodes.Node) *nodes.HTMLNode {
	node := Html(tag, attributes, children)
	node.Namespace = "http://www.w3.org/2000/svg"
	return node
}

func Text(text string) *nodes.TextNode {
	return nodes.NewTextNode(text)
}

func Fragment(children []nodes.Node) *nodes.FragmentNode {
	return nodes.NewFragmentNode(children)
}

func Component[Props map[string]interface{}](factory nodes.FunctionComponent[Props], attributes Props, children []nodes.Node) *nodes.FuncNode[Props] {
	return nodes.NewFuncNode[Props](factory, attributes, children)
}
