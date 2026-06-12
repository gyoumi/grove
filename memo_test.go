package grove_test

import (
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
)

var memoChildRenders int

type rowProps struct {
	Label string
}

func MemoRow(p rowProps) *g.Node {
	memoChildRenders++
	return g.Li(p.Label)
}

func MemoSkipApp() *g.Node {
	n, setN := g.UseState(0)
	label, setLabel := g.UseState("row")
	return g.Div(
		g.Button(g.ID("bump"), g.OnClick(func(*g.Event) { setN(n + 1) }), g.Textf("n=%d", n)),
		g.Button(g.ID("relabel"), g.OnClick(func(*g.Event) { setLabel("changed") }), "relabel"),
		g.Ul(g.Memo(g.C(MemoRow, rowProps{Label: label}))),
	)
}

func TestMemoSkipsUnchangedProps(t *testing.T) {
	memoChildRenders = 0
	r := testdom.Mount(g.C0(MemoSkipApp))
	if memoChildRenders != 1 {
		t.Fatalf("mount renders = %d", memoChildRenders)
	}

	r.Click(r.FindByAttr("id", "bump")) // parent re-renders, row props unchanged
	if memoChildRenders != 1 {
		t.Fatalf("memo child should be skipped, renders = %d", memoChildRenders)
	}

	r.Click(r.FindByAttr("id", "relabel")) // props change
	if memoChildRenders != 2 {
		t.Fatalf("memo child should re-render on prop change, renders = %d", memoChildRenders)
	}
	if r.Find("li").TextContent() != "changed" {
		t.Fatalf("memo child stale: %s", r.HTML())
	}
}

var statefulMemoRenders int

func StatefulMemo() *g.Node {
	statefulMemoRenders++
	n, setN := g.UseState(0)
	return g.Button(g.ID("inner"), g.OnClick(func(*g.Event) { setN(n + 1) }), g.Textf("inner=%d", n))
}

func MemoStateApp() *g.Node {
	n, setN := g.UseState(0)
	return g.Div(
		g.Button(g.ID("outer"), g.OnClick(func(*g.Event) { setN(n + 1) }), "outer"),
		g.Memo(g.C0(StatefulMemo)),
	)
}

func TestMemoStateStillUpdates(t *testing.T) {
	statefulMemoRenders = 0
	r := testdom.Mount(g.C0(MemoStateApp))
	r.Click(r.FindByAttr("id", "outer")) // C0 props equal → skipped
	if statefulMemoRenders != 1 {
		t.Fatalf("C0 memo should skip on parent render: %d", statefulMemoRenders)
	}
	r.Click(r.FindByAttr("id", "inner")) // own state → must render
	if statefulMemoRenders != 2 || r.FindByAttr("id", "inner").TextContent() != "inner=1" {
		t.Fatalf("memo component's own state update failed: renders=%d html=%s", statefulMemoRenders, r.HTML())
	}
}

var memoCtx = g.NewContext("light")

func MemoCtxLeaf() *g.Node {
	return g.Span(g.UseContext(memoCtx))
}

func MemoCtxMiddle() *g.Node {
	// a memo boundary between the provider and the consumer
	return g.Div(g.Memo(g.C0(MemoCtxLeaf)))
}

func MemoCtxApp() *g.Node {
	theme, setTheme := g.UseState("light")
	return g.Div(
		g.Button(g.OnClick(func(*g.Event) { setTheme("dark") }), "switch"),
		memoCtx.Provider(theme, g.Memo(g.C0(MemoCtxMiddle))),
	)
}

func TestContextCrossesMemoBoundary(t *testing.T) {
	r := testdom.Mount(g.C0(MemoCtxApp))
	r.Click(r.Find("button"))
	if got := r.Find("span").TextContent(); got != "dark" {
		t.Fatalf("context change should pierce memo boundaries, got %q", got)
	}
}

type eqProps struct {
	Values []string // incomparable: defeats plain Memo
	OnPick func(string)
}

var eqRenders int

func EqRow(p eqProps) *g.Node {
	eqRenders++
	return g.Li(g.Textf("%d items", len(p.Values)))
}

func MemoEqApp() *g.Node {
	n, setN := g.UseState(0)
	vals := []string{"a", "b"}
	return g.Div(
		g.Button(g.OnClick(func(*g.Event) { setN(n + 1) }), "bump"),
		g.Ul(g.MemoEq(EqRow, eqProps{Values: vals, OnPick: func(string) {}},
			func(old, new eqProps) bool {
				if len(old.Values) != len(new.Values) {
					return false
				}
				for i := range old.Values {
					if old.Values[i] != new.Values[i] {
						return false
					}
				}
				return true
			})),
	)
}

func TestMemoEqCustomComparison(t *testing.T) {
	eqRenders = 0
	r := testdom.Mount(g.C0(MemoEqApp))
	r.Click(r.Find("button"))
	r.Click(r.Find("button"))
	if eqRenders != 1 {
		t.Fatalf("MemoEq should skip when compared fields match: renders=%d", eqRenders)
	}
}
