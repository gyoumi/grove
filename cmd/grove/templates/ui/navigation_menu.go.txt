package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

const navTriggerClass = "inline-flex h-9 items-center justify-center gap-1 rounded-md px-3 py-2 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground focus:outline-none"

// NavMenuItem is one entry: a Label, an optional Href for a plain link, and an
// optional Panel shown on hover (which makes it a dropdown trigger instead).
type NavMenuItem struct {
	Label any
	Href  string
	Panel []any
}

type NavigationMenuProps struct {
	Class string
}

// NavigationMenu is a horizontal navigation bar whose items can reveal a rich
// panel on hover. Items without a Panel are plain links.
func NavigationMenu(p NavigationMenuProps, items ...NavMenuItem) *g.Node {
	entries := make([]any, 0, len(items))
	for _, it := range items {
		if len(it.Panel) == 0 {
			entries = append(entries, g.El("li",
				g.El("a", g.Class(navTriggerClass), g.Data("slot", "navigation-menu-link"), g.Attr("href", it.Href), it.Label),
			))
			continue
		}
		panel := append([]any{
			g.Class("invisible absolute left-0 top-full z-50 mt-1.5 min-w-[12rem] rounded-md border bg-popover p-2 text-popover-foreground opacity-0 shadow-md transition-opacity duration-150 group-hover/nav:visible group-hover/nav:opacity-100 group-focus-within/nav:visible group-focus-within/nav:opacity-100"),
			g.Data("slot", "navigation-menu-content"),
		}, it.Panel...)
		entries = append(entries, g.El("li",
			g.Class("group/nav relative"),
			g.Button(
				g.Type("button"),
				g.Class(navTriggerClass),
				g.Data("slot", "navigation-menu-trigger"),
				g.Span(it.Label),
				Icon("chevron-down", "size-3.5 transition-transform duration-200 group-hover/nav:rotate-180"),
			),
			g.Div(panel...),
		))
	}

	return g.El("nav",
		g.Class(style.CN("relative flex items-center", p.Class)),
		g.Data("slot", "navigation-menu"),
		g.Role("navigation"),
		g.El("ul", append([]any{g.Class("flex items-center gap-1")}, entries...)...),
	)
}
