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

func addTodoForm(ctx context.Context, props nodes.Props, _ nodes.Children) nodes.Child {
	currentState, ok := props["state"].(appState)
	if !ok {
		panic("addTodoForm expects the state store as its state prop")
	}

	return lander.Html("div", nodes.Attributes{}, []nodes.Child{
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
		}, []nodes.Child{}).Style("margin-right: 1rem;"),
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
		}, []nodes.Child{
			lander.Text("Add"),
		}),
	}).Style("margin-top: 1rem; display: flex")
}

func todoComponent(_ context.Context, props nodes.Props, _ nodes.Children) nodes.Child {
	onDelete, ok := props["onDelete"].(func() error)
	if !ok {
		panic("todoComponent expects a function as its onDelete prop")
	}

	onChange, ok := props["onChange"].(func() error)
	if !ok {
		panic("todoComponent expects a function as its onDelete prop")
	}

	currentTodo, ok := props["todo"].(todo)
	if !ok {
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

func todosApp(ctx context.Context, props nodes.Props, _ nodes.Children) nodes.Child {
	currentState, ok := props["state"].(appState)
	if !ok {
		panic("todosApp expects the state store as its state prop")
	}

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

	todosComponents := make([]nodes.Child, len(currentState.todos))

	for i, todo := range currentState.todos {
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
			lander.Component(store.Consumer, nodes.Props{
				"render": func(currentState appState) nodes.Child {
					return lander.Component(addTodoForm, nodes.Props{
						"state": currentState,
					}, []nodes.Child{})
				},
			}, []nodes.Child{}),
		}).Style("max-width: 300px;"),
	}).Style("padding: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(store.Consumer, nodes.Props{
			"render": func(currentState appState) nodes.Child {
				return lander.Component(todosApp, nodes.Props{"state": currentState}, []nodes.Child{})
			},
		}, []nodes.Child{}),
		"#app",
	)
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
