package ui_test

import (
	"strings"
	"testing"

	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func TestIconRendersSVG(t *testing.T) {
	r := testdom.Mount(ui.Icon("check", "size-5 text-emerald-500"))
	svg := r.Find("svg")
	if svg == nil {
		t.Fatalf("icon should render an svg: %s", r.HTML())
	}
	if svg.Attrs["data-icon"] != "check" {
		t.Fatalf("data-icon should be set: %s", svg.HTML())
	}
	cls := svg.Attrs["class"]
	if !strings.Contains(cls, "size-5") || strings.Contains(cls, "size-4") || !strings.Contains(cls, "text-emerald-500") {
		t.Fatalf("extra classes should override the default size: %q", cls)
	}
	if r.Find("path") == nil {
		t.Fatalf("check icon should have a path: %s", r.HTML())
	}

	if n := len(testdom.Mount(ui.Icon("x")).FindAll("path")); n != 2 {
		t.Fatalf("x icon should have 2 paths, got %d", n)
	}
	if testdom.Mount(ui.Icon("search")).Find("circle") == nil {
		t.Fatal("search icon should include a circle")
	}
}
