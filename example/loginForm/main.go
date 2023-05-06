package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

type loginForm struct {
	env *lander.DomEnvironment

	username string
	password string
}

func (f *loginForm) render(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Html("h1", nodes.Attributes{}, nodes.Children{
			lander.Text("Log into our app"),
		}),
		lander.Html("form", nodes.Attributes{}, nodes.Children{
			lander.Html("label", nodes.Attributes{
				"for": "username",
			}, nodes.Children{
				lander.Text("Username"),
			}).Style("font-weight: bold;"),
			lander.Html("input", nodes.Attributes{
				"name":        "username",
				"placeholder": "Enter Username",
				"change": func(event *events.DOMEvent) error {
					f.username = event.JSEvent().Get("target").Get("value").String()
					return f.env.Update()
				},
			}, nodes.Children{}),
			lander.Html("label", nodes.Attributes{
				"for": "password",
			}, nodes.Children{
				lander.Text("Password"),
			}).Style("font-weight: bold;"),
			lander.Html("input", nodes.Attributes{
				"name":        "password",
				"placeholder": "Enter Password",
				"type":        "password",
				"change": func(event *events.DOMEvent) error {
					f.password = event.JSEvent().Get("target").Get("value").String()
					return f.env.Update()
				},
			}, nodes.Children{}),
			lander.Html("button", nodes.Attributes{
				"type": "submit",
			}, nodes.Children{
				lander.Text("Submit"),
			}).Style("margin-top: 1rem;"),
		}),
	}).Style("margin: 1rem;")
}

func main() {
	c := make(chan bool)

	form := &loginForm{}

	env, err := lander.RenderInto(
		lander.Component(form.render, nodes.Props{}, nodes.Children{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	form.env = env

	<-c
}
