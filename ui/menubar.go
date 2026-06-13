package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// MenubarMenu is one top-level menu: a Label and its Items (DropdownItem,
// DropdownSeparator, DropdownLabel).
type MenubarMenu struct {
	Label string
	Items []any
}

type MenubarProps struct {
	Class string
}

// Menubar is a horizontal bar of menus, each opening below its label. At most
// one menu is open at a time; selecting an item closes it.
func Menubar(p MenubarProps, menus ...MenubarMenu) *g.Node {
	return g.C(menubarView, menubarArgs{p: p, menus: menus})
}

type menubarArgs struct {
	p     MenubarProps
	menus []MenubarMenu
}

func menubarView(a menubarArgs) *g.Node {
	openIdx, setOpenIdx := g.UseState(-1)

	entries := make([]any, 0, len(a.menus))
	for i, m := range a.menus {
		isOpen := openIdx == i
		trigger := g.Button(
			g.Type("button"),
			g.Class(style.CN(
				"flex select-none items-center rounded-sm px-3 py-1 text-sm font-medium outline-none hover:bg-accent hover:text-accent-foreground",
				map[string]bool{"bg-accent text-accent-foreground": isOpen},
			)),
			g.Data("slot", "menubar-trigger"),
			g.OnClick(func(*g.Event) {
				if isOpen {
					setOpenIdx(-1)
				} else {
					setOpenIdx(i)
				}
			}),
			m.Label,
		)
		entries = append(entries, Popover(
			PopoverProps{
				Open:    isOpen,
				OnClose: func() { setOpenIdx(-1) },
				Align:   PopoverAlignStart,
				Class:   "flex min-w-[10rem] flex-col gap-0.5 p-1",
			},
			trigger,
			menuCloseCtx.Provider(func() { setOpenIdx(-1) }, m.Items...),
		))
	}

	all := append([]any{
		g.Class(style.CN("flex items-center gap-0.5 rounded-md border bg-background p-1 shadow-sm", a.p.Class)),
		g.Data("slot", "menubar"),
		g.Role("menubar"),
	}, entries...)
	return g.Div(all...)
}
