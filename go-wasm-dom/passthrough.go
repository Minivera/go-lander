//go:build js && wasm

package go_wasm_dom

import (
	"fmt"
	"math"
	"strconv"
	realJs "syscall/js"
	"unsafe"
)

var isInFakeMode = false

type valueType uint64

const (
	valueUndefined valueType = iota
	valueNaN
	valueZero
	valueNull
	valueTrue
	valueFalse
	valueGlobal
	valueDocument

	objectConstructor
	arrayConstructor
	functionType
	numberType
	stringType

	domNode
)

// Value represents a JavaScript value. The zero value is the JavaScript value "undefined".
// Values can be checked for equality with the Equal method.
type Value struct {
	realJs.Value

	id             int
	referencedType valueType

	tag        string
	attributes map[string]string
	properties map[string]*Value
	internals  map[string]any
	listeners  map[string]Func
}

// Undefined returns the JavaScript value "undefined".
func Undefined() Value {
	if isInFakeMode {
		return Value{referencedType: valueUndefined}
	}
	return Value{Value: realJs.Undefined()}
}

// Null returns the JavaScript value "null".
func Null() Value {
	if isInFakeMode {
		return Value{referencedType: valueNull}
	}
	return Value{Value: realJs.Null()}
}

// Global returns the JavaScript global object, usually "window" or "global".
func Global() Value {
	if isInFakeMode {
		return window
	}
	return Value{Value: realJs.Global()}
}

func makeFloat(x float64) Value {
	if x == 0 {
		return Value{
			referencedType: valueZero,
		}
	}
	if x != x {
		return Value{
			referencedType: valueNaN,
		}
	}

	return Value{
		referencedType: numberType,
		internals: map[string]any{
			"floatValue": x,
		},
	}
}

func makeString(x string) Value {
	return Value{
		referencedType: stringType,
		internals: map[string]any{
			"stringValue": x,
		},
	}
}

// ValueOf returns x as a JavaScript value:
//
//	| Go                     | JavaScript             |
//	| ---------------------- | ---------------------- |
//	| js.Value               | [its value]            |
//	| js.Func                | function               |
//	| nil                    | null                   |
//	| bool                   | boolean                |
//	| integers and floats    | number                 |
//	| string                 | string                 |
//	| []interface{}          | new array              |
//	| map[string]interface{} | new object             |
//
// Panics if x is not one of the expected types.
func ValueOf(x any) Value {
	if isInFakeMode {
		switch x := x.(type) {
		case Value:
			return x
		case *Value:
			return *x
		case Func:
			return x.internalValue
		case *Func:
			return x.internalValue
		case nil:
			return Null()
		case bool:
			if x {
				return Value{
					referencedType: valueTrue,
				}
			} else {
				return Value{
					referencedType: valueFalse,
				}
			}
		case int:
			return makeFloat(float64(x))
		case int8:
			return makeFloat(float64(x))
		case int16:
			return makeFloat(float64(x))
		case int32:
			return makeFloat(float64(x))
		case int64:
			return makeFloat(float64(x))
		case uint:
			return makeFloat(float64(x))
		case uint8:
			return makeFloat(float64(x))
		case uint16:
			return makeFloat(float64(x))
		case uint32:
			return makeFloat(float64(x))
		case uint64:
			return makeFloat(float64(x))
		case uintptr:
			return makeFloat(float64(x))
		case unsafe.Pointer:
			return makeFloat(float64(uintptr(x)))
		case float32:
			return makeFloat(float64(x))
		case float64:
			return makeFloat(x)
		case string:
			return makeString(x)
		case []any:
			a := Value{
				referencedType: arrayConstructor,
				properties:     map[string]*Value{},
			}
			for i, s := range x {
				val := ValueOf(s)
				a.properties[fmt.Sprintf("%d", i)] = &val
			}
			return a
		case map[string]any:
			o := Value{
				referencedType: objectConstructor,
				properties:     map[string]*Value{},
			}
			for k, v := range x {
				val := ValueOf(v)
				o.properties[fmt.Sprintf("%v", k)] = &val
			}
			return o
		default:
			t.Fatal("ValueOf: invalid value")
		}
	}
	return Value{Value: realJs.ValueOf(x)}
}

func goValueOf(x Value) any {
	if isInFakeMode {
		switch x.referencedType {
		case valueUndefined, valueNull:
			return nil
		case valueTrue:
			return true
		case valueFalse:
			return false
		case valueZero:
			return 0
		case valueNaN:
			return math.NaN()
		case numberType:
			return x.internals["floatValue"].(float64)
		case stringType:
			return x.internals["stringValue"].(string)
		case arrayConstructor:
			s := make([]any, len(x.properties))
			for k, v := range x.properties {
				index, err := strconv.Atoi(k)
				if err != nil {
					t.Fatal("Conversion of JS array to Go slice, indexes were not numeric")
				}
				s[index] = goValueOf(*v)
			}
			return s
		case objectConstructor:
			s := make(map[string]any, len(x.properties))
			for k, v := range x.properties {
				s[k] = goValueOf(*v)
			}
		default:
			t.Fatal("Conversion of un-convertable type to Go type, no-op")
		}
	}
	return Value{Value: realJs.ValueOf(x)}
}

func valuePtr(v Value) *Value {
	return &v
}
