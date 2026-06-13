package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type HoverCardProps struct {
	Class string // extra classes for the panel
}

// HoverCard reveals a rich panel when the trigger is hovered or focused —
// like Tooltip but for more than a short label. Pure CSS: the panel is
// positioned below the trigger and fades in on hover.
//
//	ui.HoverCard(ui.HoverCardProps{}, ui.Button(...), ui.Card(...))
func HoverCard(p HoverCardProps, trigger any, content ...any) *g.Node {
	panel := append([]any{
		g.Class(style.CN(
			"invisible absolute left-1/2 top-full z-50 mt-2 w-64 -translate-x-1/2 rounded-md border bg-popover p-4 text-popover-foreground opacity-0 shadow-md outline-none transition-opacity duration-150 group-hover/hc:visible group-hover/hc:opacity-100 group-focus-within/hc:visible group-focus-within/hc:opacity-100",
			p.Class)),
		g.Data("slot", "hover-card-content"),
		g.Role("dialog"),
	}, content...)

	return g.Span(
		g.Class("group/hc relative inline-flex"),
		g.Data("slot", "hover-card"),
		g.Span(g.Data("slot", "hover-card-trigger"), trigger),
		g.Div(panel...),
	)
}
