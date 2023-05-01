package context

import (
	"fmt"

	"github.com/minivera/go-lander/internal"
)

// CurrentContext is the active context available to all consumers while the diffing/mounting process
// is underway. This context will change on each subsequent rerender and may change during a render
// process. The CurrentContext is always available is wrapping a function with WithNewContext.
var CurrentContext Context

// Context is the rendering and mounting context available to all function components. It can carry values
// across the tree and allows hooking into the various event listeners for the component's state.
//
// IMPORTANT: Contrary to how React handles context, setting values inside the context uses pointers and
// will replace the value for all other components, regardless of where they are in the tree. React instead
// uses a context hierarchy based on the component tree.
type Context interface {
	// OnMount adds an event listener for when a component is first mounted into the tree.
	// This only triggers if the new component is inserted into the tree for the first time
	// due to a complete change in the layout. Component are reused, so this may not fire when
	// adding new components into a list. Use OnRender to consistently get render updates.
	OnMount(func() error)

	// OnRender triggers every time a component is updated and its content is render into the tree
	// this will trigger even when the component has not changed in any way. Components unmounted
	// will not fire OnRender, use OnUnmount.
	OnRender(func() error)

	// OnUnmount triggers an event listener when a component is removed from the tree. This will only
	// fire once and after the component has been removed from the tree.
	OnUnmount(func() error)

	// HasValue returns if the internal context has the given value saved in memory. This does not check
	// if the value is nil or undefined, only if the context was set to something.
	HasValue(name string) bool

	// GetValue gets the value under the given key in the context and returns it as an interface. There are
	// no generics or type conversion going on here, you are responsible for tracking and casting the types
	// in the context.
	GetValue(name string) interface{}

	// SetValue sets the given value under the given key in the context and returns it as an interface.
	// There are no generics or type conversion going on here, you are responsible for tracking and
	// casting the types in the context.
	SetValue(name string, value interface{})

	// IsDirty is a utility method that will return true as soon as something changed through SetValue in the
	// context. This is helpful to set all components as dirty and force a rerender. Function components
	// may be ignored otherwise.
	IsDirty() bool

	// Update triggers an update in the virtual DOM tree. Updates are thread safe and only one update can happen
	// at a time. Only once the update has completed and all hook listeners have been notified will this resolve.
	Update() error
}

// baseContext is the implemented version of the context interface for internal use only.
type baseContext struct {
	updateFunc func() error

	previousContext *baseContext

	contextValues map[string]interface{}
	isDirty       bool

	contextPerComponent map[interface{}][]string
	currentComponent    interface{}
	componentEvents     map[interface{}]map[string]func() error
}

// WithNewContext wraps the given function with a CurrentContext. The function will keep a reference of the
// previous version of CurrentContext and will restore it once it resolves. The function expects a update
// function to set as the value of Update in the Context interface and a previous context.
//
// previousContext is not the previous version of the context, this is handled internally. Rather this is
// the context from a previous render cycle. This must be set for unmounts to work properly and for
// context values to carry over subsequent renders.
func WithNewContext(updateFunc func() error, previousContext Context, call func() error) error {
	prevContext := CurrentContext
	localContext := &baseContext{
		updateFunc:    updateFunc,
		contextValues: map[string]interface{}{},

		contextPerComponent: map[interface{}][]string{},
		currentComponent:    nil,
		componentEvents:     map[interface{}]map[string]func() error{},
	}

	// Restore the old context if it was provided
	if previousContext != nil {
		localContext.previousContext = previousContext.(*baseContext)

		for key, value := range localContext.previousContext.contextValues {
			localContext.contextValues[key] = value
		}
	}

	CurrentContext = localContext

	err := call()
	if err != nil {
		return err
	}

	// trigger this async, so we have finished restoring the tree before this happens
	go func() {
		if err := localContext.triggerEvents(); err != nil {
			panic(err)
		}
	}()

	CurrentContext = prevContext
	return nil
}

// RegisterComponent registers the given interface as the current component being rendered in the
// CurrentContext. This is needed to properly link hooks like OnMount to the given component without
// asking consumers to pass the component reference.
func RegisterComponent(component interface{}) {
	converted := CurrentContext.(*baseContext)
	converted.currentComponent = component
}

