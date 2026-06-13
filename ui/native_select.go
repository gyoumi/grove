package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

const nativeSelectClass = "flex h-9 w-full appearance-none rounded-md border border-input bg-transparent px-3 py-1 pr-8 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"

// SelectOption is one choice in a NativeSelect.
type SelectOption struct {
	Value    string
	Label    string
	Disabled bool
}

type NativeSelectProps struct {
	Value       string // selected option value (controlled)
	Options     []SelectOption
	Placeholder string // optional disabled leading option
	Class       string
	ID          string
	Disabled    bool
	OnChange    func(*g.Event)
}

// NativeSelect is a styled native <select>. The current value re-syncs from
// Value each render, and a chevron is overlaid on the right.
func NativeSelect(p NativeSelectProps) *g.Node {
	opts := make([]any, 0, len(p.Options)+1)
	if p.Placeholder != "" {
		opts = append(opts, g.El("option", g.Attr("value", ""), g.Attr("disabled", ""), p.Placeholder))
	}
	for _, o := range p.Options {
		oargs := []any{g.Attr("value", o.Value)}
		if o.Disabled {
			oargs = append(oargs, g.Attr("disabled", ""))
		}
		opts = append(opts, g.El("option", append(oargs, o.Label)...))
	}

	selArgs := []any{
		g.Class(style.CN(nativeSelectClass, p.Class)),
		g.Data("slot", "native-select"),
		g.Value(p.Value),
		g.Disabled(p.Disabled),
	}
	if p.ID != "" {
		selArgs = append(selArgs, g.ID(p.ID))
	}
	if p.OnChange != nil {
		selArgs = append(selArgs, g.OnChange(p.OnChange))
	}

	return g.Span(
		g.Class("relative inline-flex w-full items-center"),
		g.Data("slot", "native-select-wrapper"),
		g.El("select", append(selArgs, opts...)...),
		Icon("chevron-down", "pointer-events-none absolute right-3 size-4 text-muted-foreground"),
	)
}
