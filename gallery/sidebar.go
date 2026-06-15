package gallery

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/ui"
)

func sidebarSection() *g.Node {
	return section("Sidebar", sidebarDemo())
}

// sidebarDemo is a mini app shell with a collapsible sidebar.
func sidebarDemo() *g.Node {
	return g.Div(g.Class("h-64 overflow-hidden rounded-lg border"),
		ui.SidebarProvider(
			ui.Sidebar(
				ui.SidebarHeader(g.Span(g.Class("px-2 text-sm font-semibold"), "Acme")),
				ui.SidebarContent(ui.SidebarMenu(
					ui.SidebarMenuItem(ui.SidebarMenuButton(true, ui.Icon("calendar"), "Dashboard")),
					ui.SidebarMenuItem(ui.SidebarMenuButton(false, ui.Icon("search"), "Search")),
					ui.SidebarMenuItem(ui.SidebarMenuButton(false, ui.Icon("clock"), "Recent")),
				)),
			),
			ui.SidebarInset(
				g.Div(g.Class("flex items-center gap-2 border-b p-2"),
					ui.SidebarTrigger(),
					g.Span(g.Class("text-sm font-medium"), "Dashboard"),
				),
				g.Div(g.Class("p-4 text-sm text-muted-foreground"), "Toggle the sidebar with the panel button."),
			),
		),
	)
}
