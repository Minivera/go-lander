package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

type counterApp struct {
	env *lander.DomEnvironment

	count int
}

func (a *counterApp) render(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Html("h1", nodes.Attributes{}, nodes.Children{
			lander.Text("Sample counter app"),
		}),
		lander.Html("div", nodes.Attributes{}, nodes.Children{
			lander.Html("button", nodes.Attributes{
				"click": func(*events.DOMEvent) error {
					a.count -= 1
					return a.env.Update()
				},
			}, nodes.Children{
				lander.Text("-"),
			}),
			lander.Html("div", nodes.Attributes{}, nodes.Children{
				lander.Text(fmt.Sprintf("Counter is at: %d", a.count)),
			}).Style("padding-left: 1rem; padding-right: 1rem; color: red;"),
			lander.Html("button", nodes.Attributes{
				"click": func(*events.DOMEvent) error {
					a.count += 1
					return a.env.Update()
				},
			}, nodes.Children{
				lander.Text("+"),
			}),
		}).Style("display: flex;"),
		lander.Html("div", nodes.Attributes{}, nodes.Children{
			lander.Text("Testing the style updates"),
		}).Style(fmt.Sprintf("margin-top: 1rem; width: 200px; border: %dpx solid red;", a.count)),
	}).Style("padding: 1rem;")
}

func main() {
	c := make(chan bool)

	app := counterApp{}

	env, err := lander.RenderInto(
		lander.Component(app.render, nodes.Props{}, nodes.Children{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	app.env = env

	<-c
}
