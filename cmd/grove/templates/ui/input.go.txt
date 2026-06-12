package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

const inputClass = "flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"

type InputProps struct {
	Type        string // defaults to "text"
	Value       string
	Placeholder string
	Class       string
	ID          string
	Disabled    bool
	OnInput     func(*g.Event)
	OnChange    func(*g.Event)
	OnKeyDown   func(*g.Event)
}

// Input renders a themed text input. When OnInput is set the input is
// controlled: Value is re-synced to the DOM on every render.
func Input(p InputProps, args ...any) *g.Node {
	typ := p.Type
	if typ == "" {
		typ = "text"
	}
	all := []any{
		g.Class(style.CN(inputClass, p.Class)),
		g.Type(typ),
		g.Data("slot", "input"),
		g.Disabled(p.Disabled),
	}
	if p.ID != "" {
		all = append(all, g.ID(p.ID))
	}
	if p.Placeholder != "" {
		all = append(all, g.Placeholder(p.Placeholder))
	}
	if p.OnInput != nil || p.Value != "" {
		all = append(all, g.Value(p.Value))
	}
	if p.OnInput != nil {
		all = append(all, g.OnInput(p.OnInput))
	}
	if p.OnChange != nil {
		all = append(all, g.OnChange(p.OnChange))
	}
	if p.OnKeyDown != nil {
		all = append(all, g.OnKeyDown(p.OnKeyDown))
	}
	return g.Input(append(all, args...)...)
}
