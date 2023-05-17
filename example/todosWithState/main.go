package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/experimental/state"
	"github.com/minivera/go-lander/nodes"
)

type todo struct {
	id        int
	name      string
	completed bool
}

type appState struct {
	todos      []todo
	currentAdd string
}

var store = state.NewStore[appState](appState{
	todos: []todo{
		{
			id:        0,
			name:      "write more examples",
			completed: false,
		},
	},
})

type addTodoFormProps struct {
	state appState
}

func addTodoForm(ctx context.Context, props addTodoFormProps, _ nodes.Children) nodes.Child {
	currentState := props.state

	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Html("input", nodes.Attributes{
			"value": currentState.currentAdd,
			"change": func(event *events.DOMEvent) error {
				value := event.JSEvent().Get("target").Get("value").String()
				return store.SetState(ctx, func(currentState appState) appState {
					return appState{
						todos:      currentState.todos,
						currentAdd: value,
					}
				})
			},
		}, nodes.Children{}).Style("margin-right: 1rem;"),
		lander.Html("button", nodes.Attributes{
			"click": func(*events.DOMEvent) error {
				return store.SetState(ctx, func(currentState appState) appState {
					newTodos := make([]todo, len(currentState.todos))

					for i, current := range currentState.todos {
						newTodos[i] = todo{
							id:        i,
							name:      current.name,
							completed: current.completed,
						}
					}

					newTodos = append(newTodos, todo{
						id:        len(newTodos),
						name:      currentState.currentAdd,
						completed: false,
					})

					return appState{
						todos:      newTodos,
						currentAdd: "",
					}
				})
			},
		}, nodes.Children{
			lander.Text("Add"),
		}),
	}).Style("margin-top: 1rem; display: flex")
}

type todoComponentProps struct {
	onDelete    func() error
	onChange    func() error
	currentTodo todo
}

func todoComponent(_ context.Context, props todoComponentProps, _ nodes.Children) nodes.Child {
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

type todosAppProps struct {
	state appState
}

func todosApp(ctx context.Context, props todosAppProps, _ nodes.Children) nodes.Child {
	currentState := props.state

	updateTodo := func(todoId int, completed bool) error {
		return store.SetState(ctx, func(currentState appState) appState {
			newTodos := make([]todo, len(currentState.todos))

			for i, current := range currentState.todos {
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

			return appState{
				todos:      newTodos,
				currentAdd: currentState.currentAdd,
			}
		})
	}

	deleteTodo := func(todoId int) error {
		return store.SetState(ctx, func(currentState appState) appState {
			newTodos := make([]todo, len(currentState.todos)-1)

			count := 0
			for _, current := range currentState.todos {
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

			return appState{
				todos:      newTodos,
				currentAdd: currentState.currentAdd,
			}
		})
	}

	todosComponents := make(nodes.Children, len(currentState.todos))

	for i, todo := range currentState.todos {
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
			lander.Component(store.Consumer, state.ConsumerProps[appState]{
				Render: func(currentState appState) nodes.Child {
					return lander.Component(addTodoForm, addTodoFormProps{
						state: currentState,
					}, nodes.Children{})
				},
			}, nodes.Children{}),
		}).Style("max-width: 300px;"),
	}).Style("padding: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(store.Consumer, state.ConsumerProps[appState]{
			Render: func(currentState appState) nodes.Child {
				return lander.Component(todosApp, todosAppProps{state: currentState}, nodes.Children{})
			},
		}, nodes.Children{}),
		"#app",
	)
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
