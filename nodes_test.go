// +build js,wasm

package lander

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

type mockNode struct {
	baseNode
}

func (n *mockNode) Clone() Node {
	return nil
}

func TestBaseNode_ID(t *testing.T) {
	t.Run("Returns the node ID", func(t *testing.T) {
		ID := uint64(1234567)

		assert.Equal(t, ID, (&baseNode{id: ID}).ID(), "Ids should be equal")
	})
}

func TestBaseNode_SetID(t *testing.T) {
	t.Run("Sets the node ID", func(t *testing.T) {
		ID := uint64(1234567)
		node := &baseNode{id: 0}
		node.SetID(ID)

		assert.Equal(t, ID, node.ID(), "Ids should be equal")
	})
}

func TestBaseNode_Create(t *testing.T) {
	t.Run("Create and sets the node ID", func(t *testing.T) {
		ID := uint64(1234567)
		node := &baseNode{}
		err := node.Create(ID)

		assert.Nil(t, err, "Create should not trigger an error")
		assert.Equal(t, ID, node.ID(), "Ids should be equal")
	})
}

func TestBaseNode_Position(t *testing.T) {
	t.Run("Sets the siblings and parent on position", func(t *testing.T) {
		node := &baseNode{}
		parent := &mockNode{
			baseNode: baseNode{
				id: 0,
			},
		}
		next := &mockNode{
			baseNode: baseNode{
				id: 1,
			},
		}
		prev := &mockNode{
			baseNode: baseNode{
				id: 2,
			},
		}
		err := node.Position(parent, next, prev)

		assert.Nil(t, err, "Position should not trigger an error")
		assert.Equal(t, parent, node.Parent, "Node should be positioned properly")
		assert.Equal(t, next, node.NextSibling, "Node should be positioned properly")
		assert.Equal(t, prev, node.PreviousSibling, "Node should be positioned properly")
	})
}

func TestBaseNode_Update(t *testing.T) {
	t.Run("Does nothing on base update", func(t *testing.T) {
		node := &baseNode{}
		err := node.Update(nil)

		assert.Nil(t, err, "Update should not trigger an error")
	})
}

func TestBaseNode_Render(t *testing.T) {
	t.Run("Does nothing on base render", func(t *testing.T) {
		node := &baseNode{}
		err := node.Render()

		assert.Nil(t, err, "Render should not trigger an error")
	})
}

func TestBaseNode_ToString(t *testing.T) {
	t.Run("Does nothing on base toString", func(t *testing.T) {
		node := &baseNode{}
		assert.Equal(t, "", node.ToString(), "The string return should be empty")
	})
}

func TestBaseNode_GetChildren(t *testing.T) {
	t.Run("Returns an empty slice on base GetChildren", func(t *testing.T) {
		node := &baseNode{}
		assert.Equal(t, []Node{}, node.GetChildren(), "The children array should be empty")
	})
}

func TestBaseNode_InsertChildren(t *testing.T) {
	t.Run("Does nothing on base insert", func(t *testing.T) {
		node := &baseNode{}
		err := node.InsertChildren(&mockNode{}, 0)

		assert.Nil(t, err, "InsertChildren should not trigger an error")
	})
}

func TestBaseNode_ReplaceChildren(t *testing.T) {
	t.Run("Does nothing on base replace", func(t *testing.T) {
		node := &baseNode{}
		err := node.ReplaceChildren(&mockNode{}, &mockNode{})

		assert.Nil(t, err, "ReplaceChildren should not trigger an error")
	})
}

func TestBaseNode_RemoveChildren(t *testing.T) {
	t.Run("Does nothing on base remove", func(t *testing.T) {
		node := &baseNode{}
		err := node.RemoveChildren(&mockNode{})

		assert.Nil(t, err, "RemoveChildren should not trigger an error")
	})
}

func compareNodes(t *testing.T, el1 interface{}, el2 interface{}) bool {
	if !cmp.Equal(
		el1,
		el2,
		cmpopts.IgnoreUnexported(HTMLNode{}, TextNode{}, FragmentNode{}, FuncNode{}, mockNode{}),
		cmpopts.IgnoreTypes(map[string]EventListener{}),
	) {
		t.Error(
			cmp.Diff(
				el1,
				el2,
				cmpopts.IgnoreUnexported(HTMLNode{}, TextNode{}, FragmentNode{}, FuncNode{}, mockNode{}),
				cmpopts.IgnoreTypes(map[string]EventListener{}),
			),
		)
		return false
	}
	return true
}

