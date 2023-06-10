package go_wasm_dom

import (
	"fmt"
)

type screen struct {
	// All node values, fetch in this slice from unique IDs.
	allNodes []*Value

	nodesPerTag       map[string][]int
	nodesPerAttribute map[string][]int
	nodesPerID        map[string][]int
	nodesPerClass     map[string][]int

	// Assigns the node ids to their parents. Keys are the nodes, and the values
	// is a list of parents or ancestors. This means that, to move a node, we need to redefine
	// the value of its parents and ancestors.
	nodePerParent    map[int]int
	nodePerAncestors map[int][]int
}

var (
	uniqueID      = 0
	currentScreen *screen
	window        Value
	document      Value
)

func resetScreen() {
	document = Value{
		id:             uniqueID,
		referencedType: valueDocument,
		properties: map[string]*Value{
			"createElement": createFuncValue(func(this Value, args []Value) any {
				return createElement(args...)
			}),
			"createElementNS": createFuncValue(func(this Value, args []Value) any {
				return createElement(args...)
			}),
		},
	}

	window = Value{
		id:             uniqueID,
		referencedType: valueGlobal,
		properties: map[string]*Value{
			"DEBUG":    valuePtr(ValueOf(false)),
			"document": &document,
		},
	}

	uniqueID++

	currentScreen = &screen{
		allNodes:          []*Value{},
		nodesPerTag:       map[string][]int{},
		nodesPerAttribute: map[string][]int{},
		nodesPerID:        map[string][]int{},
		nodesPerClass:     map[string][]int{},
		nodePerParent:     map[int]int{},
		nodePerAncestors:  map[int][]int{},
	}

	html := createElement(makeString("html"))
	body := createElement(makeString("body"))
	head := createElement(makeString("head"))

	assignElementToParent(html, head, -1)
	assignElementToParent(html, body, -1)
}

func createElement(args ...Value) Value {
	if currentScreen == nil {
		t.Fatal("createElement tried to create an element, but tests were not initialized")
	}

	if len(args) != 1 {
		t.Fatalf("createElement expects only a single argument, %d given", len(args))
	}

	tag := args[0].String()

	element := Value{
		id:             uniqueID,
		referencedType: domNode,
		tag:            tag,
		attributes:     map[string]string{},
		properties: map[string]*Value{
			"firstChild":      valuePtr(Null()),
			"lastChild":       valuePtr(Null()),
			"nextSibling":     valuePtr(Null()),
			"previousSibling": valuePtr(Null()),
			"childNodes":      valuePtr(ValueOf([]Value{})),
			"parentElement":   valuePtr(Null()),

			// TODO: Implement the other methods from the Node API https://developer.mozilla.org/en-US/docs/Web/API/Node
			"appendChild": createFuncValue(func(this Value, args []Value) any {
				return appendChild(this, args...)
			}),
			"insertBefore": createFuncValue(func(this Value, args []Value) any {
				return insertBefore(this, args...)
			}),
			"removeChild": createFuncValue(func(this Value, args []Value) any {
				return removeChild(this, args...)
			}),
			"replaceChild": createFuncValue(func(this Value, args []Value) any {
				return replaceChild(this, args...)
			}),
			"hasAttribute": createFuncValue(func(this Value, args []Value) any {
				return hasAttribute(this, args...)
			}),
			"getAttribute": createFuncValue(func(this Value, args []Value) any {
				return getAttribute(this, args...)
			}),
			"setAttribute": createFuncValue(func(this Value, args []Value) any {
				return setAttribute(this, args...)
			}),
			"querySelector": createFuncValue(func(this Value, args []Value) any {
				return querySelector(this, args...)
			}),
			"querySelectorAll": createFuncValue(func(this Value, args []Value) any {
				return ValueOf(querySelectorAll(this, args...))
			}),
		},
		// TODO: implement add and remove event listeners
		internals: map[string]any{},
		listeners: map[string]Func{},
	}
	element.properties["classList"] = createDOMTokenList(element)

	uniqueID++

	currentScreen.nodesPerTag["tag"] = append(currentScreen.nodesPerTag["tag"], element.id)
	currentScreen.allNodes = append(currentScreen.allNodes, &element)

	return element
}

