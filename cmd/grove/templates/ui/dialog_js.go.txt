//go:build js && wasm

package ui

import (
	"syscall/js"

	g "github.com/gyoumi/grove"
)

const focusableSelector = `a[href], button:not([disabled]), input:not([disabled]), select:not([disabled]), textarea:not([disabled]), [tabindex]:not([tabindex="-1"])`

// trapFocus moves focus into the dialog content and returns a cleanup that
// restores focus to the previously focused element.
func trapFocus(ref *g.Ref[any]) func() {
	doc := js.Global().Get("document")
	prev := doc.Get("activeElement")
	if el, ok := ref.Current.(js.Value); ok && el.Type() == js.TypeObject {
		el.Call("focus")
	}
	return func() {
		if prev.Truthy() {
			prev.Call("focus")
		}
	}
}

// cycleFocus keeps Tab navigation inside the dialog content.
func cycleFocus(ref *g.Ref[any], e *g.Event) {
	el, ok := ref.Current.(js.Value)
	if !ok || el.Type() != js.TypeObject {
		return
	}
	list := el.Call("querySelectorAll", focusableSelector)
	n := list.Get("length").Int()
	if n == 0 {
		e.PreventDefault()
		return
	}
	first, last := list.Index(0), list.Index(n-1)
	active := js.Global().Get("document").Get("activeElement")
	if e.ShiftKey() {
		if active.Equal(first) || active.Equal(el) {
			e.PreventDefault()
			last.Call("focus")
		}
	} else if active.Equal(last) {
		e.PreventDefault()
		first.Call("focus")
	}
}
