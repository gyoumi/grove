package ui

import (
	"strings"

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

// sheetSideAnim/sheetSideAnimOut are the enter/leave animations per side
// (defined in the app's CSS; grove init scaffolds the keyframes).
var sheetSideAnim = map[SheetSide]string{
	SheetRight:  "animate-slide-in-right",
	SheetLeft:   "animate-slide-in-left",
	SheetTop:    "animate-slide-in-top",
	SheetBottom: "animate-slide-in-bottom",
}

var sheetSideAnimOut = map[SheetSide]string{
	SheetRight:  "animate-slide-out-right",
	SheetLeft:   "animate-slide-out-left",
	SheetTop:    "animate-slide-out-top",
	SheetBottom: "animate-slide-out-bottom",
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
	// rendered keeps the panel mounted through the close-out animation;
	// closing flips it from slide-in to slide-out, and the panel unmounts
	// when that animation ends.
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
	panelAnim, overlayAnim := sheetSideAnim[side], "animate-overlay-in"
	if closing {
		panelAnim, overlayAnim = sheetSideAnimOut[side], "animate-overlay-out"
	}
	return g.Portal(
		g.Div(
			g.Class(style.CN("fixed inset-0 z-50 bg-black/80", overlayAnim)),
			g.Data("slot", "sheet-overlay"),
			g.OnClick(func(*g.Event) { close() }),
		),
		g.Div(
			g.Class(style.CN("fixed z-50 flex flex-col gap-4 bg-background p-6 shadow-lg", sheetSideClasses[side], panelAnim, a.p.Class)),
			g.Role("dialog"),
			g.Aria("modal", "true"),
			g.TabIndex(-1),
			g.Data("slot", "sheet-content"),
			g.Data("side", string(side)),
			g.Data("state", sheetState(closing)),
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
			a.children,
		),
	)
}

func sheetState(closing bool) string {
	if closing {
		return "closing"
	}
	return "open"
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
