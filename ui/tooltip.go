package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type TooltipProps struct {
	// Label is the tooltip text shown above the trigger.
	Label string
	Class string // extra classes for the bubble
}

// Tooltip wraps its children (the trigger) and reveals a small bubble
// above them on hover or keyboard focus. Pure CSS positioning: the bubble
// is centered over the trigger, so it suits short labels on compact
// triggers (icons, buttons) rather than long prose.
func Tooltip(p TooltipProps, trigger ...any) *g.Node {
	all := []any{
		g.Class("group/tip relative inline-flex"),
		g.Data("slot", "tooltip"),
	}
	all = append(all, trigger...)
	all = append(all,
		g.Span(
			g.Class(style.CN(
				"pointer-events-none absolute bottom-full left-1/2 z-50 mb-1.5 -translate-x-1/2 whitespace-nowrap rounded-md bg-foreground px-2.5 py-1 text-xs font-medium text-background opacity-0 shadow-md transition-opacity duration-150 group-hover/tip:opacity-100 group-focus-within/tip:opacity-100",
				p.Class)),
			g.Role("tooltip"),
			g.Data("slot", "tooltip-bubble"),
			p.Label,
		),
	)
	return g.Span(all...)
}
