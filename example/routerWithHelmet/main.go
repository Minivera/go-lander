package main

import (
	"fmt"

	"github.com/minivera/go-lander/experimental/helmet"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/experimental/router"
	"github.com/minivera/go-lander/nodes"
)

var appRouter = router.NewRouter()

func routingApp(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Component(helmet.Head, nodes.Props{}, nodes.Children{
			lander.Html("title", nodes.Attributes{}, nodes.Children{
				lander.Text("Sample routing app"),
			}),
		}),
		lander.Html("h1", nodes.Attributes{}, nodes.Children{
			lander.Text("Sample routing app"),
		}),
		lander.Component(appRouter.Switch, router.SwitchProps{
			Routes: router.RouteDefinitions{
				{"/$", func(_ router.Match) nodes.Child {
					return lander.Html("div", nodes.Attributes{}, nodes.Children{
						lander.Html("h2", nodes.Attributes{}, nodes.Children{
							lander.Text("Home page"),
						}),
						lander.Html("ul", nodes.Attributes{}, nodes.Children{
							lander.Html("li", nodes.Attributes{}, nodes.Children{
								lander.Component(appRouter.Link, router.LinkProps{
									To: "/hello",
								}, nodes.Children{
									lander.Text("To /hello"),
								}),
							}),
							lander.Html("li", nodes.Attributes{}, nodes.Children{
								lander.Component(appRouter.Link, router.LinkProps{
									To: "/app",
								}, nodes.Children{
									lander.Text("To /app"),
								}),
							}),
							lander.Html("li", nodes.Attributes{}, nodes.Children{
								lander.Component(appRouter.Link, router.LinkProps{
									To: "/redirect",
								}, nodes.Children{
									lander.Text("To /redirect, which will send us back here"),
								}),
							}),
							lander.Html("li", nodes.Attributes{}, nodes.Children{
								lander.Component(appRouter.Link, router.LinkProps{
									To: "/notfound",
								}, nodes.Children{
									lander.Text("To the 404 page"),
								}),
							}),
						}),
					})
				}},
				{"/hello$", func(_ router.Match) nodes.Child {
					return lander.Html("div", nodes.Attributes{}, nodes.Children{
						lander.Component(helmet.Head, nodes.Props{}, nodes.Children{
							lander.Html("title", nodes.Attributes{}, nodes.Children{
								lander.Text("Sample routing app - Hello"),
							}),
						}),
						lander.Html("h2", nodes.Attributes{}, nodes.Children{
							lander.Text("Hello, world!"),
						}),
						lander.Html("div", nodes.Attributes{}, nodes.Children{
							lander.Component(appRouter.Link, router.LinkProps{
								To: "/",
							}, nodes.Children{
								lander.Text("Go back to Home"),
							}),
						}),
					})
				}},
				{"/app.*", func(_ router.Match) nodes.Child {
					return lander.Html("div", nodes.Attributes{}, nodes.Children{
						lander.Component(helmet.Head, nodes.Props{}, nodes.Children{
							lander.Html("title", nodes.Attributes{}, nodes.Children{
								lander.Text("Sample routing app - App"),
							}),
						}),
						lander.Html("h2", nodes.Attributes{}, nodes.Children{
							lander.Text("Welcome to the app"),
						}),
						lander.Component(appRouter.Route, router.RouteProps{
							Route: "/app/([a-zA-Z0-9]+)/(?P<subroute>[a-zA-Z0-9]+)",
							Render: func(match router.Match) nodes.Child {
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
							lander.Component(appRouter.Link, router.LinkProps{
								To: "/app/something/other",
							}, nodes.Children{
								lander.Text("Test the pattern matching"),
							}),
						}),
						lander.Html("div", nodes.Attributes{}, nodes.Children{
							lander.Component(appRouter.Link, router.LinkProps{
								To: "/",
							}, nodes.Children{
								lander.Text("Go back to Home"),
							}),
						}),
					})
				}},
				{"/redirect$", func(_ router.Match) nodes.Child {
					return lander.Component(appRouter.Redirect, router.RedirectProps{
						To: "/",
					}, nodes.Children{})

				}},
				{".*", func(match router.Match) nodes.Child {
					return lander.Html("div", nodes.Attributes{}, nodes.Children{
						lander.Component(helmet.Head, nodes.Props{}, nodes.Children{
							lander.Html("title", nodes.Attributes{}, nodes.Children{
								lander.Text("Sample routing app - Not found"),
							}),
						}),
						lander.Html("h2", nodes.Attributes{}, nodes.Children{
							lander.Text(fmt.Sprintf("404! `%s` was not found", match.Pathname)),
						}),
						lander.Html("div", nodes.Attributes{}, nodes.Children{
							lander.Component(appRouter.Link, router.LinkProps{
								To: "/",
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
			lander.Component(helmet.Provider, nodes.Props{}, nodes.Children{
				lander.Component(routingApp, nodes.Props{}, nodes.Children{}),
			}),
		}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
