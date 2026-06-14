package ui

import (
	"fmt"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type ResizableProps struct {
	Vertical        bool    // stack the panels (split top/bottom) instead of side-by-side
	DefaultFraction float64 // first panel's share, 0–1 (default 0.5)
	Class           string
}

// Resizable splits an area into two panels with a draggable divider. Drag the
// handle (or focus it and use the arrow keys) to resize; the split is clamped
// so neither panel collapses.
func Resizable(p ResizableProps, first, second *g.Node) *g.Node {
	return g.C(resizableView, resizableArgs{p: p, first: first, second: second})
}

type resizableArgs struct {
	p             ResizableProps
	first, second *g.Node
}

func resizableView(a resizableArgs) *g.Node {
	init := a.p.DefaultFraction
	if init <= 0 || init >= 1 {
		init = 0.5
	}
	fraction, setFraction := g.UseState(init)
	dragging, setDragging := g.UseState(false)
	groupRef := g.UseRef[any](nil)

	clamp := func(f float64) float64 { return max(0.08, min(0.92, f)) }

	// While dragging, listen on the document so the pointer is tracked even
	// outside the group; the listeners are torn down when the drag ends.
	g.UseEffect(func() func() {
		if !dragging {
			return nil
		}
		return dragListen(
			func(cx, cy float64) {
				if f, ok := fractionAt(groupRef, cx, cy, !a.p.Vertical); ok {
					setFraction(clamp(f))
				}
			},
			func() { setDragging(false) },
		)
	}, []any{dragging})

	dir, handleAxis := "flex-row", "w-1.5 cursor-col-resize"
	if a.p.Vertical {
		dir, handleAxis = "flex-col", "h-1.5 cursor-row-resize"
	}

	panel := func(grow float64, child *g.Node) *g.Node {
		return g.Div(
			g.Class("min-h-0 min-w-0 overflow-auto"),
			g.Data("slot", "resizable-panel"),
			g.Attr("style", fmt.Sprintf("flex: %.4f 1 0%%", grow)),
			child,
		)
	}

	handle := g.Div(
		g.Class(style.CN("flex shrink-0 items-center justify-center bg-border transition-colors hover:bg-primary/40 focus-visible:bg-primary/60 focus-visible:outline-none data-[dragging=1]:bg-primary/60", handleAxis)),
		g.Data("slot", "resizable-handle"),
		g.Role("separator"),
		g.TabIndex(0),
		g.AttrIf(dragging, "data-dragging", "1"),
		g.OnMouseDown(func(e *g.Event) { e.PreventDefault(); setDragging(true) }),
		g.OnKeyDown(func(e *g.Event) {
			const step = 0.04
			switch e.Key() {
			case "ArrowLeft", "ArrowUp":
				e.PreventDefault()
				setFraction(clamp(fraction - step))
			case "ArrowRight", "ArrowDown":
				e.PreventDefault()
				setFraction(clamp(fraction + step))
			}
		}),
	)

	return g.Div(
		g.Class(style.CN("flex h-full w-full overflow-hidden rounded-lg border", dir, a.p.Class)),
		g.Data("slot", "resizable-group"),
		g.BindRef(groupRef),
		panel(fraction, a.first),
		handle,
		panel(1-fraction, a.second),
	)
}
