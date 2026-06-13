package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// Kbd renders a keyboard key hint, e.g. ui.Kbd("⌘"), ui.Kbd("K").
func Kbd(children ...any) *g.Node {
	args := []any{
		g.Class(style.CN("pointer-events-none inline-flex h-5 min-w-5 select-none items-center justify-center gap-1 rounded border bg-muted px-1.5 font-mono text-[0.7rem] font-medium text-muted-foreground")),
		g.Data("slot", "kbd"),
	}
	return g.El("kbd", append(args, children...)...)
}

// KbdGroup lays out several Kbd hints in a row.
func KbdGroup(children ...any) *g.Node {
	args := []any{
		g.Class("inline-flex items-center gap-1"),
		g.Data("slot", "kbd-group"),
	}
	return g.Span(append(args, children...)...)
}
