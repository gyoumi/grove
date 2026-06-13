package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// ScrollArea is a scrollable viewport with a thin, themed scrollbar. Give it
// a bounded height and put content inside:
//
//	ui.ScrollArea("h-72", ...rows...)
func ScrollArea(class string, children ...any) *g.Node {
	args := []any{
		g.Class(style.CN(
			"relative overflow-auto rounded-[inherit]",
			"[scrollbar-width:thin] [scrollbar-color:var(--border)_transparent]",
			"[&::-webkit-scrollbar]:w-2 [&::-webkit-scrollbar]:h-2",
			"[&::-webkit-scrollbar-thumb]:rounded-full [&::-webkit-scrollbar-thumb]:bg-border",
			class,
		)),
		g.Data("slot", "scroll-area"),
	}
	return g.Div(append(args, children...)...)
}
