package hooks

import (
	"fmt"

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
	fmt.Println("Rendering hooks provider")
	if !context.HasValue("lander_states") {
		context.SetValue("lander_states", nil)
	}

	states, ok := context.GetValue("lander_states").(*stateChain)
	if ok {
		context.SetValue("lander_active_state", states)
	} else {
		context.SetValue("lander_active_state", nil)
	}

	if len(children) != 1 {
		panic("Provider expects only one children")
	}

	return children[0]
}
