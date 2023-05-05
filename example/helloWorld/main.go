package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

func helloWorld(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Html("h1", nodes.Attributes{}, nodes.Children{
		lander.Text("Hello, World!"),
	}).Style("margin: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(helloWorld, nodes.Props{}, nodes.Children{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
