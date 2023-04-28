package hooks

import (
	"fmt"
	"reflect"

	"github.com/minivera/go-lander/context"
)

func useInternalMemo[T any](ctx context.Context, defaultValue T,
	deps []interface{}) (bool, T, func(func(T) T) error, func() T) {

	if !ctx.HasValue("lander_states") || !ctx.HasValue("lander_active_state") {
		panic("hooks were used outside of a hook provider, make sure to wrap your entire app in a `lander.Component(hooks.Provider)`")
	}

	fmt.Printf("attempting to setup state for %v\n", defaultValue)
	activeState, activeStateOk := ctx.GetValue("lander_active_state").(*stateChain)
	states, statesOk := ctx.GetValue("lander_states").(*stateChain)
	changed := false

	if states == nil || !statesOk {
		changed = true
		states = &stateChain{
			mounted: false,
			state:   defaultValue,
			deps:    deps,
			next:    nil,
		}
		ctx.SetValue("lander_states", states)
		fmt.Printf("states were empty, creating new empty states to is %T, %v\n", states, states)
		activeState = states
		activeStateOk = true
	}

	fmt.Printf("fetched active state is %T, %v\n", activeState, activeState)
	var realActiveState *stateChain
	if activeState == nil || !activeStateOk {
		fmt.Printf("Creating new active state %v\n", stateChain{
			mounted: false,
			state:   defaultValue,
			deps:    deps,
			next:    nil,
		})
		changed = true
		currentState := states
		for currentState.next != nil {
			currentState = currentState.next
		}

		realActiveState = &stateChain{
			mounted: false,
			state:   defaultValue,
			deps:    deps,
			next:    nil,
		}
		currentState.next = realActiveState
	} else {
		fmt.Printf("Using existing active state %v, %t\n", activeState, activeState == nil)
		realActiveState = activeState
	}

	fmt.Printf("current active state is %T, %v\n", realActiveState, realActiveState)
	if realActiveState.mounted && !reflect.DeepEqual(realActiveState.deps, deps) {
		changed = true
		realActiveState.state = defaultValue
	}

	ctx.OnMount(func() error {
		realActiveState.mounted = true
		return nil
	})

	ctx.OnUnmount(func() error {
		fmt.Println("Unmounting states")
		// On unmount, remove this state from the chain, so it is not reused in the future
		if ctx.HasValue("lander_states") {
			states := ctx.GetValue("lander_states").(*stateChain)

			currentState := states
			for currentState.next != nil && currentState.next != realActiveState {
				currentState = currentState.next
			}

			if currentState.next == realActiveState {
				currentState.next = realActiveState.next
			}
		}
		return nil
	})

	fmt.Printf("setting active state to %T, %v\n", realActiveState.next, realActiveState.next)
	ctx.SetValue("lander_active_state", realActiveState.next)
	return changed, realActiveState.state.(T), func(setter func(val T) T) error {
			realActiveState.state = setter(realActiveState.state.(T))
			return ctx.Update()
		}, func() T {
			return realActiveState.state.(T)
		}
}

// UseState hooks into the context to provide some updatable state to a component. This state can be updated
// with the second parameter. Note that due to how Golang share's closure variables by reference, any state
// variable that is not a pointer will not be updated inside the event listeners. The third return value
// can be used to always get the most up-to-date state value.
func UseState[T any](ctx context.Context, defaultValue T) (T, func(func(val T) T) error, func() T) {
	_, state, stateSetter, stateGetter := useInternalMemo[T](ctx, defaultValue, nil)
	return state, stateSetter, stateGetter
}

type effectState struct {
	effect  func() (func() error, error)
	cleanup func() error
}

func UseEffect(ctx context.Context, effect func() (func() error, error), deps []interface{}) {
	state := &effectState{
		effect: effect,
		cleanup: func() error {
			return nil
		},
	}
	changed, memoizedEffect, _, _ := useInternalMemo[*effectState](ctx, state, deps)

	ctx.OnRender(func() error {
		if changed {
			receivedCleanup, err := memoizedEffect.effect()
			if receivedCleanup != nil {
				state.cleanup = receivedCleanup
			}
			return err
		}

		return nil
	})

	ctx.OnUnmount(func() error {
		if memoizedEffect.cleanup != nil {
			return memoizedEffect.cleanup()
		}

		return nil
	})
}
