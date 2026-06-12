package ui

import g "github.com/gyoumi/grove"

// Card composes a bordered, themed panel from named parts:
//
//	ui.Card(
//	    ui.CardHeader(ui.CardTitle("Title"), ui.CardDescription("...")),
//	    ui.CardContent(...),
//	    ui.CardFooter(...),
//	)
//
// Extra g.Class options append after the base classes.
func Card(args ...any) *g.Node {
	return part("rounded-xl border bg-card text-card-foreground shadow", "card", args)
}

func CardHeader(args ...any) *g.Node {
	return part("flex flex-col gap-1.5 p-6", "card-header", args)
}

func CardTitle(args ...any) *g.Node {
	return part("font-semibold leading-none tracking-tight", "card-title", args)
}

func CardDescription(args ...any) *g.Node {
	return part("text-sm text-muted-foreground", "card-description", args)
}

func CardContent(args ...any) *g.Node {
	return part("p-6 pt-0", "card-content", args)
}

func CardFooter(args ...any) *g.Node {
	return part("flex items-center p-6 pt-0", "card-footer", args)
}

func part(base, slot string, args []any) *g.Node {
	all := append([]any{g.Class(base), g.Data("slot", slot)}, args...)
	return g.Div(all...)
}
