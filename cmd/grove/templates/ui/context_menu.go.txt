package ui

import (
	"fmt"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type ContextMenuProps struct {
	Class string
}

type contextMenuArgs struct {
	p       ContextMenuProps
	trigger any
	items   []any
}

// ContextMenu opens a menu at the pointer on right-click over its trigger
// region. Compose items from DropdownItem/DropdownSeparator/DropdownLabel —
// selecting one closes the menu.
//
//	ui.ContextMenu(ui.ContextMenuProps{}, area,
//	    ui.DropdownItem(ui.DropdownItemProps{OnSelect: cut}, "Cut"),
//	)
func ContextMenu(p ContextMenuProps, trigger any, items ...any) *g.Node {
	return g.C(contextMenuView, contextMenuArgs{p: p, trigger: trigger, items: items})
}

func contextMenuView(a contextMenuArgs) *g.Node {
	open, setOpen := g.UseState(false)
	pos, setPos := g.UseState([2]float64{})

	body := []any{
		g.Data("slot", "context-menu"),
		g.Div(
			g.Data("slot", "context-menu-trigger"),
			g.On("contextmenu", func(e *g.Event) {
				e.PreventDefault()
				setPos([2]float64{e.Num("clientX"), e.Num("clientY")})
				setOpen(true)
			}),
			a.trigger,
		),
	}
	if open {
		body = append(body, g.Fragment(
			g.Div(
				g.Class("fixed inset-0 z-40"),
				g.Data("slot", "context-menu-overlay"),
				g.OnClick(func(*g.Event) { setOpen(false) }),
				g.On("contextmenu", func(e *g.Event) { e.PreventDefault(); setOpen(false) }),
			),
			g.Div(
				g.Class(style.CN("fixed z-50 flex min-w-[8rem] flex-col gap-0.5 rounded-md border bg-popover p-1 text-popover-foreground shadow-md", a.p.Class)),
				g.Role("menu"),
				g.Data("slot", "context-menu-content"),
				g.Attr("style", fmt.Sprintf("left:%gpx;top:%gpx", pos[0], pos[1])),
				menuCloseCtx.Provider(func() { setOpen(false) }, a.items...),
			),
		))
	}
	return g.Div(body...)
}