func TestNewHtmlNode(t *testing.T) {
	var event EventListener = func(Node, *DOMEvent) error { return nil }
	fakeNode := &mockNode{
		baseNode: baseNode{
			id: 0,
		},
	}

	tcs := []struct {
		scenario   string
		tag        string
		id         string
		classes    []string
		attributes map[string]interface{}
		children   []Node
		expected   *HTMLNode
		err        bool
	}{
		{
			scenario: "Succeeds when creating a node with basic info",
			tag:      "div",
			id:       "test",
			classes:  []string{"foo", "bar"},
			attributes: map[string]interface{}{
				"test":  "test",
				"foo":   1,
				"bar":   true,
				"click": event,
			},
			children: []Node{fakeNode},
			expected: &HTMLNode{
				Tag:     "div",
				DomID:   "test",
				Classes: []string{"foo", "bar"},
				Attributes: map[string]string{
					"test": "test",
					"foo":  "1",
					"bar":  "",
				},
				EventListeners: map[string]EventListener{
					"click": event,
				},
				Children: []Node{fakeNode},
				Styles:   []string{},
			},
		},
		{
			scenario: "Fail when creating a node with invalid attributes",
			tag:      "div",
			classes:  []string{},
			attributes: map[string]interface{}{
				"test": struct{}{},
			},
			children: []Node{},
			err:      true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			node, err := newHTMLNode(tc.tag, tc.id, tc.classes, tc.attributes, tc.children)

			if tc.err {
				assert.Error(t, err, "newHTMlNode should error")
			} else {
				assert.NoError(t, err, "newHTMLNode should not error")
				assert.NotNil(t, node)
				compareNodes(t, tc.expected, node)
			}
		})
	}
}

func TestHTMLNode_Update(t *testing.T) {
	var event EventListener = func(Node, *DOMEvent) error { return nil }

	tcs := []struct {
		scenario        string
		attributes      map[string]interface{}
		expected        map[string]string
		expectedID      string
		expectedClasses []string
		err             bool
	}{
		{
			scenario: "Succeeds when updating a node with valid info",
			attributes: map[string]interface{}{
				"id":    "test",
				"class": "foo bar",
				"foo":   1,
				"bar":   true,
				"click": event,
			},
			expected: map[string]string{
				"foo": "1",
				"bar": "",
			},
			expectedID:      "test",
			expectedClasses: []string{"foo", "bar"},
		},
		{
			scenario: "Fail when updating a node with invalid attributes",
			attributes: map[string]interface{}{
				"test": struct{}{},
			},
			err: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			node, _ := newHTMLNode("div", "", []string{}, map[string]interface{}{}, []Node{})
			err := node.Update(tc.attributes)

			if tc.err {
				assert.Error(t, err, "Update should error")
			} else {
				assert.NoError(t, err, "Update should not error")
				assert.Equal(t, tc.expected, node.Attributes, "Attributes should be properly changed")
				assert.Equal(t, tc.expectedID, node.DomID, "ID should be properly changed")
				assert.Equal(t, tc.expectedClasses, node.Classes, "Classes should be properly changed")
			}
		})
	}
}

func TestHTMLNode_ToString(t *testing.T) {
	tcs := []struct {
		scenario string
		node     *HTMLNode
		expected string
	}{
		{
			scenario: "Returns the printed node",
			node: Html("div#test.foo.bar", map[string]interface{}{
				"id":   "test",
				"test": "test",
				"foo":  1,
				"bar":  true,
			}, []Node{}),
			expected: `<div id="test" class="foo bar" bar="" foo="1" test="test"></div>`,
		},
		{
			scenario: "Returns the printed node with the content of its children",
			node: Html("div#test.foo.bar", map[string]interface{}{
				"test": "test",
				"foo":  1,
				"bar":  true,
			}, []Node{
				Text("test"),
			}),
			expected: `<div id="test" class="foo bar" bar="" foo="1" test="test">test</div>`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.node.ToString(), "Printed node should be equal to expected")
		})
	}
}

func TestHTMLNode_GetChildren(t *testing.T) {
	t.Run("GetChildren returns the node's children", func(t *testing.T) {
		children := []Node{&mockNode{}}
		node := &HTMLNode{
			Children: children,
		}

		assert.Equal(t, children, node.GetChildren(), "Returned children should be the same")
	})
}

