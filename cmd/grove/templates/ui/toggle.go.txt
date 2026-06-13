package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type ToggleVariant string

const (
	ToggleDefault ToggleVariant = "default"
	ToggleOutline ToggleVariant = "outline"
)

type ToggleSize string

const (
	ToggleSizeDefault ToggleSize = "default"
	ToggleSizeSm      ToggleSize = "sm"
	ToggleSizeLg      ToggleSize = "lg"
)

var toggleVariants = style.Variants{
	Base: "inline-flex items-center justify-center gap-2 rounded-md text-sm font-medium transition-colors hover:bg-muted hover:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50 data-[state=on]:bg-accent data-[state=on]:text-accent-foreground",
	Groups: map[string]map[string]string{
		"variant": {
			"default": "bg-transparent",
			"outline": "border border-input bg-transparent shadow-sm",
		},
		"size": {
			"default": "h-9 min-w-9 px-2",
			"sm":      "h-8 min-w-8 px-1.5",
			"lg":      "h-10 min-w-10 px-2.5",
		},
	},
	Defaults: map[string]string{"variant": "default", "size": "default"},
}

type ToggleProps struct {
	Pressed         bool // controlled on/off
	Variant         ToggleVariant
	Size            ToggleSize
	Disabled        bool
	OnPressedChange func(bool)
	Class           string
}

// Toggle is a controlled two-state button (a single member of a ToggleGroup
// can stand alone). data-state is "on" when Pressed.
func Toggle(p ToggleProps, children ...any) *g.Node {
	args := []any{
		g.Class(toggleVariants.Class(map[string]string{
			"variant": string(p.Variant),
			"size":    string(p.Size),
		}, p.Class)),
		g.Type("button"),
		g.Data("slot", "toggle"),
		g.Data("state", onOff(p.Pressed)),
		g.AttrIf(p.Pressed, "aria-pressed", "true"),
		g.Disabled(p.Disabled),
		g.OnClick(func(*g.Event) {
			if p.OnPressedChange != nil {
				p.OnPressedChange(!p.Pressed)
			}
		}),
	}
	return g.Button(append(args, children...)...)
}

func onOff(on bool) string {
	if on {
		return "on"
	}
	return "off"
}
