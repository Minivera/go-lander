# GO-Lander
[![Language: Go](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/minivera/go-lander)
[![Go Report Card](https://goreportcard.com/badge/github.com/minivera/go-lander)](https://goreportcard.com/report/github.com/minivera/go-lander)
[![Go Reference](https://pkg.go.dev/badge/github.com/minivera/go-lander.svg)](https://pkg.go.dev/github.com/minivera/go-lander)
![GitHub](https://img.shields.io/github/license/minivera/go-lander)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/minivera/go-lander)

GO-Lander is an experimental library to help you build frontend applications with the power of Golang and Web Assembly
(WASM). It is heavily inspired by JavaScript frontend libraries such as React and aims to provide the smallest
footprint possible.

---

**GO-Lander is _very_ unstable and not tested. If you decide to use it in a production application, please be aware
that the API might change as we work on it. The lack of tests also means that there are likely undetected bugs in
the diffing and mounting algorithms, which could cause issues in your apps. Please check the [example](./example)
directory to see what kind of behavior we currently support.**

---

## Why another WASM frontend framework?

GO-Lander was built as a passion project and a way to understand the Golang WASM implementation. It started as a clone
of [Vugu](https://github.com/vugu/vugu), with the goal of implementing an API closer to React's and experiment with the
virtual DOM, alongside [lander](https://github.com/Minivera/lander).

It evolved in a genuine attempt to offer a different experience from the existing libraries, but it still retains
its experimental and passion project nature. Our goal is not to become the next replacement for React, but rather to
offer something people can experiment with.

## Why use GO-Lander

You _shouldn't_, it's far too unstable and buggy in its current state. But if you decide to use GO-Lander, the
library has a few things to offer:

1. **Lightweight**. It is intended to have the lightest footprint possible in your application, the APIs
   are straightforward and the amount of features has been purposefully kept to a minimum. Any additional
   "nice-to-have" features can be ignored.
2. **Native**. We're using the native features of Golang to our advantage to provide a near native application
   experience, you should feel right at home. Rather than to clone React in Go, we've attempted to provide an
   experience as close to React as possible without moving away from the experience of writing Go code, especially
   when it comes to React's more "magical" features (I.E. features that work even though the common JavaScript
   knowledge would say otherwise).
3. **Component-centric**. Components can be composed to create large applications, giving you the ability to
   separate the concerns of your application in individual components, like you would for any other Golang application.
4. **Expandable**. Our goal is to keep the library simple and open, we give you the base tools to build anything
   you'd like. None of the "advanced/experimental" features such as our router use internal logic, they are all
   build with the same public APIs.

### GO-lander and Vugu

[Vugu](https://github.com/vugu/vugu) is a lot more production ready and mature than GO-Lander. Vugu is a lot closer
to Vue.js, it expects you to write `.vugu` files and your application logic in an HTML `script` tag. If you enjoy the
Vue.js experience, Vugu is likely to be much more usable. Since GO-lander is closer to vecty in terms of its
developer experience, we recommend looking
at [its own comparison to Vugu](https://github.com/hexops/vecty#vecty-vs-vugu) to make your decision.

### GO-lander and vecty

[Vecty](https://github.com/hexops/vecty) is closer to how GO-lander was designed with some key differences. Vecty
expects components to always be struct pointers with a `Render` method. GO-Lander allows you to choose how you want
to structure your app, we store components as functions, regardless of whether they are methods of a struct or not.
Vecty
also differs from GO-lander in the algorithm it uses for diffing. Vecty uses a high-performance algorithm similar to
the virtual DOM pattern where Go-lander uses a pure virtual DOM tree with a separate diffing and patching process.

## Getting started

Install the latest version of GO-lander to get started.

```bash
go get github.com/Minivera/go-lander@latest
```

Once installed, create a main file for your application, let's start with a "Hello, World!".

```go
package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

func helloWorld(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Html("h1", nodes.Attributes{}, nodes.Children{
		lander.Text("Hello, World!"),
	}).Style("margin: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(helloWorld, nodes.Props{}, nodes.Children{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
```

For a web assembly application to stay alive after its first execution, we need to create a channel and wait on it.
This ensures our application stays alive until the page is refreshed, which allows GO-lander to rerender the app on
demand.

We then call the `RenderInto` function using a component called `helloWorld` on the `#app` DOM node, which we'll
create later. GO-lander apps need to start with a component, which can then render any structure you want.

In this case, the `helloWorld` component is a function that renders a `h1` tag with the text "Hello, World!" and
some styling. Any styles you add to an HTML node with `Style` will be added to a global stylesheet and updated on
every subsequent render.

GO-lander does not provide any utilities to compile your application to Web Assembly, so we'll also need an `index.
html` file and the JavaScript file provided by your installation of Go. You may also use [`wasmserve`]
(https://github.com/hajimehoshi/wasmserve) to simplify this process.

First, we need an HTML file. We can copy the one from the official [Web Assembly wiki](https://github.
com/golang/go/wiki/WebAssembly#getting-started). Add a `<div id="app"></div>` tag to the body, the final file should
look like this.

```html

<html>
<head>
    <meta charset="utf-8"/>
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        });
    </script>
</head>
<body>
<div id="app"></div>
</body>
</html>
```

Then, copy the official Go JavaScript glue file to your current directory, this is needed to bridge the gap between
GO-lander and the JavaScript DOM environment.

```bash
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

With all that in hand, compile your application to WASM and serve it.

```bash
GOOS=js GOARCH=wasm go build -o main.wasm
# install goexec: go get -u github.com/shurcooL/goexec
goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'
```

Open your browser to `http://localhost:8080` and you should see the "Hello, World!" message. To deploy your
application, upload the `main.wasm`, `index.html`, and `wasm_exec.js` to any server capable of serving HTML files.

### Examples

For more advanced examples, check out the [example directory](./example). We use `mage` to build and run the
examples, install it with:

```bash
go install github.com/magefile/mage@latest
```

Our `magefile` contains all build targets, feel free to check them out using the command `mage` in this directory.
To run any example, use this command and replace `<example name>` with the directory name of the example you want to
run.

```bash
mage runExample <example name>
```

You can also run all examples in a single app using:

```bash
mage runExampleViewer
```

## Detailed documentation

Any GO-lander starts with the lander environment, which is created when you render an application into a DOM node.
This environment takes care of the entire render and diffing process, but you are responsible for keeping it alive
and handling any updates.

```go
env, err := lander.RenderInto(someComponent, 'Some DOM node selector')
```

When this function is called, lander will query for a DOM node in the document using the provided selector and mount
your component into it. `RenderInto` can only render components, every GO-lander app must start with a component.

Once properly mounted, the lander environment is returned. If any error happened (such as the DOM node not being
found or the library running in a server environment), an error will be returned alongside a `nil` environment.

The environment exports a single method, `Update()`. Whenever the state of your application changes, or you want to
trigger a rerender of your application, call `env.Update()`. This will rerender your application, diff it against
the previous application, and update the DOM with the changes, if any.

Both `Update` and `RenderInto` are safe to execute in parallel, only one execution of either can run at a time.
Subsequent calls will need to wait until the current render is done before they can update the tree. This means that
any update you trigger are executed in order, but also that the entire tree must have been fully mounted before you
can update. Neither are batched, if `Update` is called is quick succession, each update will execute sequentially.

### Writing and styling HTML elements

To write the HTML structure of your app, use one of the three functions provided.

- **lander.HTML(tag, attributes, children)** will create an HTML node with the given tag, which can be any valid
  HTML tag, including web components.
- **lander.SVG(tag, attributes, children)** is an override of `lander.HTML` which adds the `http://www.w3.
  org/2000/svg` to the node and creates it on the document using that namespace, which ensures that SVG tags will
  render properly.
- **lander.Text(text)** will create a text node, which can be used to add text inside any node.

HTML nodes take a `nodes.Attributes` map as their attributes parameter, which can include any valid HTML attribute
or property for the element. GO-lander will extract the attributes and assign them using the following rules:

```go
package main

import (
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

lander.HTML('div', nodes.Attributes{
    "checked": true,              // Will be converted to `checked=""` on the DOM element if true and omitted if false.
    "placeholder": "some string", // Will be kept as a string and assigned on the DOM element.
    "value": -1, // Will be converted to as string and assigned on the DOM element.
    "click": func (event *events.DOMEvent) error {} // Will be assigned using `addEventListener` on the DOM element.
}, nodes.Children{});
```

Any other type is ignored. If an attribute has an associated property on the DOM node (such as `value`), it will
also be set as a property. See
the [content vs. IDL attributes](https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes#content_versus_idl_attributes)
reference for more details.

Event listeners must be defined as either a function following the type `nodes.EventListenerFunc` or as a function
with the exact signature `func(event *events.DOMEvent) error`. Event listeners are assigned with their HTML name,
see the [event reference](https://developer.mozilla.org/en-US/docs/Web/Events) for a list of events and their names.

An event listener is called with a `events.DOMEvent` parameter, which includes three helpful methods.

- `event.JSEvent()` returns the actual JavaScript event triggers as a `js.Value`.
- `event.JSEventThis()` returns the `this` value of the event listener, if any.
- `event.PreventDefault()` triggers the prevent default call on the event, which will block any of the event's
  default behavior, such as form submits.

HTML and SVG nodes may also take a slice of children. These children can be any of the node returned by the `lander`
node factories, such as component nodes, text nodes, fragment nodes, or other HTML nodes. The `nodes.Children` type
is provided to reduce the complexity of the code when defining the slice of children.

`lander.HTML` nodes provide some styling capabilities through their `Style` and `SelectorStyle` methods. `Style`
takes any valid CSS definition and will assign it to the HTML node on render, dynamically generating class names and
a CSS file to assign in the document's `<head>`. Calling `Style` multiple times will override the previous styling
definition. Rather, if you wish to further define your CSS styles, `SelectorStyle` can be used to add styles with an
option selector, for example:

```go
lander.
    HTML('div', nodes.Attributes{}, nodes.Children{
        lander.HTML('input', nodes.Attributes{}, nodes.Children{}).
    }).
    // Will generate a random CSS class and create this CSS style in the head
    // .classname { color: red; margin: 10px }
    Style("color: red; margin: 10px").
    // Will use the previously generated CSS class and create this CSS style in the head
    // .classname input { color: red; margin: 10px }
    SelectorStyle("input", "width: 80%")
```

### Components

Components are the core of GO-lander's component pattern. In their simplest of forms, a component is a function that,
given specific inputs, returns a single lander node.

Every component's signature must match the signature of `nodes.FunctionComponent`, namely:

```go
type FunctionComponent[T any] func (ctx context.Context, props T, children nodes.Children) nodes.Child
```

The `context.Context` is covered in a later section. The component's `Props` are, unlike HTML attributes, a generic
type that can be assigned any type. The `nodes.Props` type if provided as a utility for components that take no 
props. The properties are never converted by GO-lander, they will be passed along to the function. When GO-lander 
renders the tree, it calls each component in the tree to generate its result (called a render).

You are responsible for writing your components in any way you want, GO-lander only provides you with a mean to
return a child and have it appear in the updated DOM tree.

As an example, let's explore the "Hello, World!" example with some properties.

```go
package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

type helloWorldProps struct {
	message string
}

func helloWorld(_ context.Context, props helloWorldProps, children nodes.Children) nodes.Child {
	message := props.message // Lander does not check the types, make sure you validate them yourself

	return lander.Html("h1", nodes.Attributes{}, nodes.Children{
		lander.Text(message),
		children[0],
	}).Style("margin: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		// `Component` is a generic function, the type is infered from the types of the props
		lander.Component(helloWorld, helloWorldProps{
			message: "Hello, World!"
		}, nodes.Children{
			lander.Text("From the lander README"),
		}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
```

When creating a component node through `lander.Component`, the `props` passed will be given to the component's
function, as will the children. If the component returns another component as part of its rendered node's descendant,
that component will in-turn be rendered until all components in the tree have been rendered. This is the GO-lander's 
render cycle in a nutshell.

GO-lander does not provide any type safety for props, they are passed to GO-lander's internal logic without any type 
information. You are responsible for checking the types of the values you receive, and to `panic` if the types are 
not correct or a required prop is missing.

Any component may return `nil` instead of a `nodes.Child`. The component's render result will be ignored and any
previous node it rendered will be removed. The component may return a valid node in a later render cycle, the result
will then be inserted into the DOM.

### Fragment nodes

Returning only a single child from a component may not always be practical. If a component renders a slice of
children, for example, you would need to wrap that slice in an HTML node such as a `div` to only return a single
child. GO-lander provides a final type of node to solve this problem, the `lander.Fragment` node.

A fragment takes a slice of nodes as its children and will render them as if they were a direct child of the
closest HTML node. Fragments are also very useful when you want to assign a slice of children to a node that already
has children defined, in place where the spread operator does not work. Taking our "Hello, World" example from above,
we could rewrite it with fragments like this:

```go
package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

type helloWorldProps struct {
	message string
}

func helloWorld(_ context.Context, props helloWorldProps, children nodes.Children) nodes.Child {
	return lander.Html("h1", nodes.Attributes{}, nodes.Children{
		lander.Text(props.message),
		// But now the fragment will take care of the children
		lander.Fragment(children),
	}).Style("margin: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(helloWorld, helloWorldProps{
			message: "Hello, World!"
		}, nodes.Children{
			lander.Text("From the lander README"),
			lander.Text("This node would previously have been ignored!"),
		}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
```

### Keeping state

Since GO-lander does not make any assumptions in how you want to structure your app, it also does not provide any
in-depth state management solutions built-in. Rather, you can use the native features of Go to store state for your
application with global variables or structs for example.

The most straightforward way of storing state is with struct components. To create a `struct` component, define a
`struct` type with the state values you want to manage, and create your component as a method of the `struct`. As
long as the struct is kept alive in your application and managed, your state will be tracked. For example, let,s
build a log-in form.

```go
package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

type loginForm struct {
	// We keep a reference to the environment to allow updating on form changes
	env *lander.DomEnvironment

	username string
	password string
}

// The form should be a pointer
func (f *loginForm) render(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Html("h1", nodes.Attributes{}, nodes.Children{
			lander.Text("Log into our app"),
		}),
		// We would likely want to also handle the form submit here
		lander.Html("form", nodes.Attributes{}, nodes.Children{
			lander.Html("label", nodes.Attributes{
				"for": "username",
			}, nodes.Children{
				lander.Text("Username"),
			}).Style("font-weight: bold;"),
			lander.Html("input", nodes.Attributes{
				"name":        "username",
				"placeholder": "Enter Username",
				// On change, we assign the value of username to the input value, then rerender
				"change": func(event *events.DOMEvent) error {
					f.username = event.JSEvent().Get("target").Get("value").String()
					return f.env.Update()
				},
			}, nodes.Children{}),
			lander.Html("label", nodes.Attributes{
				"for": "password",
			}, nodes.Children{
				lander.Text("Password"),
			}).Style("font-weight: bold;"),
			lander.Html("input", nodes.Attributes{
				"name":        "password",
				"placeholder": "Enter Password",
				"type":        "password",
				// On change, we assign the value of password to the input value, then rerender
				"change": func(event *events.DOMEvent) error {
					f.password = event.JSEvent().Get("target").Get("value").String()
					return f.env.Update()
				},
			}, nodes.Children{}),
			lander.Html("button", nodes.Attributes{
				"type": "submit",
			}, nodes.Children{
				lander.Text("Submit"),
			}).Style("margin-top: 1rem;"),
		}),
	}).Style("margin: 1rem;")
}

func main() {
	c := make(chan bool)

	form := &loginForm{}

	env, err := lander.RenderInto(
		// Render the form's render method instead of a function component
		lander.Component(form.render, nodes.Props{}, nodes.Children{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	form.env = env

	<-c
}
```

In this example, the `loginForm` `struct` stored our two state value, and a pointer to the GO-lander environment,
which allows us to trigger an update when the form value changes. Since the form is created in the `main` function,
it will be alive and keep our state for the entire lifecycle of the application.

A drawback of this approach is that any subsequent struct components will need to be stored either in the `main`
function, or as a property on another struct. This creates more management for your application's state, but it also
makes your state management less global and less "magical". For example, if our login form was a subcomponent
somewhere in the app, we could instead store it in another struct in a sort of tree of struct components.

```go
package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

type app struct {
	// We keep a reference to the environment to allow updating on changes
	env *lander.DomEnvironment

	loginForm *loginForm
}

// The form should be a pointer
func (a *app) render(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Html("div", nodes.Attributes{}, nodes.Children{
		lander.Component(a.loginForm.render, nodes.Props{
			"onSubmit": func(username, password string) error {
				// Do something

				// Reset the form's state
				a.loginForm = &loginForm{
					env: a.env,
				}
				return f.env.Update()
			},
		}, nodes.Children{})
	})
}

func main() {
	c := make(chan bool)

	app := &app{
		// Create our initial state for the form
		loginForm = &loginForm{}
	}

	env, err := lander.RenderInto(
		// Render the form's render method instead of a function component
		lander.Component(app.render, nodes.Props{}, nodes.Children{}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	form.env = env
	form.loginForm.env = env

	<-c
}
```

### Context

Every component function takes a GO-lander context as its first parameter. This context is similar in concept to
Golang's own [context](https://pkg.go.dev/context) and
React's [context](https://react.dev/learn/passing-data-deeply-with-context).
This context carries over data that can be accessed anywhere in the tree, it allows you to define global data that
all components in the tree can consume without needing to pass it as props throughout the entire tree.

As an example, consider a deeply nested series of components that you want to theme using a central theme. If you
wanted to carry over the theme to all components, you would need to pass it as props to all components in the tree.

```go
func FirstComponent(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
    someTheme := createTheme()
    
    return lander.Component(secondComponent, SecondComponentProps{
        theme: someTheme
    }, nodes.Children{})
}

type SecondComponentProps struct {
	
}

func SecondComponent(_ context.Context, props SecondComponentProps, _ nodes.Children) nodes.Child {
    someTheme := props.theme
    
    return lander.Component(thirdComponent, thirdComponentProps{
        theme: someTheme
    }, nodes.Children{})
}

// And so on...
```

This example might be solved by making the theme into a global struct you can import from somewhere in your app.
Another solution, depending on how you want to structure your application or library, is to store this component in
the context for all to access. It works similarly to the global struct alternative, but it instead lives inside the
application's render cycle. Let's rewrite the above example with context.

```go
func FirstComponent(ctx context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
    someTheme := createTheme()
    
    if !ctx.HasValue("theme") {
        ctx.SetValue("theme", theme)
    }
    
    return lander.Component(secondComponent, nodes.Props{}, nodes.Children{})
}

func SecondComponent(ctx context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
    someTheme := ctx.GetValue("theme").(Theme)
    
    return lander.Component(thirdComponent, nodes.Props{}, nodes.Children{})
}

// And so on...
```

The context struct provides three methods to access and set data:

- `HasValue(valueName string) bool` returns if there is a stored value under the given key in the context.
- `GetValue(valueName string) interface{}` returns the value under the given key in the context, may panic if the
  value does not exist. We recommend checking with `HasValue` to validate that the value can be found. This method
  returns a generic interface, you will need to type-case it to your type.
- `SetValue(valueName string, value interface{})` sets the value under the given key in the context, you can save
  any type in the context. This will automatically set all components for rerender in the next update.

A component that provides a context value through `SetValue` is called a provider. Once a provider is added to the
app, every other component can now see the added context value through `SetValue` (see the limitations section below)
. We recommend you add all providers at the root of your application to ensure that your app always has its context
set properly.

#### Effect on component diffing

By default, components only rerender during the diffing process if the node's change in a significant way, which
could be through props changes, children changes, or the component function changing. If the component stays the
same, but the context changes, then the component would not rerender with the new context value. For this reason,
any change to the context will set all components as "different" and trigger a rerender.

We recommend that providers do not change the value they set in context on every render. Rather, check if the value
already exists and is already set to the value you want to be set to. This can be done with `HasValue`, like in the
example below.

```go
func ThemeProvider(ctx context.Context, _ nodes.Props, children nodes.Children) nodes.Child {
    someTheme := createTheme()
    
    if !ctx.HasValue("theme") {
        // We might also want to check if the theme has changed here, for example if we provide
        // a dark mode version.
        ctx.SetValue("theme", theme)
    }
    
    return lander.Fragment(children)
}
```

Using a global theme object instead of the context in this same example would not set the components to be
rerendered when the theme changes. This is the main difference between using context and a global struct.

#### Limitations

There are a few key differences between the GO-lander context and both of its inspirations. First, GO-lander's
context does not provide any cancellation or timeout mechanisms like Golang's context. We have plans to support
render cancellation through the context to allow component to stop the diffing process from checking their result,
but this is not yet supported. Second, GO-lander's context is a single struct that is carried over the entire tree. It
is only defined once at the start of the rendering process. This differs from React's own context, as described below.

```
# In React, the context is redefined only for the descendants of any component that changes the content of the context
# for example:
App
| Some context provider -> Provides context value "foo"
| | Some child component with provider -> Consumes context value "foo", Provides context value "bar"
| | | Some descendant component -> Consumes context value "bar"
| | Some other child component -> Consumes context value "foo"
```

When a context provider redefines the context, only the descendants of that component see the new value. Other
components in the tree will see the previous value in the context.

GO-lander's context is global to the tree and is a pointer, thus any change to a context value in the tree will
affect all components, regardless of their location. The tree is visited in order, so this behavior is predictable.
Taking the same example as above, but rewriting it in GO-lander, the context would behave as described below.

```
# In GO-lander, the context is a pointer that always points to the same value regardless of positions
App
| Some context provider -> Provides context value "foo"
| | Some child component with provider -> Consumes context value "foo", Provides context value "bar"
| | | Some descendant component -> Consumes context value "bar"
| | Some other child component -> Consumes context value "bar" _different from React_
```

This limitation is important to keep in mind when you define data in the context. As soon as `SetValue` is called,
the data is available globally to all components in the tree, even components that have already been rendered.

### Lifecycle listeners

The context object also provides a set of three lifecycle listeners, which can be used to take actions when specific
things happen to your components. All three listener types take a `func() error` as their only parameter. This
function will be executed when the listening even happens. Return an error only if something critical should happen,
this will cause the entire app to stop.

- `ctx.OnMount` listens for the first time the component has been mounted, I.E. added to the DOM tree. This will
  only happens once for components and is called _after_ the component has been mounted. By this point, the DOM nodes
  of its child have been added to the tree and can be accessed. Components are reused in the tree, which may lead to
  different mounts that you would expect. See example below.
- `ctx.OnRender` listens for any full render of the component. This event happens whenever the component has changed
  in a meaningful way, which triggers a diff. This event happens _after_ the render has happened, meaning that any
  HTML nodes it returns have been updated in the DOM. This will also fire on first mount, but will not fire on unmount.
- `ctx.OnUnmount` listens for an unmount event on the component, I.E. when it is removed from the DOM tree. This
  event triggers only once _after_ the component has been removed and unmount. By this point, any child it returned
  have been removed the DOM tree. Components are reused in the tree, which may lead to different unmounts that you
  would expect. See example below.

Lifecycle listeners should be called directly in the render function, they will trigger based on the chosen event.

```go
func app(ctx context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
    ctx.OnMount(func () error {
    // do something on mount	
    })
    
    ctx.OnRender(func () error {
    // do something on render	
    })
    
    ctx.OnUnmount(func () error {
    // do something on unmount	
    })
    
    return ...
}
```

As mentioned above, components are reused. For example, if you return a list of components from a tree and remove an
element in the middle of that list, the unmounted component will be the last element of that list, not the one that
was "removed" from the developer or user perspective. For example, consider this tree:

```
App
| Todo 1
| Todo 2
| Todo 3
```

If you remove `Todo 2`, `Todo 2` will be reused and updated to the values of `Todo 3`, and `Todo 3` will trigger an
unmount. This is different from other libraries that use a "keyed" approach to lists. GO-lander does not use keys
and instead rely on node reuse. Let's visualize this:

```
App
| Todo 1
| Todo 2 <- Will now use the component function, props, and children of Todo 3
| Todo 3 <- Will unmount

# After the render cycle
App
| Todo 1
| Todo 3
```

## Experimental features

We have built a few experimental features that bridge the gap between other, more feature-rich, libraries and the
minimalist approach of GO-lander. They are very experimental (even more than the library itself), and may be subject to
change or may break in unexpected ways.

### Hooks

[React hooks](https://react.dev/reference/react) have changed the way we build frontend apps, and like many other
frontend trying to compete in the
landscape dominated by React, we have also built an alternative to hooks inside GO-lander. This experiment still
follows the core principles of GO-lander, but given the nature of the hooks api, it required some hidden magic to
properly work.

We have implemented two of the common hooks in React, `useState` and `useEffect`. Both use an internal version of
the `useMemo` hook, which is currently not available publicly.

To use the hooks API, you _must_ wrap your entire application inside the `hooks.Provider` component. While this is
not strictly required logic wise, it does help encapsulate the logic. Future versions may remove this provider.

To see the hooks in action, let's look at the [API fetch with hooks example](./example/fetchAPIWithHooks/main.go).

```go
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
		lander.Component(hooks.Provider, nodes.Props{}, nodes.Children{
			lander.Component(fetchApp, nodes.Props{}, nodes.Children{}),
		}),
		"#app",
	)
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
```

The `hooks.UseState` hook uses a generic type to do the typecasting for you, the default value passed as its second
parameter must match the generic type. The hook always returns three variables, always typed to the generic type you
provided.

- The first parameter is the current value of the state, which will be the default value on first render. Due to how
  scoping works in Go, you should not use this value in event listeners or code that is not executed in the same
  function as the hook.
- The second parameter is the state setter function, which takes a function as its argument when called. You cannot
  set the state value by passing it to the state setter, it must be a return value of a function passed to
  `setState`. This is to ensure that you always have the latest value of the state if you need to update it. Setting
  the state also automatically triggers a rerender.
- The third parameter is a state getter function. Since scoping in Go means that anonymous functions may have a
  different version of the state value depending on when they're called, this function ensures you will always get
  the most up-to-date value if you need it. This is not needed if your state value is a pointer.

The `hooks.UseEffect` hook takes a function as its second parameter, which in turn must return an error and a
cleanup function. The effect function given is called on mount, and on any subsequent render if and only if the
dependencies slice given as its third parameter change. If you do not want the hook to rerender, pass an empty or nil
slice.

The effect can return `nil`, or another function as its cleanup. This cleanup is automatically called on unmount,
which allows you to clean any asynchronous code before the component gets unmounted.

All hooks must be given the context object of the function calling them as its first parameter. All memoized values
are saved in the context, meaning that components in an application using hooks will always rerender and cannot be
optimized. This should have no effect on your app's performance, but it worth considering when looking at this
experimental feature.

### Global state management

The alternative to using the global context for state management across an entire app is to use a global struct or
store of some kind. To make that experience easier on developers, we've built a very basic version of a global state
store.

To use the store, create a package in your application and export a newly created store containing your state, for
example:

```go
package state

import "github.com/minivera/go-lander/experimental/state"

type appState struct {
	// Some state
}

var Store = state.NewStore[appState](appState{
	// State default values
})
```

That store exports two methods, which can be used to set or consume state. `Store.SetState` sets the entire
stored state to the new value. It expects the context as its first parameter, and a setter function as its second,
which should have this signature; `func(value T) T` (where `T` is the generic type given to the store). This setter
will provide the current value of the store and expects a new value. We strongly recommend creating an entirely new
value when setting the state. The app will automatically update once the state has been set.

`Store.Consumer` is a component you can use in your app to inject your state into another component. It takes a
`Render` prop, which should provides the current value of the store and must return a valid GO-lander child, like
any other component.

```go
lander.Component(store.Consumer, store.ConsumerProps{
    Render: func (currentState appState) nodes.Child {
    // Return some nodes based on the state
    },
}, nodes.Children{}),
```

The consumer component takes care of any rerendering it needs to process. At the moment, it will always rerender
even if the state has not changed between updates, which differs from more stable state management libraries.

### In-memory routing

Since WASM applications are not easily made aware of the current URL in the browser, or can easily access the
`history` API to modify it, we have built this experimental set of components and utilities to help you create a
single-page application. Please note that this only supports client-side routing, you will need to handle serving
your application under any route. The examples provided in this repository do not handle routing to any other URL
than `/`.

This experiment offers a Regex based, in-memory, router. This means that routes are defined as regular expressions,
including parameters. We do not provide any utilities to convert more traditional paths (like `/users/:username`) to
regular expressions, at the moment.

To get started, create a package in your application and export a newly created router. This router should be
available to your entire application, as it provides all the components needed to properly handle routing.

```go
package routing

import "github.com/minivera/go-lander/experimental/router"

var Router = router.NewRouter()
```

Next, wrap your entire application inside a `Router.Provider` component. Routing uses the context to store the
current location information and to listen to any changes in the URL. The router will only update the context when
the location changes and will not impact your app's performance by setting the entire app to rerender on every update.

```go
package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/experimental/router"
	"github.com/minivera/go-lander/nodes"
)

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(appRouter.Provider, nodes.Props{}, nodes.Children{
			lander.Component(yourApp, nodes.Props{}, nodes.Children{}),
		}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
```

Your application is now URL-aware and can use in-app routing. We provide three ways of changing the current location of
your application.

1. `Router.Navigate(to string, replace bool)` will immediately navigate the user to the new location defined in `to`.
   If `replace` is set to `true`, the new history entry will be replaced instead of pushed, and the user may not use
   the back button to go to the previous location.
2. `Router.Link` is a component that renders a single `<a>` anchor element. Any children passed to the component will
   render inside the anchor. It can take two props, `to` and `replace`, which behave exactly like the `Navigate`
   parameters.
3. `Router.Redirect` is a component that immediately changes the location of the browser when it renders, which will
   trigger an update once the navigation is completed. It can take two props, `to` and `replace`, which behave
   exactly like the `Navigate` parameters.

The router also provides you with two components to conditionally render content based on the location.

`Router.Route` is an "on/off" component which will only render its children if its `Route` property matches the
current location. It uses a `Render` prop, which takes a function that receives the current match when the URl matches.
Let's see it in action.

```go
package main

import "github.com/minivera/go-lander/experimental/router"

func someApp(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Component(appRouter.Route, router.RouteProps{
		Route: "/app/([a-zA-Z0-9]+)/(?P<subroute>[a-zA-Z0-9]+)",
		Render: func(match router.Match) nodes.Child {
			// Only render if the URL matches the route when compiled to a regex
			// match.Pathname includes the actual URL location
			// match.Params["0"] has the first path param, which as an unnamed capture group
			// match.Params["subroute"])) has the second path param, which as named
		},
	}, nodes.Children{}),
}
```

The `Router.Route` will extract the relevant capture groups from `Route`, which work as path parameters for your 
routes. A `/users/:username` route in a more standard routing library would translate to 
`/users/(?P<username>[a-zA-Z0-9]+)` for example. The `Render` function is called with a match struct containing the
pathname and the path parameters. Any unnamed captured groups are stored in the `match.Params` map under the index
in the regex. For example, if we had `/users/(?P<username>[a-zA-Z0-9]+)/([a-zA-Z0-9]+)`, the second path param would
be stored under the index `"1"`.

The `Router.Route` can be chained to create a complex router. However, each route is checked on render and multiple
routes may match at the same time. To make sure only one route renders, use the `Router.Switch` component.

```go
package main

import "github.com/minivera/go-lander/experimental/router"

func someApp(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	lander.Component(appRouter.Switch, router.SwitchProps{
		Routes: router.RouteDefinitions{
			{"/$", func(_ router.Match) nodes.Child {
				// Home path
			}},
			{"/hello/test$", func(_ router.Match) nodes.Child {
				// Route with subpath
			}},
			{"/hello$", func(_ router.Match) nodes.Child {
				// Hello route, must be after the more specific route as it might match
				// the subroutes.
			}},
			{".*", func(match router.Match) nodes.Child {
				// Catch all, 404 route.
			}},
		},
	}, nodes.Children{}),
}
```

The `Router.Switch` component takes a set of `router.RouteDefinitions` as its single `Routes` prop. These
definitions are identical to the `Router.Router` props. The switch will check each route in order against the
current location and render the first match it finds. For this reason, you may want to have your more specific
routes appear before first level routes as they might match against sub-routes, as explained in the code above. The
`.*` catch all regex can be added at the end to render something if no route matches.

See more in the [routing example](./example/router/main.go).


### Head tags management (Helmet)

Managing the head of the document is a common pattern in JavaScript when in-app routing is introduced. To provide a 
good routing story to developers, we have built an experiment allowing some basic manipulation of the `head` tag of 
the document directly from GO-lander's tree.

To get started, wrap your entire application inside a `helmer.Provider` component. This provider will update the 
document's head on every render. The experiment does not check if updating the head is necessary at the moment and 
will always update when the tree rerenders. 

```go
package main

import (
	"fmt"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/experimental/helmet"
	"github.com/minivera/go-lander/nodes"
)

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(helmet.Provider, nodes.Props{}, nodes.Children{
			lander.Component(yourApp, nodes.Props{}, nodes.Children{}),
		}), "#app")
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
```

Your application can now provide new `head` tags, such as `title` as HTML nodes directly in the tree. To do so, use 
the provided `helmer.Head` component. Any HTML children passed to this component will be evaluated at render time, 
if the node's tag is one of `title`, `meta`, `link`, `script`, `noscript`, or `style`, it will be saved internally 
to be added to the head at the end of the render cycle.

```go
package main

import "github.com/minivera/go-lander/experimental/router"

func someApp(_ context.Context, _ nodes.Props, _ nodes.Children) nodes.Child {
	return lander.Component(helmet.Head, router.RouteProps{}, nodes.Children{
		lander.Html("title", nodes.Attributes{}, nodes.Children{
			lander.Text("Some title"),
		}),
		lander.Html("script", nodes.Attributes{}, nodes.Children{
			lander.Text("(function() { alert('test'); })();"),
		}),
		lander.Html("style", nodes.Attributes{}, nodes.Children{
			lander.Text("* { color: red; }"),
		}),
    }),
}
```

All tags will be rendered and updated when the tree updates like you would expect with a reactive library like 
GO-lander. The `title` tag is unique however as only one tag can exist in a page. `Helmet` prioritizes the last tag 
seen in the tree, given a walk from the "top" of your app's tree to the "bottom". For example:

```
# Titles are in order of priority, if 1 is removed, 2 will be selected and so on.
SomeApp
| Head with title = "4"
| div
| | Head with title = "3"
| div
| | div
| | | Head with title = "2"
| Head with title = "1" <- This will be the final title
```

See more in the [helmet example](./example/routerWithHelmet/main.go).

## Acknowledgements

Lander would not have been possible without the massive work done by the contributors of these libraries:

- [React](https://github.com/facebook/react)
- [Vecty](https://github.com/hexops/vecty)
- [Vugu](https://github.com/vugu/vugu) (Lander started out as a React-like clone of Vugu)
