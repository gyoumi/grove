package ui_test

import (
	"strings"
	"testing"

	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func TestProgressClamps(t *testing.T) {
	for _, c := range []struct {
		in   float64
		want string
	}{{150, "width: 100%"}, {-5, "width: 0%"}, {40, "width: 40%"}} {
		r := testdom.Mount(ui.Progress(c.in))
		bar := r.FindByAttr("data-slot", "progress-bar")
		if bar == nil || !strings.Contains(bar.Attrs["style"], c.want) {
			t.Fatalf("Progress(%v) bar style = %q, want %q", c.in, bar.HTML(), c.want)
		}
	}
}

func TestNativeSelectRendersOptionsAndValue(t *testing.T) {
	r := testdom.Mount(ui.NativeSelect(ui.NativeSelectProps{
		Value:       "b",
		Placeholder: "Choose…",
		Options:     []ui.SelectOption{{Value: "a", Label: "A"}, {Value: "b", Label: "B"}},
	}))
	if n := len(r.FindAll("option")); n != 3 { // placeholder + 2
		t.Fatalf("expected 3 options, got %d: %s", n, r.HTML())
	}
	sel := r.FindByAttr("data-slot", "native-select")
	if sel == nil || sel.Props["value"] != "b" {
		t.Fatalf("select value should be controlled to b: %s", r.HTML())
	}
	// chevron icon overlaid
	if r.FindByAttr("data-icon", "chevron-down") == nil {
		t.Fatalf("native select should show a chevron: %s", r.HTML())
	}
}

func TestTableComposesSlots(t *testing.T) {
	r := testdom.Mount(ui.Table(
		ui.TableHeader(ui.TableRow(ui.TableHead("Name"))),
		ui.TableBody(ui.TableRow(ui.TableCell("Ada"))),
	))
	for _, slot := range []string{"table", "table-header", "table-body", "table-row", "table-head", "table-cell"} {
		if r.FindByAttr("data-slot", slot) == nil {
			t.Fatalf("missing %s: %s", slot, r.HTML())
		}
	}
	if r.Find("table") == nil || r.Find("thead") == nil || r.Find("td") == nil {
		t.Fatalf("table should render real table elements: %s", r.HTML())
	}
}

func TestBreadcrumbSeparatorIcon(t *testing.T) {
	r := testdom.Mount(ui.BreadcrumbSeparator())
	if r.FindByAttr("data-icon", "chevron-right") == nil {
		t.Fatalf("default separator should render a chevron: %s", r.HTML())
	}
}
