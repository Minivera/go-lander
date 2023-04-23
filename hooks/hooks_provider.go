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

func Provider(context context.Context, _ nodes.Props, children nodes.Children) nodes.Child {
	if !context.HasValue("lander_states") {
		context.SetValue("lander_states", nil)
	}

	states := context.GetValue("lander_states").(*stateChain)
	if !context.HasValue("lander_active_state") {
		context.SetValue("lander_active_state", states)
	}

	if len(children) != 1 {
		panic("Provider expects only one children")
	}

	return children[0]
}
