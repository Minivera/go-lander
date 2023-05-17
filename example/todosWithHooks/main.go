package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/experimental/hooks"
	"github.com/minivera/go-lander/nodes"
)

type addTodoFormProps struct {
	onAdd func(value string) error
}

func addTodoForm(ctx context.Context, props addTodoFormProps, _ nodes.Children) nodes.Child {
	onAdd := props.onAdd

	value, setValue, getValue := hooks.UseState[string](ctx, "")

	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Html("input", nodes.Attributes{
			"value": value,
			"change": func(event *events.DOMEvent) error {
				value = event.JSEvent().Get("target").Get("value").String()
				return setValue(func(_ string) string {
					return value
				})
			},
		}, nodes.Children{}).Style("margin-right: 1rem;"),
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
		}, nodes.Children{
			lander.Text("Add"),
		}),
	}).Style("margin-top: 1rem; display: flex")
}

type todo struct {
	id        int
	name      string
	completed bool
}

type todoComponentProps struct {
	onDelete    func() error
	onChange    func() error
	currentTodo todo
}

func todoComponent(ctx context.Context, props todoComponentProps, _ nodes.Children) nodes.Child {
	onDelete := props.onDelete
	onChange := props.onChange
	currentTodo := props.currentTodo

	return lander.Html("li", nodes.Attributes{}, nodes.Children{
		lander.Html("div", nodes.Attributes{}, nodes.Children{
			lander.Html("input", nodes.Attributes{
				"type":    "checkbox",
				"checked": currentTodo.completed,
				"change": func(*events.DOMEvent) error {
					return onChange()
				},
			}, nodes.Children{}),
			lander.Html("strong", nodes.Attributes{}, nodes.Children{
				lander.Text(currentTodo.name),
			}),
		}).Style("display: inline-flex; align-items: center; padding-right: 1rem;"),
		lander.Html("button", nodes.Attributes{
			"click": func(*events.DOMEvent) error {
				return onDelete()
			},
		}, nodes.Children{
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

	todosComponents := make(nodes.Children, len(todos))

	for i, todo := range todos {
		localTodo := todo
		todosComponents[i] = lander.Component(todoComponent, todoComponentProps{
			onDelete: func() error {
				return deleteTodo(localTodo.id)
			},
			onChange: func() error {
				return updateTodo(localTodo.id, !localTodo.completed)
			},
			currentTodo: localTodo,
		}, nodes.Children{})
	}

	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Html("h1", nodes.Attributes{}, nodes.Children{
			lander.Text("Sample todo app"),
		}),
		lander.Html("div", nodes.Attributes{}, nodes.Children{
			lander.Html("h2", nodes.Attributes{}, nodes.Children{
				lander.Text("Todos"),
			}),
			lander.Html("ul", nodes.Attributes{}, todosComponents).Style("margin-top: 1rem;"),
			lander.Component(addTodoForm, addTodoFormProps{
				onAdd: func(value string) error {
					return addTodo(value)
				},
			}, nodes.Children{}),
		}).Style("max-width: 300px;"),
	}).Style("padding: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(hooks.Provider, nodes.Props{}, nodes.Children{
			lander.Component(todosApp, nodes.Props{}, nodes.Children{}),
		}),
		"#app",
	)
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
