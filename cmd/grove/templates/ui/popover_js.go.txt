//go:build js && wasm

package ui

import (
	"math"
	"syscall/js"

	g "github.com/gyoumi/grove"
)

// measureFlip reports the side the panel should move to when its current
// placement overflows the viewport and the opposite side of the anchor
// (the popover wrapper) has room; ok is false when the placement is fine.
func measureFlip(ref *g.DOMRef, side PopoverSide) (PopoverSide, bool) {
	el, ok := ref.Current.(js.Value)
	if !ok || el.Type() != js.TypeObject {
		return "", false
	}
	anchor := el.Get("parentElement")
	if anchor.Type() != js.TypeObject {
		return "", false
	}
	rect := el.Call("getBoundingClientRect")
	arect := anchor.Call("getBoundingClientRect")
	vw := js.Global().Get("innerWidth").Float()
	vh := js.Global().Get("innerHeight").Float()
	h := rect.Get("height").Float()
	w := rect.Get("width").Float()
	const gap = 8 // matches the m*-2 offset in positionClasses

	switch side {
	case PopoverBottom:
		if rect.Get("bottom").Float() > vh && arect.Get("top").Float()-h-gap >= 0 {
			return PopoverTop, true
		}
	case PopoverTop:
		if rect.Get("top").Float() < 0 && arect.Get("bottom").Float()+h+gap <= vh {
			return PopoverBottom, true
		}
	case PopoverRight:
		if rect.Get("right").Float() > vw && arect.Get("left").Float()-w-gap >= 0 {
			return PopoverLeft, true
		}
	case PopoverLeft:
		if rect.Get("left").Float() < 0 && arect.Get("right").Float()+w+gap <= vw {
			return PopoverRight, true
		}
	}
	return "", false
}

// measureShift returns the cross-axis offset (px, applied as an inline
// transform) that keeps the panel inside the viewport. The overflow is
// computed from the panel's un-shifted position, so the offset relaxes
// back toward zero when the viewport grows; ok is false when the current
// offset is already right.
func measureShift(ref *g.DOMRef, side PopoverSide, current int) (int, bool) {
	el, ok := ref.Current.(js.Value)
	if !ok || el.Type() != js.TypeObject {
		return 0, false
	}
	rect := el.Call("getBoundingClientRect")
	const pad = 4 // breathing room against the viewport edge
	var lo, hi, limit float64
	if side == PopoverLeft || side == PopoverRight {
		lo = rect.Get("top").Float() - float64(current)
		hi = rect.Get("bottom").Float() - float64(current)
		limit = js.Global().Get("innerHeight").Float()
	} else {
		lo = rect.Get("left").Float() - float64(current)
		hi = rect.Get("right").Float() - float64(current)
		limit = js.Global().Get("innerWidth").Float()
	}
	var delta float64
	if hi > limit-pad {
		delta = (limit - pad) - hi
	}
	if lo+delta < pad {
		delta = pad - lo // when the panel fits neither way, keep the start edge visible
	}
	if n := int(math.Round(delta)); n != current {
		return n, true
	}
	return current, false
}

// onViewportResize runs fn on window resize and returns a cleanup that
// detaches the listener.
func onViewportResize(fn func()) func() {
	cb := js.FuncOf(func(js.Value, []js.Value) any { fn(); return nil })
	js.Global().Call("addEventListener", "resize", cb)
	return func() {
		js.Global().Call("removeEventListener", "resize", cb)
		cb.Release()
	}
}
