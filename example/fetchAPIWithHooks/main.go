package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/experimental/hooks"
	"github.com/minivera/go-lander/nodes"
)

type todo struct {
	Id        int    `json:"id"`
	Todo      string `json:"todo"`
	Completed bool   `json:"completed"`
	UserId    int    `json:"userId"`
}

func fetchApp(ctx context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	loading, setLoading, _ := hooks.UseState[bool](ctx, true)
	currentTodo, setTodo, _ := hooks.UseState[*todo](ctx, nil)

	hooks.UseEffect(ctx, func() (func() error, error) {
		// Simulate some loading
		time.Sleep(2 * time.Second)

		resp, err := http.Get("https://dummyjson.com/todos/1")
		if err != nil {
			return nil, err
		}

		loadedTodo := &todo{}
		err = json.NewDecoder(resp.Body).Decode(loadedTodo)
		if err != nil {
			return nil, err
		}

		err = setTodo(func(_ *todo) *todo {
			return loadedTodo
		})
		if err != nil {
			return nil, err
		}

		return nil, setLoading(func(_ bool) bool {
			return false
		})
	}, []interface{}{})

	content := lander.Html("marquee", nodes.Attributes{}, nodes.Children{
		lander.Text("Loading..."),
	}).Style("width: 150px;")
	if !loading {
		content = lander.Html("div", nodes.Attributes{}, nodes.Children{
			lander.Html("label", nodes.Attributes{
				"for": "id",
			}, nodes.Children{
				lander.Text("ID"),
			}),
			lander.Html("input", nodes.Attributes{
				"name":     "id",
				"value":    currentTodo.Id,
				"readonly": true,
			}, nodes.Children{}),
			lander.Html("label", nodes.Attributes{
				"for": "todo",
			}, nodes.Children{
				lander.Text("Todo"),
			}),
			lander.Html("input", nodes.Attributes{
				"name":     "todo",
				"value":    currentTodo.Todo,
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
				"checked":  currentTodo.Completed,
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

	_, err := lander.RenderInto(
		lander.Component(hooks.Provider, nodes.Props{}, []nodes.Child{
			lander.Component(fetchApp, nodes.Props{}, []nodes.Child{}),
		}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
