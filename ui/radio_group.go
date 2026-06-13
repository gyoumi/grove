package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// RadioItem is one choice in a RadioGroup.
type RadioItem struct {
	Value    string
	Label    any
	Disabled bool
}

type RadioGroupProps struct {
	Value    string // selected value (controlled)
	OnChange func(string)
	Class    string
}

// RadioGroup is a controlled set of mutually exclusive choices.
func RadioGroup(p RadioGroupProps, items ...RadioItem) *g.Node {
	rows := make([]any, 0, len(items))
	for _, it := range items {
		selected := it.Value == p.Value
		var dot *g.Node
		if selected {
			dot = g.Div(g.Class("size-2.5 rounded-full bg-primary"))
		}
		radio := g.Button(
			g.Type("button"),
			g.Class("flex size-4 shrink-0 items-center justify-center rounded-full border border-primary text-primary shadow focus:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"),
			g.Role("radio"),
			g.Data("slot", "radio-item"),
			g.Data("value", it.Value),
			g.AttrIf(selected, "aria-checked", "true"),
			g.Disabled(it.Disabled),
			g.OnClick(func(*g.Event) {
				if !it.Disabled && p.OnChange != nil {
					p.OnChange(it.Value)
				}
			}),
			dot,
		)
		rows = append(rows, g.Label(
			g.Class("flex items-center gap-2 text-sm font-medium leading-none"),
			g.Data("slot", "radio-row"),
			radio,
			g.Span(it.Label),
		))
	}
	all := append([]any{
		g.Class(style.CN("grid gap-2.5", p.Class)),
		g.Data("slot", "radio-group"),
		g.Role("radiogroup"),
	}, rows...)
	return g.Div(all...)
}