func TestHTMLNode_InsertChildren(t *testing.T) {
	fakeNode1 := &mockNode{
		baseNode: baseNode{
			id: 0,
		},
	}
	fakeNode2 := &mockNode{
		baseNode: baseNode{
			id: 1,
		},
	}
	fakeNode3 := &mockNode{
		baseNode: baseNode{
			id: 2,
		},
	}

	tcs := []struct {
		scenario        string
		currentChildren []Node
		newChild        Node
		position        int
		expected        []Node
	}{
		{
			scenario:        "Add to the end when inserting with a position of -1",
			currentChildren: []Node{fakeNode1},
			newChild:        fakeNode2,
			position:        -1,
			expected:        []Node{fakeNode1, fakeNode2},
		},
		{
			scenario:        "Add to the head when inserting with a position of 0",
			currentChildren: []Node{fakeNode1, fakeNode2},
			newChild:        fakeNode3,
			position:        0,
			expected:        []Node{fakeNode3, fakeNode1, fakeNode2},
		},
		{
			scenario:        "Add to the middle when inserting with a position of 1",
			currentChildren: []Node{fakeNode1, fakeNode2},
			newChild:        fakeNode3,
			position:        1,
			expected:        []Node{fakeNode1, fakeNode3, fakeNode2},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			node := &HTMLNode{
				Children: tc.currentChildren,
			}
			err := node.InsertChildren(tc.newChild, tc.position)

			assert.NoError(t, err, "InsertChildren should never error")
			assert.Equal(t, tc.expected, node.GetChildren(), "Children should be the same")
		})
	}
}

func TestHTMLNode_ReplaceChildren(t *testing.T) {
	fakeNode1 := &mockNode{
		baseNode: baseNode{
			id: 0,
		},
	}
	fakeNode2 := &mockNode{
		baseNode: baseNode{
			id: 1,
		},
	}
	fakeNode3 := &mockNode{
		baseNode: baseNode{
			id: 2,
		},
	}

	tcs := []struct {
		scenario        string
		currentChildren []Node
		oldChild        Node
		newChild        Node
		expected        []Node
	}{
		{
			scenario:        "Replace the node if it can be found in the parent",
			currentChildren: []Node{fakeNode1, fakeNode2},
			oldChild:        fakeNode2,
			newChild:        fakeNode3,
			expected:        []Node{fakeNode1, fakeNode3},
		},
		{
			scenario:        "Does not replace the node if it could not be found",
			currentChildren: []Node{fakeNode1, fakeNode2},
			oldChild:        fakeNode3,
			newChild:        fakeNode3,
			expected:        []Node{fakeNode1, fakeNode2},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			node := &HTMLNode{
				Children: tc.currentChildren,
			}
			err := node.ReplaceChildren(tc.oldChild, tc.newChild)

			assert.NoError(t, err, "ReplaceChildren should never error")
			assert.Equal(t, tc.expected, node.GetChildren(), "Children should be the same")
		})
	}
}

func TestHTMLNode_RemoveChildren(t *testing.T) {
	fakeNode1 := &mockNode{
		baseNode: baseNode{
			id: 0,
		},
	}
	fakeNode2 := &mockNode{
		baseNode: baseNode{
			id: 1,
		},
	}
	fakeNode3 := &mockNode{
		baseNode: baseNode{
			id: 2,
		},
	}

	tcs := []struct {
		scenario        string
		currentChildren []Node
		toRemove        Node
		expected        []Node
	}{
		{
			scenario:        "Remove the node if it can be found in the parent",
			currentChildren: []Node{fakeNode1, fakeNode2, fakeNode3},
			toRemove:        fakeNode2,
			expected:        []Node{fakeNode1, fakeNode3},
		},
		{
			scenario:        "Does not remove the node if it could not be found",
			currentChildren: []Node{fakeNode1, fakeNode2},
			toRemove:        fakeNode3,
			expected:        []Node{fakeNode1, fakeNode2},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			node := &HTMLNode{
				Children: tc.currentChildren,
			}
			err := node.RemoveChildren(tc.toRemove)

			assert.NoError(t, err, "RemoveChildren should never error")
			assert.Equal(t, tc.expected, node.GetChildren(), "Children should be the same")
		})
	}
}

func TestHTMLNode_Clone(t *testing.T) {
	t.Run("Clone will make a deep copy of the node", func(t *testing.T) {
		node := Html("div.test#test", map[string]interface{}{
			"test": "test",
			"foo":  true,
			"bar":  1,
		}, []Node{}).Style(`
			color: black;
			font-weight: 500;
		`).SelectorStyle(" a", `
			color: black;
			font-weight: 500;
		`)

		compareNodes(t, node, node.Clone())
	})
}

func TestNewTextNode(t *testing.T) {
	t.Run("Will return a new text node with the given text", func(t *testing.T) {
		text := "test"
		node := newTextNode(text)

		assert.Equal(t, text, node.Text, "Text should be equal")
	})
}

