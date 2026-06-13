package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// ButtonGroup joins a row of buttons into a single segmented control:
// adjacent corners and borders are collapsed so they read as one unit.
//
//	ui.ButtonGroup(
//	    ui.Button(ui.ButtonProps{Variant: ui.ButtonOutline}, "Day"),
//	    ui.Button(ui.ButtonProps{Variant: ui.ButtonOutline}, "Week"),
//	)
func ButtonGroup(children ...any) *g.Node {
	args := []any{
		g.Class(style.CN(
			"inline-flex items-center [&>*]:rounded-none [&>*:first-child]:rounded-l-md [&>*:last-child]:rounded-r-md [&>*:not(:first-child)]:-ml-px",
		)),
		g.Data("slot", "button-group"),
		g.Role("group"),
	}
	return g.Div(append(args, children...)...)
}
