package gallery

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/ui"
)

func resizableSection() *g.Node {
	return section("Resizable",
		g.Div(g.Class("h-48"),
			ui.Resizable(ui.ResizableProps{DefaultFraction: 0.4},
				g.Div(g.Class("flex h-full items-center justify-center bg-muted/30 text-sm text-muted-foreground"), "Panel"),
				g.Div(g.Class("flex h-full items-center justify-center text-sm text-muted-foreground"), "Drag the divider →"),
			),
		),
	)
}
