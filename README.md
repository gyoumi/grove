# grove

A React-style UI framework for Go, compiled to WebAssembly. Function
components render a virtual DOM, hooks hold state, and a reconciler diffs
each render and applies minimal DOM updates — the React mental model, in Go.

```go
package main

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/dom"
)

func Counter() *g.Node {
	count, setCount := g.UseState(0)
	return g.Button(
		g.Class("rounded-md bg-primary px-4 py-2 text-primary-foreground"),
		g.OnClick(func(*g.Event) { setCount(count + 1) }),
		g.Textf("count is %d", count),
	)
}

func main() {
	dom.Mount("#root", g.C0(Counter))
}
```

grove is **client-side only by design**: it renders in the browser, full
stop. There is no server rendering or hydration, and none is planned —
that keeps the engine small and the model simple.

## Quickstart

```sh
go install github.com/gyoumi/grove/cmd/grove@latest

grove init myapp
cd myapp
grove serve        # http://localhost:8080, rebuilds and reloads on save
```

`grove init` scaffolds a working app with Tailwind CSS v4 and a themeable
design system preconfigured. No Node toolchain is involved: the CLI uses the
Tailwind standalone binary (downloaded once, cached) and plain `go build`
with `GOOS=js GOARCH=wasm`.

## Components

A component is a function returning `*g.Node`. Wrap it with `g.C0` (no
props) or `g.C` (typed props) to place it in the tree:

```go
type GreetingProps struct{ Name string }

func Greeting(p GreetingProps) *g.Node {
	return g.P("hello, ", g.Strong(p.Name))
}

func App() *g.Node {
	return g.Div(
		g.C(Greeting, GreetingProps{Name: "Ada"}),
	)
}
```

Element constructors (`Div`, `Span`, `Button`, `Input`, …) accept any mix
of:

- **options** — attributes, properties, and handlers: `g.Class(…)`,
  `g.ID(…)`, `g.Value(…)`, `g.OnClick(…)`, or the generic `g.Attr`/`g.Prop`/`g.On`
- **children** — `*g.Node` values, plain strings (become text nodes), and
  `[]*g.Node` slices
- **nil** — skipped, which makes inline conditionals pleasant

`g.Fragment(…)` groups children without a wrapper element, like `<>…</>`.

### Conditionals and lists

```go
g.Div(
	g.If(loggedIn, g.Span("welcome back")),         // nil when false
	g.IfElse(busy, Spinner(), Content()),
	g.Ul(g.Map(todos, func(t Todo) *g.Node {
		return g.Li(g.Key(t.ID), t.Title)             // keys make reorders cheap
	})),
)
```

Give list children keys — `g.Key(…)` for elements, `.WithKey(…)` for
component nodes — so the reconciler can move DOM nodes instead of
rebuilding them.

### Skipping renders with Memo

`g.Memo(g.C(Row, props))` skips the component's re-render when its props
are unchanged (`==`); state updates inside it and context changes above it
still apply. For props containing slices or callbacks, `g.MemoEq` takes a
custom comparison over the fields that matter. A skipped component keeps
the handlers from its last render, so handlers inside a Memo subtree
should read changing data through a `UseRef`.

## Routing

`grove/router` is a small hash router (`#/event/42` — shareable links and
back/forward with no server config; in tests the path lives in memory):

```go
router.Routes(
	router.Route{Pattern: "/", Render: home},
	router.Route{Pattern: "/event/:id", Render: func(p router.Params) *g.Node {
		return g.C(EventPage, p["id"])
	}},
	router.Route{Pattern: "*", Render: notFound},
)

router.Link("/event/42", "open")   // anchor that navigates client-side
router.Navigate("/")               // programmatic navigation
```

## Hooks

The built-in hooks mirror React's, including the rules of hooks: call them
unconditionally, in the same order, on every render (grove panics with a
clear message if a component breaks this).

| hook | use |
| --- | --- |
| `UseState(initial)` | state + setter; equal values are a no-op |
| `UseStateLazy(func() T)` | state with a computed initial value |
| `UseReducer(reducer, initial)` | updates that derive from the latest state |
| `UseEffect(setup, deps)` | run after commit; setup may return a cleanup |
| `UseMemo(compute, deps)` | cache a computation |
| `UseCallback(fn, deps)` | cache a function value |
| `UseRef(initial)` | mutable box that survives renders |
| `UseContext(ctx)` | read the nearest `ctx.Provider` value |

`deps` works exactly like React's dependency array:

```go
g.UseEffect(fn, nil)          // after every render
g.UseEffect(fn, []any{})      // once, on mount
g.UseEffect(fn, []any{a, b})  // when a or b changes (shallow ==)
```

State updates are batched: every setter called in one event handler (or
one effect pass) produces a single re-render. Like React, `UseState`
values are snapshots of the render they came from — use `UseReducer` when
an update must read the latest state.

Context flows down the tree without prop drilling:

```go
var Theme = g.NewContext("light")

// provide
Theme.Provider(theme, g.C0(Page))

// consume, anywhere below
theme := g.UseContext(Theme)
```

## Events

Handlers receive `*g.Event` with typed accessors — `Value()`, `Checked()`,
`Key()` — plus `PreventDefault`, `StopPropagation`, and path-based access
(`e.Num("clientX")`) to anything else on the platform event. Under the
hood grove attaches **one** delegated listener per event type and bubbles
events through the virtual tree, so handlers cost nothing per element.

Controlled inputs work like React's:

```go
text, setText := g.UseState("")
g.Input(g.Value(text), g.OnInput(func(e *g.Event) { setText(e.Value()) }))
```

The `value`/`checked` properties are re-synced on every render, so the DOM
can't drift from your state.

