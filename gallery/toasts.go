package gallery

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/ui"
)

func toastsSection() *g.Node {
	return section("Toasts",
		g.Div(g.Class("flex flex-wrap items-center gap-2"),
			ui.Button(ui.ButtonProps{Variant: ui.ButtonOutline, OnClick: func(*g.Event) {
				ui.Toast("Saved", ui.ToastOptions{Description: "Your changes are saved."})
			}}, "Show toast"),
			ui.Button(ui.ButtonProps{Variant: ui.ButtonOutline, OnClick: func(*g.Event) {
				ui.Toast("Event created", ui.ToastOptions{Description: "Standup · 5 days", Variant: ui.ToastSuccess})
			}}, "Success"),
			ui.Button(ui.ButtonProps{Variant: ui.ButtonOutline, OnClick: func(*g.Event) {
				ui.Toast("Something went wrong", ui.ToastOptions{Description: "Could not save the event.", Variant: ui.ToastDestructive})
			}}, "Error"),
		),
	)
}
