package router

import (
	"fmt"
	"syscall/js"

	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

// Router contains the routing state of the application, it must be created globally in an application
// and is used to create all the other routing components.
type Router struct {
	currentURL string

	handleHistoryFunc js.Func
}

// NewRouter generates a valid router pointer with all properties set.
func NewRouter() *Router {
	return &Router{}
}

// Provider provides the context and values for the router to work properly. It must be added as one of
// the first component of the tree and all subsequent router components or logic must happen in a descendant
// of the provider. The provider also listens to the popstate events to update the application if the user
// uses the back or forward buttons. Returns a fragment node, which allows passing more than one child.
func (r *Router) Provider(ctx context.Context, _ nodes.Props, children nodes.Children) nodes.Child {
	g := js.Global()
	if !g.Truthy() {
		panic("not in browser environment, global was undefined")
	}

	if r.currentURL == "" {
		r.currentURL = g.Get("window").Get("location").Call("toString").String()
	}
	ctx.SetValue("lander_routing_url", r.currentURL)

	ctx.OnMount(func() error {
		g := js.Global()
		if !g.Truthy() {
			return fmt.Errorf("not in browser environment, global was undefined")
		}

		r.handleHistoryFunc = js.FuncOf(func(this js.Value, args []js.Value) any {
			pathname := g.Get("window").Get("location").Call("toString").String()
			r.currentURL = pathname

			return ctx.Update()
		})

		g.Get("window").Call("addEventListener", "popstate", r.handleHistoryFunc)

		return nil
	})

	ctx.OnUnmount(func() error {
		g := js.Global()
		if !g.Truthy() {
			return fmt.Errorf("not in browser environment, global was undefined")
		}

		if r.handleHistoryFunc.IsUndefined() {
			return fmt.Errorf("the history change listener was somehow undefined, critical error")
		}

		g.Get("window").Call("removeEventListener", "popstate", r.handleHistoryFunc)
		r.handleHistoryFunc.Release()

		return nil
	})

	return nodes.NewFragmentNode(children)
}
