// Package island mounts JS-rendered components — React components in
// particular — as leaf nodes inside a grove tree. The page registers
// islands on window.groveIslands before the wasm app starts; grove renders
// a container element and hands it (with JSON-decoded props) to the
// registry on mount, on every props change, and on unmount:
//
//	window.groveIslands = {
//		Greeting: {
//			mount(el, props)  { el._root = ReactDOM.createRoot(el); el._root.render(...) },
//			update(el, props) { el._root.render(...) },           // optional
//			unmount(el)       { el._root.unmount() },             // optional
//		},
//	}
//
// The JS side owns everything inside the container: grove never renders
// children into it, so an island is always a leaf. State that must survive
// updates belongs on the element (as above) or in a closure — grove holds
// no JS handles.
package island

import (
	"encoding/json"
	"fmt"

	g "github.com/gyoumi/grove"
)

type leafProps struct {
	name string
	json string
	opts []g.Option
}

// C places the island registered under name as a leaf node. props is
// marshaled to JSON and delivered to the island's mount on first commit,
// then to update whenever the encoding changes; re-renders that leave the
// props alone reach JS not at all. opts decorate the container element
// (class, id, …). Props that cannot be marshaled are a programming error
// and panic, like a rules-of-hooks violation.
func C(name string, props any, opts ...g.Option) *g.Node {
	b, err := json.Marshal(props)
	if err != nil {
		panic(fmt.Sprintf("island %q: props are not JSON-marshalable: %v", name, err))
	}
	return g.C(leaf, leafProps{name: name, json: string(b), opts: opts})
}

func leaf(p leafProps) *g.Node {
	ref := g.UseRef[any](nil)
	pushed := g.UseRef("") // the props JSON the host has already seen

	g.UseEffect(func() func() {
		hostMount(ref, p.name, p.json)
		pushed.Current = p.json
		// The cleanup runs before the container leaves the document, so
		// the JS side can still reach its DOM.
		return func() { hostUnmount(ref, p.name) }
	}, []any{})
	g.UseEffect(func() func() {
		if pushed.Current != p.json {
			hostUpdate(ref, p.name, p.json)
			pushed.Current = p.json
		}
		return nil
	}, []any{p.json})

	args := []any{g.Data("slot", "island"), g.Data("island", p.name), g.BindRef(ref)}
	for _, o := range p.opts {
		args = append(args, o)
	}
	return g.Div(args...)
}
