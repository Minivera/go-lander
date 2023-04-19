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
	fmt.Printf("Executing app render function, state is %d\n", a.count)
	return lander.Html("div", map[string]interface{}{}, []lander.Child{
		lander.Text("Hello, World!"),
		lander.Html("button", map[string]interface{}{
			"click": func(*lander.DOMEvent) error {
				fmt.Println("Called onClick handler")
				a.count += 1
				return a.env.Update()
			},
		}, []nodes.Child{
			lander.Text(fmt.Sprintf("Click me to increase by 1, counting %d", a.count)),
		}),
	}).Style("color: red; padding: 1rem;")
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
