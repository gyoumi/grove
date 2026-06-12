// The smallest grove app: one component, one piece of state, no styling.
// Run with: grove serve (from this directory).
package main

import (
	"fmt"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/dom"
)

func Counter() *g.Node {
	count, setCount := g.UseState(0)

	g.UseEffect(func() func() {
		dom.SetTitle(fmt.Sprintf("count: %d", count))
		return nil
	}, []any{count})

	return g.Div(
		g.H1("grove counter"),
		g.Button(
			g.OnClick(func(*g.Event) { setCount(count + 1) }),
			g.Textf("clicked %d times", count),
		),
	)
}

func main() {
	dom.Mount("#root", g.C0(Counter))
}
