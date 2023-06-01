package router

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	js "github.com/minivera/go-lander/go-wasm-dom"
	"github.com/minivera/go-lander/internal"
	"github.com/minivera/go-lander/nodes"
)

// Match is a struct that contains the details of a route match.
type Match struct {
	// Pathname is the matched pathname, without any transformations.
	Pathname string

	// Params is the extracted map of parameters in the pathname, based on the route regex. The
	// map will used the capture group names, or indexes if none are given.
	Params map[string]string
}

// RouteRender is the type definition for the render function when a route match. This uses the render
// prop pattern and will execute with the given match, it expects the rendered node to be returned.
type RouteRender = func(Match) nodes.Child

// RouteDefinition contains the information to define a possible route in a switch.
type RouteDefinition struct {
	// Route is a stringified regex that defines the route path. The regular expression will be executed
	// against the window's location and the Render function will execute if there is a match. Parameters
	// can be defined with capture groups and will be available in the match.
	//
	// The regular expression will execute against the full pathname, including the protocol and host.
	// For example, https://example.com/path.
	Route string

	// Render is the function to execute when there is a match on the provided Route. It will execute with
	// the match in parameters and expects a node in return.
	Render RouteRender
}

// RouteDefinitions is a slice of route definitions. We use a slice instead of a map to maintain the
// order of added definitions, which allows us to handle the switch in a cascade.
type RouteDefinitions = []RouteDefinition

// SwitchProps are the properties assigned to the Switch component, use as the generic props.
type SwitchProps struct {
	// Routes are the definitions to use for this Switch, see the Switch component documentation.
	Routes RouteDefinitions
}

// Switch is a component that expects no children and a `routes` property. The `routes` property should be
// a slice of RouteDefinitions. The Switch will loop over all definitions in order and check their routes
// against the window's location. The first route that matches will be rendered and all other routes are
// ignored. The Switch follows regex rules, so routes with less specificity should be added after more
// specific routes.
//
// Example:
// /app/test
// /app/.* <- Should be after to avoid this route matching against all possible subroutes.
//
// A catch-all route can be added at the end of the list using the `.*` regex.
func (r *Router) Switch(ctx context.Context, props SwitchProps, children nodes.Children) nodes.Child {
	if len(children) > 0 {
		panic("Router.Switch will not render any children, but a non-zero number of children were given.")
	}

	if !ctx.HasValue("lander_routing_url") {
		panic("routing components were used outside of a router provider, make sure to wrap your entire app in a `lander.Component(router.Provider)`")
	}

	pathname := ctx.GetValue("lander_routing_url").(string)

	routeDefs := props.Routes

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

// RouteProps are the properties assigned to the Route component, use as the generic props.
type RouteProps struct {
	// Route is the regex string to use as the route of this Route, see the Route component documentation.
	Route string

	// Render is the render function to execute if the route matches.
	Render RouteRender
}

// Route renders the provided render function if the route matches against the window's location.
// Route expects no children, and a `route` and `render` property. The route should be a stringified
// regex to match against the location. If the route matches, the render function will be executed with
// the match as its only parameter, it expects a node in return.
func (r *Router) Route(ctx context.Context, props RouteProps, children nodes.Children) nodes.Child {
	if len(children) > 0 {
		panic("Router.Route will not render any children, but a non-zero number of children were given.")
	}

	if !ctx.HasValue("lander_routing_url") {
		panic("routing components were used outside of a router provider, make sure to wrap your entire app in a `lander.Component(router.Provider)`")
	}
	pathname := ctx.GetValue("lander_routing_url").(string)
	internal.Debugf("Current pathname is %s\n", pathname)

	route := props.Route
	render := props.Render

	regex, err := regexp.Compile(route)
	if err != nil {
		panic(fmt.Sprintf("route %s in Router.Switch is not a valid regex, %s", route, err))
	}

	if !regex.MatchString(pathname) {
		internal.Debugf("%s did not match %s\n", pathname, route)
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

// Navigate navigates the user to the provided URL using the history browser API, if available.
// Replace can be given to replace the current state rather than pushing a new state on the history
// stack.
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

// LinkProps are the properties assigned to the Link component, use as the generic props.
type LinkProps struct {
	// To is the path to navigate to on click on the anchor, will also be used as the href attribute for the
	// anchor.
	To string

	// Decides if the path navigation should replace the current history entry or add a new entry on the
	// stack.
	Replace bool
}

// Link is an anchor component that allows routing using the history API on click. The Link expects a
// `to` property and a `replace` property, as per the Navigate API.
func (r *Router) Link(_ context.Context, props LinkProps, children nodes.Children) nodes.Child {
	to := props.To
	replace := props.Replace

	return nodes.NewHTMLNode("a", nodes.Attributes{
		"click": func(e *events.DOMEvent) error {
			r.Navigate(to, replace)
			e.PreventDefault()

			return nil
		},
		"href": to,
	}, children)
}

// RedirectProps are the properties assigned to the Redirect component, use as the generic props.
type RedirectProps struct {
	// To is the path to navigate to on render.
	To string

	// Decides if the path navigation should replace the current history entry or add a new entry on the
	// stack.
	Replace bool
}

// Redirect will automatically redirect the user on render (not in a hook) using its provided `to` and
// `replace` property, as per the Navigate API.
func (r *Router) Redirect(_ context.Context, props RedirectProps, _ nodes.Children) nodes.Child {
	to := props.To
	replace := props.Replace

	r.Navigate(to, replace)

	return nil
}
