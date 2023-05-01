package internal

import (
	"fmt"
	"syscall/js"
)

var IsDebug = false

func init() {
	if !js.Global().Truthy() || !js.Global().Get("DEBUG").Truthy() {
		return
	}

	IsDebug = js.Global().Get("DEBUG").Bool()
}

func Debugln(lines ...any) {
	if IsDebug {
		fmt.Println(lines...)
	}
}

func Debugf(format string, args ...any) {
	if IsDebug {
		fmt.Printf(format, args...)
	}
}
