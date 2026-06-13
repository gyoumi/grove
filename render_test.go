package grove_test

import (
	"strconv"
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
)

func TestStaticRender(t *testing.T) {
	root := g.Div(g.Class("a"), g.Class("b"), g.Data("x", "1"),
		g.H1("Title"),
		g.Fragment(g.Span("f1"), g.Span("f2")),
		g.Input(g.Type("text"), g.Value("hello"), g.Placeholder("p")),
	)
	r := testdom.Mount(root)
	want := `<div class="a b" data-x="1"><h1>Title</h1><span>f1</span><span>f2</span><input placeholder="p" type="text" value="hello"/></div>`
	if got := r.HTML(); got != want {
		t.Fatalf("html mismatch:\n got %s\nwant %s", got, want)
	}
}

func Counter() *g.Node {
	count, setCount := g.UseState(0)
	return g.Button(
		g.OnClick(func(*g.Event) { setCount(count + 1) }),
		g.Textf("count: %d", count),
	)
}

func TestCounter(t *testing.T) {
	r := testdom.Mount(g.C0(Counter))
	if got := r.HTML(); got != "<button>count: 0</button>" {
		t.Fatalf("initial: %s", got)
	}
	r.Click(r.Find("button"))
	if got := r.HTML(); got != "<button>count: 1</button>" {
		t.Fatalf("after click: %s", got)
	}
	r.Click(r.Find("button"))
	if got := r.HTML(); got != "<button>count: 2</button>" {
		t.Fatalf("after 2nd click: %s", got)
	}
}

// UseState returns a render-time snapshot, like React: two setCount(count+1)
// calls in one handler land on the same base value.
func SnapshotCounter() *g.Node {
	count, setCount := g.UseState(0)
	return g.Button(
		g.OnClick(func(*g.Event) { setCount(count + 1); setCount(count + 1) }),
		g.Textf("%d", count),
	)
}

var reducerRenders int

func ReducerCounter() *g.Node {
	reducerRenders++
	count, dispatch := g.UseReducer(func(s int, _ struct{}) int { return s + 1 }, 0)
	return g.Button(
		g.OnClick(func(*g.Event) { dispatch(struct{}{}); dispatch(struct{}{}) }),
		g.Textf("%d", count),
	)
}

func TestStateSnapshotAndReducerBatching(t *testing.T) {
	r := testdom.Mount(g.C0(SnapshotCounter))
	r.Click(r.Find("button"))
	if got := r.HTML(); got != "<button>1</button>" {
		t.Fatalf("snapshot semantics: want 1, got %s", got)
	}

	reducerRenders = 0
	r = testdom.Mount(g.C0(ReducerCounter))
	if reducerRenders != 1 {
		t.Fatalf("mount renders = %d, want 1", reducerRenders)
	}
	r.Click(r.Find("button"))
	if got := r.HTML(); got != "<button>2</button>" {
		t.Fatalf("reducer should see latest state: %s", got)
	}
	if reducerRenders != 2 {
		t.Fatalf("two dispatches should batch into one render, renders = %d", reducerRenders)
	}
}

func ControlledInput() *g.Node {
	text, setText := g.UseState("")
	return g.Div(
		g.Input(g.Value(text), g.OnInput(func(e *g.Event) { setText(e.Value()) })),
		g.P(g.Textf("you typed: %s", text)),
	)
}

func TestControlledInput(t *testing.T) {
	r := testdom.Mount(g.C0(ControlledInput))
	r.Input(r.Find("input"), "hi")
	if got := r.Find("p").HTML(); got != "<p>you typed: hi</p>" {
		t.Fatalf("after input: %s", got)
	}
	if v := r.Find("input").Props["value"]; v != "hi" {
		t.Fatalf("input value prop = %v, want hi", v)
	}
}

type listProps struct{ order []string }

func KeyedList(p listProps) *g.Node {
	return g.Ul(g.Map(p.order, func(s string) *g.Node {
		return g.Li(g.Key(s), s)
	}))
}

func ReorderApp() *g.Node {
	order, setOrder := g.UseState([]string{"a", "b", "c"})
	return g.Div(
		g.Button(g.OnClick(func(*g.Event) { setOrder([]string{"c", "a", "b"}) }), "shuffle"),
		g.C(KeyedList, listProps{order}),
	)
}

func TestKeyedReorderPreservesIdentity(t *testing.T) {
	r := testdom.Mount(g.C0(ReorderApp))
	before := map[string]*testdom.Elem{}
	for _, li := range r.FindAll("li") {
		before[li.HTML()] = li
	}
	r.Click(r.Find("button"))

	lis := r.FindAll("li")
	var texts []string
	for _, li := range lis {
		texts = append(texts, li.Children[0].Text)
	}
	if len(texts) != 3 || texts[0] != "c" || texts[1] != "a" || texts[2] != "b" {
		t.Fatalf("order after shuffle: %v", texts)
	}
	for _, li := range lis {
		key := li.HTML()
		if before[key] != li {
			t.Fatalf("li %s lost DOM identity across reorder", key)
		}
	}
}