func getElementPointer(element Value) *Value {
	for _, node := range currentScreen.allNodes {
		if node.id == element.id {
			return node
		}
	}

	return nil
}

func updateSiblingProperties(parent Value) {
	childNodes := parent.properties["childNodes"]

	var previousSibling *Value
	for i := 0; i < childNodes.Length(); i++ {
		child := getElementPointer(childNodes.Index(i))

		if previousSibling == nil {
			child.properties["previousSibling"] = valuePtr(Null())
		} else {
			child.properties["previousSibling"] = previousSibling
		}
		child.properties["nextSibling"] = valuePtr(Null())

		if previousSibling != nil {
			previousSibling.properties["nextSibling"] = child
		}

		childNodes.SetIndex(i, *child)
		previousSibling = child
	}
}

func setParentAndAncestors(parent, child Value) {
	currentScreen.nodePerParent[child.id] = parent.id

	ancestors, ok := currentScreen.nodePerAncestors[parent.id]
	if ok {
		currentScreen.nodePerAncestors[child.id] = append(ancestors, parent.id)
	}
}

func assignElementToParent(parent, child Value, at int) {
	setParentAndAncestors(parent, child)

	parentPointer := getElementPointer(parent)
	childPointer := getElementPointer(child)

	childNodes := parentPointer.properties["childNodes"]
	index := childNodes.Length()
	if at >= 0 && at < index {
		index = at
	}

	newChilds := make([]Value, childNodes.Length()+1)
	count := 0
	for i := 0; i < len(newChilds); i++ {
		if i == at {
			newChilds[i] = child
			continue
		}
		newChilds[i] = childNodes.Index(count)
		count++
	}
	parentPointer.properties["childNodes"] = valuePtr(ValueOf(newChilds))
	parentPointer.properties["firstChild"] = valuePtr(ValueOf(newChilds[0]))
	parentPointer.properties["lastChild"] = valuePtr(ValueOf(newChilds[len(newChilds)-1]))

	childPointer.properties["parentElement"] = parentPointer

	updateSiblingProperties(*parentPointer)
}

func getElementIndex(parent, child Value) int {
	childNodes := parent.properties["childNodes"]
	for i := 0; i < childNodes.Length(); i++ {
		if childNodes.Index(i).Equal(child) {
			return i
		}
	}

	return -1
}

func removeElementFromParent(parent, child Value) {
	delete(currentScreen.nodePerParent, child.id)
	delete(currentScreen.nodePerAncestors, child.id)

	parentPointer := getElementPointer(parent)
	childNodes := parentPointer.properties["childNodes"]

	newChilds := make([]Value, childNodes.Length()-1)
	count := 0
	for i := 0; i < len(newChilds); i++ {
		if !childNodes.Index(i).Equal(child) {
			newChilds[i] = childNodes.Index(count)
		}
		count++
	}

	parentPointer.properties["childNodes"] = valuePtr(ValueOf(newChilds))
	updateSiblingProperties(parent)
}

func replaceElementInParent(parent, child, newChild Value) {
	setParentAndAncestors(parent, child)

	childNodes := parent.properties["childNodes"]

	newChilds := make([]Value, childNodes.Length())
	for i := 0; i < len(newChilds); i++ {
		if !childNodes.Index(i).Equal(child) {
			newChilds[i] = newChild
		} else {
			newChilds[i] = childNodes.Index(i)
		}
	}

	parent.properties["childNodes"] = valuePtr(ValueOf(newChilds))
	updateSiblingProperties(parent)
}

func appendChild(caller Value, args ...Value) Value {
	if currentScreen == nil {
		t.Fatal("appendChild tried to execute, but tests were not initialized")
	}

	if len(args) != 1 {
		t.Fatalf("appendChild expects 1 argument, %d given", len(args))
	}

	element := args[0]
	if element.referencedType != domNode {
		t.Fatalf("appendChild expects argument[0] to be a DOM element, %T given", args[0])
	}

	assignElementToParent(caller, element, -1)

	return element
}

