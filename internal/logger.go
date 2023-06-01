package internal

import (
	"fmt"

	js "github.com/minivera/go-lander/go-wasm-dom"
)

// IsDebug sets if the app is in debug mode, which will allow debug logging. Set to true by adding
// a GLOBAL variable to the browser window object and setting it to true.
var IsDebug = false

func init() {
	if !js.Global().Truthy() || !js.Global().Get("DEBUG").Truthy() {
		return
	}

	IsDebug = js.Global().Get("DEBUG").Bool()
}

// Debugln executes Println with the given parameters if IsDebug is true.
func Debugln(lines ...any) {
	if IsDebug {
		fmt.Println(lines...)
	}
}

// Debugf executes Printf with the given parameters if IsDebug is true.
func Debugf(format string, args ...any) {
	if IsDebug {
		fmt.Printf(format, args...)
	}
}
