package go_lander

import (
	"fmt"
	"reflect"
	"strings"
	"syscall/js"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
)

var document js.Value

func init() {
	document = js.Global().Get("document")
}

func validElement(value js.Value) bool {
	return value != (js.Value{}) || value != js.Null() || value != js.Undefined()
}

type DomEnvironment struct {
	root        string
	rootElement js.Value
	tree        Node

	instancesToHash map[uint64]Node

	patches []patch
}

func NewLander(root string, rootNode Node) *DomEnvironment {
	return &DomEnvironment{
		root: root,
		tree: rootNode,

		instancesToHash: make(map[uint64]Node),
	}
}

func (e *DomEnvironment) Mount() error {
	err := e.buildNodeRecursively(e.tree, "[root]")
	if err != nil {
		return err
	}

	err = e.mountToDom(e.root, e.tree)
	if err != nil {
		return err
	}

	return nil
}

func (e *DomEnvironment) Update() error {
	newTree := e.tree.Clone()

	err := e.buildNodeRecursively(newTree, "[root]")
	if err != nil {
		return err
	}

	err = e.patchDom(newTree)
	if err != nil {
		return err
	}

	return nil
}

func (e *DomEnvironment) buildNodeRecursively(currentNode Node, positionString string) error {
	id := hashPosition(positionString)

	if _, ok := e.instancesToHash[id]; !ok {
		err := currentNode.Create(id)
		if err != nil {
			return err
		}
	}

	err := currentNode.Render()
	if err != nil {
		return err
	}

	for index, child := range currentNode.GetChildren() {
		err := e.buildNodeRecursively(child, positionString+fmt.Sprintf("[%d]", index))
		if err != nil {
			return err
		}

		if len(currentNode.GetChildren()) <= 1 {
			err := child.Position(currentNode, nil, nil)
			if err != nil {
				return err
			}

			continue
		}

		if index <= 0 {
			err := child.Position(currentNode, currentNode.GetChildren()[index+1], nil)
			if err != nil {
				return err
			}
		} else if index < len(currentNode.GetChildren())-1 {
			err := child.Position(currentNode, currentNode.GetChildren()[index+1], currentNode.GetChildren()[index-1])
			if err != nil {
				return err
			}
		} else {
			err := child.Position(currentNode, nil, currentNode.GetChildren()[index-1])
			if err != nil {
				return err
			}
		}
	}

	e.instancesToHash[currentNode.ID()] = currentNode

	return nil
}

func (e *DomEnvironment) mountToDom(rootElement string, vTree Node) error {
	e.rootElement = document.Call("querySelector", rootElement)
	if !validElement(e.rootElement) {
		return fmt.Errorf("failed to find mount parent using query selector %q", rootElement)
	}

	styles, err := e.recursivelyMount(e.rootElement, vTree)
	if err != nil {
		return err
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	stylesString, err := m.String("text/css", strings.Join(styles, " "))
	if err != nil {
		return err
	}

	head := document.Call("querySelector", "head")
	if !validElement(head) {
		return fmt.Errorf("failed to find heads using query selector")
	}

	styleTag := document.Call("createElement", "style")
	styleTag.Set("id", "lander-style-tag")
	styleTag.Set("innerHTML", stylesString)
	head.Call("appendChild", styleTag)

	return nil
}

func (e *DomEnvironment) recursivelyMount(lastElement js.Value, currentNode Node) ([]string, error) {
	add := false
	domElement := lastElement
	styles := make([]string, 250)

	switch typedNode := currentNode.(type) {
	case *HtmlNode:
		add = true
		domElement = newHTMLElement(document, typedNode)

		for _, style := range typedNode.Styles {
			styles = append(styles, style)
		}
	case *TextNode:
		add = true
		domElement = document.Call("createTextNode", typedNode.Text)
	}

	err := currentNode.Mount(domElement)
	if err != nil {
		return styles, err
	}

	for _, child := range currentNode.GetChildren() {
		childStyles, err := e.recursivelyMount(domElement, child)
		if err != nil {
			return styles, err
		}

		for _, style := range childStyles {
			styles = append(styles, style)
		}
	}

	if add {
		lastElement.Call("appendChild", domElement)
	}

	return styles, nil
}

func (e *DomEnvironment) patchDom(tree Node) error {
	err := e.generatePatches(nil, e.tree, tree, e.rootElement)
	if err != nil {
		return err
	}

	for _, patch := range e.patches {
		err := patch.execute(document)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *DomEnvironment) generatePatches(prev, old, new Node, lastDomElement js.Value) error {
	domElement := lastDomElement
	if node, ok := old.(*HtmlNode); ok {
		domElement = node.domNode
	}

	// If the old is missing, we need to insert the new node
	if old == nil {
		e.patches = append(e.patches, newPatchInsert(lastDomElement, prev, new))
		return nil
	}

	// If the new is missing, we need to remove the old node
	if new == nil {
		e.patches = append(e.patches, newPatchRemove(lastDomElement, prev, old))
		return nil
	}

	// If both nodes are identical, run on children
	if old.ID() == new.ID() && hashNode(old) == hashNode(new) {
		for index, child := range old.GetChildren() {
			var newChild Node
			if index <= len(new.GetChildren()) {
				newChild = new.GetChildren()[index]
			}

			err := e.generatePatches(old, child, newChild, domElement)
			if err != nil {
				return err
			}
		}

		return nil
	}

	// If both nodes are similar
	if old.ID() == new.ID() && reflect.TypeOf(old) == reflect.TypeOf(new) {
		if val, ok := new.(*TextNode); ok {
			e.patches = append(e.patches, newPatchText(lastDomElement, old, val.Text))
			return nil
		}
		if _, ok := new.(*HtmlNode); ok {
			e.patches = append(e.patches, newPatchHtml(old, new))
		}

		for index, child := range old.GetChildren() {
			var newChild Node
			if index <= len(new.GetChildren()) {
				newChild = new.GetChildren()[index]
			}

			err := e.generatePatches(old, child, newChild, domElement)
			if err != nil {
				return err
			}
		}

		return nil
	}

	e.patches = append(e.patches, newPatchReplace(lastDomElement, prev, old, new))

	return nil
}
