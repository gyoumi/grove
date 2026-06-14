package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// sidebarCtx shares the open/closed state so a SidebarTrigger anywhere in the
// tree can toggle the Sidebar.
var sidebarCtx = g.NewContext[*sidebarState](nil)

type sidebarState struct {
	open   bool
	toggle func()
}

// SidebarProvider wraps an app shell and owns the sidebar's open state. Put a
// Sidebar and a SidebarInset (the main area) inside it.
//
//	ui.SidebarProvider(
//	    ui.Sidebar(ui.SidebarContent(ui.SidebarMenu(...))),
//	    ui.SidebarInset(ui.SidebarTrigger(), ...page...),
//	)
func SidebarProvider(args ...any) *g.Node { return g.C(sidebarProviderView, args) }

func sidebarProviderView(args []any) *g.Node {
	open, setOpen := g.UseState(true)
	st := &sidebarState{open: open, toggle: func() { setOpen(!open) }}
	all := append([]any{g.Class("flex w-full"), g.Data("slot", "sidebar-wrapper")}, args...)
	return sidebarCtx.Provider(st, g.Div(all...))
}

// Sidebar is the collapsible side panel; it slides to zero width when closed.
func Sidebar(args ...any) *g.Node { return g.C(sidebarView, args) }

func sidebarView(args []any) *g.Node {
	st := g.UseContext(sidebarCtx)
	open := st == nil || st.open
	w := "w-64"
	if !open {
		w = "w-0 overflow-hidden border-r-0"
	}
	all := append([]any{
		g.Class(style.CN("flex shrink-0 flex-col gap-2 border-r bg-card transition-[width] duration-200", w)),
		g.Data("slot", "sidebar"),
		g.Data("state", openState(open)),
	}, args...)
	return g.El("aside", all...)
}

// SidebarTrigger toggles the sidebar.
func SidebarTrigger(class ...string) *g.Node { return g.C(sidebarTriggerView, class) }

func sidebarTriggerView(class []string) *g.Node {
	st := g.UseContext(sidebarCtx)
	return g.Button(
		g.Type("button"),
		g.Class(style.CN("inline-flex size-8 items-center justify-center rounded-md text-foreground transition-colors hover:bg-accent hover:text-accent-foreground", class)),
		g.Data("slot", "sidebar-trigger"),
		g.Aria("label", "Toggle sidebar"),
		g.OnClick(func(*g.Event) {
			if st != nil {
				st.toggle()
			}
		}),
		Icon("panel-left", "size-4"),
	)
}

func SidebarHeader(args ...any) *g.Node {
	return sidebarPart("flex flex-col gap-1 p-2", "sidebar-header", args)
}
func SidebarContent(args ...any) *g.Node {
	return sidebarPart("flex flex-1 flex-col gap-1 overflow-auto p-2", "sidebar-content", args)
}
func SidebarFooter(args ...any) *g.Node {
	return sidebarPart("mt-auto flex flex-col gap-1 p-2", "sidebar-footer", args)
}

func SidebarMenu(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col gap-0.5"), g.Data("slot", "sidebar-menu")}, args...)
	return g.El("ul", all...)
}

func SidebarMenuItem(args ...any) *g.Node {
	all := append([]any{g.Data("slot", "sidebar-menu-item")}, args...)
	return g.El("li", all...)
}

// SidebarMenuButton is a navigation row; active highlights the current page.
func SidebarMenuButton(active bool, args ...any) *g.Node {
	all := append([]any{
		g.Type("button"),
		g.Class(style.CN(
			"flex w-full items-center gap-2 whitespace-nowrap rounded-md px-2 py-1.5 text-left text-sm transition-colors hover:bg-accent hover:text-accent-foreground [&>svg]:size-4 [&>svg]:shrink-0",
			map[string]bool{"bg-accent font-medium text-accent-foreground": active},
		)),
		g.Data("slot", "sidebar-menu-button"),
		g.AttrIf(active, "data-active", "1"),
	}, args...)
	return g.Button(all...)
}

// SidebarInset is the main content area beside the sidebar.
func SidebarInset(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-1 flex-col"), g.Data("slot", "sidebar-inset")}, args...)
	return g.El("main", all...)
}

func sidebarPart(base, slot string, args []any) *g.Node {
	all := append([]any{g.Class(base), g.Data("slot", slot)}, args...)
	return g.Div(all...)
}
