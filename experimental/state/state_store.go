package state

import (
	"fmt"

	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

type Store[T any] struct {
	state T
}

func NewStore[T any](defaultState T) *Store[T] {
	return &Store[T]{
		state: defaultState,
	}
}

func (s *Store[T]) Consumer(_ context.Context, props nodes.Props, children nodes.Children) nodes.Child {
	render, ok := props["render"].(func(state T) nodes.Child)
	if !ok {
		panic("Store.Consumer expects a component render prop as its `render` property.")
	}

	if len(children) > 0 {
		panic("Store.Consumer will not render any children, but a non-zero number of children were given.")
	}

	return render(s.state)
}

func (s *Store[T]) SetState(ctx context.Context, setter func(value T) T) error {
	s.state = setter(s.state)
	fmt.Printf("Setting state to %v\n", s.state)
	return ctx.Update()
}
