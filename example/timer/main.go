package main

import (
	"fmt"
	"time"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	js "github.com/minivera/go-lander/go-wasm-dom"
	"github.com/minivera/go-lander/nodes"
)

type timerApp struct {
	env *lander.DomEnvironment

	run  bool
	time float64
	last time.Time
}

func (a *timerApp) runTimer() {
	if a.run {
		js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) any {
			// Track with a time and not a pure floats since renders are synchronous. The time between
			// timeouts and the real time might be different.
			a.time += time.Now().Sub(a.last).Seconds()
			a.last = time.Now()
			if err := a.env.Update(); err != nil {
				panic("something went wrong with the time")
			}

			a.runTimer()
			return nil
		}), 10)
	}
}

func (a *timerApp) render(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Html("h1", nodes.Attributes{}, nodes.Children{
			lander.Text("Sample timer app"),
		}),
		lander.Html("div", nodes.Attributes{}, nodes.Children{
			lander.Html("div", nodes.Attributes{}, nodes.Children{
				lander.Text(fmt.Sprintf("%.2fs", a.time)),
			}).Style("padding-right: 1rem; font-family: Courier New,Courier,Lucida Sans Typewriter,Lucida Typewriter,monospace;"),
			lander.Html("button", nodes.Attributes{
				"click": func(*events.DOMEvent) error {
					if a.run {
						return nil
					}

					a.last = time.Now()
					a.run = true

					a.runTimer()
					return nil
				},
			}, nodes.Children{
				lander.Text("Start"),
			}),
			lander.Html("button", nodes.Attributes{
				"click": func(*events.DOMEvent) error {
					a.run = false
					return nil
				},
			}, nodes.Children{
				lander.Text("Stop"),
			}),
			lander.Html("button", nodes.Attributes{
				"click": func(*events.DOMEvent) error {
					a.time = 0.0
					return a.env.Update()
				},
			}, nodes.Children{
				lander.Text("Reset"),
			}),
		}).Style("display: flex;"),
	}).Style("padding: 1rem;")
}

func main() {
	c := make(chan bool)

	app := timerApp{}

	env, err := lander.RenderInto(
		lander.Component(app.render, nodes.Props{}, nodes.Children{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	app.env = env

	<-c
}
