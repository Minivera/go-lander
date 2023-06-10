package go_wasm_dom

func createDOMTokenList(owner Value) *Value {
	list := &Value{
		referencedType: objectConstructor,
		properties:     map[string]*Value{},
		internals: map[string]any{
			"elements": []string{},
			"owner":    owner,
		},
	}

	list.properties["item"] = createFuncValue(func(_ Value, args []Value) any {
		if len(args) != 1 {
			t.Fatalf("DOMTokenList.item expects 1 argument, %d given", len(args))
		}

		index := args[0]
		if index.referencedType != numberType {
			t.Fatalf("DOMTokenList.item expects argument[0] to be a number, %T given", args[0])
		}

		return ValueOf(list.internals["elements"].([]string)[index.Int()])
	})

	list.properties["contains"] = createFuncValue(func(_ Value, args []Value) any {
		if len(args) != 1 {
			t.Fatalf("DOMTokenList.contains expects 1 argument, %d given", len(args))
		}

		value := args[0]
		if value.referencedType != stringType {
			t.Fatalf("DOMTokenList.contains expects argument[0] to be a string, %T given", args[0])
		}

		val := value.String()
		for _, el := range list.internals["elements"].([]string) {
			if el == val {
				return ValueOf(true)
			}
		}

		return ValueOf(false)
	})

	list.properties["add"] = createFuncValue(func(_ Value, args []Value) any {
		for i, value := range args {
			if value.referencedType != stringType {
				t.Fatalf("DOMTokenList.add expects argument[%d] to be a string, %T given", i, args[i])
			}

			val := value.String()
			list.internals["elements"] = append(list.internals["elements"].([]string), val)

			// Add the owner of this list to the classes list for nodes
			element := getElementPointer(list.internals["owner"].(Value))
			currentScreen.nodesPerClass[val] = append(currentScreen.nodesPerClass[val], element.id)
		}

		return Undefined()
	})

	list.properties["remove"] = createFuncValue(func(_ Value, args []Value) any {
		for i, value := range args {
			if value.referencedType != stringType {
				t.Fatalf("DOMTokenList.remove expects argument[%d] to be a string, %T given", i, args[i])
			}

			val := value.String()
			newList := make([]string, len(list.internals["elements"].([]string))-1)
			count := 0
			for _, el := range list.internals["elements"].([]string) {
				if el == val {
					continue
				}
				newList[count] = el
				count++
			}

			list.internals["elements"] = newList

			element := getElementPointer(list.internals["owner"].(Value))

			// Remove the owner of this list to the classes list for nodes
			newNodes := make([]int, len(currentScreen.nodesPerClass[val])-1)
			count = 0
			for _, el := range currentScreen.nodesPerClass[val] {
				if el == element.id {
					continue
				}
				newNodes[count] = el
				count++
			}

			currentScreen.nodesPerClass[val] = newNodes
		}

		return Undefined()
	})

	list.properties["replace"] = createFuncValue(func(_ Value, args []Value) any {
		if len(args) != 2 {
			t.Fatalf("DOMTokenList.replace expects 2 argument, %d given", len(args))
		}

		value := args[0]
		if value.referencedType != stringType {
			t.Fatalf("DOMTokenList.replace expects argument[0] to be a string, %T given", args[0])
		}

		newValue := args[1]
		if newValue.referencedType != stringType {
			t.Fatalf("DOMTokenList.replace expects argument[1] to be a string, %T given", args[1])
		}

		replaced := false
		val := value.String()
		newList := make([]string, 0, len(list.internals["elements"].([]string))-1)
		for _, el := range list.internals["elements"].([]string) {
			if el == val {
				newList = append(newList, newValue.String())
				replaced = true
				continue
			}
			newList = append(newList, el)
		}

		list.internals["elements"] = newList

		return ValueOf(replaced)
	})

	list.properties["toggle"] = createFuncValue(func(_ Value, args []Value) any {
		if len(args) > 2 || len(args) < 1 {
			t.Fatalf("DOMTokenList.toggle expects 1 or 2 argument, %d given", len(args))
		}

		value := args[0]
		if value.referencedType != stringType {
			t.Fatalf("DOMTokenList.replace expects argument[0] to be a string, %T given", args[0])
		}

		onlyAdd := false
		onlyRemove := false
		if len(args) == 2 {
			if args[1].Bool() {
				onlyAdd = true
			} else {
				onlyRemove = true
			}
		}

		if list.Call("contains", value).Bool() && !onlyAdd {
			list.Call("remove", value)
			return ValueOf(false)
		}
		if !list.Call("contains", value).Bool() && !onlyRemove {
			list.Call("add", value)
			return ValueOf(true)
		}

		return ValueOf(false)
	})

	// TODO
	list.properties["entries"] = createFuncValue(func(_ Value, args []Value) any {
		return Undefined()
	})
	list.properties["forEach"] = createFuncValue(func(_ Value, args []Value) any {
		return Undefined()
	})
	list.properties["keys"] = createFuncValue(func(_ Value, args []Value) any {
		return Undefined()
	})
	list.properties["values"] = createFuncValue(func(_ Value, args []Value) any {
		return Undefined()
	})

	return list
}
