package go_wasm_dom

type screen struct {
	// All node values, fetch in this slice from unique IDs.
	allNodes []Value

	nodePerTag       map[string]int
	nodePerAttribute map[string]int

	// Assigns the node ids to their parents. Keys are the nodes, values are their parents.
	// A node can have multiple children, but a node can only have one parent. This makes sure
	// that assigning a new parent moves the node.
	nodePerParent   map[int]int
	nodePerAncestor map[int]int
}

var (
	uniqueID      = 0
	currentScreen *screen
	window        Value
)

func resetScreen() {
	window = Value{
		id:             uniqueID,
		referencedType: valueGlobal,
		properties: map[string]Value{
			"DEBUG": ValueOf(false),
		},
	}

	uniqueID++

	currentScreen = &screen{
		allNodes:         []Value{},
		nodePerTag:       map[string]int{},
		nodePerAttribute: map[string]int{},
		nodePerParent:    map[int]int{},
	}

	document := createElement("document")
	body := createElement("body")
	head := createElement("head")

	assignElementToParent(document, head, -1)
	assignElementToParent(document, body, -1)
}

func createElement(args ...any) Value {
	if currentScreen == nil {
		t.Fatal("createElement tried to create an element, but tests were not initialized")
	}

	if len(args) != 1 {
		t.Fatalf("createElement expects only a single argument, %d given", len(args))
	}

	tag, ok := args[0].(string)
	if !ok {
		t.Fatalf("createElement expects argument[0] to be a string, %T given", args[0])
	}

	element := Value{
		id:             uniqueID,
		referencedType: domNode,
		tag:            tag,
		attributes:     map[string]string{},
		properties: map[string]Value{
			"firstChild":      Null(),
			"lastChild":       Null(),
			"nextSibling":     Null(),
			"previousSibling": Null(),
			"childNodes":      ValueOf([]Value{}),
			"parentElement":   Null(),
		},
		internals: map[string]any{},
		listeners: map[string]Func{},
	}
	uniqueID++

	currentScreen.nodePerTag["tag"] = element.id

	return element
}

func updateSiblingProperties(parent Value) {
	childNodes := parent.properties["childNodes"]

	previousSibling := Null()
	for i := 0; i < childNodes.Length(); i++ {
		child := childNodes.Index(i)

		child.properties["previousSibling"] = previousSibling
		child.properties["nextSibling"] = Null()

		if !previousSibling.IsNull() {
			previousSibling.properties["nextSibling"] = child
		}

		childNodes.SetIndex(i, child)
		previousSibling = child
	}
}

func assignElementToParent(parent, child Value, at int) {
	currentScreen.nodePerParent[child.id] = parent.id

	childNodes := parent.properties["childNodes"]
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
	parent.properties["childNodes"] = ValueOf(newChilds)
	parent.properties["firstChild"] = ValueOf(newChilds[0])
	parent.properties["lastChild"] = ValueOf(newChilds[len(newChilds)-1])

	child.properties["parentElement"] = parent

	updateSiblingProperties(parent)
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

	childNodes := parent.properties["childNodes"]

	newChilds := make([]Value, childNodes.Length()-1)
	count := 0
	for i := 0; i < len(newChilds); i++ {
		if !childNodes.Index(i).Equal(child) {
			newChilds[i] = childNodes.Index(count)
		}
		count++
	}

	parent.properties["childNodes"] = ValueOf(newChilds)
	updateSiblingProperties(parent)
}

func replaceElementInParent(parent, child, newChild Value) {
	currentScreen.nodePerParent[child.id] = parent.id

	childNodes := parent.properties["childNodes"]

	newChilds := make([]Value, childNodes.Length())
	for i := 0; i < len(newChilds); i++ {
		if !childNodes.Index(i).Equal(child) {
			newChilds[i] = newChild
		} else {
			newChilds[i] = childNodes.Index(i)
		}
	}

	parent.properties["childNodes"] = ValueOf(newChilds)
	updateSiblingProperties(parent)
}

func appendChild(caller Value, args ...any) Value {
	if currentScreen == nil {
		t.Fatal("appendChild tried to execute, but tests were not initialized")
	}

	if len(args) != 1 {
		t.Fatalf("appendChild expects 1 argument, %d given", len(args))
	}

	element, ok := args[0].(Value)
	if !ok && element.referencedType != domNode {
		t.Fatalf("appendChild expects argument[0] to be a DOM element, %T given", args[0])
	}

	assignElementToParent(caller, element, -1)

	return element
}