// RegisterComponentContext registers the given context type for the given component. Only when a context
// type is registered will that component trigger its listeners. This avoids calling OnMount when the
// component is unmounting for example. The given component can be different from the last component given
// to RegisterComponent.
func RegisterComponentContext(contextType string, component interface{}) {
	internal.Debugf("Registering context type %s for component %T, %v\n", contextType, component, component)
	converted := CurrentContext.(*baseContext)
	converted.contextPerComponent[component] = append(converted.contextPerComponent[component], contextType)
	converted.currentComponent = component
}

// UnregisterAllComponentContexts unregisters all context type for the given component so the context
// can restart from a clean state. The given component can be different from the last component given
// to RegisterComponent.
func UnregisterAllComponentContexts(component interface{}) {
	internal.Debugf("Removing all context types for component %T, %v\n", component, component)
	converted := CurrentContext.(*baseContext)
	converted.contextPerComponent[component] = []string{}
}

func (c *baseContext) OnMount(listener func() error) {
	c.registerListener("mount", listener)
}

func (c *baseContext) OnRender(listener func() error) {
	c.registerListener("render", listener)
}

func (c *baseContext) OnUnmount(listener func() error) {
	c.registerListener("unmount", listener)
}

func (c *baseContext) registerListener(contextType string, listener func() error) {
	internal.Debugf("Registering event type %s for component %T (%p) %v\n", contextType, c.currentComponent, c.currentComponent, c.currentComponent)
	if _, ok := c.componentEvents[c.currentComponent]; !ok {
		c.componentEvents[c.currentComponent] = map[string]func() error{}
	}

	c.componentEvents[c.currentComponent][contextType] = listener
}

func (c *baseContext) triggerEvents() error {
	internal.Debugf("Trying to trigger events %v\n", c.componentEvents)
	for component, contextEvents := range c.contextPerComponent {
		internal.Debugf("Trying to trigger events for component %T, %v\n", component, component)
		internal.Debugf("Events are %v\n", contextEvents)

		// Ignore any context listeners for contexts that are not set on this particular component
		for _, name := range contextEvents {
			if name == "unmount" {
				internal.Debugf("Searching for unmount listener of component (%p) %T in previous context\n", component, component)
				// If the context is to unmount, then find the listener in the previous context instead
				if c.previousContext == nil {
					continue
				}

				for component, events := range c.previousContext.componentEvents {
					internal.Debugf("Previous context has component %T (%p) and events %v\n", component, component, events)
				}

				listener, ok := c.previousContext.componentEvents[component]["unmount"]
				if !ok {
					internal.Debugf("Listener for unmount was not found in previous context with component %T\n", component)
					// skip if the unmounted component doesn't trigger unmount
					continue
				}

				internal.Debugf("Executing unmount with component %T\n", component)
				err := listener()
				if err != nil {
					return fmt.Errorf("error in unmount listener for component. %w", err)
				}

				continue
			}

			// Don't continue with this component if it was never registered
			events, ok := c.componentEvents[component]
			if !ok {
				internal.Debugf("Component %T was never registered\n", component)
				continue
			}

			internal.Debugf("Testing for %s with component %T\n", name, component)
			listener, found := events[name]
			if !found {
				internal.Debugf("%s with component %T was never registered\n", name, component)
				continue
			}

			internal.Debugf("Executing %s with component %T\n", name, component)
			err := listener()
			if err != nil {
				return fmt.Errorf("error in %s listener for component. %w", name, err)
			}
		}
	}

	return nil
}

func (c *baseContext) HasValue(name string) bool {
	_, ok := c.contextValues[name]
	return ok
}

func (c *baseContext) GetValue(name string) interface{} {
	return c.contextValues[name]
}

func (c *baseContext) SetValue(name string, value interface{}) {
	c.contextValues[name] = value
	c.isDirty = true
}

func (c *baseContext) Update() error {
	return c.updateFunc()
}

func (c *baseContext) IsDirty() bool {
	return c.isDirty || (c.previousContext != nil && c.previousContext.isDirty)
}
