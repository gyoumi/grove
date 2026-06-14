package ui

import (
	"strings"

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
	// rendered keeps the panel mounted while it slides back down on close.
	rendered, setRendered := g.UseState(a.p.Open)
	closing, setClosing := g.UseState(false)

	g.UseEffect(func() func() {
		if a.p.Open {
			setRendered(true)
			setClosing(false)
		} else if rendered {
			setClosing(true)
		}
		return nil
	}, []any{a.p.Open})

	g.UseEffect(func() func() {
		if a.p.Open && rendered {
			return trapFocus(contentRef)
		}
		return nil
	}, []any{a.p.Open, rendered})

	if !rendered {
		return nil
	}
	close := func() {
		if a.p.OnClose != nil {
			a.p.OnClose()
		}
	}
	panelAnim, overlayAnim, state := "animate-slide-in-bottom", "animate-overlay-in", "open"
	if closing {
		panelAnim, overlayAnim, state = "animate-slide-out-bottom", "animate-overlay-out", "closing"
	}
	return g.Portal(
		g.Div(
			g.Class(style.CN("fixed inset-0 z-50 bg-black/80", overlayAnim)),
			g.Data("slot", "drawer-overlay"),
			g.OnClick(func(*g.Event) { close() }),
		),
		g.Div(
			g.Class(style.CN("fixed inset-x-0 bottom-0 z-50 mt-24 flex max-h-[80vh] flex-col gap-4 rounded-t-xl border bg-background p-6 shadow-lg", panelAnim, a.p.Class)),
			g.Role("dialog"),
			g.Aria("modal", "true"),
			g.TabIndex(-1),
			g.Data("slot", "drawer-content"),
			g.Data("state", state),
			g.BindRef(contentRef),
			g.OnClick(func(e *g.Event) { e.StopPropagation() }),
			g.On("animationend", func(e *g.Event) {
				if closing && strings.HasPrefix(e.Str("animationName"), "slide-out") {
					setRendered(false)
					setClosing(false)
				}
			}),
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
