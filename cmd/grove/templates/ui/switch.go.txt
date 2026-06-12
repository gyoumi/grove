package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type SwitchProps struct {
	Checked  bool
	Disabled bool
	Class    string
	ID       string
	OnChange func(checked bool)
}

// Switch renders a toggle: a sliding thumb on a track, controlled by the
// caller's state. It is a real button with role="switch", so it is
// keyboard- and screen-reader-friendly out of the box.
func Switch(p SwitchProps, args ...any) *g.Node {
	track := "bg-input"
	thumb := "translate-x-0"
	state := "unchecked"
	if p.Checked {
		track = "bg-primary"
		thumb = "translate-x-4"
		state = "checked"
	}
	all := []any{
		g.Class(style.CN(
			"inline-flex h-5 w-9 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50",
			track, p.Class)),
		g.Type("button"),
		g.Role("switch"),
		g.Aria("checked", boolStr(p.Checked)),
		g.Data("slot", "switch"),
		g.Data("state", state),
		g.Disabled(p.Disabled),
		g.Span(g.Class("pointer-events-none block size-4 rounded-full bg-background shadow-lg transition-transform", thumb)),
	}
	if p.ID != "" {
		all = append(all, g.ID(p.ID))
	}
	if p.OnChange != nil {
		onChange := p.OnChange
		next := !p.Checked
		all = append(all, g.OnClick(func(*g.Event) { onChange(next) }))
	}
	return g.Button(append(all, args...)...)
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
