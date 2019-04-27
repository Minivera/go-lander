package main

import (
	"fmt"

	lander "github.com/minivera/go-lander"
)

func container(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("div", map[string]string{}, children).Style(`
			margin-left: auto;
			margin-right: auto;
			max-width: 24rem;
		`),
	}
}

func heading(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("div", map[string]string{}, children).Style(`
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

func title(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("h1", map[string]string{}, children).Style(`
			font-size: 1.1447142426em;
			font-weight: 400;
			margin-bottom: 1rem;
			margin-top: 1rem;
			text-align: center;
		`),
	}
}

func separator(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("div", map[string]string{}, children).Style(`
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

func form(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("form", map[string]string{}, children).Style(`
			padding: 1rem 0;
		`),
	}
}

func formGroup(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("div", map[string]string{}, children).Style(`
			display: block;
			padding: 1rem 0;
		`),
	}
}

func passwordLabel(attrs map[string]string, children []lander.Node) []lander.Node {
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

func passwordHelp(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("a", map[string]string{}, children).Style(`
			margin-left: 0.5rem;
			text-decoration: none;
			color: #575B5F;
		`).SelectorStyle(":hover", `
			color: #1E50DA;
		`),
	}
}

func footer(_ map[string]string, children []lander.Node) []lander.Node {
	return []lander.Node{
		lander.Html("footer", map[string]string{}, children).Style(`
			padding: 1rem 0;
		`),
	}
}

func main() {
	tree := lander.Component(wrapper, map[string]string{}, []lander.Node{
		lander.Component(container, map[string]string{}, []lander.Node{
			lander.Component(heading, map[string]string{}, []lander.Node{
				lander.Component(manifold, map[string]string{}, []lander.Node{}),
				lander.Component(title, map[string]string{}, []lander.Node{
					lander.Html("span", map[string]string{}, []lander.Node{
						lander.Text("Log in to Manifold"),
					}),
				}),
			}),
			lander.Component(oAuthButton, map[string]string{}, []lander.Node{
				lander.Html("div", map[string]string{}, []lander.Node{
					lander.Text("Continue with github"),
				}),
			}),
			lander.Component(separator, map[string]string{}, []lander.Node{
				lander.Html("span", map[string]string{}, []lander.Node{
					lander.Text("or use your email address"),
				}),
			}),
			lander.Component(form, map[string]string{}, []lander.Node{
				lander.Component(formGroup, map[string]string{}, []lander.Node{
					lander.Component(inputLabel, map[string]string{
						"for": "email",
					}, []lander.Node{
						lander.Text("Email address"),
					}),
					lander.Component(input, map[string]string{
						"id":       "email",
						"name":     "email",
						"type":     "email",
						"required": "true",
					}, []lander.Node{}),
				}),
				lander.Component(formGroup, map[string]string{}, []lander.Node{
					lander.Component(passwordLabel, map[string]string{
						"for": "password",
					}, []lander.Node{
						lander.Text("Password"),
						lander.Component(passwordHelp, map[string]string{
							"tab-index": "-1",
						}, []lander.Node{
							lander.Text("Forgot?"),
						}),
					}),
					lander.Component(input, map[string]string{
						"id":        "password",
						"name":      "password",
						"type":      "password",
						"required":  "true",
						"minLength": "8",
					}, []lander.Node{}),
				}),
				lander.Component(footer, map[string]string{}, []lander.Node{
					lander.Component(button, map[string]string{
						"type": "submit",
					}, []lander.Node{
						lander.Text("Log in to continue"),
					}),
				}),
				lander.Component(text, map[string]string{}, []lander.Node{
					lander.Text("Don’t have an account? "),
					lander.Html("a", map[string]string{}, []lander.Node{
						lander.Text("Sign up now."),
					}),
				}),
			}),
		}),
	})

	err := lander.NewLander("#app", tree).Mount()
	if err != nil {
		fmt.Println(err)
	}
}
