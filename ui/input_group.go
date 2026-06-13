package ui

import g "github.com/gyoumi/grove"

// InputGroup wraps an input with leading/trailing addons (icons, text,
// buttons). The group carries the border and focus ring; any input dropped
// inside is rendered borderless:
//
//	ui.InputGroup(
//	    ui.InputGroupAddon(ui.Icon("search")),
//	    ui.Input(ui.InputProps{Placeholder: "Search…"}),
//	    ui.InputGroupAddon(ui.Kbd("⌘K")),
//	)
func InputGroup(args ...any) *g.Node {
	all := append([]any{
		g.Class("flex h-9 w-full items-center rounded-md border border-input bg-transparent text-sm shadow-sm transition-colors focus-within:ring-1 focus-within:ring-ring [&_input]:h-full [&_input]:border-0 [&_input]:bg-transparent [&_input]:shadow-none [&_input]:focus-visible:ring-0"),
		g.Data("slot", "input-group"),
	}, args...)
	return g.Div(all...)
}

// InputGroupAddon is a non-editable slot on either side of the input.
func InputGroupAddon(args ...any) *g.Node {
	all := append([]any{
		g.Class("flex items-center gap-1 px-3 text-muted-foreground [&>svg]:size-4"),
		g.Data("slot", "input-group-addon"),
	}, args...)
	return g.Div(all...)
}

// InputGroupText is muted text inside an addon.
func InputGroupText(children ...any) *g.Node {
	all := append([]any{g.Class("text-sm text-muted-foreground"), g.Data("slot", "input-group-text")}, children...)
	return g.Span(all...)
}