## Styling

grove treats Tailwind class strings as the styling primitive, and ships a
themeable design system on top:

1. **Theme** — `grove init` writes a CSS-variable theme (`--background`,
   `--primary`, `--radius`, … with a `.dark` variant), so classes like
   `bg-primary text-primary-foreground` work out of the box and restyling
   an app means editing variables, not components. Toggle dark mode with
   `dom.SetRootClass("dark", on)`.
2. **Class scanning** — the generated `styles/input.css` tells Tailwind to
   scan `**/*.go`, so utilities used in Go string literals are compiled in.
3. **`style.CN`** — conditional class composition plus Tailwind conflict
   resolution (`CN("p-4 bg-muted", userClass)` lets the caller's classes
   win), and `style.Variants` for components whose look is picked by named
   variants. Conflict resolution covers the full utility surface —
   shorthand/longhand hierarchies (`p`/`px`/`pt`, `inset`, `rounded`
   corners, border sides…), value-classified groups (text size vs color,
   shadow size vs color, gradient stop position vs color), variant and
   important scoping, and arbitrary values/properties.
4. **Components** — the `ui` package: Button, Badge, Card, Input, Label,
   Checkbox, Separator, Alert, Avatar, Switch, Tooltip, a Calendar (single
   date or date range, month navigation, min/max bounds), dropdown
   DatePicker and TimePicker (an input-styled trigger that opens the
   calendar, or hour/minute columns, in a popover), a modal Dialog
   (Escape/overlay dismissal, focus trapping), and anchored Popover +
   Dropdown menus (side/align placement, outside-click and Escape
   dismissal, arrow-key focus, and viewport collision handling — panels
   flip to the other side, shift until they fit, and re-measure on window
   resize). All of it is plain Tailwind on the theme variables:

```go
ui.Card(
	ui.CardHeader(ui.CardTitle("Create account")),
	ui.CardContent(
		ui.Input(ui.InputProps{Value: name, OnInput: onName}),
	),
	ui.CardFooter(
		ui.Button(ui.ButtonProps{Variant: ui.ButtonDestructive}, "Delete"),
	),
)
```

Components are meant to be **owned, not imported**:

```sh
grove add button card dialog    # copies the source into ./ui/, edit freely
```

(Importing `github.com/gyoumi/grove/ui` directly also works — add an extra
`@source` line for grove's module path to your CSS so Tailwind sees those
class strings.)

## React islands

A grove app can hand individual leaf nodes to React (or any JS renderer)
with the `island` package. The page registers islands on
`window.groveIslands` before the wasm starts:

```js
import { createRoot } from "https://esm.sh/react-dom@18/client";

window.groveIslands = {
  Greeting: {
    mount(el, props)  { el._root = createRoot(el); el._root.render(e(Greeting, props)) },
    update(el, props) { el._root.render(e(Greeting, props)) }, // optional
    unmount(el)       { el._root.unmount() },                  // optional
  },
};
```

and grove places them like any other node, re-delivering props (as JSON)
whenever they change:

```go
island.C("Greeting", map[string]any{"count": count}, g.Class("rounded-lg border p-4"))
```

The JS side owns everything inside the container element — islands are
leaves, with no grove children. Outside the browser the lifecycle is
observable in tests via `island.SetHost`. A runnable app lives in
[`examples/islands`](examples/islands).

## CLI

| command | what it does |
| --- | --- |
| `grove init <app>` | scaffold an app (Tailwind + themeable design system, no Node) |
| `grove serve` | dev server: rebuild on save, SSE live reload |
| `grove build` | release build: `-s -w`, minified CSS, optional `wasm-opt`, size report |
| `grove add <component>` | copy a ui component's source into your app |

`serve` and `build` take `-tinygo` to compile with TinyGo instead of the
standard toolchain (it must be on PATH); the matching `wasm_exec.js` is
placed in `dist/` automatically.

## Testing

The engine never touches `syscall/js` directly — it renders through a
small `Renderer` interface. `testdom` implements it in memory, so
components are testable with plain `go test`, no browser:

```go
func TestCounter(t *testing.T) {
	r := testdom.Mount(g.C0(Counter))
	r.Click(r.Find("button"))
	if got := r.HTML(); got != "<button>count is 1</button>" {
		t.Fatalf("got %s", got)
	}
}
```

`testdom` can fire events (`Click`, `Input`, `KeyDown`, `SetChecked`),
query the fake DOM (`Find`, `FindByAttr`, `FindText`), and snapshot it as
HTML.

## Bundle size

A hello-world app is ~2.6 MB of wasm (~740 KB gzipped) with the standard
toolchain — that's the cost of shipping the Go runtime, and it's a flat
cost, not per-component. `grove build` strips symbols, runs `wasm-opt`
when available, and prints both raw and gzipped sizes. Serve wasm with
gzip or brotli enabled.

With `grove build -tinygo` the same app is **~295 KB raw (~120 KB
gzipped)**. grove's engine runs its full test suite under the TinyGo
compiler, so the engine itself is safe there; the trade-offs are TinyGo's:
slower compiles, and reflection-heavy stdlib packages (`encoding/json`
works but costs size) or unsupported ones may not carry over. The standard
toolchain stays the default.

## Roadmap

- Batched DOM patch protocol to cut wasm↔JS call overhead
- Longest-increasing-subsequence move optimization for keyed lists

Not on the roadmap: server-side rendering and hydration — grove stays a
client-side framework.

## Examples

[`examples/`](examples/) contains four runnable apps — `counter` (smallest
possible), `todo` (state, lists, keys), `showcase` (every ui component,
dark mode, dialog), and `islands` (a React component living inside a grove
app). Run any of them with `grove serve` from its directory.
