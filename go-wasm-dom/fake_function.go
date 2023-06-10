//go:build js && wasm

package go_wasm_dom

import realJs "syscall/js"

// Func is a wrapped Go function to be called by JavaScript.
type Func struct {
	realJs.Func

	fn            func(this Value, args []Value) any
	internalValue Value
}

func createFuncValue(fn func(this Value, args []Value) any) *Value {
	val := Value{
		referencedType: functionType,
		internals: map[string]any{
			"prototype": uniqueID,
			"fn":        fn,
		},
	}
	uniqueID++
	return &val
}

// FuncOf returns a function to be used by JavaScript.
//
// The Go function fn is called with the value of JavaScript's "this" keyword and the
// arguments of the invocation. The return value of the invocation is
// the result of the Go function mapped back to JavaScript according to ValueOf.
//
// Invoking the wrapped Go function from JavaScript will
// pause the event loop and spawn a new goroutine.
// Other wrapped functions which are triggered during a call from Go to JavaScript
// get executed on the same goroutine.
//
// As a consequence, if one wrapped function blocks, JavaScript's event loop
// is blocked until that function returns. Hence, calling any async JavaScript
// API, which requires the event loop, like fetch (http.Client), will cause an
// immediate deadlock. Therefore a blocking function should explicitly start a
// new goroutine.
//
// Func.Release must be called to free up resources when the function will not be invoked any more.
func FuncOf(fn func(this Value, args []Value) any) Func {
	if isInFakeMode {
		return Func{
			fn:            fn,
			internalValue: *createFuncValue(fn),
		}
	}
	return Func{
		Func: realJs.FuncOf(func(this realJs.Value, args []realJs.Value) any {
			intermediaryValues := make([]Value, len(args))
			for i, arg := range args {
				intermediaryValues[i] = Value{Value: arg}
			}

			return fn(Value{Value: this}, intermediaryValues)
		}),
	}
}

// Release frees up resources allocated for the function.
// The function must not be invoked after calling Release.
// It is allowed to call Release while the function is still running.
func (c Func) Release() {
	if isInFakeMode {
		return
	}
	c.Func.Release()
}

// Equal reports whether v and w are equal according to JavaScript's === operator.
func (c Func) Equal(w Value) bool {
	if isInFakeMode {
		return c.internalValue.Equal(w)
	}
	return Value{Value: c.Func.Value}.Equal(w)
}

// IsUndefined reports whether v is the JavaScript value "undefined".
func (c Func) IsUndefined() bool {
	if isInFakeMode {
		return c.internalValue.IsUndefined()
	}
	return Value{Value: c.Func.Value}.IsUndefined()
}

// IsNull reports whether v is the JavaScript value "null".
func (c Func) IsNull() bool {
	if isInFakeMode {
		return c.internalValue.IsNull()
	}
	return Value{Value: c.Func.Value}.IsNull()
}

// IsNaN reports whether v is the JavaScript value "NaN".
func (c Func) IsNaN() bool {
	if isInFakeMode {
		return c.internalValue.IsNaN()
	}
	return Value{Value: c.Func.Value}.IsNaN()
}

// Type returns the JavaScript type of the value v. It is similar to JavaScript's typeof operator,
// except that it returns TypeNull instead of TypeObject for null.
func (c Func) Type() realJs.Type {
	if isInFakeMode {
		return c.internalValue.Type()
	}
	return Value{Value: c.Func.Value}.Type()
}

// Get returns the JavaScript property p of value v.
// It panics if v is not a JavaScript object.
func (c Func) Get(p string) Value {
	if isInFakeMode {
		return c.internalValue.Get(p)
	}
	return Value{Value: c.Func.Value}.Get(p)
}

// Set sets the JavaScript property p of value v to ValueOf(x).
// It panics if v is not a JavaScript object.
func (c Func) Set(p string, x any) {
	if isInFakeMode {
		c.internalValue.Set(p, x)
		return
	}
	Value{Value: c.Func.Value}.Set(p, x)
}

