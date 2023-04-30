package router

import (
	"fmt"
	"regexp"
	"strconv"
	"syscall/js"

	"github.com/minivera/go-lander/events"

	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

type Match struct {
	Pathname string
	Params   map[string]string
}

type RouteRender = func(Match) nodes.Child

type RouteDefinition struct {
	Route  string
	Render RouteRender
}

type RouteDefinitions = []RouteDefinition

func (r *Router) Switch(ctx context.Context, props nodes.Props, children nodes.Children) nodes.Child {
	if len(children) > 0 {
		panic("Router.Switch will not render any children, but a non-zero number of children were given.")
	}

	if !ctx.HasValue("lander_routing_url") {
		panic("routing components were used outside of a router provider, make sure to wrap your entire app in a `lander.Component(router.Provider)`")
	}

	pathname := ctx.GetValue("lander_routing_url").(string)

	routeDefs, ok := props["routes"].(RouteDefinitions)
	if !ok {
		panic("Router.Route expects a render function prop as its `render` property.")
	}

	for _, definition := range routeDefs {
		regex, err := regexp.Compile(definition.Route)
		if err != nil {
			panic(fmt.Sprintf("route %s in Router.Switch is not a valid regex, %s", definition.Route, err))
		}

		if !regex.MatchString(pathname) {
			continue
		}

		submatch := regex.FindStringSubmatch(pathname)
		groupNames := regex.SubexpNames()

		currentMatch := Match{
			Pathname: pathname,
			Params:   map[string]string{},
		}

		for i, val := range submatch[1:] {
			if groupNames[i+1] != "" {
				currentMatch.Params[groupNames[i+1]] = val
			} else {
				currentMatch.Params[strconv.Itoa(i)] = val
			}
		}

		return definition.Render(currentMatch)
	}

	return nil
}

func (r *Router) Route(ctx context.Context, props nodes.Props, children nodes.Children) nodes.Child {
	if len(children) > 0 {
		panic("Router.Route will not render any children, but a non-zero number of children were given.")
	}

	if !ctx.HasValue("lander_routing_url") {
		panic("routing components were used outside of a router provider, make sure to wrap your entire app in a `lander.Component(router.Provider)`")
	}
	pathname := ctx.GetValue("lander_routing_url").(string)
	fmt.Printf("Current pathname is %s\n", pathname)

	route, ok := props["route"].(string)
	if !ok {
		panic("Router.Route expects a string prop as its `route` property.")
	}

	render, ok := props["render"].(RouteRender)
	if !ok {
		panic("Router.Route expects a render function prop as its `render` property.")
	}

	regex, err := regexp.Compile(route)
	if err != nil {
		panic(fmt.Sprintf("route %s in Router.Switch is not a valid regex, %s", route, err))
	}

	if !regex.MatchString(pathname) {
		fmt.Printf("%s did not match %s\n", pathname, route)
		return nil
	}

	submatch := regex.FindStringSubmatch(pathname)
	groupNames := regex.SubexpNames()

	currentMatch := Match{
		Pathname: pathname,
		Params:   map[string]string{},
	}

	for i, val := range submatch[1:] {
		if groupNames[i+1] != "" {
			currentMatch.Params[groupNames[i+1]] = val
		} else {
			currentMatch.Params[strconv.Itoa(i)] = val
		}
	}

	return render(currentMatch)
}

func (r *Router) Navigate(to string, replace bool) {
	g := js.Global()
	if g.Truthy() {
		if replace {
			g.Get("window").Get("history").Call("replaceState", nil, "", to)
		} else {
			g.Get("window").Get("history").Call("pushState", nil, "", to)
		}

		r.handleHistoryFunc.Invoke()
	}
}

func (r *Router) Link(_ context.Context, props nodes.Props, children nodes.Children) nodes.Child {
	to, ok := props["to"].(string)
	if !ok {
		panic("Link expects a string prop as its `to` property.")
	}

	replace, ok := props["replace"].(bool)
	if !ok {
		replace = false
	}

	return nodes.NewHTMLNode("a", nodes.Attributes{
		"click": func(*events.DOMEvent) error {
			r.Navigate(to, replace)

			return nil
		},
	}, children)
}

func (r *Router) Redirect(_ context.Context, props nodes.Props, children nodes.Children) nodes.Child {
	to, ok := props["to"].(string)
	if !ok {
		panic("Link expects a string prop as its `to` property.")
	}

	replace, ok := props["replace"].(bool)
	if !ok {
		replace = false
	}

	r.Navigate(to, replace)

	return nil
}
