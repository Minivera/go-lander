package hooks

import (
	"fmt"
	"reflect"

	"github.com/minivera/go-lander/context"
)

func useInternalMemo[T any](ctx context.Context, defaultValue T,
	deps []interface{}) (bool, T, func(val T) error) {

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
		activeState = states
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
		fmt.Printf("Using nil active state %v, %t\n", activeState, activeState == nil)
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

	ctx.SetValue("lander_active_state", realActiveState.next)
	return changed, realActiveState.state.(T), func(val T) error {
		realActiveState.state = val
		return ctx.Update()
	}
}

func UseState[T any](ctx context.Context, defaultValue T) (T, func(val T) error) {
	_, state, stateSetter := useInternalMemo[T](ctx, defaultValue, []interface{}{})
	return state, stateSetter
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
	changed, memoizedEffect, _ := useInternalMemo[*effectState](ctx, state, deps)

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

	ctx.OnUnmount(memoizedEffect.cleanup)
}
