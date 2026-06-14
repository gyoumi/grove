//go:build js && wasm

package dom

import (
	"syscall/js"

	"github.com/gyoumi/grove/renderer"
)

// groveIDProp tags grove-created elements so the delegated event listeners
// can find the nearest grove element above an event target.
const groveIDProp = "__groveID"

// eventAlias maps grove event names to the bubbling DOM event actually
// listened for (focus and blur themselves don't bubble).
var eventAlias = map[string]string{
	"focus": "focusin",
	"blur":  "focusout",
}

type jsRenderer struct {
	doc       js.Value
	container js.Value
	dispatch  renderer.Dispatch

	listened map[string]bool
	funcs    []js.Func // kept alive for the app's lifetime

	queue   []func()
	drainFn js.Func
}

func newRenderer(doc, container js.Value) *jsRenderer {
	r := &jsRenderer{doc: doc, container: container, listened: map[string]bool{}}
	r.drainFn = js.FuncOf(func(js.Value, []js.Value) any {
		q := r.queue
		r.queue = nil
		for _, f := range q {
			f()
		}
		return nil
	})
	return r
}

func (r *jsRenderer) SetDispatch(d renderer.Dispatch) { r.dispatch = d }

// PortalRoot is the mount container; portal children sit beside the app's
// root content, outside any transformed subtree, so their fixed positioning
// is relative to the viewport.
func (r *jsRenderer) PortalRoot() renderer.Node { return r.container }

// svgNS is the SVG namespace; svgTags are the element names that must be
// created in it (createElement would make inert HTMLUnknownElements). Once an
// element is in the SVG namespace its descendants inherit it, so tag-name
// detection is enough to render inline icons and simple vector graphics.
const svgNS = "http://www.w3.org/2000/svg"

var svgTags = map[string]bool{
	"svg": true, "path": true, "g": true, "circle": true, "rect": true,
	"line": true, "polyline": true, "polygon": true, "ellipse": true,
	"defs": true, "use": true, "text": true, "tspan": true, "title": true,
	"linearGradient": true, "radialGradient": true, "stop": true, "clipPath": true,
}

func (r *jsRenderer) CreateElement(tag string, id int) renderer.Node {
	var el js.Value
	if svgTags[tag] {
		el = r.doc.Call("createElementNS", svgNS, tag)
	} else {
		el = r.doc.Call("createElement", tag)
	}
	el.Set(groveIDProp, id)
	return el
}

func (r *jsRenderer) CreateText(text string) renderer.Node {
	return r.doc.Call("createTextNode", text)
}

func (r *jsRenderer) SetText(n renderer.Node, text string) {
	n.(js.Value).Set("nodeValue", text)
}

func (r *jsRenderer) SetAttr(n renderer.Node, name, value string) {
	n.(js.Value).Call("setAttribute", name, value)
}

func (r *jsRenderer) RemoveAttr(n renderer.Node, name string) {
	n.(js.Value).Call("removeAttribute", name)
}

func (r *jsRenderer) SetProp(n renderer.Node, name string, value any) {
	el := n.(js.Value)
	if value == nil {
		// Clearing: empty string for value (null would render as "null"),
		// null for everything else.
		if name == "value" {
			el.Set(name, "")
		} else {
			el.Set(name, js.Null())
		}
		return
	}
	el.Set(name, js.ValueOf(value))
}

func (r *jsRenderer) InsertBefore(parent, child, before renderer.Node) {
	b := js.Null()
	if before != nil {
		b = before.(js.Value)
	}
	parent.(js.Value).Call("insertBefore", child.(js.Value), b)
}

func (r *jsRenderer) Remove(parent, child renderer.Node) {
	p := parent.(js.Value)
	c := child.(js.Value)
	// The child may already be gone if outside code mutated the DOM.
	if c.Get("parentNode").Equal(p) {
		p.Call("removeChild", c)
	}
}

func (r *jsRenderer) Listen(event string) {
	if r.listened[event] {
		return
	}
	r.listened[event] = true
	domEvent := event
	if alias, ok := eventAlias[event]; ok {
		domEvent = alias
	}
	cb := js.FuncOf(func(_ js.Value, args []js.Value) any {
		ev := args[0]
		t := ev.Get("target")
		for !t.IsNull() && !t.IsUndefined() {
			if id := t.Get(groveIDProp); !id.IsUndefined() && !id.IsNull() {
				r.dispatch(id.Int(), event, ev)
				return nil
			}
			if t.Equal(r.container) {
				return nil
			}
			t = t.Get("parentNode")
		}
		return nil
	})
	r.funcs = append(r.funcs, cb)
	r.container.Call("addEventListener", domEvent, cb)
}

// Flush is a no-op: this renderer applies each DOM op immediately. The
// batched renderer (MountBatched) is the one that defers and flushes here.
func (r *jsRenderer) Flush() {}

func (r *jsRenderer) Schedule(f func()) {
	r.queue = append(r.queue, f)
	js.Global().Call("queueMicrotask", r.drainFn)
}

// --- renderer.EventOps ---

func walk(raw any, path []string) js.Value {
	v, ok := raw.(js.Value)
	if !ok {
		return js.Undefined()
	}
	for _, p := range path {
		if v.Type() != js.TypeObject {
			return js.Undefined()
		}
		v = v.Get(p)
	}
	return v
}

func (r *jsRenderer) Str(raw any, path ...string) string {
	v := walk(raw, path)
	switch v.Type() {
	case js.TypeString:
		return v.String()
	case js.TypeNumber, js.TypeBoolean:
		return js.Global().Get("String").Invoke(v).String()
	}
	return ""
}

func (r *jsRenderer) Bool(raw any, path ...string) bool {
	v := walk(raw, path)
	return v.Type() == js.TypeBoolean && v.Bool()
}

func (r *jsRenderer) Num(raw any, path ...string) float64 {
	v := walk(raw, path)
	if v.Type() == js.TypeNumber {
		return v.Float()
	}
	return 0
}

func (r *jsRenderer) PreventDefault(raw any) {
	if v, ok := raw.(js.Value); ok {
		v.Call("preventDefault")
	}
}

func (r *jsRenderer) StopPropagation(raw any) {
	if v, ok := raw.(js.Value); ok {
		v.Call("stopPropagation")
	}
}
