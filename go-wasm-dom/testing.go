package go_wasm_dom

import "testing"

var t *testing.T

func StartTestMode(givenT *testing.T) {
	t = givenT
	isInFakeMode = true
}

func Reset() {
	isInFakeMode = false
	resetScreen()
}
