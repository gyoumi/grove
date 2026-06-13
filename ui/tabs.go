package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// Tab is one tab: a Value, the Label shown on its trigger, and the Content
// shown when it is active.
type Tab struct {
	Value   string
	Label   any
	Content []any
}

type TabsProps struct {
	DefaultValue string
	Class        string
}

// Tabs shows one panel at a time, switched by a row of triggers. It tracks
// the active tab itself (defaulting to DefaultValue, else the first tab).
func Tabs(p TabsProps, tabs ...Tab) *g.Node {
	return g.C(tabsView, tabsArgs{p: p, tabs: tabs})
}

type tabsArgs struct {
	p    TabsProps
	tabs []Tab
}

func tabsView(a tabsArgs) *g.Node {
	def := a.p.DefaultValue
	if def == "" && len(a.tabs) > 0 {
		def = a.tabs[0].Value
	}
	active, setActive := g.UseState(def)

	triggers := make([]any, 0, len(a.tabs))
	for _, t := range a.tabs {
		on := t.Value == active
		triggers = append(triggers, g.Button(
			g.Type("button"),
			g.Class(style.CN(
				"inline-flex items-center justify-center whitespace-nowrap rounded-md px-3 py-1 text-sm font-medium transition-all focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring",
				map[string]bool{
					"bg-background text-foreground shadow":        on,
					"text-muted-foreground hover:text-foreground": !on,
				})),
			g.Data("slot", "tabs-trigger"),
			g.Data("value", t.Value),
			g.Data("state", activeState(on)),
			g.Role("tab"),
			g.OnClick(func(*g.Event) { setActive(t.Value) }),
			t.Label,
		))
	}
	list := g.Div(append([]any{
		g.Class("inline-flex h-9 w-fit items-center justify-center rounded-lg bg-muted p-1 text-muted-foreground"),
		g.Data("slot", "tabs-list"),
		g.Role("tablist"),
	}, triggers...)...)

	var content *g.Node
	for _, t := range a.tabs {
		if t.Value == active {
			content = g.Div(append([]any{
				g.Class("mt-2"),
				g.Data("slot", "tabs-content"),
				g.Data("value", t.Value),
				g.Role("tabpanel"),
			}, t.Content...)...)
			break
		}
	}

	return g.Div(
		g.Class(style.CN("flex flex-col gap-1", a.p.Class)),
		g.Data("slot", "tabs"),
		list,
		content,
	)
}

func activeState(on bool) string {
	if on {
		return "active"
	}
	return "inactive"
}
