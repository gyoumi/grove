package grove_test

import (
	"strconv"
	"strings"
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
)

var effLog []string

type effProps struct{ dep int }

func EffChild(p effProps) *g.Node {
	g.UseEffect(func() func() {
		effLog = append(effLog, "setup:"+strconv.Itoa(p.dep))
		return func() { effLog = append(effLog, "cleanup:"+strconv.Itoa(p.dep)) }
	}, []any{p.dep})
	g.UseEffect(func() func() {
		effLog = append(effLog, "once")
		return nil
	}, []any{})
	return g.Span("child")
}

func EffApp() *g.Node {
	dep, setDep := g.UseState(0)
	show, setShow := g.UseState(true)
	return g.Div(
		g.Button(g.ID("inc"), g.OnClick(func(*g.Event) { setDep(dep + 1) }), "inc"),
		g.Button(g.ID("same"), g.OnClick(func(*g.Event) { setDep(dep) }), "same"),
		g.Button(g.ID("hide"), g.OnClick(func(*g.Event) { setShow(false) }), "hide"),
		g.If(show, g.C(EffChild, effProps{dep: dep})),
	)
}

func TestEffectLifecycle(t *testing.T) {
	effLog = nil
	r := testdom.Mount(g.C0(EffApp))
	assertLog := func(want ...string) {
		t.Helper()
		if strings.Join(effLog, ",") != strings.Join(want, ",") {
			t.Fatalf("effect log:\n got %v\nwant %v", effLog, want)
		}
	}
	assertLog("setup:0", "once")

	r.Click(r.FindByAttr("id", "inc")) // dep 0→1: cleanup+setup, "once" untouched
	assertLog("setup:0", "once", "cleanup:0", "setup:1")

	r.Click(r.FindByAttr("id", "same")) // equal state: no render, no effects
	assertLog("setup:0", "once", "cleanup:0", "setup:1")

	r.Click(r.FindByAttr("id", "hide")) // unmount: cleanup runs
	assertLog("setup:0", "once", "cleanup:0", "setup:1", "cleanup:1")
	if r.Find("span") != nil {
		t.Fatal("child should be unmounted")
	}
}

// An effect that sets state must trigger a follow-up render in the same
// settle cycle.
func EffectSetsState() *g.Node {
	n, setN := g.UseState(0)
	g.UseEffect(func() func() {
		setN(42)
		return nil
	}, []any{})
	return g.Span(g.Textf("%d", n))
}

func TestEffectStateUpdate(t *testing.T) {
	r := testdom.Mount(g.C0(EffectSetsState))
	if got := r.HTML(); got != "<span>42</span>" {
		t.Fatalf("got %s", got)
	}
}

var memoComputes int

func MemoApp() *g.Node {
	a, setA := g.UseState(1)
	b, setB := g.UseState(1)
	doubled := g.UseMemo(func() int { memoComputes++; return a * 2 }, []any{a})
	return g.Div(
		g.Button(g.ID("a"), g.OnClick(func(*g.Event) { setA(a + 1) }), "a"),
		g.Button(g.ID("b"), g.OnClick(func(*g.Event) { setB(b + 1) }), "b"),
		g.Span(g.Textf("%d-%d", doubled, b)),
	)
}

func TestUseMemo(t *testing.T) {
	memoComputes = 0
	r := testdom.Mount(g.C0(MemoApp))
	if memoComputes != 1 {
		t.Fatalf("computes after mount = %d", memoComputes)
	}
	r.Click(r.FindByAttr("id", "b")) // unrelated state: no recompute
	if memoComputes != 1 {
		t.Fatalf("memo recomputed on unrelated update: %d", memoComputes)
	}
	r.Click(r.FindByAttr("id", "a"))
	if memoComputes != 2 {
		t.Fatalf("memo should recompute when dep changes: %d", memoComputes)
	}
	if got := r.Find("span").TextContent(); got != "4-2" {
		t.Fatalf("got %s", got)
	}
}

func RefApp() *g.Node {
	renders := g.UseRef(0)
	renders.Current++
	count, setCount := g.UseState(0)
	return g.Button(
		g.OnClick(func(*g.Event) { setCount(count + 1) }),
		g.Textf("%d renders, count %d", renders.Current, count),
	)
}

func TestUseRefPersists(t *testing.T) {
	r := testdom.Mount(g.C0(RefApp))
	r.Click(r.Find("button"))
	if got := r.HTML(); got != "<button>2 renders, count 1</button>" {
		t.Fatalf("got %s", got)
	}
}

var themeCtx = g.NewContext("light")

func ThemeLabel() *g.Node {
	theme := g.UseContext(themeCtx)
	return g.Span(theme)
}

func ThemeApp() *g.Node {
	theme, setTheme := g.UseState("light")
	return g.Div(
		g.Button(g.OnClick(func(*g.Event) { setTheme("dark") }), "switch"),
		themeCtx.Provider(theme,
			g.C0(ThemeLabel),
		),
		g.C0(ThemeLabel), // outside the provider: sees the default
	)
}

func TestContext(t *testing.T) {
	r := testdom.Mount(g.C0(ThemeApp))
	spans := r.FindAll("span")
	if spans[0].TextContent() != "light" || spans[1].TextContent() != "light" {
		t.Fatalf("initial: %s", r.HTML())
	}
	r.Click(r.Find("button"))
	spans = r.FindAll("span")
	if spans[0].TextContent() != "dark" {
		t.Fatalf("provider consumer should see dark: %s", r.HTML())
	}
	if spans[1].TextContent() != "light" {
		t.Fatalf("outside consumer should keep default: %s", r.HTML())
	}
}

func BadHooks() *g.Node {
	n, setN := g.UseState(0)
	if n == 0 {
		g.UseRef(0) // conditional hook: count changes on next render
	}
	return g.Button(g.OnClick(func(*g.Event) { setN(n + 1) }), "x")
}

func TestRulesOfHooksPanic(t *testing.T) {
	r := testdom.Mount(g.C0(BadHooks))
	defer func() {
		rec := recover()
		if rec == nil {
			t.Fatal("expected a rules-of-hooks panic")
		}
		if !strings.Contains(rec.(string), "rules of hooks") {
			t.Fatalf("unexpected panic: %v", rec)
		}
	}()
	r.Click(r.Find("button"))
}

func TestHookOutsideComponentPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when calling a hook outside render")
		}
	}()
	g.UseState(0)
}

func DOMRefApp() *g.Node {
	ref := g.UseRef[any](nil)
	g.UseEffect(func() func() {
		if el, ok := ref.Current.(*testdom.Elem); ok {
			el.Attrs["data-focused"] = "yes"
		}
		return nil
	}, []any{})
	return g.Input(g.BindRef(ref))
}

func TestBindRef(t *testing.T) {
	r := testdom.Mount(g.C0(DOMRefApp))
	if r.Find("input").Attrs["data-focused"] != "yes" {
		t.Fatalf("ref should expose the DOM handle in effects: %s", r.HTML())
	}
}
