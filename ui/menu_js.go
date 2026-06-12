//go:build js && wasm

package ui

import (
	"syscall/js"

	g "github.com/gyoumi/grove"
)

// menuFocusMove moves keyboard focus through a menu's enabled items: dir
// +1 for ArrowDown, -1 for ArrowUp, wrapping at the ends.
func menuFocusMove(ref *g.DOMRef, e *g.Event, dir int) {
	menu, ok := ref.Current.(js.Value)
	if !ok || menu.Type() != js.TypeObject {
		return
	}
	items := menu.Call("querySelectorAll", `[role="menuitem"]:not([disabled])`)
	n := items.Get("length").Int()
	if n == 0 {
		return
	}
	e.PreventDefault()
	active := js.Global().Get("document").Get("activeElement")
	cur := -1
	for i := 0; i < n; i++ {
		if items.Index(i).Equal(active) {
			cur = i
			break
		}
	}
	next := cur + dir
	switch {
	case cur == -1:
		next = 0
		if dir < 0 {
			next = n - 1
		}
	case next < 0:
		next = n - 1
	case next >= n:
		next = 0
	}
	items.Index(next).Call("focus")
}