func insertBefore(caller Value, args ...any) Value {
	if currentScreen == nil {
		t.Fatal("insertBefore tried to execute, but tests were not initialized")
	}

	if len(args) != 2 {
		t.Fatalf("insertBefore expects 2 argument, %d given", len(args))
	}

	element, ok := args[0].(Value)
	if !ok && element.referencedType != domNode {
		t.Fatalf("insertBefore expects argument[0] to be a DOM element, %T given", args[0])
	}

	previousSibling, ok := args[1].(Value)
	if !ok && element.referencedType != domNode {
		t.Fatalf("insertBefore expects argument[1] to be a DOM element, %T given", args[1])
	}

	assignElementToParent(caller, element, getElementIndex(caller, previousSibling)-1)

	return element
}

func removeChild(caller Value, args ...any) Value {
	if currentScreen == nil {
		t.Fatal("removeChild tried to execute, but tests were not initialized")
	}

	if len(args) != 1 {
		t.Fatalf("removeChild expects 1 argument, %d given", len(args))
	}

	element, ok := args[0].(Value)
	if !ok && element.referencedType != domNode {
		t.Fatalf("insertBefore expects argument[0] to be a DOM element, %T given", args[0])
	}

	removeElementFromParent(caller, element)

	return element
}

func replaceChild(caller Value, args ...any) Value {
	if currentScreen == nil {
		t.Fatal("replaceChild tried to execute, but tests were not initialized")
	}

	if len(args) != 2 {
		t.Fatalf("replaceChild expects 2 argument, %d given", len(args))
	}

	element, ok := args[0].(Value)
	if !ok && element.referencedType != domNode {
		t.Fatalf("replaceChild expects argument[0] to be a DOM element, %T given", args[0])
	}

	previousElement, ok := args[1].(Value)
	if !ok && element.referencedType != domNode {
		t.Fatalf("insertBefore expects argument[1] to be a DOM element, %T given", args[1])
	}

	replaceElementInParent(caller, previousElement, element)

	return element
}

func hasAttribute(caller Value, args ...any) Value {
	if currentScreen == nil {
		t.Fatal("hasAttribute tried to execute, but tests were not initialized")
	}

	if len(args) != 1 {
		t.Fatalf("hasAttribute expects 1 argument, %d given", len(args))
	}

	key, ok := args[0].(Value)
	if !ok && key.referencedType != stringType {
		t.Fatalf("hasAttribute expects argument[0] to be a string, %T given", args[0])
	}

	if _, ok = caller.attributes[goValueOf(key).(string)]; !ok {
		return Value{
			referencedType: valueFalse,
		}
	}

	return Value{
		referencedType: valueTrue,
	}
}

func getAttribute(caller Value, args ...any) Value {
	if currentScreen == nil {
		t.Fatal("getAttribute tried to execute, but tests were not initialized")
	}

	if len(args) != 1 {
		t.Fatalf("setAttribute expects 1 argument, %d given", len(args))
	}

	key, ok := args[0].(Value)
	if !ok && key.referencedType != stringType {
		t.Fatalf("getAttribute expects argument[0] to be a string, %T given", args[0])
	}

	if val, ok := caller.attributes[goValueOf(key).(string)]; ok {
		return ValueOf(val)
	}

	return Null()
}

func setAttribute(caller Value, args ...any) Value {
	if currentScreen == nil {
		t.Fatal("setAttribute tried to execute, but tests were not initialized")
	}

	if len(args) != 2 {
		t.Fatalf("setAttribute expects 2 argument, %d given", len(args))
	}

	key, ok := args[0].(Value)
	if !ok && key.referencedType != stringType {
		t.Fatalf("setAttribute expects argument[0] to be a string, %T given", args[0])
	}

	value, ok := args[0].(Value)
	if !ok && key.referencedType != stringType {
		t.Fatalf("setAttribute expects argument[1] to be a string, %T given", args[0])
	}

	caller.attributes[goValueOf(key).(string)] = goValueOf(value).(string)

	return Undefined()
}
