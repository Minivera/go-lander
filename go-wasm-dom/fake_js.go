//go:build js && wasm

package go_wasm_dom

import (
	"fmt"
	"strconv"
	"strings"
	realJs "syscall/js"
)

func convertArg(arg any) any {
	// When using args that can be of type Value, we must make sure to convert them back
	// to realJs values, otherwise ValueOf will throw.
	switch fakeValue := arg.(type) {
	case Value:
		return fakeValue.Value
	case Func:
		return fakeValue.Func
	default:
		return arg
	}
}

func convertArgs(args []any) []any {
	converted := make([]any, len(args))

	for i, arg := range args {
		converted[i] = convertArg(arg)
	}

	return converted
}

// Equal reports whether v and w are equal according to JavaScript's === operator.
func (v Value) Equal(w Value) bool {
	if isInFakeMode {
		switch v.referencedType {
		case valueUndefined:
			return w.referencedType == valueUndefined
		case valueNull:
			return w.referencedType == valueNull
		case valueTrue:
			return w.referencedType == valueTrue
		case valueFalse:
			return w.referencedType == valueFalse
		case valueZero:
			return w.referencedType == valueZero
		case valueNaN:
			return w.referencedType == valueNaN
		case numberType:
			return w.referencedType == numberType && v.internals["floatValue"] == w.internals["floatValue"]
		case stringType:
			return w.referencedType == stringType && v.internals["stringValue"] == w.internals["stringValue"]
		case objectConstructor:
			return w.referencedType == objectConstructor && compareObjectValues(v, w)
		case arrayConstructor:
			return w.referencedType == arrayConstructor && compareObjectValues(v, w)
		case domNode:
			return w.referencedType == domNode && compareDOMValues(v, w)
		case functionType:
			// TODO
			return w.referencedType == functionType
		default:
			t.Fatal("Type: bad value type")
		}
	}
	return v.Value.Equal(w.Value)
}

// IsUndefined reports whether v is the JavaScript value "undefined".
func (v Value) IsUndefined() bool {
	if isInFakeMode {
		return v.referencedType == valueUndefined
	}
	return v.Value.IsUndefined()
}

// IsNull reports whether v is the JavaScript value "null".
func (v Value) IsNull() bool {
	if isInFakeMode {
		return v.referencedType == valueNull
	}
	return v.Value.IsNull()
}

// IsNaN reports whether v is the JavaScript value "NaN".
func (v Value) IsNaN() bool {
	if isInFakeMode {
		return v.referencedType == valueNaN
	}
	return v.Value.IsNaN()
}

// Type returns the JavaScript type of the value v. It is similar to JavaScript's typeof operator,
// except that it returns TypeNull instead of TypeObject for null.
func (v Value) Type() realJs.Type {
	if isInFakeMode {
		switch v.referencedType {
		case valueUndefined:
			return realJs.TypeUndefined
		case valueNull:
			return realJs.TypeNull
		case valueTrue, valueFalse:
			return realJs.TypeBoolean
		case numberType, valueZero, valueNaN:
			return realJs.TypeNumber
		case stringType:
			return realJs.TypeString
		case objectConstructor, arrayConstructor, domNode:
			return realJs.TypeObject
		case functionType:
			return realJs.TypeFunction
		default:
			t.Fatal("Type: bad value type")
		}
	}
	return v.Value.Type()
}

// Get returns the JavaScript property p of value v.
// It panics if v is not a JavaScript object.
func (v Value) Get(p string) Value {
	if isInFakeMode {
		if v.properties == nil || (v.referencedType != objectConstructor && v.referencedType != domNode &&
			v.referencedType != valueDocument && v.referencedType != valueGlobal) {
			t.Fatal("Get: Value is not a JavaScript Object")
		}

		if _, ok := v.properties[p]; !ok {
			t.Fatalf("Get: Value does not have a property called %s", p)
		}

		return ValueOf(v.properties[p])
	}
	return Value{Value: v.Value.Get(p)}
}

