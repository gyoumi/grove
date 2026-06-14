package grove_test

import (
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
)

func portalApp() *g.Node {
	open, set := g.UseState(false)
	return g.Div(g.Class("transformed"),
		g.Button(g.Data("slot", "portal-open"), g.OnClick(func(*g.Event) { set(true) }), "open"),
		g.If(open, g.Portal(g.Div(g.Data("slot", "portaled"),
			g.Button(g.Data("slot", "portal-close"), g.OnClick(func(*g.Event) { set(false) }), "close"),
		))),
	)
}

// Portal children mount directly under the container — escaping the
// transformed app root — yet stay wired for events through the virtual tree.
func TestPortalEscapesToContainerAndCleansUp(t *testing.T) {
	r := testdom.Mount(g.C0(portalApp))
	if r.FindByAttr("data-slot", "portaled") != nil {
		t.Fatal("portal content should not render while closed")
	}

	r.Click(r.FindByAttr("data-slot", "portal-open"))
	p := r.FindByAttr("data-slot", "portaled")
	if p == nil {
		t.Fatalf("portal should render when open: %s", r.HTML())
	}
	// It must be a direct child of the container, not nested inside the
	// transformed app root (which would trap fixed positioning).
	direct := false
	for _, c := range r.Container.Children {
		if c == p {
			direct = true
		}
	}
	if !direct {
		t.Fatalf("portal content should mount directly under the container: %s", r.HTML())
	}

	// A click inside the portal reaches its handler via the virtual tree.
	r.Click(r.FindByAttr("data-slot", "portal-close"))
	if r.FindByAttr("data-slot", "portaled") != nil {
		t.Fatalf("closing should unmount the portal content: %s", r.HTML())
	}
}
