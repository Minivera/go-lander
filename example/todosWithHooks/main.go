package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/experimental/hooks"
	"github.com/minivera/go-lander/nodes"
)

func addTodoForm(ctx context.Context, props nodes.Props, _ nodes.Children) nodes.Child {
	onAdd, ok := props["onAdd"].(func(value string) error)
	if !ok {
		fmt.Println("addTodoForm expects a function as its onAdd prop")
		// TODO: This is pretty terrible, improve. Maybe make props a struct?
		panic("addTodoForm expects a function as its onAdd prop")
	}

	value, setValue, getValue := hooks.UseState[string](ctx, "")

	return lander.Html("div", nodes.Attributes{}, []nodes.Child{
		lander.Html("input", nodes.Attributes{
			"value": value,
			"change": func(event *events.DOMEvent) error {
				value = event.JSEvent().Get("target").Get("value").String()
				return setValue(func(_ string) string {
					return value
				})
			},
		}, []nodes.Child{}).Style("margin-right: 1rem;"),
		lander.Html("button", nodes.Attributes{
			"click": func(*events.DOMEvent) error {
				err := onAdd(getValue())
				if err != nil {
					return err
				}

				return setValue(func(_ string) string {
					return ""
				})
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

func todoComponent(ctx context.Context, props nodes.Props, _ nodes.Children) nodes.Child {
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
	todos, setTodos, _ := hooks.UseState[[]todo](ctx, []todo{
		{
			id:        0,
			name:      "write more examples",
			completed: false,
		},
	})

	updateTodo := func(todoId int, completed bool) error {
		return setTodos(func(todos []todo) []todo {
			newTodos := make([]todo, len(todos))

			for i, current := range todos {
				if todoId == current.id {
					newTodos[i] = todo{
						id:        i,
						name:      current.name,
						completed: completed,
					}
				} else {
					newTodos[i] = todo{
						id:        i,
						name:      current.name,
						completed: current.completed,
					}
				}
			}

			return newTodos
		})
	}

	deleteTodo := func(todoId int) error {
		return setTodos(func(todos []todo) []todo {
			newTodos := make([]todo, len(todos)-1)

			count := 0
			for _, current := range todos {
				if current.id == todoId {
					continue
				}

				newTodos[count] = todo{
					id:        count,
					name:      current.name,
					completed: current.completed,
				}
				count += 1
			}

			return newTodos
		})
	}

	addTodo := func(name string) error {
		return setTodos(func(todos []todo) []todo {
			newTodos := make([]todo, len(todos))

			for i, current := range todos {
				newTodos[i] = todo{
					id:        i,
					name:      current.name,
					completed: current.completed,
				}
			}

			newTodos = append(newTodos, todo{
				id:        len(newTodos),
				name:      name,
				completed: false,
			})

			return newTodos
		})
	}

	todosComponents := make([]nodes.Child, len(todos))

	for i, todo := range todos {
		localTodo := todo
		todosComponents[i] = lander.Component(todoComponent, nodes.Props{
			"onDelete": func() error {
				return deleteTodo(localTodo.id)
			},
			"onChange": func() error {
				fmt.Printf("Triggering on change of todo %d, %v\n", localTodo.id, localTodo)
				return updateTodo(localTodo.id, !localTodo.completed)
			},
			"todo": localTodo,
		}, []nodes.Child{})
	}

	return lander.Html("div", nodes.Attributes{}, []nodes.Child{
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
	}).Style("padding: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(hooks.Provider, nodes.Props{}, []nodes.Child{
			lander.Component(todosApp, nodes.Props{}, []nodes.Child{}),
		}),
		"#app",
	)
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
