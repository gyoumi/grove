//go:build js && wasm

package ui

import (
	"syscall/js"

	g "github.com/gyoumi/grove"
)

// dragListen attaches document mousemove/mouseup handlers for the duration of
// a drag and returns a cleanup that detaches them.
func dragListen(onMove func(cx, cy float64), onUp func()) func() {
	doc := js.Global().Get("document")
	moveFn := js.FuncOf(func(_ js.Value, args []js.Value) any {
		e := args[0]
		onMove(e.Get("clientX").Float(), e.Get("clientY").Float())
		return nil
	})
	upFn := js.FuncOf(func(js.Value, []js.Value) any {
		onUp()
		return nil
	})
	doc.Call("addEventListener", "mousemove", moveFn)
	doc.Call("addEventListener", "mouseup", upFn)
	return func() {
		doc.Call("removeEventListener", "mousemove", moveFn)
		doc.Call("removeEventListener", "mouseup", upFn)
		moveFn.Release()
		upFn.Release()
	}
}

// fractionAt converts a pointer position to the first panel's fraction of the
// group along the resize axis.
func fractionAt(ref *g.Ref[any], cx, cy float64, horizontal bool) (float64, bool) {
	el, ok := ref.Current.(js.Value)
	if !ok || el.Type() != js.TypeObject {
		return 0, false
	}
	r := el.Call("getBoundingClientRect")
	if horizontal {
		w := r.Get("width").Float()
		if w == 0 {
			return 0, false
		}
		return (cx - r.Get("left").Float()) / w, true
	}
	h := r.Get("height").Float()
	if h == 0 {
		return 0, false
	}
	return (cy - r.Get("top").Float()) / h, true
}
