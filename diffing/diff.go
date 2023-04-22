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
	prev, old, new nodes.Node) ([]Patch, []string, error) {

	var patches []Patch
	var currentStyles []string

	var oldChildren []nodes.Node
	var newChildren []nodes.Node

	fmt.Printf("Diffing %T, %v against %T, %v\n", old, old, new, new)
	if new == nil {
		fmt.Println("New was missing, removing")
		// If the new is missing, then we should remove unneeded children
		patches = append(patches, newPatchRemove(prev, old))
		return patches, currentStyles, nil
	} else if old == nil {
		fmt.Println("Old was missing, inserting")
		// If the old node is missing, then we are mounting for the first time
		patches = append(patches, newPatchInsert(listenerFunc, prev, new))
		switch typedNode := new.(type) {
		case *nodes.HTMLNode:
			newChildren = typedNode.Children

			currentStyles = append(currentStyles, typedNode.Styles...)
		}

		return patches, currentStyles, nil
	} else if reflect.TypeOf(old) != reflect.TypeOf(new) {
		fmt.Println("Types were different, replacing")
		// If both nodes exist, but they are of a different type, replace and patch
		patches = append(patches, newPatchReplace(listenerFunc, prev, old, new))

		switch typedNode := old.(type) {
		case *nodes.HTMLNode:
			newChildren = typedNode.Children

			currentStyles = append(currentStyles, new.(*nodes.HTMLNode).Styles...)
		}

		return patches, currentStyles, nil
	} else if old.Diff(new) {
		fmt.Println("Nodes were different, updating")
		// If both nodes have the same type, but have differences
		switch typedNode := old.(type) {
		case *nodes.HTMLNode:
			patches = append(patches, newPatchHTML(listenerFunc, typedNode, new.(*nodes.HTMLNode)))
			oldChildren = typedNode.Children
			newConverted := new.(*nodes.HTMLNode)
			newChildren = newConverted.Children

			currentStyles = append(currentStyles, new.(*nodes.HTMLNode).Styles...)
		case *nodes.TextNode:
			patches = append(patches, newPatchText(prev, typedNode, new.(*nodes.TextNode).Text))
		default:
			return nil, []string{}, fmt.Errorf("somehow got neither a text, nor a HTML node during patching, cannot process node")
		}
	} else {
		fmt.Println("No changes")
		// If the two nodes are the same, still run on the children
		if oldConverted, ok := old.(*nodes.HTMLNode); ok {
			oldChildren = oldConverted.Children
			currentStyles = append(currentStyles, oldConverted.Styles...)
		}
		if newConverted, ok := new.(*nodes.HTMLNode); ok {
			newChildren = newConverted.Children
		}
	}

	// Start by running through the old children and patch individually
	count := 0
	for _, child := range oldChildren {
		var newChild nodes.Node
		if count < len(newChildren) {
			newChild = newChildren[count]
		}

		childPatches, styles, err := GeneratePatches(listenerFunc, old, child, newChild)
		if err != nil {
			return nil, []string{}, err
		}
		patches = append(patches, childPatches...)
		currentStyles = append(currentStyles, styles...)

		count += 1
	}

	// If we still have new nodes left, then loop over them and insert
	if count >= len(newChildren) {
		return patches, currentStyles, nil
	}

	for _, child := range newChildren[count:] {
		childPatches, styles, err := GeneratePatches(listenerFunc, old, nil, child)
		if err != nil {
			return nil, []string{}, err
		}
		patches = append(patches, childPatches...)
		currentStyles = append(currentStyles, styles...)
	}

	return patches, currentStyles, nil
}
