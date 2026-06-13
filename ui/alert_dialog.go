package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type AlertDialogProps struct {
	Open    bool
	OnClose func() // called on Escape; the overlay does not dismiss
	Class   string
}

type alertDialogArgs struct {
	p        AlertDialogProps
	children []any
}

// AlertDialog is a modal that interrupts for a confirmation: unlike Dialog it
// does not dismiss on an overlay click, so the user must pick an action.
// Compose it from AlertDialogHeader/Title/Description/Footer with
// AlertDialogAction and AlertDialogCancel.
func AlertDialog(p AlertDialogProps, children ...any) *g.Node {
	return g.C(alertDialogView, alertDialogArgs{p: p, children: children})
}

func alertDialogView(a alertDialogArgs) *g.Node {
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
		g.Div(g.Class("fixed inset-0 z-50 bg-black/80"), g.Data("slot", "alert-dialog-overlay")),
		g.Div(
			g.Class(style.CN("fixed left-1/2 top-1/2 z-50 grid w-full max-w-lg -translate-x-1/2 -translate-y-1/2 gap-4 border bg-background p-6 shadow-lg sm:rounded-lg", a.p.Class)),
			g.Role("alertdialog"),
			g.Aria("modal", "true"),
			g.TabIndex(-1),
			g.Data("slot", "alert-dialog-content"),
			g.BindRef(contentRef),
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

func AlertDialogHeader(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col gap-2 text-center sm:text-left"), g.Data("slot", "alert-dialog-header")}, args...)
	return g.Div(all...)
}

func AlertDialogTitle(args ...any) *g.Node {
	all := append([]any{g.Class("text-lg font-semibold"), g.Data("slot", "alert-dialog-title")}, args...)
	return g.H2(all...)
}

func AlertDialogDescription(args ...any) *g.Node {
	all := append([]any{g.Class("text-sm text-muted-foreground"), g.Data("slot", "alert-dialog-description")}, args...)
	return g.P(all...)
}

func AlertDialogFooter(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"), g.Data("slot", "alert-dialog-footer")}, args...)
	return g.Div(all...)
}

// AlertDialogAction is the confirming button (primary styling).
func AlertDialogAction(onClick func(*g.Event), children ...any) *g.Node {
	all := []any{
		g.Class("inline-flex h-9 items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow transition-colors hover:bg-primary/90"),
		g.Type("button"),
		g.Data("slot", "alert-dialog-action"),
		g.OnClick(onClick),
	}
	return g.Button(append(all, children...)...)
}

// AlertDialogCancel is the dismissing button (outline styling).
func AlertDialogCancel(onClick func(*g.Event), children ...any) *g.Node {
	all := []any{
		g.Class("inline-flex h-9 items-center justify-center rounded-md border border-input bg-background px-4 py-2 text-sm font-medium shadow-sm transition-colors hover:bg-accent hover:text-accent-foreground"),
		g.Type("button"),
		g.Data("slot", "alert-dialog-cancel"),
		g.OnClick(onClick),
	}
	return g.Button(append(all, children...)...)
}
