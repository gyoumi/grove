package ui

import (
	"fmt"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// Time is a wall-clock time of day with no date — the value a TimePicker
// holds.
type Time struct {
	Hour   int // 0-23
	Minute int // 0-59
}

// String formats the time as zero-padded HH:MM.
func (t Time) String() string { return fmt.Sprintf("%02d:%02d", t.Hour, t.Minute) }

// Minutes is the time as minutes since midnight.
func (t Time) Minutes() int { return t.Hour*60 + t.Minute }

// TimeFromMinutes builds a Time from minutes since midnight.
func TimeFromMinutes(m int) Time { return Time{Hour: m / 60, Minute: m % 60} }

// timeTriggerClass matches dateTriggerClass; kept separate so each picker
// copies in as a self-contained file.
const timeTriggerClass = "flex h-9 w-full items-center justify-between gap-2 rounded-md border border-input bg-transparent px-3 py-1 text-left text-sm shadow-sm transition-colors hover:bg-accent/40 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 data-[empty=1]:text-muted-foreground"

type TimePickerProps struct {
	Value       *Time // nil shows the placeholder
	OnChange    func(Time)
	Step        int // minute increment for the minute column (default 30)
	Placeholder string
	Disabled    bool
	Class       string // extra classes for the trigger
}

// TimePicker is a dropdown time field: the trigger reads like a text input
// showing HH:MM (or a placeholder), and clicking it opens scrollable hour
// and minute columns. Picking from a column sets that part of the value;
// the popover closes on outside click or Escape, so both parts can be set
// without reopening.
func TimePicker(p TimePickerProps) *g.Node { return g.C(timePickerView, p) }

func timePickerView(p TimePickerProps) *g.Node {
	open, setOpen := g.UseState(false)
	step := p.Step
	if step <= 0 {
		step = 30
	}

	have := p.Value != nil
	cur := Time{}
	if have {
		cur = *p.Value
	}
	label := p.Placeholder
	if have {
		label = cur.String()
	} else if label == "" {
		label = "Pick a time"
	}

	set := func(t Time) {
		if p.OnChange != nil {
			p.OnChange(t)
		}
	}

	trigger := g.Button(
		g.Type("button"),
		g.Class(style.CN(timeTriggerClass, p.Class)),
		g.Data("slot", "time-picker-trigger"),
		g.AttrIf(!have, "data-empty", "1"),
		g.Disabled(p.Disabled),
		g.OnClick(func(*g.Event) { setOpen(!open) }),
		g.Span(g.Class("truncate"), label),
		Icon("chevron-down", "size-4 shrink-0 text-muted-foreground"),
	)

	hours := make([]*g.Node, 0, 24)
	for h := range 24 {
		hh := h
		hours = append(hours, timeOption("hour", hh, have && cur.Hour == hh, func() {
			set(Time{Hour: hh, Minute: cur.Minute})
		}))
	}
	mins := make([]*g.Node, 0, 60/step+1)
	for m := 0; m < 60; m += step {
		mm := m
		mins = append(mins, timeOption("minute", mm, have && cur.Minute == mm, func() {
			set(Time{Hour: cur.Hour, Minute: mm})
		}))
	}

	panel := g.Div(g.Class("flex h-48 gap-1"),
		g.Div(g.Class("flex flex-col gap-0.5 overflow-y-auto pr-1"), g.Data("slot", "time-hours"), hours),
		g.Div(g.Class("w-px bg-border")),
		g.Div(g.Class("flex flex-col gap-0.5 overflow-y-auto pr-1"), g.Data("slot", "time-minutes"), mins),
	)

	return Popover(
		PopoverProps{
			Open:         open,
			OnClose:      func() { setOpen(false) },
			Align:        PopoverAlignStart,
			Class:        "w-auto p-2",
			WrapperClass: "relative block w-full",
		},
		trigger,
		panel,
	)
}

func timeOption(kind string, val int, selected bool, onClick func()) *g.Node {
	cls := style.CN(
		"rounded px-3 py-1 text-center text-sm tabular-nums transition-colors",
		map[string]bool{
			"hover:bg-accent hover:text-accent-foreground": !selected,
			"bg-primary text-primary-foreground":           selected,
		},
	)
	return g.Button(
		g.Type("button"),
		g.Class(cls),
		g.Data("slot", "time-option"),
		g.Data(kind, fmt.Sprint(val)),
		g.AttrIf(selected, "data-selected", "1"),
		g.OnClick(func(*g.Event) { onClick() }),
		fmt.Sprintf("%02d", val),
	)
}
