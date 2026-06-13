package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// dateTriggerClass styles a button to read like an Input but behave as a
// dropdown trigger; data-empty dims it when nothing is selected.
const dateTriggerClass = "flex h-9 w-full items-center justify-between gap-2 rounded-md border border-input bg-transparent px-3 py-1 text-left text-sm shadow-sm transition-colors hover:bg-accent/40 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 data-[empty=1]:text-muted-foreground"

type DatePickerProps struct {
	Mode CalendarMode // single (default) or range

	// Single (CalendarSingle): Value is the chosen day, OnChange fires on pick.
	Value    *Date
	OnChange func(Date)

	// Range (CalendarRange): Start/End are the chosen range, OnRange fires
	// once both ends are picked.
	Start, End *Date
	OnRange    func(start, end Date)

	Placeholder string            // trigger text when nothing is chosen
	Format      func(Date) string // day formatting in the trigger (default ISO)
	Min, Max    *Date
	Disabled    bool
	Class       string // extra classes for the trigger
}

// DatePicker is a dropdown date field: the trigger reads like a text input
// showing the selection (or a placeholder), and clicking it opens a
// Calendar in a popover. It selects a single day or a date range; in range
// mode the popover stays open until both ends are chosen, then closes.
func DatePicker(p DatePickerProps) *g.Node { return g.C(datePickerView, p) }

func datePickerView(p DatePickerProps) *g.Node {
	open, setOpen := g.UseState(false)

	format := p.Format
	if format == nil {
		format = func(d Date) string { return d.ISO() }
	}

	label, empty := datePickerLabel(p, format)
	trigger := g.Button(
		g.Type("button"),
		g.Class(style.CN(dateTriggerClass, p.Class)),
		g.Data("slot", "date-picker-trigger"),
		g.AttrIf(empty, "data-empty", "1"),
		g.Disabled(p.Disabled),
		g.OnClick(func(*g.Event) { setOpen(!open) }),
		g.Span(g.Class("truncate"), label),
		g.Span(g.Class("shrink-0 text-muted-foreground"), "▾"),
	)

	cal := Calendar(CalendarProps{
		Mode:     p.Mode,
		Selected: p.Value,
		OnSelect: func(d Date) {
			if p.OnChange != nil {
				p.OnChange(d)
			}
			setOpen(false)
		},
		Start: p.Start, End: p.End,
		OnRange: func(s, e Date) {
			if p.OnRange != nil {
				p.OnRange(s, e)
			}
			setOpen(false)
		},
		Min: p.Min, Max: p.Max,
		Class: "border-0 bg-transparent p-0 shadow-none",
	})

	return Popover(
		PopoverProps{
			Open:         open,
			OnClose:      func() { setOpen(false) },
			Align:        PopoverAlignStart,
			Class:        "w-auto p-2",
			WrapperClass: "relative block w-full",
		},
		trigger,
		cal,
	)
}

// datePickerLabel is the trigger text and whether the picker is empty.
func datePickerLabel(p DatePickerProps, format func(Date) string) (string, bool) {
	if p.Mode == CalendarRange {
		switch {
		case p.Start != nil && p.End != nil:
			return format(*p.Start) + " – " + format(*p.End), false
		case p.Start != nil:
			return format(*p.Start) + " – …", false
		}
	} else if p.Value != nil {
		return format(*p.Value), false
	}
	if p.Placeholder != "" {
		return p.Placeholder, true
	}
	if p.Mode == CalendarRange {
		return "Pick a date range", true
	}
	return "Pick a date", true
}
