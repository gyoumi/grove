package gallery_test

import (
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/gallery"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func TestGalleryRenders(t *testing.T) {
	ui.DismissAllToasts() // the toast store is global; isolate the mount
	r := testdom.Mount(g.C0(gallery.Page))

	if r.FindText("Component gallery") == nil {
		t.Fatal("gallery should render its title")
	}
	// a representative spread across every section, including the heavy set
	for _, slot := range []string{
		"button", "badge", "tabs", "accordion", "select-trigger", "table",
		"radio-group", "slider", "breadcrumb", "menubar", "card", "chart",
		"carousel", "sidebar", "resizable-handle",
	} {
		if r.FindByAttr("data-slot", slot) == nil {
			t.Fatalf("gallery is missing the %q component", slot)
		}
	}
}
