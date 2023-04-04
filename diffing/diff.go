package diffing

import (
	"fmt"
	"reflect"

	"github.com/minivera/go-lander/nodes"
)

func GeneratePatches(prev, old, new nodes.Node) ([]Patch, error) {
	var patches []Patch

	// If the new is missing, then nothing should happen. Removing unneeded children
	// is part of the update patch.
	if new == nil {
		return patches, nil
	}

	// If the old node is missing, then we are mounting for the first time
	if old == nil {
		patches = append(patches, newPatchInsert(prev, new))
		return patches, nil
	}

	// If both nodes don't have the same type, we should replace and mount the new node
	if reflect.TypeOf(old) != reflect.TypeOf(new) {
		patches = append(patches, newPatchReplace(prev, old, new))
		return patches, nil
	}

	// If both nodes have the same type, but have differences
	if old.Diff(new) {
		if val, ok := new.(*nodes.TextNode); ok {
			patches = append(patches, newPatchText(prev, old, val.Text))
			return patches, nil
		} else if _, ok := new.(*nodes.HTMLNode); ok {
			patches = append(patches, newPatchHTML(prev, old, new))
		}
	}

	converted, ok := old.(*nodes.HTMLNode)
	if !ok {
		return nil, fmt.Errorf("somehow got neither a text nor a HTML node during patching, cannot process node")
	}
	newConverted, ok := new.(*nodes.HTMLNode)

	// Start by running through the old children and patch individually
	count := 0
	for _, child := range converted.Children {
		var newChild nodes.Node
		if count < len(newConverted.Children) {
			newChild = newConverted.Children[count]
		}

		childPatches, err := GeneratePatches(old, child, newChild)
		if err != nil {
			return nil, err
		}
		patches = append(patches, childPatches...)

		count += 1
	}

	// If we still have new nodes left, then loop over them and insert
	if count >= len(newConverted.Children) {
		return patches, nil
	}

	for _, child := range newConverted.Children {
		patches = append(patches, newPatchInsert(old, child))
	}

	return patches, nil
}
