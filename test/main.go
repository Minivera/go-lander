package main

import (
	"fmt"

	lander "github.com/minivera/go-lander"
)

func container(_ map[string]interface{}, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("div", map[string]interface{}{}, children).Style(`
			margin-left: auto;
			margin-right: auto;
			max-width: 24rem;
		`),
	}
}

func heading(_ map[string]interface{}, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("div", map[string]interface{}{}, children).Style(`
			align-content: center;
			display: flex;
			flex-direction: column;
			margin-top: 1.5rem;
		`).SelectorStyle(" > *", `
			margin-left: auto;
			margin-right: auto;
		`),
	}
}

func title(_ map[string]interface{}, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("h1", map[string]interface{}{}, children).Style(`
			font-size: 1.1447142426em;
			font-weight: 400;
			margin-bottom: 1rem;
			margin-top: 1rem;
			text-align: center;
		`),
	}
}

func separator(_ map[string]interface{}, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("div", map[string]interface{}{}, children).Style(`
			text-align: center;
			margin: 1rem auto;
			border-bottom: 1px solid #C4C7CA;
		`).SelectorStyle(" span", `
			position: relative;
			padding: 0 1rem;
			top: 11px;
			font-size: 14px;
			color: #575B5F;
			background: white;
		`),
	}
}

func form(attrs map[string]interface{}, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("form", attrs, children).Style(`
			padding: 1rem 0;
		`),
	}
}

func formGroup(_ map[string]interface{}, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("div", map[string]interface{}{}, children).Style(`
			display: block;
			padding: 1rem 0;
		`),
	}
}

func passwordLabel(attrs map[string]interface{}, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("label", attrs, children).Style(`
			display: block;
			padding-bottom: 2px;
			font-size: 14px;
			cursor: pointer;
			display: flex;
			justify-content: space-between;
		`).SelectorStyle(" > span", `
			color: #575B5F;
		`),
	}
}

func passwordHelp(_ map[string]interface{}, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("a", map[string]interface{}{}, children).Style(`
			margin-left: 0.5rem;
			text-decoration: none;
			color: #575B5F;
		`).SelectorStyle(":hover", `
			color: #1E50DA;
		`),
	}
}

func footer(_ map[string]interface{}, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("footer", map[string]interface{}{}, children).Style(`
			padding: 1rem 0;
		`),
	}
}

var errorState []string

func loginPage(_ map[string]interface{}, _ []lander.Node) []lander.Node {
	var errorNode lander.Node

	if len(errorState) > 0 {
		listNodes := make([]lander.Node, len(errorState))

		for index, err := range errorState {
			listNodes[index] = lander.Html("li", map[string]interface{}{}, []lander.Node{
				lander.Text(err),
			})
		}

		errorNode = lander.Component(toast, map[string]interface{}{}, []lander.Node{
			lander.Html("ul", map[string]interface{}{}, listNodes),
		})
	}

	var onSubmit lander.EventListener = func(currentNode lander.Node, event *lander.DOMEvent) error {
		event.PreventDefault()

		errorState = []string{
			"error",
		}
		return nil
	}

	return []lander.Node{
		lander.Component(container, map[string]interface{}{}, []lander.Node{
			lander.Component(heading, map[string]interface{}{}, []lander.Node{
				lander.Component(manifold, map[string]interface{}{}, []lander.Node{}),
				lander.Component(title, map[string]interface{}{}, []lander.Node{
					lander.Html("span", map[string]interface{}{}, []lander.Node{
						lander.Text("Log in to Manifold"),
					}),
				}),
			}),
			lander.Component(oAuthButton, map[string]interface{}{}, []lander.Node{
				lander.Html("div", map[string]interface{}{}, []lander.Node{
					lander.Text("Continue with github"),
				}),
			}),
			lander.Component(separator, map[string]interface{}{}, []lander.Node{
				lander.Html("span", map[string]interface{}{}, []lander.Node{
					lander.Text("or use your email address"),
				}),
			}),
			lander.Component(form, map[string]interface{}{
				"submit": onSubmit,
			}, []lander.Node{
				lander.Component(formGroup, map[string]interface{}{}, []lander.Node{
					lander.Component(inputLabel, map[string]interface{}{
						"for": "email",
					}, []lander.Node{
						lander.Text("Email address"),
					}),
					lander.Component(input, map[string]interface{}{
						"id":   "email",
						"name": "email",
						"type": "email",
					}, []lander.Node{}),
				}),
				lander.Component(formGroup, map[string]interface{}{}, []lander.Node{
					lander.Component(passwordLabel, map[string]interface{}{
						"for": "password",
					}, []lander.Node{
						lander.Text("Password"),
						lander.Component(passwordHelp, map[string]interface{}{
							"tab-index": -1,
						}, []lander.Node{
							lander.Text("Forgot?"),
						}),
					}),
					lander.Component(input, map[string]interface{}{
						"id":   "password",
						"name": "password",
						"type": "password",
					}, []lander.Node{}),
				}),
				errorNode,
				lander.Component(footer, map[string]interface{}{}, []lander.Node{
					lander.Component(button, map[string]interface{}{
						"type": "submit",
					}, []lander.Node{
						lander.Text("Log in to continue"),
					}),
				}),
				lander.Component(text, map[string]interface{}{}, []lander.Node{
					lander.Text("Don’t have an account? "),
					lander.Html("a", map[string]interface{}{}, []lander.Node{
						lander.Text("Sign up now."),
					}),
				}),
			}),
		}),
	}
}

var landerEnv *lander.DomEnvironment

func main() {
	c := make(chan bool)

	tree := lander.Component(wrapper, map[string]interface{}{}, []lander.Node{
		lander.Component(loginPage, map[string]interface{}{}, []lander.Node{}),
	})

	landerEnv = lander.NewLander("#app", tree)

	err := landerEnv.Mount()
	if err != nil {
		fmt.Println(err)
	}

	<-c
}
