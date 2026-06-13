package ui

import (
	"fmt"

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
	// WrapperClass overrides the anchoring wrapper's classes (default
	// "relative inline-flex"); pass "relative block w-full" to let the
	// trigger fill its container, as the field pickers do.
	WrapperClass string
}

// Popover anchors a floating panel to its trigger: the panel is positioned
// with CSS relative to the trigger per Side/Align, a transparent overlay
// underneath closes it on outside clicks, and Escape closes it from the
// keyboard. In the browser the panel measures itself after opening and
// corrects viewport collisions: it flips to the opposite side when that
// side has room (the rendered side is exposed as data-side), shifts along
// the cross axis until it fits, and re-measures when the window resizes
// while open.
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

	return g.Span(g.Class(style.CN("relative inline-flex", p.WrapperClass)), g.Data("slot", "popover"),
		trigger,
		g.If(p.Open, g.Fragment(
			g.Div(g.Class("fixed inset-0 z-40"), g.Data("slot", "popover-overlay"),
				g.OnClick(func(*g.Event) { close() })),
			g.C(popoverPanel, panelArgs{side: side, align: align, class: p.Class, onClose: close, content: content}),
		)),
	)
}

type panelArgs struct {
	side    PopoverSide
	align   PopoverAlign
	class   string
	onClose func()
	content []any
}

// popoverPanel renders the floating panel; it re-measures on every open
// (the panel remounts when Open flips), after each side change, and on
// window resize. Collision handling is flip first, then shift: move to the
// opposite side when the placement overflows and that side has room, then
// slide along the cross axis until the panel fits the viewport.
func popoverPanel(a panelArgs) *g.Node {
	measureRef := g.UseRef[any](nil)
	flipped, setFlipped := g.UseState(PopoverSide(""))
	shift, setShift := g.UseState(0)
	resizes, bump := g.UseReducer(func(n int, _ struct{}) int { return n + 1 }, 0)

	side := a.side
	if flipped != "" {
		side = flipped
	}

	g.UseEffect(func() func() {
		return onViewportResize(func() {
			setFlipped("") // re-evaluate from the requested side
			bump(struct{}{})
		})
	}, []any{})
	g.UseEffect(func() func() {
		if s, ok := measureFlip(measureRef, side); ok {
			setFlipped(s)
			return nil // shift is measured once the new side is committed
		}
		if px, ok := measureShift(measureRef, side, shift); ok {
			setShift(px)
		}
		return nil
	}, []any{side, resizes})

	axis := "X"
	if side == PopoverLeft || side == PopoverRight {
		axis = "Y"
	}
	panel := []any{
		g.Class(style.CN(
			"absolute z-50 min-w-[8rem] rounded-md border bg-popover p-4 text-popover-foreground shadow-md outline-none",
			positionClasses[side][a.align], a.class)),
		g.Role("dialog"),
		g.TabIndex(-1),
		g.Data("slot", "popover-content"),
		g.Data("side", string(side)),
		// The inline transform composes with the translate utilities in
		// positionClasses (CSS translate and transform are distinct
		// properties), so centered panels shift correctly too.
		g.AttrIf(shift != 0, "style", fmt.Sprintf("transform: translate%s(%dpx)", axis, shift)),
		g.BindRef(measureRef),
		g.OnKeyDown(func(e *g.Event) {
			if e.Key() == "Escape" {
				a.onClose()
			}
		}),
	}
	return g.Div(append(panel, a.content...)...)
}
