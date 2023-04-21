package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/nodes"
)

type counterApp struct {
	env *lander.DomEnvironment

	count int
}

func (a *counterApp) render(_ lander.Props, _ lander.Children) lander.Child {
	return lander.Html("div", map[string]interface{}{}, []lander.Child{
		lander.Html("h1", map[string]interface{}{}, []nodes.Child{
			lander.Text("Sample counter app"),
		}),
		lander.Html("div", map[string]interface{}{}, []nodes.Child{
			lander.Html("button", map[string]interface{}{
				"click": func(*lander.DOMEvent) error {
					a.count -= 1
					return a.env.Update()
				},
			}, []nodes.Child{
				lander.Text("-"),
			}),
			lander.Html("div", map[string]interface{}{}, []nodes.Child{
				lander.Text(fmt.Sprintf("Counter is at: %d", a.count)),
			}).Style("padding-left: 1rem; padding-right: 1rem; color: red;"),
			lander.Html("button", map[string]interface{}{
				"click": func(*lander.DOMEvent) error {
					a.count += 1
					return a.env.Update()
				},
			}, []nodes.Child{
				lander.Text("+"),
			}),
		}).Style("display: flex;"),
		lander.Html("div", map[string]interface{}{}, []nodes.Child{
			lander.Text("Testing the style updates"),
		}).Style(fmt.Sprintf("margin-top: 1rem; width: 200px; border: %dpx solid red;", a.count)),
	}).Style("padding: 1rem;")
}

func main() {
	c := make(chan bool)

	app := counterApp{}

	env, err := lander.RenderInto(
		lander.Component(app.render, map[string]interface{}{}, []lander.Child{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	app.env = env

	<-c
}
