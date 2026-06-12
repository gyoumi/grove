package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type PopoverSide string

const (
	PopoverBottom PopoverSide = "bottom"
	PopoverTop    PopoverSide = "top"
	PopoverLeft   PopoverSide = "left"
	PopoverRight  PopoverSide = "right"
)

type PopoverAlign string

const (
	PopoverAlignStart  PopoverAlign = "start"
	PopoverAlignCenter PopoverAlign = "center"
	PopoverAlignEnd    PopoverAlign = "end"
)

// positionClasses anchors the panel relative to the trigger, per side and
// alignment. Literal class strings keep the Tailwind scanner aware of them.
var positionClasses = map[PopoverSide]map[PopoverAlign]string{
	PopoverBottom: {
		PopoverAlignStart:  "top-full left-0 mt-2",
		PopoverAlignCenter: "top-full left-1/2 -translate-x-1/2 mt-2",
		PopoverAlignEnd:    "top-full right-0 mt-2",
	},
	PopoverTop: {
		PopoverAlignStart:  "bottom-full left-0 mb-2",
		PopoverAlignCenter: "bottom-full left-1/2 -translate-x-1/2 mb-2",
		PopoverAlignEnd:    "bottom-full right-0 mb-2",
	},
	PopoverRight: {
		PopoverAlignStart:  "left-full top-0 ml-2",
		PopoverAlignCenter: "left-full top-1/2 -translate-y-1/2 ml-2",
		PopoverAlignEnd:    "left-full bottom-0 ml-2",
	},
	PopoverLeft: {
		PopoverAlignStart:  "right-full top-0 mr-2",
		PopoverAlignCenter: "right-full top-1/2 -translate-y-1/2 mr-2",
		PopoverAlignEnd:    "right-full bottom-0 mr-2",
	},
}

type PopoverProps struct {
	// Open controls visibility; the popover is fully controlled.
	Open bool
	// OnClose fires on Escape or a click outside the panel.
	OnClose func()
	Side    PopoverSide  // default bottom
	Align   PopoverAlign // default center
	Class   string       // extra classes for the panel
}

// Popover anchors a floating panel to its trigger: the panel is positioned
// with CSS relative to the trigger per Side/Align, a transparent overlay
// underneath closes it on outside clicks, and Escape closes it from the
// keyboard. (Viewport collision flipping is roadmap — pick a Side that has
// room.)
func Popover(p PopoverProps, trigger *g.Node, content ...any) *g.Node {
	side := p.Side
	if side == "" {
		side = PopoverBottom
	}
	align := p.Align
	if align == "" {
		align = PopoverAlignCenter
	}
	close := func() {
		if p.OnClose != nil {
			p.OnClose()
		}
	}

	panel := []any{
		g.Class(style.CN(
			"absolute z-50 min-w-[8rem] rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none",
			positionClasses[side][align], p.Class)),
		g.Role("dialog"),
		g.TabIndex(-1),
		g.Data("slot", "popover-content"),
		g.Data("side", string(side)),
		g.OnKeyDown(func(e *g.Event) {
			if e.Key() == "Escape" {
				close()
			}
		}),
	}
	panel = append(panel, content...)

	return g.Span(g.Class("relative inline-flex"), g.Data("slot", "popover"),
		trigger,
		g.If(p.Open, g.Fragment(
			g.Div(g.Class("fixed inset-0 z-40"), g.Data("slot", "popover-overlay"),
				g.OnClick(func(*g.Event) { close() })),
			g.Div(panel...),
		)),
	)
}
