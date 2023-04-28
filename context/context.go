package context

import (
	"fmt"
)

var CurrentContext Context

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

	HasValue(name string) bool
	GetValue(name string) interface{}
	SetValue(name string, value interface{})
	IsDirty() bool

	Update() error
}

type baseContext struct {
	updateFunc func() error

	previousContext *baseContext

	contextValues map[string]interface{}
	isDirty       bool

	contextPerComponent map[interface{}][]string
	currentComponent    interface{}
	componentEvents     map[interface{}]map[string]func() error
}

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

func RegisterComponent(component interface{}) {
	converted := CurrentContext.(*baseContext)
	converted.currentComponent = component
}

func RegisterComponentContext(contextType string, component interface{}) {
	fmt.Printf("Registering context type %s for component %T, %v\n", contextType, component, component)
	converted := CurrentContext.(*baseContext)
	converted.contextPerComponent[component] = append(converted.contextPerComponent[component], contextType)
	converted.currentComponent = component
}

func UnregisterAllComponentContexts(component interface{}) {
	fmt.Printf("Removing all context types for component %T, %v\n", component, component)
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
	fmt.Printf("Registering event type %s for component %T (%p) %v\n", contextType, c.currentComponent, c.currentComponent, c.currentComponent)
	if _, ok := c.componentEvents[c.currentComponent]; !ok {
		c.componentEvents[c.currentComponent] = map[string]func() error{}
	}

	c.componentEvents[c.currentComponent][contextType] = listener
}

func (c *baseContext) triggerEvents() error {
	fmt.Printf("Trying to trigger events %v\n", c.componentEvents)
	for component, contextEvents := range c.contextPerComponent {
		fmt.Printf("Trying to trigger events for component %T, %v\n", component, component)
		fmt.Printf("Events are %v\n", contextEvents)

		// Ignore any context listeners for contexts that are not set on this particular component
		for _, name := range contextEvents {
			if name == "unmount" {
				fmt.Printf("Searching for unmount listener of component (%p) %T in previous context\n", component, component)
				// If the context is to unmount, then find the listener in the previous context instead
				if c.previousContext == nil {
					continue
				}

				for component, events := range c.previousContext.componentEvents {
					fmt.Printf("Previous context has component %T (%p) and events %v\n", component, component, events)
				}

				listener, ok := c.previousContext.componentEvents[component]["unmount"]
				if !ok {
					fmt.Printf("Listener for unmount was not found in previous context with component %T\n", component)
					// skip if the unmounted component doesn't trigger unmount
					continue
				}

				fmt.Printf("Executing unmount with component %T\n", component)
				err := listener()
				if err != nil {
					return fmt.Errorf("error in unmount listener for component. %w", err)
				}

				continue
			}

			// Don't continue with this component if it was never registered
			events, ok := c.componentEvents[component]
			if !ok {
				fmt.Printf("Component %T was never registered\n", component)
				continue
			}

			fmt.Printf("Testing for %s with component %T\n", name, component)
			listener, found := events[name]
			if !found {
				fmt.Printf("%s with component %T was never registered\n", name, component)
				continue
			}

			fmt.Printf("Executing %s with component %T\n", name, component)
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
