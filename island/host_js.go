//go:build js && wasm

package island

import (
	"syscall/js"

	g "github.com/gyoumi/grove"
)

// entry returns window.groveIslands[name]; a missing registry or entry is
// a wiring mistake that would otherwise fail silently, so mount warns on
// the console.
func entry(name string, warn bool) js.Value {
	reg := js.Global().Get("groveIslands")
	if reg.Type() != js.TypeObject {
		if warn {
			js.Global().Get("console").Call("warn",
				"grove: window.groveIslands is not set; island "+name+" will not render")
		}
		return js.Undefined()
	}
	e := reg.Get(name)
	if e.Type() != js.TypeObject {
		if warn {
			js.Global().Get("console").Call("warn", "grove: no island registered under "+name)
		}
		return js.Undefined()
	}
	return e
}

func call(ref *g.DOMRef, name, op, propsJSON string, warn bool) {
	el, ok := ref.Current.(js.Value)
	if !ok || el.Type() != js.TypeObject {
		return
	}
	e := entry(name, warn)
	if e.IsUndefined() {
		return
	}
	fn := e.Get(op)
	if fn.Type() != js.TypeFunction {
		return // update and unmount are optional
	}
	if op == "unmount" {
		fn.Invoke(el)
		return
	}
	fn.Invoke(el, js.Global().Get("JSON").Call("parse", propsJSON))
}

func hostMount(ref *g.DOMRef, name, propsJSON string)  { call(ref, name, "mount", propsJSON, true) }
func hostUpdate(ref *g.DOMRef, name, propsJSON string) { call(ref, name, "update", propsJSON, false) }
func hostUnmount(ref *g.DOMRef, name string)           { call(ref, name, "unmount", "", false) }
