//go:build js && wasm

package dom

import (
	"syscall/js"

	grove "github.com/gyoumi/grove"
)

// Mount renders root into the element matching the CSS selector and blocks
// forever so event handlers stay alive. Call it last in main.
func Mount(selector string, root *grove.Node) {
	doc := js.Global().Get("document")
	container := doc.Call("querySelector", selector)
	if container.IsNull() {
		panic("grove/dom: no element matches selector " + selector)
	}
	container.Set("innerHTML", "")
	r := newRenderer(doc, container)
	grove.Mount(r, container, root)
	select {}
}

// Document returns the browser document as a js.Value, as an escape hatch
// for APIs grove doesn't wrap.
func Document() js.Value { return js.Global().Get("document") }

// Window returns the browser window object.
func Window() js.Value { return js.Global() }

// SetTitle sets the page title.
func SetTitle(title string) { Document().Set("title", title) }

// Focus focuses a DOM handle obtained via g.BindRef.
func Focus(handle any) {
	if v, ok := handle.(js.Value); ok && v.Type() == js.TypeObject {
		v.Call("focus")
	}
}

// SetRootClass toggles a class on <html> — handy for shadcn-style dark mode:
// dom.SetRootClass("dark", enabled).
func SetRootClass(name string, on bool) {
	cl := Document().Get("documentElement").Get("classList")
	if on {
		cl.Call("add", name)
	} else {
		cl.Call("remove", name)
	}
}