// Set sets the JavaScript property p of value v to ValueOf(x).
// It panics if v is not a JavaScript object.
func (v Value) Set(p string, x any) {
	if isInFakeMode {
		if v.properties == nil || (v.referencedType != objectConstructor && v.referencedType != domNode &&
			v.referencedType != valueDocument && v.referencedType != valueGlobal) {
			t.Fatal("Set: Value is not a JavaScript Object")
		}

		val := ValueOf(x)
		v.properties[p] = &val
		if p == "id" {
			// Assign the element to the nodes list per ID if needed
			currentScreen.nodesPerID[v.properties[p].String()] = append(
				currentScreen.nodesPerID[v.properties[p].String()],
				v.id,
			)
		}
		return
	}
	v.Value.Set(p, convertArg(x))
}

// Delete deletes the JavaScript property p of value v.
// It panics if v is not a JavaScript object.
func (v Value) Delete(p string) {
	if isInFakeMode {
		if v.properties == nil || (v.referencedType != objectConstructor && v.referencedType != domNode &&
			v.referencedType != valueDocument && v.referencedType != valueGlobal) {
			t.Fatal("Delete: Value is not a JavaScript Object")
		}

		delete(v.properties, p)
		return
	}
	v.Value.Delete(p)
}

// Index returns JavaScript index i of value v.
// It panics if v is not a JavaScript object.
func (v Value) Index(i int) Value {
	if isInFakeMode {
		if v.properties == nil || v.referencedType != arrayConstructor {
			t.Fatal("Index: Value is not a JavaScript Array")
		}

		if _, ok := v.properties[strconv.Itoa(i)]; !ok {
			t.Fatalf("Index: Value index %d is out of bounds", i)
		}

		return ValueOf(v.properties[strconv.Itoa(i)])
	}
	return Value{Value: v.Value.Index(i)}
}

// SetIndex sets the JavaScript index i of value v to ValueOf(x).
// It panics if v is not a JavaScript object.
func (v Value) SetIndex(i int, x any) {
	if isInFakeMode {
		if v.properties == nil || v.referencedType != arrayConstructor {
			t.Fatal("SetIndex: Value is not a JavaScript Array")
		}

		val := ValueOf(x)
		v.properties[strconv.Itoa(i)] = &val
		return
	}
	v.Value.SetIndex(i, convertArg(x))
}

// Length returns the JavaScript property "length" of v.
// It panics if v is not a JavaScript object.
func (v Value) Length() int {
	if isInFakeMode {
		if v.referencedType == arrayConstructor || v.referencedType == objectConstructor {
			return len(v.properties)
		}
		if v.referencedType == stringType {
			return len(v.internals["stringValue"].(string))
		}
		return 0
	}
	return v.Value.Length()
}

// Call does a JavaScript call to the method m of value v with the given arguments.
// It panics if v has no method m.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Value) Call(m string, args ...any) Value {
	if isInFakeMode {
		method, ok := v.properties["m"]
		if !ok {
			t.Fatalf("Call: Object doesn't have method %s, it might not be implemented yet", m)
		}

		args := convertArgs(args)
		valueArgs := make([]Value, len(args))
		for i, a := range args {
			valueArgs[i] = ValueOf(a)
		}
		return method.Invoke(append([]Value{v}, valueArgs...))
	}
	return Value{Value: v.Value.Call(m, convertArgs(args)...)}
}

// Invoke does a JavaScript call of the value v with the given arguments.
// It panics if v is not a JavaScript function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Value) Invoke(args ...any) Value {
	if isInFakeMode {
		if v.referencedType != functionType {
			t.Fatal("Invoke: Value is not a JavaScript function")
		}

		args := convertArgs(args)
		valueArgs := make([]Value, len(args)-1)
		thisArg := Value{}
		for i, a := range args {
			if i == 0 {
				thisArg = ValueOf(a)
			}
			valueArgs[i] = ValueOf(a)
		}

		return ValueOf(v.internals["fn"].(func(this Value, args []Value) any)(thisArg, valueArgs))
	}
	return Value{Value: v.Value.Invoke(convertArgs(args)...)}
}

