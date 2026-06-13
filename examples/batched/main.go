// The batched example renders through grove's experimental batched patch
// protocol: dom.MountBatched ships each commit's DOM mutations to a JS
// applier in a single call instead of one syscall/js call per mutation.
// It exercises events, a keyed list (insert/remove/move), and a value
// property under batching — open it in a browser to validate the protocol.
package main

import (
	"strconv"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/dom"
)

func App() *g.Node {
	items, setItems := g.UseState([]string{"alpha", "beta", "gamma"})
	count, setCount := g.UseState(0)

	add := func(*g.Event) {
		setItems(append([]string{"item " + strconv.Itoa(len(items))}, items...))
	}
	rotate := func(*g.Event) {
		if len(items) > 1 {
			setItems(append(items[1:], items[0]))
		}
	}
	drop := func(s string) func(*g.Event) {
		return func(*g.Event) {
			out := make([]string, 0, len(items))
			for _, it := range items {
				if it != s {
					out = append(out, it)
				}
			}
			setItems(out)
		}
	}

	return g.Div(g.Class("mx-auto max-w-md space-y-6 p-8"),
		g.H1(g.Class("text-2xl font-bold tracking-tight"), "batched rendering"),
		g.Div(g.Class("flex items-center gap-3"),
			g.Button(
				g.Class("rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground"),
				g.OnClick(func(*g.Event) { setCount(count + 1) }),
				g.Textf("clicked %d times", count),
			),
			g.Input(g.Class("h-9 w-20 rounded-md border px-2 text-sm"), g.Value(strconv.Itoa(count))),
		),
		g.Div(g.Class("flex gap-2"),
			g.Button(g.Class("rounded-md border px-3 py-1.5 text-sm hover:bg-accent"), g.OnClick(add), "add"),
			g.Button(g.Class("rounded-md border px-3 py-1.5 text-sm hover:bg-accent"), g.OnClick(rotate), "rotate"),
		),
		g.Ul(g.Class("space-y-1"),
			g.Map(items, func(s string) *g.Node {
				return g.Li(g.Key(s), g.Class("flex items-center justify-between rounded-md border px-3 py-2 text-sm"),
					g.Span(s),
					g.Button(g.Class("text-muted-foreground hover:text-destructive"), g.OnClick(drop(s)), "remove"),
				)
			}),
		),
	)
}

func main() { dom.MountBatched("#root", g.C0(App)) }
