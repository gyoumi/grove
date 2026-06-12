package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// menuCloseCtx hands the menu's close function down to its items, so
// selecting an item closes the menu without wiring it through every prop.
var menuCloseCtx = g.NewContext[func()](nil)

type DropdownProps struct {
	// Open controls visibility; the menu is fully controlled.
	Open bool
	// OnClose fires on Escape, outside clicks, and after item selection.
	OnClose func()
	Align   PopoverAlign // default start
	Class   string       // extra classes for the menu panel
}

// Dropdown anchors a menu below its trigger. Compose it from DropdownItem,
// DropdownSeparator, and DropdownLabel; items close the menu when picked,
// and ArrowUp/ArrowDown move focus through the items in the browser.
//
//	ui.Dropdown(ui.DropdownProps{Open: open, OnClose: func() { setOpen(false) }},
//	    ui.Button(ui.ButtonProps{OnClick: func(*g.Event) { setOpen(!open) }}, "actions"),
//	    ui.DropdownItem(ui.DropdownItemProps{OnSelect: rename}, "Rename"),
//	    ui.DropdownSeparator(),
//	    ui.DropdownItem(ui.DropdownItemProps{OnSelect: del, Class: "text-destructive"}, "Delete"),
//	)
func Dropdown(p DropdownProps, trigger *g.Node, items ...any) *g.Node {
	align := p.Align
	if align == "" {
		align = PopoverAlignStart
	}
	close := func() {
		if p.OnClose != nil {
			p.OnClose()
		}
	}
	// a plain ref (not UseRef — Dropdown is not a component): handlers are
	// re-bound on every render, so each render's handlers see its own ref,
	// which BindRef points at the mounted menu element
	menuRef := &g.DOMRef{}

	content := []any{
		g.Role("menu"),
		g.BindRef(menuRef),
		// replaces the popover's keydown so it must handle Escape too
		g.OnKeyDown(func(e *g.Event) {
			switch e.Key() {
			case "Escape":
				close()
			case "ArrowDown":
				menuFocusMove(menuRef, e, 1)
			case "ArrowUp":
				menuFocusMove(menuRef, e, -1)
			}
		}),
		menuCloseCtx.Provider(close, items...),
	}

	return Popover(PopoverProps{
		Open:    p.Open,
		OnClose: p.OnClose,
		Side:    PopoverBottom,
		Align:   align,
		Class:   style.CN("flex min-w-[10rem] flex-col gap-0.5 p-1", p.Class),
	}, trigger, content...)
}

type DropdownItemProps struct {
	// OnSelect runs when the item is picked; the menu then closes itself.
	OnSelect func()
	Disabled bool
	Class    string // e.g. "text-destructive" for dangerous actions
}

type dropdownItemArgs struct {
	p        DropdownItemProps
	children []any
}

// DropdownItem is one selectable row of a Dropdown.
func DropdownItem(p DropdownItemProps, children ...any) *g.Node {
	return g.C(dropdownItem, dropdownItemArgs{p: p, children: children})
}

func dropdownItem(a dropdownItemArgs) *g.Node {
	closeMenu := g.UseContext(menuCloseCtx)
	all := []any{
		g.Class(style.CN(
			"flex w-full cursor-pointer select-none items-center gap-2 rounded-sm px-2 py-1.5 text-left text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground disabled:pointer-events-none disabled:opacity-50",
			a.p.Class)),
		g.Type("button"),
		g.Role("menuitem"),
		g.Data("slot", "dropdown-item"),
		g.Disabled(a.p.Disabled),
		g.OnClick(func(*g.Event) {
			if a.p.Disabled {
				return
			}
			if a.p.OnSelect != nil {
				a.p.OnSelect()
			}
			if closeMenu != nil {
				closeMenu()
			}
		}),
	}
	return g.Button(append(all, a.children...)...)
}

// DropdownSeparator draws a thin rule between item groups.
func DropdownSeparator() *g.Node {
	return g.Div(g.Class("-mx-1 my-1 h-px bg-border"), g.Role("none"), g.Data("slot", "dropdown-separator"))
}

// DropdownLabel is a non-interactive heading inside the menu.
func DropdownLabel(args ...any) *g.Node {
	all := append([]any{g.Class("px-2 py-1.5 text-xs font-medium text-muted-foreground"), g.Data("slot", "dropdown-label")}, args...)
	return g.Div(all...)
}
