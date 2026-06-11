package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// Separator renders a horizontal rule; pass vertical=true for a vertical
// one (the parent needs a height, e.g. flex with h-*).
func Separator(vertical bool, args ...any) *g.Node {
	cls := "bg-border shrink-0 h-px w-full"
	orientation := "horizontal"
	if vertical {
		cls = "bg-border shrink-0 w-px h-full"
		orientation = "vertical"
	}
	all := []any{
		g.Class(style.CN(cls)),
		g.Role("none"),
		g.Data("slot", "separator"),
		g.Data("orientation", orientation),
	}
	return g.Div(append(all, args...)...)
}
