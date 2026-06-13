//go:build js && wasm

package dom

import (
	"syscall/js"

	grove "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/patch"
	"github.com/gyoumi/grove/renderer"
)

// MountBatched is an experimental alternative to Mount that renders through
// grove's batched patch protocol: DOM mutations for a commit are recorded
// into an op buffer and shipped to a JS-side applier in a single call,
// instead of one syscall/js call per mutation. The applier mirrors the
// host-verified reference in package patch.
//
// It is opt-in and not yet the default: ref-based components (anything using
// g.BindRef, e.g. the Popover/DatePicker/Dialog/island helpers) need the
// node handle to be a live js.Value, but under batching a handle is an
// integer id, so those need the default Mount until the resolver lands.
// Plain element/text/attribute/property/event/list rendering works.
func MountBatched(selector string, root *grove.Node) {
	doc := js.Global().Get("document")
	container := doc.Call("querySelector", selector)
	if container.IsNull() {
		panic("grove/dom: no element matches selector " + selector)
	}
	container.Set("innerHTML", "")
	r := newBatchRenderer(container)
	grove.Mount(r, patch.ContainerID, root)
	select {}
}

// batchRenderer embeds the host-verified patch.Recorder for the recording
// half — so its mutation encoding is exactly the code the protocol test
// covers — and overrides only the browser-specific parts: shipping the
// buffer across the boundary (Flush), microtask scheduling, and reading
// real event objects.
type batchRenderer struct {
	*patch.Recorder

	dispatch   renderer.Dispatch
	apply      js.Value // the per-flush JS closure: apply(payload string)
	dispatchFn js.Func  // Go callback the JS listeners invoke

	queue   []func()
	drainFn js.Func
}

func newBatchRenderer(container js.Value) *batchRenderer {
	r := &batchRenderer{Recorder: patch.NewRecorder()}
	r.drainFn = js.FuncOf(func(js.Value, []js.Value) any {
		q := r.queue
		r.queue = nil
		for _, f := range q {
			f()
		}
		return nil
	})
	r.dispatchFn = js.FuncOf(func(_ js.Value, args []js.Value) any {
		if r.dispatch != nil {
			r.dispatch(args[0].Int(), args[1].String(), args[2])
		}
		return nil
	})
	// Build the applier once, bound to this container and dispatch callback.
	ctor := js.Global().Get("Function").New("container", "dispatch", applierSource)
	r.apply = ctor.Invoke(container, r.dispatchFn)
	return r
}

func (r *batchRenderer) SetDispatch(d renderer.Dispatch) { r.dispatch = d }

// Flush ships the whole batch across the boundary in one call.
func (r *batchRenderer) Flush() {
	if r.Buf.Len() == 0 {
		return
	}
	r.apply.Invoke(r.Buf.Encode())
	r.Buf.Reset()
}

func (r *batchRenderer) Schedule(f func()) {
	r.queue = append(r.queue, f)
	js.Global().Call("queueMicrotask", r.drainFn)
}

// --- EventOps over the raw js.Value event (same semantics as the direct
// renderer; kept separate so that renderer stays untouched) ---

func (r *batchRenderer) PreventDefault(raw any) {
	if v, ok := raw.(js.Value); ok {
		v.Call("preventDefault")
	}
}

func (r *batchRenderer) StopPropagation(raw any) {
	if v, ok := raw.(js.Value); ok {
		v.Call("stopPropagation")
	}
}

func (r *batchRenderer) walk(raw any, path []string) js.Value {
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

func (r *batchRenderer) Str(raw any, path ...string) string {
	v := r.walk(raw, path)
	switch v.Type() {
	case js.TypeString:
		return v.String()
	case js.TypeNumber, js.TypeBoolean:
		return js.Global().Get("String").Invoke(v).String()
	}
	return ""
}

func (r *batchRenderer) Bool(raw any, path ...string) bool {
	v := r.walk(raw, path)
	return v.Type() == js.TypeBoolean && v.Bool()
}

func (r *batchRenderer) Num(raw any, path ...string) float64 {
	v := r.walk(raw, path)
	if v.Type() == js.TypeNumber {
		return v.Float()
	}
	return 0
}

// applierSource is the body of a JS function (container, dispatch) that
// returns a per-flush applier. It maintains an id→node map (the container is
// id 1) and delegated listeners, and mirrors patch.Apply op for op. Kept as
// source so the dom package stays self-contained; constructed once per mount.
const applierSource = `
var map = {1: container};
var listened = {};
var alias = {focus: "focusin", blur: "focusout"};

function unesc(s) {
  if (s.indexOf("\\") < 0) return s;
  var out = "", i = 0;
  while (i < s.length) {
    var c = s.charAt(i);
    if (c === "\\" && i + 1 < s.length) {
      var n = s.charAt(i + 1);
      out += n === "t" ? "\t" : n === "n" ? "\n" : n;
      i += 2;
    } else { out += c; i++; }
  }
  return out;
}

function setProp(node, name, kind, v) {
  if (kind === "z") { node[name] = name === "value" ? "" : null; return; }
  if (kind === "b") { node[name] = v === "true"; return; }
  if (kind === "n") { node[name] = +v; return; }
  node[name] = v;
}

function listen(ev) {
  if (listened[ev]) return;
  listened[ev] = true;
  container.addEventListener(alias[ev] || ev, function (e) {
    var t = e.target;
    while (t) {
      if (t.__groveID !== undefined && t.__groveID !== null) { dispatch(t.__groveID, ev, e); return; }
      if (t === container) return;
      t = t.parentNode;
    }
  });
}

return function (payload) {
  var lines = payload.split("\n");
  for (var i = 0; i < lines.length; i++) {
    var f = lines[i].split("\t");
    for (var j = 0; j < f.length; j++) f[j] = unesc(f[j]);
    var id = +f[1];
    switch (f[0]) {
      case "e": { var el = document.createElement(f[2]); el.__groveID = +f[3]; map[id] = el; break; }
      case "t": map[id] = document.createTextNode(f[2]); break;
      case "s": map[id].nodeValue = f[2]; break;
      case "a": map[id].setAttribute(f[2], f[3]); break;
      case "r": map[id].removeAttribute(f[2]); break;
      case "p": setProp(map[id], f[2], f[3], f[4]); break;
      case "i": map[id].insertBefore(map[+f[2]], +f[3] ? map[+f[3]] : null); break;
      case "x": map[id].removeChild(map[+f[2]]); break;
      case "l": listen(f[1]); break;
    }
  }
};
`
