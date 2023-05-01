package hooks

import (
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

type stateChain struct {
	mounted bool
	state   interface{}
	deps    []interface{}

	next *stateChain
}

// Provider provides teh context for hooks to work properly. This Provider must be added as the first
// component of the app. It takes care of setting up and tracking the states on every render. Returns
// a fragment node, which allows passing more than one child.
func Provider(context context.Context, _ nodes.Props, children nodes.Children) nodes.Child {
	if !context.HasValue("lander_states") {
		context.SetValue("lander_states", nil)
	}

	states, ok := context.GetValue("lander_states").(*stateChain)
	if ok {
		context.SetValue("lander_active_state", states)
	} else {
		context.SetValue("lander_active_state", nil)
	}

	return nodes.NewFragmentNode(children)
}
