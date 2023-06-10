package go_wasm_dom

func compareObjectValues(a Value, b Value) bool {
	if len(a.properties) != len(b.properties) {
		return false
	}

	for k, v := range a.properties {
		otherVal, ok := b.properties[k]
		if !ok {
			return false
		}

		if !v.Equal(*otherVal) {
			return false
		}
	}

	return true
}

func compareDOMValues(a Value, b Value) bool {
	if a.tag != b.tag {
		return false
	}

	if len(a.attributes) != len(b.attributes) {
		return false
	}

	if len(a.properties) != len(b.properties) {
		return false
	}

	for k, v := range a.attributes {
		otherVal, ok := b.attributes[k]
		if !ok {
			return false
		}

		if v != otherVal {
			return false
		}
	}

	return compareObjectValues(a, b)
}
