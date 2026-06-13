// The island lifecycle test drives the off-browser stub host (SetHost),
// which only exists outside js/wasm; go test runs on the host, so the
// build tag keeps `GOOS=js GOARCH=wasm go vet` from compiling a file that
// references a symbol absent on that platform.
//go:build !(js && wasm)

package island_test

import (
	"fmt"
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/island"
	"github.com/gyoumi/grove/testdom"
)

type recHost struct {
	calls []string
	els   []any
}

func (h *recHost) Mount(el any, name, props string) {
	h.calls = append(h.calls, fmt.Sprintf("mount %s %s", name, props))
	h.els = append(h.els, el)
}

func (h *recHost) Update(el any, name, props string) {
	h.calls = append(h.calls, fmt.Sprintf("update %s %s", name, props))
	h.els = append(h.els, el)
}

func (h *recHost) Unmount(el any, name string) {
	h.calls = append(h.calls, "unmount "+name)
	h.els = append(h.els, el)
}

var (
	setCount func(int)
	setShow  func(bool)
	setTick  func(int)
)

func app() *g.Node {
	count, sc := g.UseState(1)
	show, ss := g.UseState(true)
	tick, st := g.UseState(0)
	setCount, setShow, setTick = sc, ss, st
	return g.Div(g.Data("tick", fmt.Sprint(tick)),
		g.If(show, island.C("Chart", map[string]any{"count": count}, g.Class("h-40"))),
	)
}

func TestIslandLifecycle(t *testing.T) {
	h := &recHost{}
	island.SetHost(h)
	defer island.SetHost(nil)

	r := testdom.Mount(g.C0(app))
	el := r.FindByAttr("data-island", "Chart")
	if el == nil {
		t.Fatalf("island container missing: %s", r.HTML())
	}
	if el.Attrs["class"] != "h-40" || el.Attrs["data-slot"] != "island" {
		t.Fatalf("container options not applied: %s", el.HTML())
	}
	if len(h.calls) != 1 || h.calls[0] != `mount Chart {"count":1}` {
		t.Fatalf("mount call wrong: %v", h.calls)
	}
	if h.els[0] == nil {
		t.Fatal("mount should receive the container handle")
	}

	// changed props reach the host exactly once, as an update
	setCount(2)
	r.Settle()
	if len(h.calls) != 2 || h.calls[1] != `update Chart {"count":2}` {
		t.Fatalf("update call wrong: %v", h.calls)
	}

	// re-renders that leave the props alone don't reach the host
	setTick(1)
	r.Settle()
	if len(h.calls) != 2 {
		t.Fatalf("props-preserving re-render must not call the host: %v", h.calls)
	}

	// removal unmounts the island
	setShow(false)
	r.Settle()
	if len(h.calls) != 3 || h.calls[2] != "unmount Chart" {
		t.Fatalf("unmount call wrong: %v", h.calls)
	}
	if r.FindByAttr("data-island", "Chart") != nil {
		t.Fatalf("container should be gone: %s", r.HTML())
	}

	// showing again is a fresh mount
	setShow(true)
	r.Settle()
	if len(h.calls) != 4 || h.calls[3] != `mount Chart {"count":2}` {
		t.Fatalf("remount call wrong: %v", h.calls)
	}
}

func TestIslandBadPropsPanic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("non-marshalable props should panic")
		}
	}()
	island.C("Chart", map[string]any{"fn": func() {}})
}
