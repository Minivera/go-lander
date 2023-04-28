package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

type todo struct {
	Id        int    `json:"id"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
	UserId    int    `json:"userId"`
}

type fetchApp struct {
	env *lander.DomEnvironment

	loaded     bool
	loadedTodo todo
}

func (a *fetchApp) render(ctx context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	ctx.OnMount(func() error {
		// Simulate some loading
		time.Sleep(2 * time.Second)

		resp, err := http.Get("https://dummyjson.com/todos/1")
		if err != nil {
			return err
		}

		err = json.NewDecoder(resp.Body).Decode(&a.loadedTodo)
		if err != nil {
			return err
		}

		a.loaded = true

		return a.env.Update()
	})

	content := lander.Html("marquee", nodes.Attributes{}, nodes.Children{
		lander.Text("Loading..."),
	}).Style("width: 150px;")
	if a.loaded {
		content = lander.Html("div", nodes.Attributes{}, nodes.Children{
			lander.Html("label", nodes.Attributes{
				"for": "id",
			}, nodes.Children{
				lander.Text("ID"),
			}),
			lander.Html("input", nodes.Attributes{
				"name":     "id",
				"value":    a.loadedTodo.Id,
				"readonly": true,
			}, nodes.Children{}),
			lander.Html("label", nodes.Attributes{
				"for": "todo",
			}, nodes.Children{
				lander.Text("Todo"),
			}),
			lander.Html("input", nodes.Attributes{
				"name":     "todo",
				"value":    a.loadedTodo.Todo,
				"readonly": true,
			}, nodes.Children{}),
			lander.Html("label", nodes.Attributes{
				"for": "completed",
			}, nodes.Children{
				lander.Text("Completed?"),
			}),
			lander.Html("input", nodes.Attributes{
				"name":     "completed",
				"type":     "checkbox",
				"checked":  a.loadedTodo.Completed,
				"readonly": true,
			}, nodes.Children{}),
		}).Style("width: 150px;")
	}

	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Html("h1", nodes.Attributes{}, nodes.Children{
			lander.Text("Sample loading app"),
		}),
		lander.Html("div", nodes.Attributes{}, nodes.Children{
			content,
		}),
	}).Style("padding: 1rem;")
}

func main() {
	c := make(chan bool)

	app := fetchApp{}

	env, err := lander.RenderInto(
		lander.Component(app.render, nodes.Props{}, []nodes.Child{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	app.env = env

	<-c
}
