package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

const selectTriggerClass = "flex h-9 w-full items-center justify-between gap-2 rounded-md border border-input bg-transparent px-3 py-1 text-left text-sm shadow-sm transition-colors hover:bg-accent/40 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 data-[empty=1]:text-muted-foreground"

type SelectProps struct {
	Value       string // selected option value (controlled)
	Options     []SelectOption
	Placeholder string
	OnChange    func(string)
	Disabled    bool
	Class       string // extra classes for the trigger
}

// Select is a dropdown listbox: the trigger shows the chosen option, and
// clicking it opens a popover of options with a check on the current one.
// Unlike NativeSelect it is a fully styled menu, so it themes consistently
// across browsers.
func Select(p SelectProps) *g.Node { return g.C(selectView, p) }

func selectView(p SelectProps) *g.Node {
	open, setOpen := g.UseState(false)

	label, empty := p.Placeholder, true
	for _, o := range p.Options {
		if o.Value == p.Value {
			label, empty = o.Label, false
		}
	}
	if empty && label == "" {
		label = "Select…"
	}

	trigger := g.Button(
		g.Type("button"),
		g.Class(style.CN(selectTriggerClass, p.Class)),
		g.Data("slot", "select-trigger"),
		g.AttrIf(empty, "data-empty", "1"),
		g.Disabled(p.Disabled),
		g.OnClick(func(*g.Event) { setOpen(!open) }),
		g.Span(g.Class("truncate"), label),
		Icon("chevron-down", "size-4 shrink-0 text-muted-foreground"),
	)

	rows := make([]any, 0, len(p.Options)+1)
	rows = append(rows, g.Role("listbox"))
	for _, o := range p.Options {
		selected := o.Value == p.Value
		var check *g.Node
		if selected {
			check = Icon("check", "size-4")
		}
		rows = append(rows, g.Button(
			g.Type("button"),
			g.Role("option"),
			g.Class("flex w-full cursor-pointer select-none items-center justify-between gap-2 rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground disabled:pointer-events-none disabled:opacity-50"),
			g.Data("slot", "select-item"),
			g.Data("value", o.Value),
			g.AttrIf(selected, "aria-selected", "true"),
			g.Disabled(o.Disabled),
			g.OnClick(func(*g.Event) {
				if o.Disabled {
					return
				}
				if p.OnChange != nil {
					p.OnChange(o.Value)
				}
				setOpen(false)
			}),
			g.Span(o.Label),
			check,
		))
	}

	return Popover(PopoverProps{
		Open:         open,
		OnClose:      func() { setOpen(false) },
		Align:        PopoverAlignStart,
		Class:        "flex max-h-72 min-w-[8rem] flex-col gap-0.5 overflow-auto p-1",
		WrapperClass: "relative block w-full",
	}, trigger, rows...)
}
