//go:build js && wasm

package ui

import (
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
