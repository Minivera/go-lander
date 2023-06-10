package go_wasm_dom

import "testing"

var t *testing.T

func StartTestMode(givenT *testing.T) {
	t = givenT
	isInFakeMode = true
	resetScreen()
}

func EnableDebug(givenT *testing.T) {
	if currentScreen == nil {
		givenT.Fatalf("tried to enable debug without first starting the test mode")
	}

	window.properties["DEBUG"] = valuePtr(ValueOf(true))
}

func Reset() {
	resetScreen()
}

func EndTestMode() {
	isInFakeMode = false
}
