package go_lander

import "strings"

func Html(tag string, attributes map[string]string, children []Node) *HtmlNode {
	tagname, id, classes := hyperscript(tag)

	if val, ok := attributes["id"]; ok {
		id = val
	}

	if val, ok := attributes["class"]; ok {
		classes = strings.Split(val, " ")
	}

	return newHtmlNode(tagname, id, classes, attributes, children)
}

func Svg(tag string, attributes map[string]string, children []Node) *HtmlNode {
	node := Html(tag, attributes, children)
	node.namespace = "http://www.w3.org/2000/svg"
	return node
}

func Text(text string) *TextNode {
	return newTextNode(text)
}

func Fragment(children []Node) *FragmentNode {
	return newFragmentNode(children)
}

func Component(factory FunctionComponent, attributes map[string]string, children []Node) *FuncNode {
	return newFuncNode(factory, attributes, children)
}
