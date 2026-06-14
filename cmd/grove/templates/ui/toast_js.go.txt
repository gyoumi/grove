//go:build js && wasm

package ui

import "syscall/js"

// scheduleToastDismiss removes toast id after ms via the browser timer.
func scheduleToastDismiss(id, ms int) {
	var cb js.Func
	cb = js.FuncOf(func(js.Value, []js.Value) any {
		dismissToast(id)
		cb.Release()
		return nil
	})
	js.Global().Call("setTimeout", cb, ms)
}
