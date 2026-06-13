package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type CollapsibleProps struct {
	DefaultOpen bool
	Class       string
}

// Collapsible shows or hides its content when the trigger is clicked. It
// manages its own open state; trigger is the clickable label.
//
//	ui.Collapsible(ui.CollapsibleProps{}, g.Span("Details"), g.P("…"))
func Collapsible(p CollapsibleProps, trigger any, content ...any) *g.Node {
	return g.C(collapsibleView, collapsibleArgs{p: p, trigger: trigger, content: content})
}

type collapsibleArgs struct {
	p       CollapsibleProps
	trigger any
	content []any
}

func collapsibleView(a collapsibleArgs) *g.Node {
	open, setOpen := g.UseState(a.p.DefaultOpen)
	body := []any{
		g.Class(style.CN("flex flex-col gap-2", a.p.Class)),
		g.Data("slot", "collapsible"),
		g.Data("state", openState(open)),
		g.Button(
			g.Type("button"),
			g.Class("flex items-center gap-2"),
			g.Data("slot", "collapsible-trigger"),
			g.OnClick(func(*g.Event) { setOpen(!open) }),
			a.trigger,
		),
	}
	if open {
		body = append(body, g.Div(g.Data("slot", "collapsible-content"), a.content))
	}
	return g.Div(body...)
}

// openState maps a bool to the data-state value shared by the disclosure
// components.
func openState(open bool) string {
	if open {
		return "open"
	}
	return "closed"
}
