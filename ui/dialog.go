package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type DialogProps struct {
	// Open controls visibility; the dialog is fully controlled by the
	// parent's state.
	Open bool
	// OnClose is called when the user dismisses the dialog (Escape, overlay
	// click) — set your open state to false here.
	OnClose func()
	// Class adds classes to the content panel.
	Class string
}

type dialogArgs struct {
	p        DialogProps
	children []any
}

// Dialog renders a modal dialog: overlay, centered panel, Escape and
// overlay-click dismissal, focus capture while open, and focus restore on
// close. Compose the panel from DialogHeader/DialogTitle/DialogDescription/
// DialogFooter like shadcn:
//
//	ui.Dialog(ui.DialogProps{Open: open, OnClose: func() { setOpen(false) }},
//	    ui.DialogHeader(ui.DialogTitle("Edit profile")),
//	    ui.DialogFooter(ui.Button(ui.ButtonProps{OnClick: save}, "Save")),
//	)
func Dialog(p DialogProps, children ...any) *g.Node {
	return g.C(dialog, dialogArgs{p: p, children: children})
}

func dialog(a dialogArgs) *g.Node {
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
			g.Data("slot", "dialog-overlay"),
			g.OnClick(func(*g.Event) { close() }),
		),
		g.Div(
			g.Class(style.CN("fixed left-1/2 top-1/2 z-50 grid w-full max-w-lg -translate-x-1/2 -translate-y-1/2 gap-4 border bg-background p-6 shadow-lg sm:rounded-lg", a.p.Class)),
			g.Role("dialog"),
			g.Aria("modal", "true"),
			g.TabIndex(-1),
			g.Data("slot", "dialog-content"),
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

func DialogHeader(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col gap-1.5 text-center sm:text-left"), g.Data("slot", "dialog-header")}, args...)
	return g.Div(all...)
}

func DialogTitle(args ...any) *g.Node {
	all := append([]any{g.Class("text-lg font-semibold leading-none tracking-tight"), g.Data("slot", "dialog-title")}, args...)
	return g.H2(all...)
}

func DialogDescription(args ...any) *g.Node {
	all := append([]any{g.Class("text-sm text-muted-foreground"), g.Data("slot", "dialog-description")}, args...)
	return g.P(all...)
}

func DialogFooter(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col-reverse gap-2 sm:flex-row sm:justify-end"), g.Data("slot", "dialog-footer")}, args...)
	return g.Div(all...)
}
