// The islands example mounts a real React component inside a grove app.
// index.html registers the island on window.groveIslands (React via esm.sh)
// before starting the wasm; grove passes it props and React re-renders on
// every change, while keeping its own internal state.
package main

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/dom"
	"github.com/gyoumi/grove/island"
)

func App() *g.Node {
	count, setCount := g.UseState(0)
	return g.Div(g.Class("mx-auto max-w-xl space-y-6 p-8"),
		g.H1(g.Class("text-2xl font-bold tracking-tight"), "grove + React islands"),
		g.P(g.Class("text-muted-foreground"),
			"The bordered card below is a React component rendered by react-dom. ",
			"grove owns the page and passes the count down as props; ",
			"the island also keeps its own React state."),
		g.Button(
			g.Class("rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"),
			g.OnClick(func(*g.Event) { setCount(count + 1) }),
			g.Textf("count in grove: %d", count),
		),
		island.C("Greeting", map[string]any{"count": count},
			g.Class("block rounded-lg border bg-card p-4 shadow-sm")),
	)
}

func main() { dom.Mount("#root", g.C0(App)) }
