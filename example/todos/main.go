package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/nodes"
)

type addTodoForm struct {
	env *lander.DomEnvironment

	value string
}

func (a *addTodoForm) render(props lander.Props, _ lander.Children) lander.Child {
	onAdd, ok := props["onAdd"].(func(value string) error)
	if !ok {
		fmt.Println("addTodoForm expects a function as its onAdd prop")
		// TODO: This is pretty terrible, improve. Maybe make props a struct?
		panic("addTodoForm expects a function as its onAdd prop")
	}

	return lander.Html("div", map[string]interface{}{}, []lander.Child{
		lander.Html("input", map[string]interface{}{
			"value": a.value,
			"change": func(event *lander.DOMEvent) error {
				a.value = event.JSEvent().Get("target").Get("value").String()
				return a.env.Update()
			},
		}, []nodes.Child{}).Style("margin-right: 1rem;"),
		lander.Html("button", map[string]interface{}{
			"click": func(*lander.DOMEvent) error {
				return onAdd(a.value)
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

type todosApp struct {
	env *lander.DomEnvironment

	todos []todo
	form  addTodoForm
}

func (a *todosApp) updateTodo(todoId int, completed bool) {
	todos := make([]todo, len(a.todos))

	for i, current := range a.todos {
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

	a.todos = todos
}

func (a *todosApp) deleteTodo(todoId int) {
	todos := make([]todo, len(a.todos)-1)

	count := 0
	for _, current := range a.todos {
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

	a.todos = todos
}

func (a *todosApp) addTodo(name string) {
	todos := make([]todo, len(a.todos))

	for i, current := range a.todos {
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

	a.todos = todos
}

func (a *todosApp) render(_ lander.Props, _ lander.Children) lander.Child {
	fmt.Printf("Todos are %v\n", a.todos)
	todos := make([]nodes.Child, len(a.todos))

	for i, todo := range a.todos {
		todos[i] = lander.Html("li", map[string]interface{}{}, []nodes.Child{
			lander.Html("div", map[string]interface{}{}, []nodes.Child{
				lander.Html("input", map[string]interface{}{
					"type":    "checkbox",
					"checked": todo.completed,
					"change": func(*lander.DOMEvent) error {
						a.updateTodo(todo.id, !todo.completed)
						return a.env.Update()
					},
				}, []nodes.Child{}),
				lander.Html("strong", map[string]interface{}{}, []nodes.Child{
					lander.Text(todo.name),
				}),
			}).Style("display: inline-flex; align-items: center; padding-right: 1rem;"),
			lander.Html("button", map[string]interface{}{
				"click": func(*lander.DOMEvent) error {
					a.deleteTodo(todo.id)
					return a.env.Update()
				},
			}, []nodes.Child{
				lander.Text("X"),
			}).Style("display: inline;"),
		})
	}

	return lander.Html("div", map[string]interface{}{}, []lander.Child{
		lander.Html("h1", map[string]interface{}{}, []nodes.Child{
			lander.Text("Sample todo app"),
		}),
		lander.Html("div", map[string]interface{}{}, []nodes.Child{
			lander.Html("h2", map[string]interface{}{}, []nodes.Child{
				lander.Text("Todos"),
			}),
			lander.Html("ul", map[string]interface{}{}, todos).Style("margin-top: 1rem;"),
			lander.Component(a.form.render, map[string]interface{}{
				"onAdd": func(value string) error {
					a.addTodo(value)
					// Reset the form's state
					a.form = addTodoForm{
						env: a.env,
					}
					return a.env.Update()
				},
			}, []nodes.Child{}),
		}).Style("max-width: 300px;"),
	}).Style("padding: 1rem;")
}

func main() {
	c := make(chan bool)

	app := todosApp{
		todos: []todo{
			{
				id:        0,
				name:      "write more examples",
				completed: false,
			},
		},
		form: addTodoForm{},
	}

	env, err := lander.RenderInto(
		lander.Component(app.render, map[string]interface{}{}, []lander.Child{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	app.env = env
	app.form.env = env

	<-c
}
