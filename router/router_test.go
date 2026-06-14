package router_test

import (
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/router"
	"github.com/gyoumi/grove/testdom"
)

func routedApp() *g.Node {
	return router.Routes(
		router.Route{Pattern: "/", Render: func(router.Params) *g.Node {
			return g.Div(
				g.H1("home"),
				router.Link("/event/42", g.ID("to-event"), "open event 42"),
			)
		}},
		router.Route{Pattern: "/event/:id", Render: func(p router.Params) *g.Node {
			return g.Div(
				g.H1(g.Textf("event %s", p["id"])),
				router.Link("/", g.ID("to-home"), "back"),
			)
		}},
		router.Route{Pattern: "*", Render: func(router.Params) *g.Node {
			return g.H1("not found")
		}},
	)
}

func TestRouting(t *testing.T) {
	router.Navigate("/")
	r := testdom.Mount(routedApp())
	if r.Find("h1").TextContent() != "home" {
		t.Fatalf("should start at home: %s", r.HTML())
	}

	r.Click(r.FindByAttr("id", "to-event"))
	r.Settle()
	if got := r.Find("h1").TextContent(); got != "event 42" {
		t.Fatalf("link should navigate with params, got %q", got)
	}
	if router.Path() != "/event/42" {
		t.Fatalf("path = %q", router.Path())
	}

	router.Navigate("/nope/extra/deep")
	r.Settle()
	if got := r.Find("h1").TextContent(); got != "not found" {
		t.Fatalf("fallback route should match, got %q", got)
	}

	router.Navigate("/")
	r.Settle()
	if got := r.Find("h1").TextContent(); got != "home" {
		t.Fatalf("navigate home failed, got %q", got)
	}
}

func TestLinkHref(t *testing.T) {
	router.Navigate("/")
	r := testdom.Mount(routedApp())
	link := r.FindByAttr("id", "to-event")
	if link.Attrs["href"] != "/event/42" {
		t.Fatalf("link href = %q", link.Attrs["href"])
	}
}

func toggleApp() *g.Node {
	on, setOn := g.UseState(true)
	return g.Div(
		g.Button(g.ID("toggle"), g.OnClick(func(*g.Event) { setOn(!on) }), "toggle"),
		g.If(on, routedApp()),
	)
}

func TestUnmountUnsubscribes(t *testing.T) {
	router.Navigate("/")
	r := testdom.Mount(g.C0(toggleApp))
	if r.Find("h1") == nil {
		t.Fatal("routes should render while shown")
	}
	r.Click(r.FindByAttr("id", "toggle")) // unmounts the Routes component
	if r.Find("h1") != nil {
		t.Fatalf("routes should be unmounted: %s", r.HTML())
	}
	// navigating after unmount must not panic via a stale subscription
	router.Navigate("/event/1")
	r.Settle()
}
