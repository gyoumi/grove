package ui_test

import (
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func TestAccordionSingleOpen(t *testing.T) {
	r := testdom.Mount(ui.Accordion(ui.AccordionProps{},
		ui.AccordionItem{Value: "a", Title: "Section A", Content: []any{g.P("body a")}},
		ui.AccordionItem{Value: "b", Title: "Section B", Content: []any{g.P("body b")}},
	))
	if r.FindByAttr("data-slot", "accordion-content") != nil {
		t.Fatalf("should start collapsed: %s", r.HTML())
	}
	r.Click(r.FindText("Section A"))
	if r.FindText("body a") == nil {
		t.Fatalf("clicking A should reveal its content: %s", r.HTML())
	}
	// single mode: opening B closes A
	r.Click(r.FindText("Section B"))
	if r.FindText("body b") == nil || r.FindText("body a") != nil {
		t.Fatalf("single accordion should keep one section open: %s", r.HTML())
	}
}

func TestTabsSwitch(t *testing.T) {
	r := testdom.Mount(ui.Tabs(ui.TabsProps{},
		ui.Tab{Value: "1", Label: "One", Content: []any{g.P("first")}},
		ui.Tab{Value: "2", Label: "Two", Content: []any{g.P("second")}},
	))
	if r.FindText("first") == nil || r.FindText("second") != nil {
		t.Fatalf("first tab should be active by default: %s", r.HTML())
	}
	r.Click(r.FindText("Two"))
	if r.FindText("second") == nil || r.FindText("first") != nil {
		t.Fatalf("clicking Two should switch panels: %s", r.HTML())
	}
}

func radioApp() *g.Node {
	v, set := g.UseState("a")
	return ui.RadioGroup(ui.RadioGroupProps{Value: v, OnChange: set},
		ui.RadioItem{Value: "a", Label: "A"},
		ui.RadioItem{Value: "b", Label: "B"},
	)
}

func TestRadioGroupSelects(t *testing.T) {
	r := testdom.Mount(g.C0(radioApp))
	if r.FindByAttr("data-value", "a").Attrs["aria-checked"] != "true" {
		t.Fatalf("a should start selected: %s", r.HTML())
	}
	r.Click(r.FindByAttr("data-value", "b"))
	if r.FindByAttr("data-value", "b").Attrs["aria-checked"] != "true" {
		t.Fatalf("b should become selected: %s", r.HTML())
	}
	if r.FindByAttr("data-value", "a").Attrs["aria-checked"] == "true" {
		t.Fatal("a should no longer be selected")
	}
}

func sliderApp() *g.Node {
	v, set := g.UseState(20.0)
	return ui.Slider(ui.SliderProps{Value: v, OnChange: set})
}

func TestSliderControlled(t *testing.T) {
	r := testdom.Mount(g.C0(sliderApp))
	sl := r.FindByAttr("data-slot", "slider")
	if sl.Props["value"] != "20" {
		t.Fatalf("slider should start at 20: %v", sl.Props["value"])
	}
	r.Input(sl, "55")
	if got := r.FindByAttr("data-slot", "slider").Props["value"]; got != "55" {
		t.Fatalf("slider should update to 55, got %v", got)
	}
}

func otpApp() *g.Node {
	v, set := g.UseState("")
	return ui.InputOTP(ui.InputOTPProps{Value: v, Length: 4, OnChange: set})
}

func TestInputOTPFiltersAndFills(t *testing.T) {
	r := testdom.Mount(g.C0(otpApp))
	if n := len(r.FindAll("div")); n == 0 {
		t.Fatal("otp should render slot divs")
	}
	in := r.FindByAttr("data-slot", "input-otp-input")
	r.Input(in, "12ab345") // digits only, capped at length 4 → "1234"
	if got := r.FindByAttr("data-index", "0").TextContent(); got != "1" {
		t.Fatalf("first slot should show 1, got %q: %s", got, r.HTML())
	}
	if got := r.FindByAttr("data-index", "3").TextContent(); got != "4" {
		t.Fatalf("fourth slot should show 4, got %q", got)
	}
}

func TestCollapsibleToggles(t *testing.T) {
	r := testdom.Mount(ui.Collapsible(ui.CollapsibleProps{}, g.Span("Toggle"), g.P("hidden body")))
	if r.FindText("hidden body") != nil {
		t.Fatalf("collapsible should start closed: %s", r.HTML())
	}
	r.Click(r.FindByAttr("data-slot", "collapsible-trigger"))
	if r.FindText("hidden body") == nil {
		t.Fatalf("clicking the trigger should reveal content: %s", r.HTML())
	}
}