func TestTextNode_Update(t *testing.T) {
	tcs := []struct {
		scenario   string
		attributes map[string]interface{}
		expected   string
		err        bool
	}{
		{
			scenario: "Succeeds when updating a node with valid text",
			attributes: map[string]interface{}{
				"text": "test",
			},
			expected: "test",
		},
		{
			scenario: "Fail when updating a node with invalid attributes",
			attributes: map[string]interface{}{
				"test": struct{}{},
			},
			err: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			node := newTextNode("test")
			err := node.Update(tc.attributes)

			if tc.err {
				assert.Error(t, err, "Update should error")
			} else {
				assert.NoError(t, err, "Update should not error")
				assert.Equal(t, tc.expected, node.Text, "Text should be properly changed")
			}
		})
	}
}

func TestTextNode_ToString(t *testing.T) {
	t.Run("ToString on Text node will return the text", func(t *testing.T) {
		text := "test"
		node := TextNode{
			Text: text,
		}

		assert.Equal(t, text, node.ToString(), "ToString should return the text")
	})
}

func TestTextNode_Clone(t *testing.T) {
	t.Run("Clone will make a deep copy of the node", func(t *testing.T) {
		node := Text("test")

		compareNodes(t, node, node.Clone())
	})
}

func TestNewFragmentNode(t *testing.T) {
	t.Run("Will return a new fragment node with the given children", func(t *testing.T) {
		children := []Node{&mockNode{}}
		node := newFragmentNode(children)

		assert.Equal(t, children, node.Children, "children should be equal")
	})
}

func TestFragmentNode_ToString(t *testing.T) {
	t.Run("ToString on Fragment node will return the children", func(t *testing.T) {
		node := FragmentNode{
			Children: []Node{
				Text("test"),
			},
		}

		assert.Equal(t, `test`, node.ToString(), "ToString should return the children's toString")
	})
}

func TestFragmentNode_GetChildren(t *testing.T) {
	t.Run("GetChildren returns the node's children", func(t *testing.T) {
		children := []Node{&mockNode{}}
		node := &FragmentNode{
			Children: children,
		}

		assert.Equal(t, children, node.GetChildren(), "Returned children should be the same")
	})
}

func TestFragmentNode_InsertChildren(t *testing.T) {
	fakeNode1 := &mockNode{
		baseNode: baseNode{
			id: 0,
		},
	}
	fakeNode2 := &mockNode{
		baseNode: baseNode{
			id: 1,
		},
	}
	fakeNode3 := &mockNode{
		baseNode: baseNode{
			id: 2,
		},
	}

	tcs := []struct {
		scenario        string
		currentChildren []Node
		newChild        Node
		position        int
		expected        []Node
	}{
		{
			scenario:        "Add to the end when inserting with a position of -1",
			currentChildren: []Node{fakeNode1},
			newChild:        fakeNode2,
			position:        -1,
			expected:        []Node{fakeNode1, fakeNode2},
		},
		{
			scenario:        "Add to the head when inserting with a position of 0",
			currentChildren: []Node{fakeNode1, fakeNode2},
			newChild:        fakeNode3,
			position:        0,
			expected:        []Node{fakeNode3, fakeNode1, fakeNode2},
		},
		{
			scenario:        "Add to the middle when inserting with a position of 1",
			currentChildren: []Node{fakeNode1, fakeNode2},
			newChild:        fakeNode3,
			position:        1,
			expected:        []Node{fakeNode1, fakeNode3, fakeNode2},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			node := &FragmentNode{
				Children: tc.currentChildren,
			}
			err := node.InsertChildren(tc.newChild, tc.position)

			assert.NoError(t, err, "InsertChildren should never error")
			assert.Equal(t, tc.expected, node.GetChildren(), "Children should be the same")
		})
	}
}

func TestFragmentNode_ReplaceChildren(t *testing.T) {
	fakeNode1 := &mockNode{
		baseNode: baseNode{
			id: 0,
		},
	}
	fakeNode2 := &mockNode{
		baseNode: baseNode{
			id: 1,
		},
	}
	fakeNode3 := &mockNode{
		baseNode: baseNode{
			id: 2,
		},
	}

	tcs := []struct {
		scenario        string
		currentChildren []Node
		oldChild        Node
		newChild        Node
		expected        []Node
	}{
		{
			scenario:        "Replace the node if it can be found in the parent",
			currentChildren: []Node{fakeNode1, fakeNode2},
			oldChild:        fakeNode2,
			newChild:        fakeNode3,
			expected:        []Node{fakeNode1, fakeNode3},
		},
		{
			scenario:        "Does not replace the node if it could not be found",
			currentChildren: []Node{fakeNode1, fakeNode2},
			oldChild:        fakeNode3,
			newChild:        fakeNode3,
			expected:        []Node{fakeNode1, fakeNode2},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			node := &FragmentNode{
				Children: tc.currentChildren,
			}
			err := node.ReplaceChildren(tc.oldChild, tc.newChild)

			assert.NoError(t, err, "ReplaceChildren should never error")
			assert.Equal(t, tc.expected, node.GetChildren(), "Children should be the same")
		})
	}
}

