package ui

import g "github.com/gyoumi/grove"

// Item and its parts render a horizontal list row with optional leading media
// and trailing actions:
//
//	ui.Item(
//	    ui.ItemMedia(ui.Avatar(...)),
//	    ui.ItemContent(ui.ItemTitle("Ada"), ui.ItemDescription("ada@x.io")),
//	    ui.ItemActions(ui.Button(...)),
//	)
func Item(args ...any) *g.Node {
	all := append([]any{
		g.Class("flex items-center gap-3 rounded-md border bg-card p-3 text-card-foreground"),
		g.Data("slot", "item"),
	}, args...)
	return g.Div(all...)
}

func ItemMedia(args ...any) *g.Node {
	all := append([]any{g.Class("flex shrink-0 items-center justify-center text-muted-foreground"), g.Data("slot", "item-media")}, args...)
	return g.Div(all...)
}

func ItemContent(args ...any) *g.Node {
	all := append([]any{g.Class("flex min-w-0 flex-1 flex-col gap-0.5"), g.Data("slot", "item-content")}, args...)
	return g.Div(all...)
}

func ItemTitle(args ...any) *g.Node {
	all := append([]any{g.Class("truncate text-sm font-medium leading-none"), g.Data("slot", "item-title")}, args...)
	return g.Div(all...)
}

func ItemDescription(args ...any) *g.Node {
	all := append([]any{g.Class("truncate text-sm text-muted-foreground"), g.Data("slot", "item-description")}, args...)
	return g.Div(all...)
}

func ItemActions(args ...any) *g.Node {
	all := append([]any{g.Class("ml-auto flex shrink-0 items-center gap-2"), g.Data("slot", "item-actions")}, args...)
	return g.Div(all...)
}
