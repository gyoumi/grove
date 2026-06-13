package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type SheetSide string

const (
	SheetRight  SheetSide = "right"
	SheetLeft   SheetSide = "left"
	SheetTop    SheetSide = "top"
	SheetBottom SheetSide = "bottom"
)

var sheetSideClasses = map[SheetSide]string{
	SheetRight:  "inset-y-0 right-0 h-full w-3/4 border-l sm:max-w-sm",
	SheetLeft:   "inset-y-0 left-0 h-full w-3/4 border-r sm:max-w-sm",
	SheetTop:    "inset-x-0 top-0 border-b",
	SheetBottom: "inset-x-0 bottom-0 border-t",
}

type SheetProps struct {
	Open    bool
	OnClose func()
	Side    SheetSide // default right
	Class   string
}

type sheetArgs struct {
	p        SheetProps
	children []any
}

// Sheet is a panel that slides in from an edge: a modal overlay plus a
// side/top/bottom panel with focus capture and Escape/overlay dismissal.
func Sheet(p SheetProps, children ...any) *g.Node {
	return g.C(sheetView, sheetArgs{p: p, children: children})
}

func sheetView(a sheetArgs) *g.Node {
	side := a.p.Side
	if side == "" {
		side = SheetRight
	}
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
			g.Data("slot", "sheet-overlay"),
			g.OnClick(func(*g.Event) { close() }),
		),
		g.Div(
			g.Class(style.CN("fixed z-50 flex flex-col gap-4 bg-background p-6 shadow-lg", sheetSideClasses[side], a.p.Class)),
			g.Role("dialog"),
			g.Aria("modal", "true"),
			g.TabIndex(-1),
			g.Data("slot", "sheet-content"),
			g.Data("side", string(side)),
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
			a.children,
		),
	)
}

func SheetHeader(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col gap-1.5 text-center sm:text-left"), g.Data("slot", "sheet-header")}, args...)
	return g.Div(all...)
}

func SheetTitle(args ...any) *g.Node {
	all := append([]any{g.Class("text-lg font-semibold text-foreground"), g.Data("slot", "sheet-title")}, args...)
	return g.H2(all...)
}

func SheetDescription(args ...any) *g.Node {
	all := append([]any{g.Class("text-sm text-muted-foreground"), g.Data("slot", "sheet-description")}, args...)
	return g.P(all...)
}

func SheetFooter(args ...any) *g.Node {
	all := append([]any{g.Class("mt-auto flex flex-col gap-2"), g.Data("slot", "sheet-footer")}, args...)
	return g.Div(all...)
}
