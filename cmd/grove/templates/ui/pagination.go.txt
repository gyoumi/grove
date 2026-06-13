package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

const paginationLinkBase = "inline-flex h-9 min-w-9 items-center justify-center gap-1 rounded-md px-3 text-sm font-medium transition-colors hover:bg-accent hover:text-accent-foreground"

// Pagination and its parts render page navigation:
//
//	ui.Pagination(ui.PaginationContent(
//	    ui.PaginationItem(ui.PaginationPrevious("#prev")),
//	    ui.PaginationItem(ui.PaginationLink("#1", true, "1")),
//	    ui.PaginationItem(ui.PaginationNext("#next")),
//	))
func Pagination(args ...any) *g.Node {
	all := append([]any{
		g.Class("mx-auto flex w-full justify-center"),
		g.Data("slot", "pagination"),
		g.Role("navigation"),
		g.Attr("aria-label", "pagination"),
	}, args...)
	return g.El("nav", all...)
}

func PaginationContent(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-row items-center gap-1"), g.Data("slot", "pagination-content")}, args...)
	return g.El("ul", all...)
}

func PaginationItem(args ...any) *g.Node {
	all := append([]any{g.Data("slot", "pagination-item")}, args...)
	return g.El("li", all...)
}

// PaginationLink is a page number link; active marks the current page.
func PaginationLink(href string, active bool, children ...any) *g.Node {
	cls := style.CN(paginationLinkBase, map[string]bool{"border border-input bg-background": active})
	all := append([]any{
		g.Class(cls),
		g.Data("slot", "pagination-link"),
		g.AttrIf(active, "aria-current", "page"),
		g.Attr("href", href),
	}, children...)
	return g.El("a", all...)
}

// PaginationPrevious and PaginationNext are the directional controls.
func PaginationPrevious(href string) *g.Node {
	return g.El("a",
		g.Class(style.CN(paginationLinkBase, "px-2.5")),
		g.Data("slot", "pagination-previous"),
		g.Attr("href", href),
		g.Attr("aria-label", "Go to previous page"),
		Icon("chevron-left", "size-4"),
		g.Span("Previous"),
	)
}

func PaginationNext(href string) *g.Node {
	return g.El("a",
		g.Class(style.CN(paginationLinkBase, "px-2.5")),
		g.Data("slot", "pagination-next"),
		g.Attr("href", href),
		g.Attr("aria-label", "Go to next page"),
		g.Span("Next"),
		Icon("chevron-right", "size-4"),
	)
}

func PaginationEllipsis() *g.Node {
	return g.Span(
		g.Class("flex size-9 items-center justify-center"),
		g.Data("slot", "pagination-ellipsis"),
		g.Attr("aria-hidden", "true"),
		Icon("more-horizontal", "size-4"),
	)
}
