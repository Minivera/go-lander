package diffing

import (
	"fmt"
	"reflect"
	"syscall/js"

	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

// GeneratePatches generate a set of patches to update the real DOM and the virtual dom passed as the
// old node. It will run recursively on all nodes of the tree and return the patches in a slice to be
// executed sequentially. GeneratePatches expects the tree it is given to be made exclusively of HTML
// nodes (HTML and text), and no components.
func GeneratePatches(listenerFunc func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{},
	prev, old, new nodes.Node) ([]Patch, error) {

	var patches []Patch

	var oldChildren []nodes.Node
	var newChildren []nodes.Node

	if new == nil {
		// If the new is missing, then we should remove unneeded children
		newPatchRemove(prev, new)
	} else if old == nil {
		// If the old node is missing, then we are mounting for the first time
		patches = append(patches, newPatchInsert(prev, new))
		switch typedNode := new.(type) {
		case *nodes.HTMLNode:
			patches = append(
				patches,
				newPatchHTML(listenerFunc, typedNode, new.(*nodes.HTMLNode)),
			)
			newChildren = typedNode.Children
		case *nodes.TextNode:
			patches = append(
				patches,
				newPatchText(prev, old.(*nodes.TextNode), typedNode.Text),
			)
		default:
			return nil, fmt.Errorf("somehow got neither a text, nor a HTML node during patching, cannot process node")
		}
	} else if reflect.TypeOf(old) != reflect.TypeOf(new) {
		// If both nodes exist, but they are of a different type, replace and patch
		patches = append(patches, newPatchReplace(prev, old, new))

		switch typedNode := new.(type) {
		case *nodes.HTMLNode:
			patches = append(
				patches,
				newPatchHTML(listenerFunc, typedNode, new.(*nodes.HTMLNode)),
			)
			newChildren = typedNode.Children
		case *nodes.TextNode:
			patches = append(
				patches,
				newPatchText(prev, old.(*nodes.TextNode), typedNode.Text),
			)
		default:
			return nil, fmt.Errorf("somehow got neither a text, nor a HTML node during patching, cannot process node")
		}
	} else if old.Diff(new) {
		// If both nodes have the same type, but have differences
		switch typedNode := new.(type) {
		case *nodes.HTMLNode:
			patches = append(patches, newPatchHTML(listenerFunc, typedNode, new.(*nodes.HTMLNode)))
			oldChildren = typedNode.Children
			newConverted := new.(*nodes.HTMLNode)
			newChildren = newConverted.Children
		case *nodes.TextNode:
			patches = append(patches, newPatchText(prev, old.(*nodes.TextNode), typedNode.Text))
		default:
			return nil, fmt.Errorf("somehow got neither a text, nor a HTML node during patching, cannot process node")
		}
	}

	// Start by running through the old children and patch individually
	count := 0
	for _, child := range oldChildren {
		var newChild nodes.Node
		if count < len(newChildren) {
			newChild = newChildren[count]
		}

		child.Position(old)

		childPatches, err := GeneratePatches(listenerFunc, old, child, newChild)
		if err != nil {
			return nil, err
		}
		patches = append(patches, childPatches...)

		count += 1
	}

	// If we still have new nodes left, then loop over them and insert
	if count >= len(newChildren) {
		return patches, nil
	}

	for _, child := range newChildren {
		child.Position(new)

		patches = append(patches, newPatchInsert(old, child))
	}

	return patches, nil
}
