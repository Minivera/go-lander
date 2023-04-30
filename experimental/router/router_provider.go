package router

import (
	"fmt"
	"syscall/js"

	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

type Router struct {
	currentURL string

	handleHistoryFunc js.Func
}

func NewRouter() *Router {
	return &Router{}
}

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
