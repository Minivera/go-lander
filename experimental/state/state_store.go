package state

import (
	"github.com/minivera/go-lander/context"
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

// ConsumerProps are the properties to assign to the Consumer component, use as the generic type.
type ConsumerProps[T any] struct {
	// Render is the function that will be executed to render this component's content. It takes in the entire
	// state as parameters and will render every time the state is updated.
	Render func(state T) nodes.Child
}

// Consumer is a component that will use its `render` prop to render a component with the requested state.
// The render function will always execute with the most up-to-date version of the state. Consumer
// expects no children and only renders its render prop.
func (s *Store[T]) Consumer(_ context.Context, props ConsumerProps[T], children nodes.Children) nodes.Child {
	render := props.Render

	if len(children) > 0 {
		panic("Store.Consumer will not render any children, but a non-zero number of children were given.")
	}

	// TODO: Force this component to rerender somehow, since we're not using context, it might not
	// TODO: always rerender?
	return render(s.state)
}

// SetState is a utility function to set the global state to a new value. This will trigger a full tree
// update and set all components to be rerendered, not only the consumer components. SetState expects
// a function as its setter parameter to ensure the value is updated with the latest version of the
// state. This does not merge the old and new state together, the new state is expected to include
// the new and old state merged.
func (s *Store[T]) SetState(ctx context.Context, setter func(value T) T) error {
	s.state = setter(s.state)
	return ctx.Update()
}
