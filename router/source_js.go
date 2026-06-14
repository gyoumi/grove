//go:build js && wasm

package router

import (
	"syscall/js"
)

// historySource keeps the path in the real URL via the History API, so links
// are clean (/event/42, no #) and back/forward work. Hosting the built app
// needs a fallback that serves index.html for unknown paths; grove serve
// does this, and grove build documents it.
type historySource struct {
	subs     subscribers
	listened bool
	cb       js.Func
}

var src source = &historySource{}

func (h *historySource) path() string {
	return normalize(js.Global().Get("location").Get("pathname").String())
}

func (h *historySource) navigate(p string) {
	// pushState changes the URL without a reload and does not fire popstate,
	// so notify subscribers directly.
	js.Global().Get("history").Call("pushState", js.Null(), "", p)
	h.subs.notify()
}

func (h *historySource) subscribe(fn func()) func() {
	if !h.listened {
		h.listened = true
		// popstate fires on back/forward navigation.
		h.cb = js.FuncOf(func(js.Value, []js.Value) any {
			h.subs.notify()
			return nil
		})
		js.Global().Call("addEventListener", "popstate", h.cb)
	}
	return h.subs.add(fn)
}
