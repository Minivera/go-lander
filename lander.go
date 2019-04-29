package go_lander

func Html(tag string, attributes map[string]interface{}, children []Node) *HtmlNode {
	tagname, id, classes := hyperscript(tag)

	node, err := newHtmlNode(tagname, id, classes, attributes, children)
	if err != nil {
		// We don't want to return an error for ease of use, panic instead
		panic(err)
	}

	return node
}

func Svg(tag string, attributes map[string]interface{}, children []Node) *HtmlNode {
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

func Component(factory FunctionComponent, attributes map[string]interface{}, children []Node) *FuncNode {
	return newFuncNode(factory, attributes, children)
}
