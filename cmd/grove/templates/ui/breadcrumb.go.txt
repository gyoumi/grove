package ui

import g "github.com/gyoumi/grove"

// Breadcrumb and its parts render a navigation trail:
//
//	ui.Breadcrumb(ui.BreadcrumbList(
//	    ui.BreadcrumbItem(ui.BreadcrumbLink("/", "Home")),
//	    ui.BreadcrumbSeparator(),
//	    ui.BreadcrumbItem(ui.BreadcrumbPage("Settings")),
//	))
func Breadcrumb(args ...any) *g.Node {
	all := append([]any{g.Data("slot", "breadcrumb"), g.Attr("aria-label", "breadcrumb")}, args...)
	return g.El("nav", all...)
}

func BreadcrumbList(args ...any) *g.Node {
	all := append([]any{
		g.Class("flex flex-wrap items-center gap-1.5 break-words text-sm text-muted-foreground"),
		g.Data("slot", "breadcrumb-list"),
	}, args...)
	return g.El("ol", all...)
}

func BreadcrumbItem(args ...any) *g.Node {
	all := append([]any{g.Class("inline-flex items-center gap-1.5"), g.Data("slot", "breadcrumb-item")}, args...)
	return g.El("li", all...)
}

// BreadcrumbLink is a navigable crumb. href is set on an anchor.
func BreadcrumbLink(href string, children ...any) *g.Node {
	all := append([]any{
		g.Class("transition-colors hover:text-foreground"),
		g.Data("slot", "breadcrumb-link"),
		g.Attr("href", href),
	}, children...)
	return g.El("a", all...)
}

// BreadcrumbPage is the current (non-link) crumb.
func BreadcrumbPage(children ...any) *g.Node {
	all := append([]any{
		g.Class("font-normal text-foreground"),
		g.Data("slot", "breadcrumb-page"),
		g.Role("link"),
		g.Attr("aria-current", "page"),
	}, children...)
	return g.Span(all...)
}

// BreadcrumbSeparator is the divider between crumbs (a chevron by default).
func BreadcrumbSeparator(children ...any) *g.Node {
	body := children
	if len(body) == 0 {
		body = []any{Icon("chevron-right", "size-3.5")}
	}
	all := append([]any{
		g.Class("[&>svg]:size-3.5 text-muted-foreground/60"),
		g.Data("slot", "breadcrumb-separator"),
		g.Role("presentation"),
		g.Attr("aria-hidden", "true"),
	}, body...)
	return g.El("li", all...)
}

// BreadcrumbEllipsis collapses overflowing crumbs.
func BreadcrumbEllipsis() *g.Node {
	return g.Span(
		g.Class("flex size-9 items-center justify-center"),
		g.Data("slot", "breadcrumb-ellipsis"),
		g.Role("presentation"),
		g.Attr("aria-hidden", "true"),
		Icon("more-horizontal", "size-4"),
	)
}
