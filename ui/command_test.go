package ui_test

import (
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func commandApp(onPick func(string)) func() *g.Node {
	return func() *g.Node {
		return ui.Command(ui.CommandProps{},
			ui.CommandGroup{Heading: "Fruits", Items: []ui.CommandItem{
				{Value: "Apple", OnSelect: func() { onPick("apple") }},
				{Value: "Banana", OnSelect: func() { onPick("banana") }},
				{Value: "Cherry", OnSelect: func() { onPick("cherry") }},
			}},
		)
	}
}

func TestCommandFiltersClicksAndKeyboard(t *testing.T) {
	var picked string
	r := testdom.Mount(g.C0(commandApp(func(s string) { picked = s })))

	// all three show, the first is active
	if r.FindByAttr("data-value", "Apple") == nil || r.FindByAttr("data-value", "Cherry") == nil {
		t.Fatalf("all items should show initially: %s", r.HTML())
	}
	if r.FindByAttr("data-value", "Apple").Attrs["data-active"] != "1" {
		t.Fatal("first item should start active")
	}

	// typing filters by substring (case-insensitive)
	in := r.FindByAttr("data-slot", "command-input")
	r.Input(in, "an")
	if r.FindByAttr("data-value", "Banana") == nil {
		t.Fatalf("Banana should match 'an': %s", r.HTML())
	}
	if r.FindByAttr("data-value", "Apple") != nil || r.FindByAttr("data-value", "Cherry") != nil {
		t.Fatal("non-matching items should be hidden")
	}

	// no matches → empty state
	r.Input(in, "zzz")
	if r.FindByAttr("data-slot", "command-empty") == nil {
		t.Fatalf("no matches should show the empty state: %s", r.HTML())
	}

	// keyboard: clear filter, arrow down to Banana, Enter selects it
	r.Input(in, "")
	r.KeyDown(in, "ArrowDown")
	if r.FindByAttr("data-value", "Banana").Attrs["data-active"] != "1" {
		t.Fatalf("ArrowDown should move active to Banana: %s", r.HTML())
	}
	r.KeyDown(in, "Enter")
	if picked != "banana" {
		t.Fatalf("Enter should select the active item, got %q", picked)
	}

	// hovering moves the active highlight to that row
	r.Fire(r.FindByAttr("data-value", "Apple"), "mouseover", nil)
	if r.FindByAttr("data-value", "Apple").Attrs["data-active"] != "1" {
		t.Fatalf("hovering should activate the row: %s", r.HTML())
	}
	if r.FindByAttr("data-value", "Banana").Attrs["data-active"] == "1" {
		t.Fatal("only one row should be active at a time")
	}

	// clicking selects directly
	picked = ""
	r.Click(r.FindByAttr("data-value", "Cherry"))
	if picked != "cherry" {
		t.Fatalf("click should select, got %q", picked)
	}
}

func comboboxApp() *g.Node {
	v, set := g.UseState("")
	return ui.Combobox(ui.ComboboxProps{
		Value: v, OnChange: set, Placeholder: "Pick a fruit",
		Options: []ui.SelectOption{{Value: "a", Label: "Apple"}, {Value: "b", Label: "Banana"}},
	})
}

func TestComboboxOpensFiltersSelects(t *testing.T) {
	r := testdom.Mount(g.C0(comboboxApp))
	trig := r.FindByAttr("data-slot", "combobox-trigger")
	if trig == nil || trig.Attrs["data-empty"] != "1" {
		t.Fatalf("combobox should start empty: %s", r.HTML())
	}
	if r.FindByAttr("data-slot", "command") != nil {
		t.Fatal("the command list should be closed until opened")
	}

	r.Click(trig)
	if r.FindByAttr("data-slot", "command") == nil {
		t.Fatalf("clicking the trigger should open the command list: %s", r.HTML())
	}
	r.Input(r.FindByAttr("data-slot", "command-input"), "ban")
	if r.FindByAttr("data-value", "Banana") == nil || r.FindByAttr("data-value", "Apple") != nil {
		t.Fatalf("filter should narrow to Banana: %s", r.HTML())
	}
	r.Click(r.FindByAttr("data-value", "Banana"))
	if r.FindByAttr("data-slot", "command") != nil {
		t.Fatal("picking should close the popover")
	}
	if got := r.FindByAttr("data-slot", "combobox-trigger").TextContent(); !contains(got, "Banana") {
		t.Fatalf("trigger should show the choice, got %q", got)
	}
}
