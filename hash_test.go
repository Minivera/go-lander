// +build js,wasm

package lander

import (
	"testing"

	"github.com/cespare/xxhash"
	"github.com/stretchr/testify/assert"
)

func TestHashPosition(t *testing.T) {
	t.Run("Will return a hashed version of the position string given", func(t *testing.T) {
		pos := "testststst"
		assert.Equal(t, xxhash.Sum64String(pos), hashPosition(pos), "The string should be hashed with xxhash")
	})
}

func TestHashNode(t *testing.T) {
	htmlNode, _ := newHTMLNode("div", "test", []string{"test"}, map[string]interface{}{
		"test": "test",
	}, []Node{})
	textNode := newTextNode("test")
	funcNode := newFuncNode(nil, map[string]interface{}{
		"test": "test",
	}, []Node{})
	arrayNode := newFragmentNode([]Node{})

	tcs := []struct {
		scenario string
		node     Node
		expected uint64
	}{
		{
			scenario: "Hashes a valid HTML node",
			node:     htmlNode,
			expected: xxhash.Sum64String(`[test="test"][tag="div"][id="test"][class="test"][children="0"]`),
		},
		{
			scenario: "Hashes a valid Text node",
			node:     textNode,
			expected: xxhash.Sum64String(`[text="test"]`),
		},
		{
			scenario: "Hashes a valid Func node",
			node:     funcNode,
			expected: xxhash.Sum64String(`[test="test"][children="0"]`),
		},
		{
			scenario: "Hashes a valid Fragment node",
			node:     arrayNode,
			expected: xxhash.Sum64String(`[children="0"]`),
		},
		{
			scenario: "Returns 0 on an invalid node type",
			node:     &mockNode{},
			expected: 0,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			assert.Equal(t, tc.expected, hashNode(tc.node), "The hashed node should be valid")
		})
	}
}
