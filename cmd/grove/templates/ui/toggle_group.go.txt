package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// ToggleGroupItem is one toggle in a ToggleGroup.
type ToggleGroupItem struct {
	Value    string
	Children []any
	Disabled bool
}

type ToggleGroupProps struct {
	// Multiple lets several items be pressed at once; the default is a
	// single-selection group.
	Multiple bool
	Value    []string // pressed values (controlled)
	OnChange func([]string)
	Variant  ToggleVariant
	Size     ToggleSize
	Class    string
}

// ToggleGroup is a controlled row of segmented toggles sharing one selection.
func ToggleGroup(p ToggleGroupProps, items ...ToggleGroupItem) *g.Node {
	sel := sliceSet(p.Value)

	next := func(v string) []string {
		if p.Multiple {
			out := make([]string, 0, len(p.Value)+1)
			found := false
			for _, x := range p.Value {
				if x == v {
					found = true
					continue
				}
				out = append(out, x)
			}
			if !found {
				out = append(out, v)
			}
			return out
		}
		if sel[v] {
			return nil
		}
		return []string{v}
	}

	rows := make([]any, 0, len(items))
	for _, it := range items {
		on := sel[it.Value]
		args := []any{
			g.Class(toggleVariants.Class(map[string]string{
				"variant": string(p.Variant),
				"size":    string(p.Size),
			}, "rounded-none shadow-none first:rounded-l-md last:rounded-r-md focus:z-10")),
			g.Type("button"),
			g.Data("slot", "toggle-group-item"),
			g.Data("value", it.Value),
			g.Data("state", onOff(on)),
			g.AttrIf(on, "aria-pressed", "true"),
			g.Disabled(it.Disabled),
			g.OnClick(func(*g.Event) {
				if !it.Disabled && p.OnChange != nil {
					p.OnChange(next(it.Value))
				}
			}),
		}
		rows = append(rows, g.Button(append(args, it.Children...)...))
	}

	all := append([]any{
		g.Class(style.CN("inline-flex items-center [&>*:not(:first-child)]:-ml-px", p.Class)),
		g.Data("slot", "toggle-group"),
		g.Role("group"),
	}, rows...)
	return g.Div(all...)
}
