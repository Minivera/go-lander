package state

import (
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/internal"
	"github.com/minivera/go-lander/nodes"
)

// Store stores an arbitrary state type in-memory. It does not use context to save the store value
// and instead keeps the store in the struct. Make sure the store stays alive during the entire lifecycle
// of the app to avoid losing state.
type Store[T any] struct {
	state T
}

// NewStore creates a new global state store to be kept in the global memory of the application. The
// store will instantiate with the default state.
func NewStore[T any](defaultState T) *Store[T] {
	return &Store[T]{
		state: defaultState,
	}
}

// Consumer is a component that will use its `render` prop to render a component with the requested state.
// The render function will always execute with the most up-to-date version of the state. Consumer
// expects no children and only renders its render prop.
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

// SetState is a utility function to set the global state to a new value. This will trigger a full tree
// update and set all components to be rerendered, not only the consumer components. SetState expects
// a function as its setter parameter to ensure the value is updated with the latest version of the
// state. This does not merge the old and new state together, the new state is expected to include
// the new and old state merged.
func (s *Store[T]) SetState(ctx context.Context, setter func(value T) T) error {
	s.state = setter(s.state)
	internal.Debugf("Setting state to %v\n", s.state)
	return ctx.Update()
}