func TestFragmentNode_RemoveChildren(t *testing.T) {
	fakeNode1 := &mockNode{
		baseNode: baseNode{
			id: 0,
		},
	}
	fakeNode2 := &mockNode{
		baseNode: baseNode{
			id: 1,
		},
	}
	fakeNode3 := &mockNode{
		baseNode: baseNode{
			id: 2,
		},
	}

	tcs := []struct {
		scenario        string
		currentChildren []Node
		toRemove        Node
		expected        []Node
	}{
		{
			scenario:        "Remove the node if it can be found in the parent",
			currentChildren: []Node{fakeNode1, fakeNode2, fakeNode3},
			toRemove:        fakeNode2,
			expected:        []Node{fakeNode1, fakeNode3},
		},
		{
			scenario:        "Does not remove the node if it could not be found",
			currentChildren: []Node{fakeNode1, fakeNode2},
			toRemove:        fakeNode3,
			expected:        []Node{fakeNode1, fakeNode2},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			node := &FragmentNode{
				Children: tc.currentChildren,
			}
			err := node.RemoveChildren(tc.toRemove)

			assert.NoError(t, err, "RemoveChildren should never error")
			assert.Equal(t, tc.expected, node.GetChildren(), "Children should be the same")
		})
	}
}

func TestFragmentNode_Clone(t *testing.T) {
	t.Run("Clone will make a deep copy of the node", func(t *testing.T) {
		node := Fragment([]Node{&mockNode{}})

		compareNodes(t, node, node.Clone())
	})
}

func TestNewFuncNode(t *testing.T) {
	t.Run("Will return a new function node with the given parameters", func(t *testing.T) {
		children := []Node{&mockNode{}}
		attributes := map[string]interface{}{
			"test": "test",
		}
		var comp FunctionComponent = func(_ map[string]interface{}, _ []Node) []Node {
			return nil
		}
		node := newFuncNode(comp, attributes, children)

		assert.Equal(t, attributes, node.Attributes, "Attributes should be equal")
		assert.Equal(t, children, node.givenChildren, "Children should be equal")
	})
}

func TestFuncNode_Update(t *testing.T) {
	t.Run("Func nodes updates the attributes when updated", func(t *testing.T) {
		attributes := map[string]interface{}{
			"test": "test",
		}
		node := FuncNode{
			Attributes: map[string]interface{}{},
		}

		err := node.Update(attributes)

		assert.NoError(t, err, "Update should not error")
		assert.Equal(t, attributes, node.Attributes, "Attributes should be equal")
	})
}

func TestFuncNode_Render(t *testing.T) {
	t.Run("Render will execute the factory with parameters", func(t *testing.T) {
		children := []Node{&mockNode{}}
		var comp FunctionComponent = func(_ map[string]interface{}, _ []Node) []Node {
			return children
		}
		node := FuncNode{
			Attributes:    map[string]interface{}{},
			givenChildren: []Node{},
			factory:       comp,
		}

		err := node.Render()

		assert.NoError(t, err, "Update should not error")
		assert.Equal(t, children, node.Children, "Children should be equal")
	})
}

func TestFuncNode_ToString(t *testing.T) {
	t.Run("ToString on Func node will return the children", func(t *testing.T) {
		node := FuncNode{
			Children: []Node{
				Text("test"),
			},
		}

		assert.Equal(t, `test`, node.ToString(), "ToString should return the children's toString")
	})
}

func TestFuncNode_GetChildren(t *testing.T) {
	t.Run("GetChildren returns the node's children", func(t *testing.T) {
		children := []Node{&mockNode{}}
		node := &FuncNode{
			Children: children,
		}

		assert.Equal(t, children, node.GetChildren(), "Returned children should be the same")
	})
}

func TestFuncNode_Clone(t *testing.T) {
	t.Run("Clone will make a deep copy of the node", func(t *testing.T) {
		node := &FuncNode{
			Attributes: map[string]interface{}{
				"test": "test",
			},
			givenChildren: []Node{&mockNode{}},
			Children:      []Node{},
		}

		compareNodes(t, node, node.Clone())
	})
}
