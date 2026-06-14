package ui_test

import (
	"strings"
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func TestChartRendersSeries(t *testing.T) {
	r := testdom.Mount(ui.Chart(ui.ChartProps{Type: ui.ChartLine, Series: []ui.ChartSeries{
		{Values: []float64{1, 3, 2, 5}},
		{Values: []float64{2, 1, 4, 3}},
	}}))
	if r.Find("svg") == nil {
		t.Fatalf("chart should render an svg: %s", r.HTML())
	}
	paths := 0
	for _, p := range r.FindAll("path") {
		if p.Attrs["data-slot"] == "chart-series" {
			paths++
		}
	}
	if paths != 2 {
		t.Fatalf("two line series should render two paths, got %d", paths)
	}
}

func carouselApp() *g.Node {
	return ui.Carousel(ui.CarouselProps{Dots: true}, g.Div("A"), g.Div("B"), g.Div("C"))
}

func TestCarouselNavigates(t *testing.T) {
	r := testdom.Mount(g.C0(carouselApp))
	track := r.FindByAttr("data-slot", "carousel-track")
	if !strings.Contains(track.Attrs["style"], "translateX(-0%)") {
		t.Fatalf("should start on the first slide: %s", track.Attrs["style"])
	}
	if _, disabled := r.FindByAttr("data-slot", "carousel-previous").Attrs["disabled"]; !disabled {
		t.Fatal("previous should be disabled on the first slide")
	}
	r.Click(r.FindByAttr("data-slot", "carousel-next"))
	if got := r.FindByAttr("data-slot", "carousel-track").Attrs["style"]; !strings.Contains(got, "translateX(-100%)") {
		t.Fatalf("next should advance to slide 2: %s", got)
	}
	if r.FindByAttr("data-index", "1").Attrs["class"] == "" || !strings.Contains(r.FindByAttr("data-index", "1").Attrs["class"], "bg-primary") {
		t.Fatalf("the second dot should be active: %s", r.HTML())
	}
}

func TestSidebarToggles(t *testing.T) {
	r := testdom.Mount(ui.SidebarProvider(
		ui.Sidebar(ui.SidebarContent(ui.SidebarMenu(
			ui.SidebarMenuItem(ui.SidebarMenuButton(true, "Home")),
		))),
		ui.SidebarInset(ui.SidebarTrigger(), g.Div("content")),
	))
	sb := r.FindByAttr("data-slot", "sidebar")
	if sb == nil || sb.Attrs["data-state"] != "open" {
		t.Fatalf("sidebar should start open: %s", r.HTML())
	}
	r.Click(r.FindByAttr("data-slot", "sidebar-trigger"))
	if r.FindByAttr("data-slot", "sidebar").Attrs["data-state"] != "closed" {
		t.Fatal("the trigger should collapse the sidebar")
	}
	r.Click(r.FindByAttr("data-slot", "sidebar-trigger"))
	if r.FindByAttr("data-slot", "sidebar").Attrs["data-state"] != "open" {
		t.Fatal("the trigger should reopen the sidebar")
	}
}

var formValidated bool

func heavyFormApp() *g.Node {
	f := ui.UseForm(map[string]string{})
	return g.Div(
		ui.FormField(ui.FormFieldProps{Form: f, Name: "email", Label: "Email"}),
		g.Button(g.Data("slot", "submit"), g.OnClick(func(*g.Event) {
			formValidated = f.Validate(map[string]func(string) string{
				"email": func(v string) string {
					if v == "" {
						return "Required"
					}
					return ""
				},
			})
		}), "Submit"),
	)
}

func TestFormValidateBindClear(t *testing.T) {
	formValidated = false
	r := testdom.Mount(g.C0(heavyFormApp))

	r.Click(r.FindByAttr("data-slot", "submit"))
	if r.FindText("Required") == nil || formValidated {
		t.Fatalf("empty submit should fail validation: %s", r.HTML())
	}
	// typing clears the field error
	r.Input(r.FindByAttr("id", "email"), "ada@x.io")
	if r.FindText("Required") != nil {
		t.Fatal("editing should clear the error")
	}
	r.Click(r.FindByAttr("data-slot", "submit"))
	if !formValidated {
		t.Fatal("a filled field should pass validation")
	}
}

func TestResizableKeyboard(t *testing.T) {
	r := testdom.Mount(ui.Resizable(ui.ResizableProps{}, g.Div("L"), g.Div("R")))
	first := r.FindByAttr("data-slot", "resizable-panel")
	if !strings.Contains(first.Attrs["style"], "flex: 0.5000") {
		t.Fatalf("should start at an even split: %s", first.Attrs["style"])
	}
	r.KeyDown(r.FindByAttr("data-slot", "resizable-handle"), "ArrowRight")
	if got := r.FindByAttr("data-slot", "resizable-panel").Attrs["style"]; !strings.Contains(got, "flex: 0.5400") {
		t.Fatalf("ArrowRight should grow the first panel: %s", got)
	}
}
