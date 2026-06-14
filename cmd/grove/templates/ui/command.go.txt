package ui

import (
	"strings"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// CommandItem is one selectable row. Value is what the query filters on;
// Label is what's shown (defaults to Value when nil).
type CommandItem struct {
	Value    string
	Label    any
	OnSelect func()
	Disabled bool
}

// CommandGroup is a labelled set of items.
type CommandGroup struct {
	Heading string
	Items   []CommandItem
}

type CommandProps struct {
	Placeholder string // search box placeholder
	EmptyText   string // shown when nothing matches (default "No results.")
	Class       string
}

// Command is a searchable, keyboard-driven list: type to filter, Arrow keys
// to move, Enter to pick. It powers Combobox and command palettes (drop it
// in a Dialog). Filtering is a case-insensitive substring match on each
// item's Value.
func Command(p CommandProps, groups ...CommandGroup) *g.Node {
	return g.C(commandView, commandArgs{p: p, groups: groups})
}

type commandArgs struct {
	p      CommandProps
	groups []CommandGroup
}

type flatItem struct {
	item  CommandItem
	group int
}

func commandView(a commandArgs) *g.Node {
	query, setQuery := g.UseState("")
	active, setActive := g.UseState(0)

	// Flatten the matching items in display order for keyboard navigation.
	var visible []flatItem
	q := strings.ToLower(strings.TrimSpace(query))
	for gi, grp := range a.groups {
		for _, it := range grp.Items {
			if q == "" || strings.Contains(strings.ToLower(it.Value), q) {
				visible = append(visible, flatItem{item: it, group: gi})
			}
		}
	}
	if active > len(visible)-1 {
		active = max(0, len(visible)-1)
	}

	pick := func(it CommandItem) {
		if !it.Disabled && it.OnSelect != nil {
			it.OnSelect()
		}
	}

	placeholder := a.p.Placeholder
	if placeholder == "" {
		placeholder = "Type a command or search…"
	}
	input := g.Div(
		g.Class("flex items-center gap-2 border-b px-3"),
		g.Data("slot", "command-input-wrapper"),
		Icon("search", "size-4 shrink-0 text-muted-foreground"),
		g.Input(
			g.Class("flex h-10 w-full bg-transparent py-3 text-sm outline-none placeholder:text-muted-foreground"),
			g.Data("slot", "command-input"),
			g.Value(query),
			g.Attr("placeholder", placeholder),
			g.OnInput(func(e *g.Event) { setQuery(e.Value()); setActive(0) }),
			g.OnKeyDown(func(e *g.Event) {
				switch e.Key() {
				case "ArrowDown":
					e.PreventDefault()
					setActive(min(active+1, len(visible)-1))
				case "ArrowUp":
					e.PreventDefault()
					setActive(max(active-1, 0))
				case "Enter":
					e.PreventDefault()
					if active >= 0 && active < len(visible) {
						pick(visible[active].item)
					}
				}
			}),
		),
	)

	var rows []any
	if len(visible) == 0 {
		empty := a.p.EmptyText
		if empty == "" {
			empty = "No results."
		}
		rows = append(rows, g.Div(
			g.Class("py-6 text-center text-sm text-muted-foreground"),
			g.Data("slot", "command-empty"),
			empty,
		))
	} else {
		lastGroup := -1
		for i, fi := range visible {
			if fi.group != lastGroup {
				lastGroup = fi.group
				if h := a.groups[fi.group].Heading; h != "" {
					rows = append(rows, g.Div(
						g.Class("px-2 py-1.5 text-xs font-medium text-muted-foreground"),
						g.Data("slot", "command-group-heading"),
						h,
					))
				}
			}
			label := fi.item.Label
			if label == nil {
				label = fi.item.Value
			}
			rows = append(rows, g.Button(
				g.Type("button"),
				g.Class(style.CN(
					"flex w-full cursor-pointer select-none items-center gap-2 rounded-sm px-2 py-1.5 text-left text-sm outline-none disabled:pointer-events-none disabled:opacity-50",
					map[string]bool{"bg-accent text-accent-foreground": i == active},
				)),
				g.Data("slot", "command-item"),
				g.Data("value", fi.item.Value),
				g.AttrIf(i == active, "data-active", "1"),
				g.Disabled(fi.item.Disabled),
				g.OnClick(func(*g.Event) { pick(fi.item) }),
				label,
			))
		}
	}

	return g.Div(
		g.Class(style.CN("flex w-full flex-col overflow-hidden rounded-md bg-popover text-popover-foreground", a.p.Class)),
		g.Data("slot", "command"),
		input,
		g.Div(
			g.Class("max-h-72 overflow-y-auto p-1"),
			g.Data("slot", "command-list"),
			g.Role("listbox"),
			rows,
		),
	)
}
