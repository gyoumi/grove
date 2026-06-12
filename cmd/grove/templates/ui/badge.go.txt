package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type BadgeVariant string

const (
	BadgeDefault     BadgeVariant = "default"
	BadgeSecondary   BadgeVariant = "secondary"
	BadgeDestructive BadgeVariant = "destructive"
	BadgeOutline     BadgeVariant = "outline"
)

var badgeVariants = style.Variants{
	Base: "inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
	Groups: map[string]map[string]string{
		"variant": {
			"default":     "border-transparent bg-primary text-primary-foreground shadow hover:bg-primary/80",
			"secondary":   "border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80",
			"destructive": "border-transparent bg-destructive text-destructive-foreground shadow hover:bg-destructive/80",
			"outline":     "text-foreground",
		},
	},
	Defaults: map[string]string{"variant": "default"},
}

type BadgeProps struct {
	Variant BadgeVariant
	Class   string
}

// Badge renders a small themed status label.
func Badge(p BadgeProps, children ...any) *g.Node {
	args := []any{
		g.Class(badgeVariants.Class(map[string]string{"variant": string(p.Variant)}, p.Class)),
		g.Data("slot", "badge"),
	}
	return g.Span(append(args, children...)...)
}
