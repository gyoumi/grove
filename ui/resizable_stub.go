//go:build !(js && wasm)

package ui

import g "github.com/gyoumi/grove"

// Pointer dragging needs a real document; outside the browser (tests) the
// handle is resized with the keyboard instead.
func dragListen(onMove func(cx, cy float64), onUp func()) func() { return func() {} }

func fractionAt(ref *g.Ref[any], cx, cy float64, horizontal bool) (float64, bool) {
	return 0, false
}
