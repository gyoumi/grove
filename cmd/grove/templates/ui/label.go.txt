package ui

import g "github.com/gyoumi/grove"

// Label renders a shadcn-style form label. forID associates it with a
// control by id ("" for none).
func Label(forID string, args ...any) *g.Node {
	all := []any{
		g.Class("flex items-center gap-2 text-sm font-medium leading-none select-none peer-disabled:cursor-not-allowed peer-disabled:opacity-50"),
		g.Data("slot", "label"),
	}
	if forID != "" {
		all = append(all, g.For(forID))
	}
	return g.Label(append(all, args...)...)
}
