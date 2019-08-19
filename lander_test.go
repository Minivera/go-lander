package lander

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHtml(t *testing.T) {
	t.Run("Should succeed if given valid node data", func(t *testing.T) {
		attributes := map[string]interface{}{
			"test": "test",
			"foo":  true,
		}
		tag := "div"
		id := "test"
		classes := []string{"test", "foo"}

		assert.Equal(
			t,
			&HTMLNode{
				Tag:            tag,
				Attributes:     map[string]string{"test": "test", "foo": ""},
				EventListeners: map[string]EventListener{},
				Classes:        classes,
				Children:       []Node{},
				DomID:          id,
				Styles:         []string{},
			},
			Html(fmt.Sprintf("%s#%s.%s", tag, id, strings.Join(classes, ".")), attributes, []Node{}),
			"The node should be equal to the expected node",
		)
	})

	t.Run("Should panic if given unknown attributes", func(t *testing.T) {
		attributes := map[string]interface{}{
			"test": struct{}{},
		}

		assert.Panics(t, func() { Html("", attributes, []Node{}) }, "The Html creator should panic")
	})
}

func TestSvg(t *testing.T) {
	t.Run("Should create an SVG namespaced HTML node", func(t *testing.T) {
		assert.Equal(
			t,
			&HTMLNode{
				Tag:            "d",
				Attributes:     map[string]string{},
				EventListeners: map[string]EventListener{},
				Classes:        []string{},
				Children:       []Node{},
				DomID:          "",
				namespace:      "http://www.w3.org/2000/svg",
				Styles:         []string{},
			},
			Svg("d", map[string]interface{}{}, []Node{}),
			"The node should be equal to the expected node",
		)
	})
}

func TestText(t *testing.T) {
	t.Run("Should create a text node with the text", func(t *testing.T) {
		assert.Equal(
			t,
			&TextNode{
				Text: "test",
			},
			Text("test"),
			"The node should be equal to the expected node",
		)
	})
}

func TestFragment(t *testing.T) {
	t.Run("Should create a fragment node with the children", func(t *testing.T) {
		children := []Node{
			&TextNode{
				Text: "test",
			},
		}

		assert.Equal(
			t,
			&FragmentNode{
				Children: children,
			},
			Fragment(children),
			"The node should be equal to the expected node",
		)
	})
}

func TestComponent(t *testing.T) {
	t.Run("Should create a function node with the factory and attributes", func(t *testing.T) {
		factory := func(map[string]interface{}, []Node) []Node { return nil }
		attributes := map[string]interface{}{}
		children := []Node{
			&TextNode{
				Text: "test",
			},
		}

		assert.Equal(
			t,
			children,
			Component(factory, attributes, children).givenChildren,
			"The node should be equal to the expected node",
		)
		assert.Equal(
			t,
			attributes,
			Component(factory, attributes, children).Attributes,
			"The node should be equal to the expected node",
		)
	})
}
