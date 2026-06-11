package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type ButtonVariant string

const (
	ButtonDefault     ButtonVariant = "default"
	ButtonDestructive ButtonVariant = "destructive"
	ButtonOutline     ButtonVariant = "outline"
	ButtonSecondary   ButtonVariant = "secondary"
	ButtonGhost       ButtonVariant = "ghost"
	ButtonLink        ButtonVariant = "link"
)

type ButtonSize string

const (
	ButtonSizeDefault ButtonSize = "default"
	ButtonSizeSm      ButtonSize = "sm"
	ButtonSizeLg      ButtonSize = "lg"
	ButtonSizeIcon    ButtonSize = "icon"
)

var buttonVariants = style.Variants{
	Base: "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50",
	Groups: map[string]map[string]string{
		"variant": {
			"default":     "bg-primary text-primary-foreground shadow hover:bg-primary/90",
			"destructive": "bg-destructive text-destructive-foreground shadow-sm hover:bg-destructive/90",
			"outline":     "border border-input bg-background shadow-sm hover:bg-accent hover:text-accent-foreground",
			"secondary":   "bg-secondary text-secondary-foreground shadow-sm hover:bg-secondary/80",
			"ghost":       "hover:bg-accent hover:text-accent-foreground",
			"link":        "text-primary underline-offset-4 hover:underline",
		},
		"size": {
			"default": "h-9 px-4 py-2",
			"sm":      "h-8 rounded-md px-3 text-xs",
			"lg":      "h-10 rounded-md px-8",
			"icon":    "h-9 w-9",
		},
	},
	Defaults: map[string]string{"variant": "default", "size": "default"},
}

type ButtonProps struct {
	Variant  ButtonVariant
	Size     ButtonSize
	Class    string // extra classes, merged with style.CN
	Type     string // defaults to "button"
	Disabled bool
	OnClick  func(*g.Event)
}

// Button renders a shadcn-style button. Children may be strings, nodes, or
// further grove options.
func Button(p ButtonProps, children ...any) *g.Node {
	typ := p.Type
	if typ == "" {
		typ = "button"
	}
	args := []any{
		g.Class(buttonVariants.Class(map[string]string{
			"variant": string(p.Variant),
			"size":    string(p.Size),
		}, p.Class)),
		g.Type(typ),
		g.Disabled(p.Disabled),
		g.Data("slot", "button"),
	}
	if p.OnClick != nil {
		args = append(args, g.OnClick(p.OnClick))
	}
	return g.Button(append(args, children...)...)
}
