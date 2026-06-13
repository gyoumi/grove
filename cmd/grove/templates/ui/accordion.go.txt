package ui

import (
	"maps"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// AccordionItem is one section: a Value that identifies it, a Title shown on
// the trigger, and Content revealed when open.
type AccordionItem struct {
	Value   string
	Title   any
	Content []any
}

type AccordionProps struct {
	// Multiple lets several sections be open at once; the default opens at
	// most one.
	Multiple    bool
	DefaultOpen []string
	Class       string
}

// Accordion is a vertical stack of collapsible sections. It manages which
// sections are open itself.
func Accordion(p AccordionProps, items ...AccordionItem) *g.Node {
	return g.C(accordionView, accordionArgs{p: p, items: items})
}

type accordionArgs struct {
	p     AccordionProps
	items []AccordionItem
}

func accordionView(a accordionArgs) *g.Node {
	openSet, setOpen := g.UseState(sliceSet(a.p.DefaultOpen))

	toggle := func(value string) {
		next := map[string]bool{}
		if a.p.Multiple {
			maps.Copy(next, openSet)
		}
		next[value] = !openSet[value]
		setOpen(next)
	}

	rows := make([]any, 0, len(a.items))
	for _, it := range a.items {
		open := openSet[it.Value]
		header := g.El("h3", g.Class("flex"),
			g.Button(
				g.Type("button"),
				g.Class("flex flex-1 items-center justify-between py-4 text-sm font-medium transition-all hover:underline [&[data-state=open]>svg]:rotate-180"),
				g.Data("slot", "accordion-trigger"),
				g.Data("state", openState(open)),
				g.OnClick(func(*g.Event) { toggle(it.Value) }),
				g.Span(it.Title),
				Icon("chevron-down", "size-4 shrink-0 text-muted-foreground transition-transform duration-200"),
			),
		)
		section := []any{
			g.Class("border-b"),
			g.Data("slot", "accordion-item"),
			g.Data("value", it.Value),
			g.Data("state", openState(open)),
			header,
		}
		if open {
			section = append(section, g.Div(
				g.Class("overflow-hidden pb-4 text-sm text-muted-foreground"),
				g.Data("slot", "accordion-content"),
				it.Content,
			))
		}
		rows = append(rows, g.Div(section...))
	}

	all := append([]any{g.Class(style.CN("w-full", a.p.Class)), g.Data("slot", "accordion")}, rows...)
	return g.Div(all...)
}

func sliceSet(items []string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, s := range items {
		m[s] = true
	}
	return m
}
