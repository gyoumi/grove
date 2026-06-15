package gallery

import (
	"fmt"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/ui"
)

func togglesSection() *g.Node {
	toggles, setToggles := g.UseState([]string{"bold"})
	return section("Toggles",
		g.Div(g.Class("flex flex-wrap items-center gap-4"),
			ui.ToggleGroup(ui.ToggleGroupProps{
				Multiple: true, Variant: ui.ToggleOutline, Value: toggles,
				OnChange: setToggles,
			},
				ui.ToggleGroupItem{Value: "bold", Children: []any{g.Strong("B")}},
				ui.ToggleGroupItem{Value: "italic", Children: []any{g.Em("I")}},
				ui.ToggleGroupItem{Value: "underline", Children: []any{g.Span(g.Class("underline"), "U")}},
			),
			ui.TypographyMuted(fmt.Sprintf("on: %v", toggles)),
		),
	)
}
