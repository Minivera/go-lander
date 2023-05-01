package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/experimental/router"
	"github.com/minivera/go-lander/nodes"
)

var appRouter = router.NewRouter()

func routingApp(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Html("h1", nodes.Attributes{}, nodes.Children{
			lander.Text("Sample routing app"),
		}),
		lander.Component(appRouter.Switch, nodes.Props{
			"routes": router.RouteDefinitions{
				{"/$", func(_ router.Match) nodes.Child {
					return lander.Html("div", nodes.Attributes{}, nodes.Children{
						lander.Html("h2", nodes.Attributes{}, nodes.Children{
							lander.Text("Home page"),
						}),
						lander.Html("ul", nodes.Attributes{}, nodes.Children{
							lander.Html("li", nodes.Attributes{}, nodes.Children{
								lander.Component(appRouter.Link, nodes.Props{
									"to": "/hello",
								}, nodes.Children{
									lander.Text("To /hello"),
								}),
							}),
							lander.Html("li", nodes.Attributes{}, nodes.Children{
								lander.Component(appRouter.Link, nodes.Props{
									"to": "/app",
								}, nodes.Children{
									lander.Text("To /app"),
								}),
							}),
							lander.Html("li", nodes.Attributes{}, nodes.Children{
								lander.Component(appRouter.Link, nodes.Props{
									"to": "/redirect",
								}, nodes.Children{
									lander.Text("To /redirect, which will send us back here"),
								}),
							}),
							lander.Html("li", nodes.Attributes{}, nodes.Children{
								lander.Component(appRouter.Link, nodes.Props{
									"to": "/notfound",
								}, nodes.Children{
									lander.Text("To the 404 page"),
								}),
							}),
						}),
					})
				}},
				{"/hello$", func(_ router.Match) nodes.Child {
					return lander.Html("div", nodes.Attributes{}, nodes.Children{
						lander.Html("h2", nodes.Attributes{}, nodes.Children{
							lander.Text("Hello, world!"),
						}),
						lander.Html("div", nodes.Attributes{}, nodes.Children{
							lander.Component(appRouter.Link, nodes.Props{
								"to": "/",
							}, nodes.Children{
								lander.Text("Go back to Home"),
							}),
						}),
					})
				}},
				{"/app.*", func(_ router.Match) nodes.Child {
					return lander.Html("div", nodes.Attributes{}, nodes.Children{
						lander.Html("h2", nodes.Attributes{}, nodes.Children{
							lander.Text("Welcome to the app"),
						}),
						lander.Component(appRouter.Route, nodes.Props{
							"route": "/app/([a-zA-Z0-9]+)/(?P<subroute>[a-zA-Z0-9]+)",
							"render": func(match router.Match) nodes.Child {
								return lander.Html("div", nodes.Attributes{}, nodes.Children{
									lander.Html("b", nodes.Attributes{}, nodes.Children{
										lander.Text("Matched:"),
									}),
									lander.Html("span", nodes.Attributes{}, nodes.Children{
										lander.Text(fmt.Sprintf("Pathname: %s", match.Pathname)),
									}),
									lander.Html("span", nodes.Attributes{}, nodes.Children{
										lander.Text(fmt.Sprintf("Path %s", match.Params["0"])),
									}),
									lander.Html("span", nodes.Attributes{}, nodes.Children{
										lander.Text(fmt.Sprintf("Subpath %s", match.Params["subroute"])),
									}),
								}).Style("display: flex; flex-direction: column; margin: 1rem; border: 1px solid black;")
							},
						}, nodes.Children{}),
						lander.Html("div", nodes.Attributes{}, nodes.Children{
							lander.Component(appRouter.Link, nodes.Props{
								"to": "/app/something/other",
							}, nodes.Children{
								lander.Text("Test the pattern matching"),
							}),
						}),
						lander.Html("div", nodes.Attributes{}, nodes.Children{
							lander.Component(appRouter.Link, nodes.Props{
								"to": "/",
							}, nodes.Children{
								lander.Text("Go back to Home"),
							}),
						}),
					})
				}},
				{"/redirect$", func(_ router.Match) nodes.Child {
					return lander.Component(appRouter.Redirect, nodes.Props{
						"to": "/",
					}, nodes.Children{})

				}},
				{".*", func(match router.Match) nodes.Child {
					return lander.Html("div", nodes.Attributes{}, nodes.Children{
						lander.Html("h2", nodes.Attributes{}, nodes.Children{
							lander.Text(fmt.Sprintf("404! `%s` was not found", match.Pathname)),
						}),
						lander.Html("div", nodes.Attributes{}, nodes.Children{
							lander.Component(appRouter.Link, nodes.Props{
								"to": "/",
							}, nodes.Children{
								lander.Text("Go back to Home"),
							}),
						}),
					})
				}},
			},
		}, nodes.Children{}),
	}).Style("padding: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(appRouter.Provider, nodes.Props{}, nodes.Children{
			lander.Component(routingApp, nodes.Props{}, nodes.Children{}),
		}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