var applyOrder func([]string)

func reorderProbe() *g.Node {
	order, set := g.UseState([]string{"a", "b", "c", "d", "e"})
	applyOrder = set
	return g.Ul(g.Map(order, func(s string) *g.Node {
		return g.Li(g.Key(s), s)
	}))
}

func liOrder(r *testdom.R) []string {
	var out []string
	for _, li := range r.FindAll("li") {
		out = append(out, li.Children[0].Text)
	}
	return out
}

func sameStrs(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// minimalMoves is the fewest DOM moves a keyed reorder needs — the count
// outside the longest increasing subsequence — computed here with a plain
// O(n^2) LIS so the test doesn't lean on the reconciler's own routine.
func minimalMoves(oldOrder, newOrder []string) int {
	pos := map[string]int{}
	for i, s := range oldOrder {
		pos[s] = i
	}
	seq := make([]int, len(newOrder))
	for i, s := range newOrder {
		seq[i] = pos[s]
	}
	best := 0
	dp := make([]int, len(seq))
	for i := range seq {
		dp[i] = 1
		for j := range i {
			if seq[j] < seq[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > best {
			best = dp[i]
		}
	}
	return len(newOrder) - best
}

// The reconciler should move the minimum number of DOM nodes for a keyed
// reorder: only those outside the longest increasing subsequence. The old
// greedy heuristic moved up to N-1 for a single-element rotation.
func TestKeyedReorderMinimalMoves(t *testing.T) {
	r := testdom.Mount(g.C0(reorderProbe))
	id := map[string]*testdom.Elem{}
	for _, li := range r.FindAll("li") {
		id[li.Children[0].Text] = li
	}

	cur := []string{"a", "b", "c", "d", "e"}
	cases := [][]string{
		{"b", "c", "d", "e", "a"}, // rotate left   → 1 move
		{"d", "e", "a", "b", "c"}, // rotate right×2 → 2 moves
		{"e", "d", "c", "b", "a"}, // reverse        → 4 moves
		{"e", "a", "b", "c", "d"}, // rotate right   → 1 move
		{"e", "c", "a", "d", "b"}, // scramble
		{"a", "b", "c", "d", "e"}, // identity       → 0 moves
	}
	for _, next := range cases {
		want := minimalMoves(cur, next)
		r.ResetMoves()
		applyOrder(next)
		r.Settle()

		if got := liOrder(r); !sameStrs(got, next) {
			t.Fatalf("order after %v→%v: got %v", cur, next, got)
		}
		if r.Moves() != want {
			t.Fatalf("reorder %v→%v: %d moves, want minimal %d", cur, next, r.Moves(), want)
		}
		cur = next
	}

	// every key kept its original DOM node across all the reorders
	for _, li := range r.FindAll("li") {
		text := li.Children[0].Text
		if id[text] != li {
			t.Fatalf("li %q lost DOM identity across reorders", text)
		}
	}
}

// Reorders mixed with insertions and removals must still land in the right
// order; new keys mount, missing keys unmount, the rest move minimally.
func TestKeyedReorderWithInsertAndRemove(t *testing.T) {
	r := testdom.Mount(g.C0(reorderProbe))
	applyOrder([]string{"x", "c", "a", "y", "e"}) // drop b,d; add x,y; reorder
	r.Settle()
	if got := liOrder(r); !sameStrs(got, []string{"x", "c", "a", "y", "e"}) {
		t.Fatalf("order: %v", got)
	}
}

func ToggleInner() *g.Node {
	on, setOn := g.UseState(false)
	return g.Fragment(
		g.Button(g.OnClick(func(*g.Event) { setOn(!on) }), "t"),
		g.If(on, g.Span("inner")),
	)
}

func ToggleApp() *g.Node {
	return g.Div(g.Span("a"), g.C0(ToggleInner), g.Span("z"))
}

// A self-re-rendering component sandwiched between siblings must insert its
// new DOM at the right position (anchor computation across fragments).
func TestComponentRerenderPosition(t *testing.T) {
	r := testdom.Mount(g.C0(ToggleApp))
	r.Click(r.Find("button"))
	want := "<div><span>a</span><button>t</button><span>inner</span><span>z</span></div>"
	if got := r.HTML(); got != want {
		t.Fatalf("after show:\n got %s\nwant %s", got, want)
	}
	r.Click(r.Find("button"))
	want = "<div><span>a</span><button>t</button><span>z</span></div>"
	if got := r.HTML(); got != want {
		t.Fatalf("after hide:\n got %s\nwant %s", got, want)
	}
}

func SwapApp() *g.Node {
	on, setOn := g.UseState(false)
	child := g.Em("x")
	if on {
		child = g.Strong("x")
	}
	return g.Div(
		g.Button(g.OnClick(func(*g.Event) { setOn(!on) }), "go"),
		child,
	)
}

func TestTypeChangeReplaces(t *testing.T) {
	r := testdom.Mount(g.C0(SwapApp))
	if r.Find("em") == nil {
		t.Fatal("expected <em> initially")
	}
	r.Click(r.Find("button"))
	if r.Find("em") != nil || r.Find("strong") == nil {
		t.Fatalf("expected <strong> after toggle: %s", r.HTML())
	}
	want := "<div><button>go</button><strong>x</strong></div>"
	if got := r.HTML(); got != want {
		t.Fatalf("replacement position wrong:\n got %s\nwant %s", got, want)
	}
}

var bubbleLog []string

func BubbleApp() *g.Node {
	return g.Div(g.OnClick(func(*g.Event) { bubbleLog = append(bubbleLog, "outer") }),
		g.Button(g.ID("stop"), g.OnClick(func(e *g.Event) {
			bubbleLog = append(bubbleLog, "inner")
			e.StopPropagation()
		}), "in"),
		g.Button(g.ID("pass"), g.OnClick(func(*g.Event) {
			bubbleLog = append(bubbleLog, "inner2")
		}), "in2"),
	)
}

func TestEventBubblingAndStopPropagation(t *testing.T) {
	bubbleLog = nil
	r := testdom.Mount(BubbleApp())
	r.Click(r.FindByAttr("id", "stop"))
	if len(bubbleLog) != 1 || bubbleLog[0] != "inner" {
		t.Fatalf("stopPropagation failed: %v", bubbleLog)
	}
	bubbleLog = nil
	r.Click(r.FindByAttr("id", "pass"))
	if len(bubbleLog) != 2 || bubbleLog[0] != "inner2" || bubbleLog[1] != "outer" {
		t.Fatalf("bubbling failed: %v", bubbleLog)
	}
}

type todoItem struct {
	id   int
	text string
	done bool
}

func TodoApp() *g.Node {
	todos, setTodos := g.UseState([]todoItem(nil))
	nextID, setNextID := g.UseState(1)
	input, setInput := g.UseState("")

	add := func(*g.Event) {
		if input == "" {
			return
		}
		setTodos(append(todos[:len(todos):len(todos)], todoItem{nextID, input, false}))
		setNextID(nextID + 1)
		setInput("")
	}
	toggle := func(id int) func(*g.Event) {
		return func(*g.Event) {
			out := make([]todoItem, len(todos))
			copy(out, todos)
			for i := range out {
				if out[i].id == id {
					out[i].done = !out[i].done
				}
			}
			setTodos(out)
		}
	}
	del := func(id int) func(*g.Event) {
		return func(*g.Event) {
			var out []todoItem
			for _, td := range todos {
				if td.id != id {
					out = append(out, td)
				}
			}
			setTodos(out)
		}
	}

	return g.Div(
		g.Input(g.Type("text"), g.Value(input), g.OnInput(func(e *g.Event) { setInput(e.Value()) })),
		g.Button(g.Data("action", "add"), g.OnClick(add), "add"),
		g.Ul(g.Map(todos, func(td todoItem) *g.Node {
			id := strconv.Itoa(td.id)
			return g.Li(g.Key(id), g.ClassIf(td.done, "done"),
				g.Span(td.text),
				g.Button(g.Data("action", "toggle-"+id), g.OnClick(toggle(td.id)), "toggle"),
				g.Button(g.Data("action", "del-"+id), g.OnClick(del(td.id)), "x"),
			)
		})),
	)
}

func TestTodoAppSimulation(t *testing.T) {
	r := testdom.Mount(g.C0(TodoApp))

	addTodo := func(text string) {
		r.Input(r.Find("input"), text)
		r.Click(r.FindByAttr("data-action", "add"))
	}

	addTodo("alpha")
	addTodo("beta")
	lis := r.FindAll("li")
	if len(lis) != 2 {
		t.Fatalf("want 2 todos, html: %s", r.HTML())
	}
	if v := r.Find("input").Props["value"]; v != "" {
		t.Fatalf("controlled input should be cleared after add, got %v", v)
	}

	betaLi := lis[1]
	r.Click(r.FindByAttr("data-action", "toggle-1"))
	if cls := r.FindAll("li")[0].Attrs["class"]; cls != "done" {
		t.Fatalf("toggle should add done class, got %q (html %s)", cls, r.HTML())
	}
	r.Click(r.FindByAttr("data-action", "toggle-1"))
	if cls := r.FindAll("li")[0].Attrs["class"]; cls != "" {
		t.Fatalf("untoggle should remove done class, got %q", cls)
	}

	r.Click(r.FindByAttr("data-action", "del-1"))
	lis = r.FindAll("li")
	if len(lis) != 1 || lis[0].Children[0].TextContent() != "beta" {
		t.Fatalf("after delete: %s", r.HTML())
	}
	if lis[0] != betaLi {
		t.Fatal("surviving todo lost DOM identity when its sibling was removed")
	}
}