func insertBefore(caller Value, args ...Value) Value {
	if currentScreen == nil {
		t.Fatal("insertBefore tried to execute, but tests were not initialized")
	}

	if len(args) != 2 {
		t.Fatalf("insertBefore expects 2 argument, %d given", len(args))
	}

	element := args[0]
	if element.referencedType != domNode {
		t.Fatalf("insertBefore expects argument[0] to be a DOM element, %T given", args[0])
	}

	previousSibling := args[1]
	if element.referencedType != domNode {
		t.Fatalf("insertBefore expects argument[1] to be a DOM element, %T given", args[1])
	}

	assignElementToParent(caller, element, getElementIndex(caller, previousSibling)-1)

	return element
}

func removeChild(caller Value, args ...Value) Value {
	if currentScreen == nil {
		t.Fatal("removeChild tried to execute, but tests were not initialized")
	}

	if len(args) != 1 {
		t.Fatalf("removeChild expects 1 argument, %d given", len(args))
	}

	element := args[0]
	if element.referencedType != domNode {
		t.Fatalf("insertBefore expects argument[0] to be a DOM element, %T given", args[0])
	}

	removeElementFromParent(caller, element)

	return element
}

func replaceChild(caller Value, args ...Value) Value {
	if currentScreen == nil {
		t.Fatal("replaceChild tried to execute, but tests were not initialized")
	}

	if len(args) != 2 {
		t.Fatalf("replaceChild expects 2 argument, %d given", len(args))
	}

	element := args[0]
	if element.referencedType != domNode {
		t.Fatalf("replaceChild expects argument[0] to be a DOM element, %T given", args[0])
	}

	previousElement := args[1]
	if element.referencedType != domNode {
		t.Fatalf("insertBefore expects argument[1] to be a DOM element, %T given", args[1])
	}

	replaceElementInParent(caller, previousElement, element)

	return element
}

func hasAttribute(caller Value, args ...Value) Value {
	if currentScreen == nil {
		t.Fatal("hasAttribute tried to execute, but tests were not initialized")
	}

	if len(args) != 1 {
		t.Fatalf("hasAttribute expects 1 argument, %d given", len(args))
	}

	key := args[0]
	if key.referencedType != stringType {
		t.Fatalf("hasAttribute expects argument[0] to be a string, %T given", args[0])
	}

	if _, ok := caller.attributes[goValueOf(key).(string)]; !ok {
		return Value{
			referencedType: valueFalse,
		}
	}

	return Value{
		referencedType: valueTrue,
	}
}

func getAttribute(caller Value, args ...Value) Value {
	if currentScreen == nil {
		t.Fatal("getAttribute tried to execute, but tests were not initialized")
	}

	if len(args) != 1 {
		t.Fatalf("setAttribute expects 1 argument, %d given", len(args))
	}

	key := args[0]
	if key.referencedType != stringType {
		t.Fatalf("getAttribute expects argument[0] to be a string, %T given", args[0])
	}

	if val, ok := caller.attributes[goValueOf(key).(string)]; ok {
		return ValueOf(val)
	}

	return Null()
}

func setAttribute(caller Value, args ...Value) Value {
	if currentScreen == nil {
		t.Fatal("setAttribute tried to execute, but tests were not initialized")
	}

	if len(args) != 2 {
		t.Fatalf("setAttribute expects 2 argument, %d given", len(args))
	}

	key := args[0]
	if key.referencedType != stringType {
		t.Fatalf("setAttribute expects argument[0] to be a string, %T given", args[0])
	}

	value := args[0]
	if key.referencedType != stringType {
		t.Fatalf("setAttribute expects argument[1] to be a string, %T given", args[0])
	}

	callerPointer := getElementPointer(caller)
	callerPointer.attributes[goValueOf(key).(string)] = goValueOf(value).(string)

	// Assign the element under its attributes with either double or single quote
	// we keep each twice so users can use either in the selectors.
	attributeWithDoubleQuotes := fmt.Sprintf(
		"%s=\"%s\"", goValueOf(key).(string), key.String(),
	)
	currentScreen.nodesPerAttribute[attributeWithDoubleQuotes] = append(
		currentScreen.nodesPerAttribute[attributeWithDoubleQuotes],
		caller.id,
	)

	attributeWithSingleQuotes := fmt.Sprintf(
		"%s='%s'", goValueOf(key).(string), key.String(),
	)
	currentScreen.nodesPerAttribute[attributeWithSingleQuotes] = append(
		currentScreen.nodesPerAttribute[attributeWithSingleQuotes],
		caller.id,
	)

	return Undefined()
}
