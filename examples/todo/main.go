// A todo list demonstrating controlled inputs, keyed lists, and derived
// state — styled with plain Tailwind classes on the shadcn theme.
// Run with: grove serve (from this directory).
package main

import (
	"strconv"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/dom"
)

type todo struct {
	id   int
	text string
	done bool
}

func App() *g.Node {
	todos, setTodos := g.UseState([]todo(nil))
	nextID, setNextID := g.UseState(1)
	input, setInput := g.UseState("")

	add := func() {
		if input == "" {
			return
		}
		setTodos(append(todos[:len(todos):len(todos)], todo{nextID, input, false}))
		setNextID(nextID + 1)
		setInput("")
	}
	toggle := func(id int) func(*g.Event) {
		return func(*g.Event) {
			out := make([]todo, len(todos))
			copy(out, todos)
			for i := range out {
				if out[i].id == id {
					out[i].done = !out[i].done
				}
			}
			setTodos(out)
		}
	}
	remove := func(id int) func(*g.Event) {
		return func(*g.Event) {
			var out []todo
			for _, td := range todos {
				if td.id != id {
					out = append(out, td)
				}
			}
			setTodos(out)
		}
	}

	remaining := 0
	for _, td := range todos {
		if !td.done {
			remaining++
		}
	}

	return g.Div(g.Class("mx-auto flex min-h-svh max-w-md flex-col gap-4 p-8"),
		g.H1(g.Class("text-2xl font-semibold tracking-tight"), "todos"),
		g.Div(g.Class("flex gap-2"),
			g.Input(
				g.Class("flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"),
				g.Placeholder("what needs doing?"),
				g.Value(input),
				g.OnInput(func(e *g.Event) { setInput(e.Value()) }),
				g.OnKeyDown(func(e *g.Event) {
					if e.Key() == "Enter" {
						add()
					}
				}),
			),
			g.Button(
				g.Class("inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground shadow hover:bg-primary/90"),
				g.OnClick(func(*g.Event) { add() }),
				"add",
			),
		),
		g.Ul(g.Class("flex flex-col gap-1"),
			g.Map(todos, func(td todo) *g.Node {
				return g.Li(g.Key(strconv.Itoa(td.id)),
					g.Class("group flex items-center gap-3 rounded-md border px-3 py-2"),
					g.Input(
						g.Type("checkbox"),
						g.Class("size-4 accent-primary"),
						g.Checked(td.done),
						g.OnChange(toggle(td.id)),
					),
					g.Span(
						g.Class("flex-1 text-sm"),
						g.ClassIf(td.done, "text-muted-foreground line-through"),
						td.text,
					),
					g.Button(
						g.Class("text-sm text-muted-foreground opacity-0 transition-opacity hover:text-destructive group-hover:opacity-100"),
						g.OnClick(remove(td.id)),
						"✕",
					),
				)
			}),
		),
		g.If(len(todos) > 0,
			g.P(g.Class("text-sm text-muted-foreground"), g.Textf("%d remaining", remaining)),
		),
	)
}

func main() {
	dom.Mount("#root", g.C0(App))
}
