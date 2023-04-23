package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/hooks"
	"github.com/minivera/go-lander/nodes"
)

func addTodoForm(ctx context.Context, props nodes.Props, _ nodes.Children) nodes.Child {
	onAdd, ok := props["onAdd"].(func(value string) error)
	if !ok {
		fmt.Println("addTodoForm expects a function as its onAdd prop")
		// TODO: This is pretty terrible, improve. Maybe make props a struct?
		panic("addTodoForm expects a function as its onAdd prop")
	}

	value, setValue := hooks.UseState[string](ctx, "")

	return lander.Html("div", nodes.Attributes{}, []nodes.Child{
		lander.Html("input", nodes.Attributes{
			"value": value,
			"change": func(event *events.DOMEvent) error {
				value = event.JSEvent().Get("target").Get("value").String()
				return setValue(value)
			},
		}, []nodes.Child{}).Style("margin-right: 1rem;"),
		lander.Html("button", nodes.Attributes{
			"click": func(*events.DOMEvent) error {
				err := onAdd(value)
				if err != nil {
					return err
				}

				return setValue("")
			},
		}, []nodes.Child{
			lander.Text("Add"),
		}),
	}).Style("margin-top: 1rem; display: flex")
}

type todo struct {
	id        int
	name      string
	completed bool
}

func todoComponent(_ context.Context, props nodes.Props, _ nodes.Children) nodes.Child {
	onDelete, ok := props["onDelete"].(func() error)
	if !ok {
		fmt.Println("todoComponent expects a function as its onDelete prop")
		// TODO: This is pretty terrible, improve. Maybe make props a struct?
		panic("todoComponent expects a function as its onDelete prop")
	}

	onChange, ok := props["onChange"].(func() error)
	if !ok {
		fmt.Println("todoComponent expects a function as its onDelete prop")
		// TODO: This is pretty terrible, improve. Maybe make props a struct?
		panic("todoComponent expects a function as its onDelete prop")
	}

	currentTodo, ok := props["todo"].(todo)
	if !ok {
		fmt.Println("todoComponent expects a todo as its todo prop")
		// TODO: This is pretty terrible, improve. Maybe make props a struct?
		panic("todoComponent expects a todo as its todo prop")
	}

	return lander.Html("li", nodes.Attributes{}, []nodes.Child{
		lander.Html("div", nodes.Attributes{}, []nodes.Child{
			lander.Html("input", nodes.Attributes{
				"type":    "checkbox",
				"checked": currentTodo.completed,
				"change": func(*events.DOMEvent) error {
					return onChange()
				},
			}, []nodes.Child{}),
			lander.Html("strong", nodes.Attributes{}, []nodes.Child{
				lander.Text(currentTodo.name),
			}),
		}).Style("display: inline-flex; align-items: center; padding-right: 1rem;"),
		lander.Html("button", nodes.Attributes{
			"click": func(*events.DOMEvent) error {
				return onDelete()
			},
		}, []nodes.Child{
			lander.Text("X"),
		}).Style("display: inline;"),
	})

}

func todosApp(ctx context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	todos, setTodos := hooks.UseState[[]todo](ctx, []todo{
		{
			id:        0,
			name:      "write more examples",
			completed: false,
		},
	})

	updateTodo := func(todoId int, completed bool) error {
		todos := make([]todo, len(todos))

		for i, current := range todos {
			if todoId == current.id {
				todos[i] = todo{
					id:        i,
					name:      current.name,
					completed: completed,
				}
			} else {
				todos[i] = todo{
					id:        i,
					name:      current.name,
					completed: current.completed,
				}
			}
		}

		return setTodos(todos)
	}

	deleteTodo := func(todoId int) error {
		todos := make([]todo, len(todos)-1)

		count := 0
		for _, current := range todos {
			if current.id == todoId {
				continue
			}

			todos[count] = todo{
				id:        count,
				name:      current.name,
				completed: current.completed,
			}
			count += 1
		}

		return setTodos(todos)
	}

	addTodo := func(name string) error {
		todos := make([]todo, len(todos))

		for i, current := range todos {
			todos[i] = todo{
				id:        i,
				name:      current.name,
				completed: current.completed,
			}
		}

		todos = append(todos, todo{
			id:        len(todos),
			name:      name,
			completed: false,
		})

		return setTodos(todos)
	}

	fmt.Printf("Todos are %v\n", todos)
	todosComponents := make([]nodes.Child, len(todos))

	for i, todo := range todos {
		todosComponents[i] = lander.Component(todoComponent, nodes.Props{
			"onDelete": func() error {
				return deleteTodo(todo.id)
			},
			"onChange": func() error {
				return updateTodo(todo.id, !todo.completed)
			},
			"todo": todo,
		}, []nodes.Child{})
	}

	return lander.Component(hooks.Provider, nodes.Props{}, []nodes.Child{
		lander.Html("div", nodes.Attributes{}, []nodes.Child{
			lander.Html("h1", nodes.Attributes{}, []nodes.Child{
				lander.Text("Sample todo app"),
			}),
			lander.Html("div", nodes.Attributes{}, []nodes.Child{
				lander.Html("h2", nodes.Attributes{}, []nodes.Child{
					lander.Text("Todos"),
				}),
				lander.Html("ul", nodes.Attributes{}, todosComponents).Style("margin-top: 1rem;"),
				lander.Component(addTodoForm, nodes.Props{
					"onAdd": func(value string) error {
						return addTodo(value)
					},
				}, []nodes.Child{}),
			}).Style("max-width: 300px;"),
		}).Style("padding: 1rem;"),
	})
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(todosApp, nodes.Props{}, []nodes.Child{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
