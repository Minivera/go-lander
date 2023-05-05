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
   experience, you should feel right at home. Rather than to clone React in Go, we've attempted to provide the React 
   experience without moving away from the experience of writing Go code.
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

In this case, the `helloWorld` component is a function that renders a `h1` tag with the text "Hello, World!" ans 
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

### Writing and styling HTML

### Components

### Fragment nodes

### Keeping state

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