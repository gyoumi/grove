package gallery

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/ui"
)

func disclosureSection() *g.Node {
	return section("Disclosure",
		ui.Tabs(ui.TabsProps{},
			ui.Tab{Value: "overview", Label: "Overview", Content: []any{ui.TypographyMuted("The overview panel.")}},
			ui.Tab{Value: "activity", Label: "Activity", Content: []any{ui.TypographyMuted("Recent activity here.")}},
			ui.Tab{Value: "settings", Label: "Settings", Content: []any{ui.TypographyMuted("Settings live here.")}},
		),
		ui.Accordion(ui.AccordionProps{DefaultOpen: []string{"a"}},
			ui.AccordionItem{Value: "a", Title: "What is this?", Content: []any{g.Span("A grove component gallery.")}},
			ui.AccordionItem{Value: "b", Title: "How does it work?", Content: []any{g.Span("Every section renders a live component.")}},
		),
		ui.Collapsible(ui.CollapsibleProps{},
			g.Span(g.Class("text-sm font-medium"), "Advanced options ", ui.Icon("chevron-down", "size-4")),
			ui.TypographyMuted("Hidden details revealed on toggle."),
		),
	)
}
