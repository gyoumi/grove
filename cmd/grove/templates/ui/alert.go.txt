package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type AlertVariant string

const (
	AlertDefault     AlertVariant = "default"
	AlertDestructive AlertVariant = "destructive"
)

var alertVariants = style.Variants{
	Base: "relative w-full rounded-lg border px-4 py-3 text-sm",
	Groups: map[string]map[string]string{
		"variant": {
			"default":     "bg-background text-foreground",
			"destructive": "border-destructive/50 text-destructive",
		},
	},
	Defaults: map[string]string{"variant": "default"},
}

type AlertProps struct {
	Variant AlertVariant
	Class   string
}

// Alert renders a callout box:
//
//	ui.Alert(ui.AlertProps{}, ui.AlertTitle("Heads up"), ui.AlertDescription("..."))
func Alert(p AlertProps, children ...any) *g.Node {
	args := []any{
		g.Class(alertVariants.Class(map[string]string{"variant": string(p.Variant)}, p.Class)),
		g.Role("alert"),
		g.Data("slot", "alert"),
	}
	return g.Div(append(args, children...)...)
}

func AlertTitle(args ...any) *g.Node {
	all := append([]any{g.Class("mb-1 font-medium leading-none tracking-tight"), g.Data("slot", "alert-title")}, args...)
	return g.H5(all...)
}

func AlertDescription(args ...any) *g.Node {
	all := append([]any{g.Class("text-sm [&_p]:leading-relaxed"), g.Data("slot", "alert-description")}, args...)
	return g.Div(all...)
}