// Delete deletes the JavaScript property p of value v.
// It panics if v is not a JavaScript object.
func (c Func) Delete(p string) {
	if isInFakeMode {
		c.internalValue.Delete(p)
		return
	}
	Value{Value: c.Func.Value}.Delete(p)
}

// Index returns JavaScript index i of value v.
// It panics if v is not a JavaScript object.
func (c Func) Index(i int) Value {
	if isInFakeMode {
		return c.internalValue.Index(i)
	}
	return Value{Value: c.Func.Value}.Index(i)
}

// SetIndex sets the JavaScript index i of value v to ValueOf(x).
// It panics if v is not a JavaScript object.
func (c Func) SetIndex(i int, x any) {
	if isInFakeMode {
		c.internalValue.SetIndex(i, x)
		return
	}
	Value{Value: c.Func.Value}.SetIndex(i, x)
}

// Length returns the JavaScript property "length" of v.
// It panics if v is not a JavaScript object.
func (c Func) Length() int {
	if isInFakeMode {
		return c.internalValue.Length()
	}
	return Value{Value: c.Func.Value}.Length()
}

// Call does a JavaScript call to the method m of value v with the given arguments.
// It panics if v has no method m.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (c Func) Call(m string, args ...any) Value {
	if isInFakeMode {
		return c.internalValue.Call(m, args)
	}
	return Value{Value: c.Func.Value}.Call(m, args...)
}

// Invoke does a JavaScript call of the value v with the given arguments.
// It panics if v is not a JavaScript function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (c Func) Invoke(args ...any) Value {
	if isInFakeMode {
		return c.internalValue.Invoke(args)
	}
	return Value{Value: c.Func.Value}.Invoke(args...)
}

// New uses JavaScript's "new" operator with value v as constructor and the given arguments.
// It panics if v is not a JavaScript function.
// The arguments get mapped to JavaScript values according to the ValueOf function.
func (c Func) New(args ...any) Value {
	if isInFakeMode {
		return c.internalValue.New(args)
	}
	return Value{Value: c.Func.Value}.New(args...)
}

// Float returns the value v as a float64.
// It panics if v is not a JavaScript number.
func (c Func) Float() float64 {
	if isInFakeMode {
		return c.internalValue.Float()
	}
	return Value{Value: c.Func.Value}.Float()
}

// Int returns the value v truncated to an int.
// It panics if v is not a JavaScript number.
func (c Func) Int() int {
	if isInFakeMode {
		return c.internalValue.Int()
	}
	return Value{Value: c.Func.Value}.Int()
}

// Bool returns the value v as a bool.
// It panics if v is not a JavaScript boolean.
func (c Func) Bool() bool {
	if isInFakeMode {
		return c.internalValue.Bool()
	}
	return Value{Value: c.Func.Value}.Bool()
}

// Truthy returns the JavaScript "truthiness" of the value v. In JavaScript,
// false, 0, "", null, undefined, and NaN are "falsy", and everything else is
// "truthy". See https://developer.mozilla.org/en-US/docs/Glossary/Truthy.
func (c Func) Truthy() bool {
	if isInFakeMode {
		return c.internalValue.Truthy()
	}
	return Value{Value: c.Func.Value}.Truthy()
}

// String returns the value v as a string.
// String is a special case because of Go's String method convention. Unlike the other getters,
// it does not panic if v's Type is not TypeString. Instead, it returns a string of the form "<T>"
// or "<T: V>" where T is v's type and V is a string representation of v's value.
func (c Func) String() string {
	if isInFakeMode {
		return c.internalValue.String()
	}
	return Value{Value: c.Func.Value}.String()
}

// InstanceOf reports whether v is an instance of type t according to JavaScript's instanceof operator.
func (c Func) InstanceOf(t Value) bool {
	if isInFakeMode {
		return c.internalValue.InstanceOf(t)
	}
	return Value{Value: c.Func.Value}.InstanceOf(t)
}
