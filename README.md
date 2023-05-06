# GO-Lander

GO-Lander is an experimental library to help you build frontend applications with the power of Golang and Web Assembly
(WASM). It is heavily inspired by JavaScript frontend libraries such as React and aims to provide the smallest
footprint possible.

---

**GO-Lander is _very_ unstable and not tested. If you decide to use it in a production application, please be aware
that the API might change as we work on it. The lack of tests also means that there are likely undetected bugs in
the diffing and mounting algorithms, which could cause issues in your apps. Please check the [example](./example)
directory to see what kind of behaviour we currently support.**

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
   experience as clode to React as possible without moving away from the experience of writing Go code, especially 
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
developer experience, we recommend looking at [its own comparison to Vugu](https://github.com/hexops/vecty#vecty-vs-vugu) to make your decision.

### GO-lander and vecty

[Vecty](https://github.com/hexops/vecty) is closer to how GO-lander was designed with some key differences. Vecty 
expects components to always be struct pointers with a `Render` method. GO-Lander allows you to choose how you want 
to structure your app, we store components as functions, regardless of it they are methods of a struct or not. Vecty 
also differs from GO-lander in the algorithm it uses for diffing. Vecty uses a high-performance algorithm similar to 
the virtual DOM pattern where Go-lander uses a pure virtual DOM tree with a separate diffing and patching process.

## Getting started

Install the latest version GO-lander to get started.

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

HTML nodes take  a `nodes.Attributes` map as their attributes parameter, which can include any valid HTML attribute 
or property for the element. GO-lander will extract the attributes and assign them using the following rules:

```go
package main

import (
    "github.com/minivera/go-lander/events"
	"github.com/minivera/go-lander/nodes"
)

lander.HTML('div', nodes.Attributes{
	"checked": true, // Will be converted to `checked=""` on the DOM element if true and omitted if false.
	"placeholder": "some string", // Will be kept as a string and assigned on the DOM element.
	"value": -1, // Will be converted to as string and assigned on the DOM element.
	"click": func(event *events.DOMEvent) error {} // Will be assigned using `addEventListener` on the DOM element.
}, nodes.Children{});
```

Any other type is ignored. If an attribute has an associated property on the DOM node (such as `value`), it will 
also be set as a property. See the [content vs. IDL attributes](https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes#content_versus_idl_attributes)  reference for more details.

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
type FunctionComponent func(ctx context.Context, props nodes.Props, children nodes.Children) nodes.Child
```

The `context.Context` is covered in a latter section. The component's `Props` are, like HTML attributes, a map of 
string keys to any type. The properties are never converted by GO-lander, they will be passed along to the function. 
When GO-lander renders the tree, it calls each component in the tree to generate its result (called a render).

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

func helloWorld(_ context.Context, props nodes.Props, children nodes.Children) nodes.Child {
	message := props["message"].(string) // Lander does not provide any prop type safety for now
	
	return lander.Html("h1", nodes.Attributes{}, nodes.Children{
		lander.Text(message),
		children[0],
	}).Style("margin: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(helloWorld, nodes.Props{
			"message": "Hello, World!"
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
that component will in-turn be rendered until all components in the tree have been. This is the entire render 
process of GO-lander.

GO-lander does not provide any type safety for props, they are assigned in a `map[string]interface{}` variable and 
passed without any type information. You are responsible for checking the types of the values you receive, and to 
`panic` if the types are not correct or a required prop is missing.

Any component may return `nil` instead of a `nodes.Child`. The component's render result will be ignored and any 
previous node it rendered will be removed. The component may return a valid node in a later render cycle, the result 
will then be inserted into the DOM.

### Fragment nodes

Returning only a single child from a component may not always be practical. If a component renders a slice of 
children, for example, you would need to wrap that slice in an HTML node such as a `div` to only return a single 
child. GO-lander provide a final type of node to solve this problem, the `lander.Fragment` node.

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

func helloWorld(_ context.Context, props nodes.Props, children nodes.Children) nodes.Child {
	message := props["message"].(string) // Lander does not provide any prop type safety for now
	
	return lander.Html("h1", nodes.Attributes{}, nodes.Children{
		lander.Text(message),
		// But now the fragment will take care of the children
		lander.Fragment(children),
	}).Style("margin: 1rem;")
}

func main() {
	c := make(chan bool)

	_, err := lander.RenderInto(
		lander.Component(helloWorld, nodes.Props{
			"message": "Hello, World!"
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

### Context and listeners

## Experimental features

### Hooks

### Global state management

### In-memory routing

## Acknowledgements

Lander would not have been possible without the massive work done by the contributors of these libraries:

- [React](https://github.com/facebook/react)
- [Vecty](https://github.com/hexops/vecty)
- [Vugu](https://github.com/vugu/vugu) (Lander started out as a React-like clone of Vugu)