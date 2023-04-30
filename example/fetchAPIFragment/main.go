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

type todos struct {
	Todos []struct {
		Id        int    `json:"id"`
		Todo      string `json:"todo"`
		Completed bool   `json:"completed"`
		UserId    int    `json:"userId"`
	} `json:"todos"`
	Total int `json:"total"`
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
}

type fetchApp struct {
	env *lander.DomEnvironment

	loaded      bool
	loadedTodos todos
}

func (a *fetchApp) todosList(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	todos := make([]nodes.Node, len(a.loadedTodos.Todos))
	for _, todo := range a.loadedTodos.Todos {
		todos = append(todos, lander.Html("fieldset", nodes.Attributes{}, nodes.Children{
			lander.Html("legend", nodes.Attributes{}, nodes.Children{
				lander.Text(todo.Todo),
			}),
			lander.Html("label", nodes.Attributes{
				"for": "id",
			}, nodes.Children{
				lander.Text("ID"),
			}),
			lander.Html("input", nodes.Attributes{
				"name":     "id",
				"value":    todo.Id,
				"readonly": true,
			}, nodes.Children{}),
			lander.Html("label", nodes.Attributes{
				"for": "todo",
			}, nodes.Children{
				lander.Text("Todo"),
			}),
			lander.Html("input", nodes.Attributes{
				"name":     "todo",
				"value":    todo.Todo,
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
				"checked":  todo.Completed,
				"disabled": true,
			}, nodes.Children{}),
		}).Style("width: 150px;"))
	}

	return lander.Fragment(todos)
}

func (a *fetchApp) render(ctx context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	ctx.OnMount(func() error {
		// Simulate some loading
		time.Sleep(2 * time.Second)

		resp, err := http.Get("https://dummyjson.com/todos")
		if err != nil {
			return err
		}

		err = json.NewDecoder(resp.Body).Decode(&a.loadedTodos)
		if err != nil {
			return err
		}

		a.loaded = true

		return a.env.Update()
	})

	var content nodes.Node = lander.Html("marquee", nodes.Attributes{}, nodes.Children{
		lander.Text("Loading..."),
	}).Style("width: 150px;")
	if a.loaded {
		content = lander.Component(a.todosList, nodes.Props{}, nodes.Children{})
	}

	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Html("h1", nodes.Attributes{}, nodes.Children{
			lander.Text("Sample loading app with Fragments"),
		}),
		lander.Html("div", nodes.Attributes{}, nodes.Children{
			content,
		}).Style("display: flex;flex-wrap: wrap;gap: 1rem;"),
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
