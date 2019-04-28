package go_lander

import (
	"fmt"
	"strings"

	"github.com/cespare/xxhash"
)

func hashPosition(positionString string) uint64 {
	return xxhash.Sum64String(positionString)
}

func hashNode(node Node) uint64 {
	switch typedNode := node.(type) {
	case *HtmlNode:
		stringToHash := ""

		for key, value := range typedNode.Attributes {
			stringToHash += fmt.Sprintf(`[%s="%s"]`, key, value)
		}

		stringToHash += fmt.Sprintf(`[id="%s"]`, typedNode.DomID)
		stringToHash += fmt.Sprintf(`[class="%s"]`, strings.Join(typedNode.Classes, " "))
		stringToHash += fmt.Sprintf(`[children="%d"]`, len(typedNode.Children))

		return xxhash.Sum64String(stringToHash)
	case *TextNode:
		return xxhash.Sum64String(fmt.Sprintf(`[id="%s"]`, typedNode.Text))
	case *FuncNode:
		stringToHash := ""

		for key, value := range typedNode.Attributes {
			stringToHash += fmt.Sprintf(`[%s="%s"]`, key, value)
		}

		stringToHash += fmt.Sprintf(`[children="%d"]`, len(typedNode.givenChildren))

		return xxhash.Sum64String(stringToHash)
	case *FragmentNode:
		return xxhash.Sum64String(fmt.Sprintf(`[children="%d"]`, len(typedNode.Children)))
	}
	return 0
}
