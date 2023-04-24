package diffing

import (
	"fmt"
	"reflect"
	"syscall/js"

	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

// GeneratePatches generate a set of patches to update the real DOM and the virtual dom passed as the
// old node. It will run recursively on all nodes of the tree and return the patches in a slice to be
// executed sequentially. GeneratePatches expects the tree it is given to be made exclusively of HTML
// nodes (HTML and text), and no components.
func GeneratePatches(listenerFunc func(listener events.EventListenerFunc, this js.Value, args []js.Value) interface{},
	lastComponent *nodes.FuncNode, prev, old, new nodes.Node) ([]Patch, []string, error) {

	var patches []Patch
	var currentStyles []string

	var oldChildren []nodes.Node
	var newChildren []nodes.Node

	fmt.Printf("Diffing %T, %v against %T, %v\n", old, old, new, new)
	if new == nil {
		if typedNode, ok := old.(*nodes.FuncNode); ok {
			// If we hit a function not when there's a need to replace, then we should render
			// that function node and keep going as-is. The final tree should never include
			// any component.
			fmt.Println("New was missing and old node is a component, rendering and keep going")
			context.RegisterComponent(lastComponent)
			return GeneratePatches(listenerFunc, typedNode, prev, typedNode.Render(context.CurrentContext), new)
		}

		fmt.Println("New was missing, removing")
		// If the new is missing, then we should remove unneeded children
		patches = append(patches, newPatchRemove(prev, old))

		// Register the last seen component as now a unmount component if the first encountered dom node was
		// to be removed
		if lastComponent != nil {
			context.UnregisterAllComponentContexts(lastComponent)
			context.RegisterComponentContext("unmount", lastComponent)
		}

		return patches, currentStyles, nil
	} else if old == nil {
		if typedNode, ok := new.(*nodes.FuncNode); ok {
			// If we hit a function not when there's a need to insert, then we should render
			// that function node and keep going as-is. The final tree should never include
			// any component.
			fmt.Println("Old is missing and new node is a component, rendering and keep going")
			context.RegisterComponentContext("render", typedNode)
			return GeneratePatches(listenerFunc, typedNode, prev, old, typedNode.Render(context.CurrentContext))
		}

		fmt.Println("Old was missing, inserting")
		// If the old node is missing, then we are mounting for the first time
		patches = append(patches, newPatchInsert(listenerFunc, prev, new))

		// Register the last seen component as now a mount component if the first encountered dom node was
		// to be inserted
		if lastComponent != nil {
			context.RegisterComponentContext("mount", lastComponent)
		}

		if typedNode, ok := new.(*nodes.HTMLNode); ok {
			newChildren = typedNode.Children

			currentStyles = append(currentStyles, typedNode.Styles...)
		}

		return patches, currentStyles, nil
	} else if reflect.TypeOf(old) != reflect.TypeOf(new) {
		if typedNode, ok := new.(*nodes.FuncNode); ok {
			// If we hit a function not when there's a need to replace, then we should render
			// that function node and keep going as-is. The final tree should never include
			// any component.
			fmt.Println("Types were different and new node is a component, rendering and keep going")
			context.RegisterComponentContext("render", typedNode)
			return GeneratePatches(listenerFunc, typedNode, prev, old, typedNode.Render(context.CurrentContext))
		}

		fmt.Println("Types were different, replacing")
		// If both nodes exist, but they are of a different type, replace and patch
		patches = append(patches, newPatchReplace(listenerFunc, prev, old, new))

		// Register the last seen component as now a mount component if the first encountered dom node was
		// to be replaced
		if lastComponent != nil {
			context.RegisterComponentContext("mount", lastComponent)
		}

		if typedNode, ok := new.(*nodes.HTMLNode); ok {
			newChildren = typedNode.Children

			currentStyles = append(currentStyles, typedNode.Styles...)
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

		childPatches, styles, err := GeneratePatches(listenerFunc, nil, old, child, newChild)
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
		childPatches, styles, err := GeneratePatches(listenerFunc, nil, old, nil, child)
		if err != nil {
			return nil, []string{}, err
		}
		patches = append(patches, childPatches...)
		currentStyles = append(currentStyles, styles...)
	}

	return patches, currentStyles, nil
}
