// +build js,wasm

package lander

import (
	"syscall/js"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedJsValue struct {
	mock.Mock
}

func (d mockedJsValue) Call(name string, args ...interface{}) js.Value {
	ret := d.Called(name, args)
	return ret.Get(0).(js.Value)
}

func (d mockedJsValue) Get(name string) js.Value {
	ret := d.Called(name)
	return ret.Get(0).(js.Value)
}

func (d mockedJsValue) Index(index int) js.Value {
	ret := d.Called(index)
	return ret.Get(0).(js.Value)
}

func (d mockedJsValue) Set(name string, val interface{}) {
	d.Called(name, val)
}

func (d mockedJsValue) Truthy() bool {
	ret := d.Called()
	return ret.Bool(0)
}

func TestPatchText_Execute(t *testing.T) {
	fakeNode1 := &mockNode{
		baseNode: baseNode{
			id: 1,
		},
	}
	fakeNode2 := &mockNode{
		baseNode: baseNode{
			id: 2,
		},
	}
	textNode := newTextNode("foo")
	htmlNode, _ := newHTMLNode("div", "", []string{}, map[string]interface{}{}, []Node{fakeNode1, textNode, fakeNode2})

	tcs := []struct {
		scenario string
		oldNode  Node
		parent   Node
		newText  string
		err      bool
	}{
		{
			scenario: "Works when given valid values",
			oldNode:  textNode,
			parent:   htmlNode,
			newText:  "test",
		},
		{
			scenario: "Fails when the child cannot be found",
			oldNode:  &mockNode{},
			parent:   htmlNode,
			newText:  "test",
			err:      true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			parentDom := mockedJsValue{}
			array := js.ValueOf([]interface{}{map[string]interface{}{}, map[string]interface{}{}, map[string]interface{}{}})

			parentDom.On("Get", "childNodes").Return(array)

			patch := newPatchText(parentDom, tc.parent, tc.oldNode, tc.newText)
			err := patch.execute(mockedJsValue{})

			if tc.err {
				assert.Error(t, err, "execute should error")
			} else {
				assert.Equal(t, tc.newText, tc.oldNode.ToString(), "The old text node should have is text changed")
				parentDom.AssertExpectations(t)
				assert.Equal(t, tc.newText, array.Index(1).Get("nodeValue").String(), "Text should be changed on the text node")
			}
		})
	}
}

/* func TestPatchHTML_Execute(t *testing.T) {
	htmlNode1, _ := newHTMLNode("div", "", []string{}, map[string]interface{}{
		"test": "test",
	}, []Node{})
	htmlNode2, _ := newHTMLNode("div", "", []string{}, map[string]interface{}{
		"test":  "foo",
		"id":    "test",
		"class": "foo bar",
	}, []Node{})

	tcs := []struct {
		scenario string
		oldNode  Node
		newNode  Node
		err      bool
	}{
		{
			scenario: "Works when given valid values",
			oldNode:  htmlNode1,
			newNode:  htmlNode2,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			node := js.ValueOf(map[string]interface{}{
				"setAttribute": true,
			})
			oldNode := js.ValueOf(map[string]interface{}{})
			document := mockedJsValue{}

			document.On("Call", "createElement", mock.Anything).Return(node)
			document.On("Call", "querySelector", fmt.Sprintf(`[data-lander-id="%d"]`, tc.oldNode.ID())).Return(oldNode)
			document.On("Call", "replaceChild", oldNode, node).Once()

			patch := newPatchHTML(tc.oldNode, tc.newNode)
			err := patch.execute(document)

			if tc.err {
				assert.Error(t, err, "execute should error")
			} else {
				newHtml, ok := tc.newNode.(*HTMLNode)
				assert.True(t, ok)

				oldHtml, ok := tc.oldNode.(*HTMLNode)
				assert.True(t, ok)

				assert.Equal(t, newHtml.Attributes, oldHtml.Attributes, "The attributes of both nodes should be equal")
				assert.Equal(t, newHtml.Classes, oldHtml.Classes, "The class of both nodes should be equal")
				assert.Equal(t, newHtml.DomID, oldHtml.DomID, "The dom id of both nodes should be equal")
				document.AssertExpectations(t)
			}
		})
	}
} */
