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
	prev nodes.Node, prevDOMNode js.Value, old, new nodes.Node) ([]Patch, []string, error) {

	var patches []Patch
	var currentStyles []string

	var oldChildren []nodes.Node
	var newChildren []nodes.Node

	fmt.Printf("Diffing %T, %v against %T, %v\n", old, old, new, new)
	if new == nil {
		if typedNode, ok := old.(*nodes.FuncNode); ok {
			// If we hit a function as the old node when there's a need to remove, then we
			// should do nothing and trigger an unmount on the old node, then keep going so we
			// can remove the HTML nodes.
			fmt.Println("New was missing and old node is a component, triggering a unmount")
			context.RegisterComponent(typedNode)
			context.UnregisterAllComponentContexts(typedNode)
			context.RegisterComponentContext("unmount", typedNode)
		}

		fmt.Println("New was missing, removing")
		// If the new is missing, then we should remove unneeded children
		patches = append(patches, newPatchRemove(prev, prevDOMNode, old))

		return patches, currentStyles, nil
	} else if old == nil {
		fmt.Println("Old was missing, inserting")
		// If the old node is missing, then we are mounting for the first time
		patches = append(patches, newPatchInsert(listenerFunc, prevDOMNode, prev, new))

		if typedNode, ok := new.(*nodes.HTMLNode); ok {
			currentStyles = append(currentStyles, typedNode.Styles...)
		}

		return patches, currentStyles, nil
	} else if reflect.TypeOf(old) != reflect.TypeOf(new) {
		fmt.Println("Types were different, replacing")
		// If both nodes exist, but they are of a different type, replace and patch
		patches = append(patches, newPatchReplace(listenerFunc, prevDOMNode, prev, old, new))

		if typedNode, ok := new.(*nodes.HTMLNode); ok {
			currentStyles = append(currentStyles, typedNode.Styles...)
		}

		if typedNode, ok := old.(*nodes.FuncNode); ok {
			// If we hit a function node as the old when there's a need to replace, then we should
			// trigger an unmount on the old node and not render. We don't care about the old node here
			// as we should never rerender it.
			fmt.Println("Types were different and old node is a component, keep going on the new children")
			context.UnregisterAllComponentContexts(typedNode)
			context.RegisterComponentContext("unmount", typedNode)
		}

		return patches, currentStyles, nil
	} else if old.Diff(new) {
		fmt.Println("Nodes were different, updating")
		// If both nodes have the same type, but have differences
		switch typedNode := old.(type) {
		case *nodes.FuncNode:
			// If we hit a function node for both nodes, and they are different, then we should render the
			// new node and assign its result as the result of the old node. We can then keep going on
			// both children. No need to update the function nodes themselves
			oldChildren = nodes.Children{typedNode.RenderResult}
			newConverted := new.(*nodes.FuncNode)

			// Registering with old node so we can keep the references of the current
			// tree alive. Otherwise, the context will track the wrong nodes.
			context.RegisterComponent(typedNode)
			context.RegisterComponentContext("render", typedNode)
			newChildren = nodes.Children{newConverted.Clone().Render(context.CurrentContext)}
		case *nodes.FragmentNode:
			// If we hit a function node for both nodes, and they are different, then we should render the
			// new node and assign its result as the result of the old node. We can then keep going on
			// both children. No need to update the function nodes themselves
			oldChildren = typedNode.Children
			newConverted := new.(*nodes.FragmentNode)
			newChildren = newConverted.Children
		case *nodes.HTMLNode:
			newConverted := new.(*nodes.HTMLNode)
			if typedNode.Tag != newConverted.Tag {
				// If the tags are different, this is not a diff, this is a replace
				patches = append(patches, newPatchReplace(listenerFunc, prevDOMNode, prev, old, new))
				currentStyles = append(currentStyles, newConverted.Styles...)
			} else {
				patches = append(patches, newPatchHTML(listenerFunc, typedNode, new.(*nodes.HTMLNode)))
				oldChildren = typedNode.Children
				newChildren = newConverted.Children

				currentStyles = append(currentStyles, new.(*nodes.HTMLNode).Styles...)
				prevDOMNode = typedNode.DomNode
			}
		case *nodes.TextNode:
			patches = append(patches, newPatchText(prev, typedNode, new.(*nodes.TextNode).Text))
		default:
			return nil, []string{}, fmt.Errorf("somehow got neither a text, nor a HTML node during patching, cannot process node")
		}
	} else {
		fmt.Println("No changes")
		switch oldConverted := old.(type) {
		case *nodes.FuncNode:
			// For function nodes, use the previous render result as the old
			// children and update with the new children. Even if they are the same,
			// they may render differently due to state changes.
			oldChildren = nodes.Children{oldConverted.RenderResult}
			newConverted := new.(*nodes.FuncNode)

			// Registering with old node so we can keep the references of the current
			// tree alive. Otherwise, the context will track the wrong nodes.
			context.RegisterComponent(oldConverted)
			context.RegisterComponentContext("render", oldConverted)
			newChildren = nodes.Children{newConverted.Clone().Render(context.CurrentContext)}
		case *nodes.FragmentNode:
			oldChildren = oldConverted.Children
			newConverted := new.(*nodes.FragmentNode)
			newChildren = newConverted.Children
		case *nodes.HTMLNode:
			oldChildren = oldConverted.Children
			currentStyles = append(currentStyles, oldConverted.Styles...)
			patches = append(patches, newPatchListeners(listenerFunc, oldConverted))

			newConverted := new.(*nodes.HTMLNode)
			newChildren = newConverted.Children
			prevDOMNode = oldConverted.DomNode
		}
	}

	// Start by running through the old children and patch individually
	count := 0
	for _, child := range oldChildren {
		var newChild nodes.Node
		if count < len(newChildren) {
			newChild = newChildren[count]
		}

		childPatches, styles, err := GeneratePatches(listenerFunc, old, prevDOMNode, child, newChild)
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
		childPatches, styles, err := GeneratePatches(listenerFunc, old, prevDOMNode, nil, child)
		if err != nil {
			return nil, []string{}, err
		}
		patches = append(patches, childPatches...)
		currentStyles = append(currentStyles, styles...)
	}

	return patches, currentStyles, nil
}
