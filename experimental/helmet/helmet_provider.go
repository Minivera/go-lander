package helmet

import (
	"fmt"
	"syscall/js"

	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

const helmetTag = "lander-helmet"

type helmetDef struct {
	tag  string
	node *nodes.HTMLNode
}

// Provider provides the tools for the head tag to be properly inserted. This Provider must be added as
// the first component of the app. It takes care of setting up and tracking the changes to the head, then
// making sure the head is properly updated on every render. Helmet is fairly basic at the moment and will
// update the head on every render.
func Provider(ctx context.Context, _ nodes.Props, children nodes.Children) nodes.Child {
	ctx.SetValue("lander_helmet_defs", []helmetDef{})

	ctx.OnRender(func() error {
		defs := ctx.GetValue("lander_helmet_defs").([]helmetDef)
		document := js.Global().Get("document")

		head := document.Call("querySelector", "head")
		if !head.Truthy() {
			return fmt.Errorf("failed to find head using query selector")
		}

		previousTags := document.Call("querySelectorAll", fmt.Sprintf("[data-tag=\"%s\"]", helmetTag))
		for i := 0; i < previousTags.Length(); i++ {
			head.Call("removeChild", previousTags.Index(i))
		}

		for _, def := range defs {
			if def.tag == "title" {
				// For title, remove every title we can find. not only the helmet generated tags
				previousTags := document.Call("querySelectorAll", def.tag)
				for i := 0; i < previousTags.Length(); i++ {
					head.Call("removeChild", previousTags.Index(i))
				}
			}

			element := nodes.NewHTMLElement(document, def.node)
			for _, child := range def.node.Children {
				element.Set("innerHTML", element.Get("innerHTML").String()+child.ToString())
			}
			head.Call("appendChild", element)
		}

		return nil
	})

	return nodes.NewFragmentNode(children)
}
