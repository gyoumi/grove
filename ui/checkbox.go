package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type CheckboxProps struct {
	Checked  bool
	Disabled bool
	Class    string
	ID       string
	OnChange func(checked bool)
}

// Checkbox renders a styled native checkbox (grove's port uses the real
// <input type="checkbox"> with accent-color theming rather than a custom
// control, keeping it dependency-free and accessible by default).
func Checkbox(p CheckboxProps, args ...any) *g.Node {
	all := []any{
		g.Class(style.CN("size-4 shrink-0 rounded-[4px] border border-input accent-primary shadow-sm disabled:cursor-not-allowed disabled:opacity-50", p.Class)),
		g.Type("checkbox"),
		g.Data("slot", "checkbox"),
		g.Checked(p.Checked),
		g.Disabled(p.Disabled),
	}
	if p.ID != "" {
		all = append(all, g.ID(p.ID))
	}
	if p.OnChange != nil {
		onChange := p.OnChange
		all = append(all, g.OnChange(func(e *g.Event) { onChange(e.Checked()) }))
	}
	return g.Input(append(all, args...)...)
}
