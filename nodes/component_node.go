package nodes

type FunctionComponent[Props map[string]interface{}] func(attributes Props, children []Node) Node

type FuncNode[Props map[string]interface{}] struct {
	baseNode

	factory       FunctionComponent[Props]
	givenChildren []Node

	Properties   Props
	RenderResult Node
}

func NewFuncNode[Props map[string]interface{}](factory FunctionComponent[Props], attributes map[string]interface{}, givenChildren []Node) *FuncNode[Props] {
	return &FuncNode[Props]{
		Properties:    attributes,
		factory:       factory,
		givenChildren: givenChildren,
	}
}

func (n *FuncNode[Props]) Update(newAttributes Props, newChildren []Node) {
	n.Properties = newAttributes
	n.givenChildren = newChildren
}

func (n *FuncNode[Props]) Render() Node {
	n.RenderResult = n.factory(n.Properties, n.givenChildren)

	n.RenderResult.Render()

	return n.RenderResult
}

func (n *FuncNode[Props]) ToString() string {
	return n.RenderResult.ToString()
}

func (n *FuncNode[Props]) Diff(other Node) bool {
	otherAsFunc, ok := other.(*FuncNode[Props])
	if !ok {
		return false
	}

	if &otherAsFunc.factory != &n.factory {
		return false
	}

	if len(otherAsFunc.Properties) != len(n.Properties) {
		return false
	}

	for key, val := range n.Properties {
		otherVal, ok := otherAsFunc.Properties[key]
		if !ok {
			return false
		}

		if val != otherVal {
			return false
		}
	}

	// We check if any of the given children were dirty in the general diff code
	if len(otherAsFunc.givenChildren) != len(n.givenChildren) {
		return false
	}

	return true
}
