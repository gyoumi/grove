package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// Skeleton is a pulsing placeholder block for content that is still loading.
// Size it with classes: ui.Skeleton("h-4 w-32").
func Skeleton(class ...string) *g.Node {
	return g.Div(
		g.Class(style.CN("animate-pulse rounded-md bg-muted", []string(class))),
		g.Data("slot", "skeleton"),
	)
}
