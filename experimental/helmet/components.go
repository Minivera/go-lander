package helmet

import (
	"github.com/minivera/go-lander/context"
	"github.com/minivera/go-lander/nodes"
)

// Head is the main component of Helmet. Any valid children passed to this component will be rendered inside
// the head, removing any unnecessary tags along the way. This component expects html nodes as its first level
// children. It will read each child and, if they are a valid head tag (`title`, `meta`, `link`, `script`,
// `noscript`, or `style`), it will insert them into the head of the document. `title` is special and only the
// deepest `title` in the tree will render, all other tags will be removed.
func Head(ctx context.Context, _ nodes.Props, children nodes.Children) nodes.Child {
	if !ctx.HasValue("lander_helmet_defs") {
		panic("helmet.Head were used outside of a helmet provider, make sure to wrap your entire app in a `lander.Component(helmet.Provider)`")
	}

	defs := ctx.GetValue("lander_helmet_defs").([]helmetDef)

	for _, child := range children {
		if child.Type() != nodes.HTMLNodeType {
			continue
		}

		htmlChild := child.(*nodes.HTMLNode)

		// Make sure we only process the tags we care about
		if htmlChild.Tag != "title" && htmlChild.Tag != "meta" && htmlChild.Tag != "link" &&
			htmlChild.Tag != "script" && htmlChild.Tag != "noscript" && htmlChild.Tag != "style" {
			continue
		}

		htmlChild.Attributes["data-tag"] = helmetTag

		defs = append(defs, helmetDef{
			tag:  htmlChild.Tag,
			node: htmlChild,
		})
	}

	ctx.SetValue("lander_helmet_defs", defs)

	return nil
}
