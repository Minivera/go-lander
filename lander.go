//go:build js && wasm

package lander

import (
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

type Child = nodes.Node
type Children = nodes.Children

type Props = nodes.Props
type EventListener = events.EventListenerFunc
type DOMEvent = events.DOMEvent

type FunctionComponent = nodes.FunctionComponent

func Html(tag string, attributes map[string]interface{}, children Children) *nodes.HTMLNode {
	return nodes.NewHTMLNode(tag, attributes, children)
}

func Svg(tag string, attributes map[string]interface{}, children Children) *nodes.HTMLNode {
	node := Html(tag, attributes, children)
	node.Namespace = "http://www.w3.org/2000/svg"
	return node
}

func Text(text string) *nodes.TextNode {
	return nodes.NewTextNode(text)
}

func Component(factory FunctionComponent, attributes Props, children Children) *nodes.FuncNode {
	return nodes.NewFuncNode(factory, attributes, children)
}
