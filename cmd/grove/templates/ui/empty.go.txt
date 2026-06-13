package ui

import g "github.com/gyoumi/grove"

// Empty and its parts render an empty-state placeholder:
//
//	ui.Empty(
//	    ui.EmptyHeader(
//	        ui.EmptyMedia(ui.Icon("search", "size-6")),
//	        ui.EmptyTitle("No results"),
//	        ui.EmptyDescription("Try a different search."),
//	    ),
//	    ui.EmptyContent(ui.Button(ui.ButtonProps{}, "Clear filters")),
//	)
func Empty(args ...any) *g.Node {
	all := append([]any{
		g.Class("flex min-h-48 flex-col items-center justify-center gap-4 rounded-lg border border-dashed p-8 text-center"),
		g.Data("slot", "empty"),
	}, args...)
	return g.Div(all...)
}

func EmptyHeader(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col items-center gap-1.5"), g.Data("slot", "empty-header")}, args...)
	return g.Div(all...)
}

func EmptyMedia(args ...any) *g.Node {
	all := append([]any{
		g.Class("mb-2 flex size-12 items-center justify-center rounded-full bg-muted text-muted-foreground"),
		g.Data("slot", "empty-media"),
	}, args...)
	return g.Div(all...)
}

func EmptyTitle(args ...any) *g.Node {
	all := append([]any{g.Class("text-base font-medium"), g.Data("slot", "empty-title")}, args...)
	return g.Div(all...)
}

func EmptyDescription(args ...any) *g.Node {
	all := append([]any{g.Class("text-sm text-muted-foreground"), g.Data("slot", "empty-description")}, args...)
	return g.El("p", all...)
}

func EmptyContent(args ...any) *g.Node {
	all := append([]any{g.Class("flex items-center gap-2"), g.Data("slot", "empty-content")}, args...)
	return g.Div(all...)
}
