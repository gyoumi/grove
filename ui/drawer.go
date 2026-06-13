package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type DrawerProps struct {
	Open    bool
	OnClose func()
	Class   string
}

type drawerArgs struct {
	p        DrawerProps
	children []any
}

// Drawer is a panel that rises from the bottom with a grab handle — a
// touch-friendly bottom sheet. Overlay and Escape dismiss it.
func Drawer(p DrawerProps, children ...any) *g.Node {
	return g.C(drawerView, drawerArgs{p: p, children: children})
}

func drawerView(a drawerArgs) *g.Node {
	contentRef := g.UseRef[any](nil)
	g.UseEffect(func() func() {
		if !a.p.Open {
			return nil
		}
		return trapFocus(contentRef)
	}, []any{a.p.Open})

	if !a.p.Open {
		return nil
	}
	close := func() {
		if a.p.OnClose != nil {
			a.p.OnClose()
		}
	}
	return g.Fragment(
		g.Div(
			g.Class("fixed inset-0 z-50 bg-black/80"),
			g.Data("slot", "drawer-overlay"),
			g.OnClick(func(*g.Event) { close() }),
		),
		g.Div(
			g.Class(style.CN("fixed inset-x-0 bottom-0 z-50 mt-24 flex max-h-[80vh] flex-col gap-4 rounded-t-xl border bg-background p-6 shadow-lg", a.p.Class)),
			g.Role("dialog"),
			g.Aria("modal", "true"),
			g.TabIndex(-1),
			g.Data("slot", "drawer-content"),
			g.BindRef(contentRef),
			g.OnClick(func(e *g.Event) { e.StopPropagation() }),
			g.OnKeyDown(func(e *g.Event) {
				switch e.Key() {
				case "Escape":
					close()
				case "Tab":
					cycleFocus(contentRef, e)
				}
			}),
			g.Div(g.Class("mx-auto -mt-2 mb-2 h-1.5 w-12 shrink-0 rounded-full bg-muted"), g.Data("slot", "drawer-handle")),
			a.children,
		),
	)
}

func DrawerHeader(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col gap-1.5 text-center"), g.Data("slot", "drawer-header")}, args...)
	return g.Div(all...)
}

func DrawerTitle(args ...any) *g.Node {
	all := append([]any{g.Class("text-lg font-semibold"), g.Data("slot", "drawer-title")}, args...)
	return g.H2(all...)
}

func DrawerDescription(args ...any) *g.Node {
	all := append([]any{g.Class("text-sm text-muted-foreground"), g.Data("slot", "drawer-description")}, args...)
	return g.P(all...)
}

func DrawerFooter(args ...any) *g.Node {
	all := append([]any{g.Class("mt-auto flex flex-col gap-2"), g.Data("slot", "drawer-footer")}, args...)
	return g.Div(all...)
}
