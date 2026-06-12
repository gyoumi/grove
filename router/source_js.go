//go:build js && wasm

package router

import (
	"strings"
	"syscall/js"
)

// hashSource keeps the path in location.hash so links are shareable and
// back/forward work.
type hashSource struct {
	subs     subscribers
	listened bool
	cb       js.Func
}

var src source = &hashSource{}

func (h *hashSource) path() string {
	hash := js.Global().Get("location").Get("hash").String()
	return normalize(strings.TrimPrefix(hash, "#"))
}

func (h *hashSource) navigate(p string) {
	// Triggers a hashchange event, which notifies subscribers.
	js.Global().Get("location").Set("hash", p)
}

func (h *hashSource) subscribe(fn func()) func() {
	if !h.listened {
		h.listened = true
		h.cb = js.FuncOf(func(js.Value, []js.Value) any {
			h.subs.notify()
			return nil
		})
		js.Global().Call("addEventListener", "hashchange", h.cb)
	}
	return h.subs.add(fn)
}
