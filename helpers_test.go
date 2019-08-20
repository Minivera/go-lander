// +build js,wasm

package lander

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractAttributes(t *testing.T) {
	var sampleEventListener EventListener = func(Node, *DOMEvent) error {
		return nil
	}

	tcs := []struct {
		scenario               string
		attributes             map[string]interface{}
		expectedAttributes     map[string]string
		expectedEventListeners []string
		err                    error
	}{
		{
			scenario: "Work when given simple attributes",
			attributes: map[string]interface{}{
				"value":    "test",
				"min":      2,
				"required": true,
			},
			expectedAttributes: map[string]string{
				"value":    "test",
				"min":      "2",
				"required": "",
			},
			expectedEventListeners: []string{},
		},
		{
			scenario: "Work when given attributes with event listeners",
			attributes: map[string]interface{}{
				"click": sampleEventListener,
			},
			expectedAttributes:     map[string]string{},
			expectedEventListeners: []string{"click"},
		},
		{
			scenario: "Fails if given random attributes",
			attributes: map[string]interface{}{
				"test": struct{}{},
			},
			err: fmt.Errorf("attributes only support vars of type string, bool, int or EventListener, %T received", struct{}{}),
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			attributes, listeners, err := extractAttributes(tc.attributes)

			if tc.err != nil {
				assert.Equal(t, tc.err, err, "Should return an error of type", tc.err)
			} else {
				assert.Nil(t, err, "Should not return an error")
				assert.Equal(t, tc.expectedAttributes, attributes, "Attributes should equal the expected scenario")

				for _, value := range tc.expectedEventListeners {
					assert.NotNil(t, listeners[value], "Listeners should equal the expected scenario")
				}
			}
		})
	}
}

func TestHyperscript(t *testing.T) {
	tcs := []struct {
		scenario        string
		given           string
		expectedTag     string
		expectedID      string
		expectedClasses []string
	}{
		{
			scenario:        "Returns a div by default",
			given:           "",
			expectedTag:     "div",
			expectedClasses: []string{},
		},
		{
			scenario:        "Returns the tag if only containing a tag",
			given:           "input",
			expectedTag:     "input",
			expectedClasses: []string{},
		},
		{
			scenario:        "Works with a normal tag without id nor classes",
			given:           "div",
			expectedTag:     "div",
			expectedClasses: []string{},
		},
		{
			scenario:        "Works with a tag and an id",
			given:           "div#test",
			expectedTag:     "div",
			expectedID:      "test",
			expectedClasses: []string{},
		},
		{
			scenario:        "Works with a tag, an id and one class",
			given:           "div#test.test",
			expectedTag:     "div",
			expectedID:      "test",
			expectedClasses: []string{"test"},
		},
		{
			scenario:        "Works with a tag, an id and multiple classes",
			given:           "div#test.test.foo.bar",
			expectedTag:     "div",
			expectedID:      "test",
			expectedClasses: []string{"test", "foo", "bar"},
		},
		{
			scenario:        "Works with id and class randomly placed in the hyperscript tag",
			given:           "div.test#test.foo.bar",
			expectedTag:     "div",
			expectedID:      "test",
			expectedClasses: []string{"test", "foo", "bar"},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			tag, id, classes := hyperscript(tc.given)

			assert.Equal(t, tc.expectedTag, tag)
			assert.Equal(t, tc.expectedID, id)
			assert.Equal(t, tc.expectedClasses, classes)
		})
	}
}