// New uses JavaScript's "new" operator with value v as constructor and the given arguments.
// It panics if v is not a JavaScript function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (v Value) New(args ...any) Value {
	if isInFakeMode {
		if v.referencedType != functionType {
			t.Fatal("New: Value is not a JavaScript function")
		}

		args := convertArgs(args)
		valueArgs := make([]Value, len(args))
		for i, a := range args {
			valueArgs[i] = ValueOf(a)
		}

		return ValueOf(v.internals["fn"].(func(this Value, args []Value) any)(Value{}, valueArgs))
	}
	return Value{Value: v.Value.New(convertArgs(args)...)}
}

// Float returns the value v as a float64.
// It panics if v is not a JavaScript number.
func (v Value) Float() float64 {
	if isInFakeMode {
		if v.referencedType != numberType && v.referencedType != valueNaN && v.referencedType != valueZero {
			t.Fatal("Float: Value is not a JavaScript number")
		}

		return goValueOf(v).(float64)
	}
	return v.Value.Float()
}

// Int returns the value v truncated to an int.
// It panics if v is not a JavaScript number.
func (v Value) Int() int {
	if isInFakeMode {
		return int(v.Float())
	}
	return v.Value.Int()
}

// Bool returns the value v as a bool.
// It panics if v is not a JavaScript boolean.
func (v Value) Bool() bool {
	if isInFakeMode {
		if v.referencedType == valueTrue {
			return true
		} else if v.referencedType == valueFalse {
			return false
		} else {
			t.Fatal("Bool: Value is not a JavaScript boolean")
		}
	}
	return v.Value.Bool()
}

// Truthy returns the JavaScript "truthiness" of the value v. In JavaScript,
// false, 0, "", null, undefined, and NaN are "falsy", and everything else is
// "truthy". See https://developer.mozilla.org/en-US/docs/Glossary/Truthy.
func (v Value) Truthy() bool {
	if isInFakeMode {
		return !(v.referencedType == valueFalse ||
			v.referencedType == valueZero ||
			v.String() == "" ||
			v.referencedType == valueNull ||
			v.referencedType == valueUndefined ||
			v.referencedType == valueNaN)
	}
	return v.Value.Truthy()
}

// String returns the value v as a string.
// String is a special case because of Go's String method convention. Unlike the other getters,
// it does not panic if v's Type is not TypeString. Instead, it returns a string of the form "<T>"
// or "<T: V>" where T is v's type and V is a string representation of v's value.
func (v Value) String() string {
	if isInFakeMode {
		switch v.referencedType {
		case valueUndefined:
			return "undefined"
		case valueNull:
			return "null"
		case valueTrue:
			return "true"
		case valueFalse:
			return "false"
		case valueZero:
			return "0"
		case valueNaN:
			return "NaN"
		case numberType:
			return fmt.Sprintf("%f", v.internals["floatValue"].(float64))
		case stringType:
			return v.internals["stringValue"].(string)
		case arrayConstructor:
			content := make([]string, len(v.properties))
			for k, v := range v.properties {
				index, err := strconv.Atoi(k)
				if err != nil {
					t.Fatal("Conversion of JS array to Go slice, indexes were not numeric")
				}
				content[index] = v.String()
			}
			return strings.Join(content, ",")
		case objectConstructor:
			return "[object Object]"
		case functionType:
			// TODO: This would normally pring the content of the string
			return "function() {}"
		default:
			t.Fatal("Conversion of un-convertable type to Go type, no-op")
		}
	}
	return v.Value.String()
}

// InstanceOf reports whether v is an instance of type t according to JavaScript's instanceof operator.
func (v Value) InstanceOf(t Value) bool {
	if isInFakeMode {
		prototypeA, okA := v.internals["prototype"]
		prototypeB, okB := t.internals["prototype"]
		return okA && okB && prototypeA == prototypeB
	}
	return v.Value.InstanceOf(t.Value)
}
