package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

const comboboxTriggerClass = "flex h-9 w-full items-center justify-between gap-2 rounded-md border border-input bg-transparent px-3 py-1 text-left text-sm shadow-sm transition-colors hover:bg-accent/40 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 data-[empty=1]:text-muted-foreground"

type ComboboxProps struct {
	Value             string // selected option value (controlled)
	Options           []SelectOption
	Placeholder       string // trigger text when nothing is chosen
	SearchPlaceholder string // command input placeholder
	EmptyText         string
	OnChange          func(string)
	Disabled          bool
	Class             string // extra classes for the trigger
}

// Combobox is a searchable Select: the trigger shows the chosen option, and
// clicking it opens a Command (filterable list) in a popover. The current
// option is checked.
func Combobox(p ComboboxProps) *g.Node { return g.C(comboboxView, p) }

func comboboxView(p ComboboxProps) *g.Node {
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
		g.Class(style.CN(comboboxTriggerClass, p.Class)),
		g.Data("slot", "combobox-trigger"),
		g.AttrIf(empty, "data-empty", "1"),
		g.Disabled(p.Disabled),
		g.OnClick(func(*g.Event) { setOpen(!open) }),
		g.Span(g.Class("truncate"), label),
		Icon("chevrons-up-down", "size-4 shrink-0 text-muted-foreground"),
	)

	items := make([]CommandItem, 0, len(p.Options))
	for _, o := range p.Options {
		selected := o.Value == p.Value
		var labelNode any = o.Label
		if selected {
			labelNode = g.Span(
				g.Class("flex w-full items-center justify-between gap-2"),
				g.Span(o.Label),
				Icon("check", "size-4"),
			)
		}
		value := o.Value
		items = append(items, CommandItem{
			Value:    o.Label,
			Label:    labelNode,
			Disabled: o.Disabled,
			OnSelect: func() {
				if p.OnChange != nil {
					p.OnChange(value)
				}
				setOpen(false)
			},
		})
	}

	cmd := Command(CommandProps{
		Placeholder: p.SearchPlaceholder,
		EmptyText:   p.EmptyText,
		Class:       "min-w-[12rem]",
	}, CommandGroup{Items: items})

	return Popover(PopoverProps{
		Open:         open,
		OnClose:      func() { setOpen(false) },
		Align:        PopoverAlignStart,
		Class:        "w-auto p-0",
		WrapperClass: "relative block w-full",
	}, trigger, cmd)
}
