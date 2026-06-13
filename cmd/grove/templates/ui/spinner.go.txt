package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// Spinner is a spinning loader icon that inherits the current text color.
// Size it with classes: ui.Spinner("size-6").
func Spinner(class ...string) *g.Node {
	return g.Span(
		g.Class(style.CN("inline-flex text-muted-foreground", []string(class))),
		g.Data("slot", "spinner"),
		g.Role("status"),
		g.Attr("aria-label", "Loading"),
		Icon("loader", "size-4 animate-spin"),
	)
}
