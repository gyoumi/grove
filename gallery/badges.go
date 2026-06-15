package gallery

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/ui"
)

func badgesSection() *g.Node {
	return section("Badges",
		g.Div(g.Class("flex flex-wrap items-center gap-2"),
			ui.Badge(ui.BadgeProps{}, "Default"),
			ui.Badge(ui.BadgeProps{Variant: ui.BadgeSecondary}, "Secondary"),
			ui.Badge(ui.BadgeProps{Variant: ui.BadgeDestructive}, "Destructive"),
			ui.Badge(ui.BadgeProps{Variant: ui.BadgeOutline}, "Outline"),
			ui.Kbd("⌘"), ui.Kbd("K"),
		),
	)
}
