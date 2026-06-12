package ui_test

import (
	"strings"
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func popoverApp() *g.Node {
	open, setOpen := g.UseState(false)
	return ui.Popover(
		ui.PopoverProps{Open: open, OnClose: func() { setOpen(false) }, Side: ui.PopoverTop, Align: ui.PopoverAlignEnd},
		ui.Button(ui.ButtonProps{OnClick: func(*g.Event) { setOpen(!open) }}, "toggle"),
		g.P("popover body"),
	)
}

func TestPopoverOpenCloseAndPosition(t *testing.T) {
	r := testdom.Mount(g.C0(popoverApp))
	if r.FindByAttr("data-slot", "popover-content") != nil {
		t.Fatal("popover should start closed")
	}

	r.Click(r.Find("button"))
	content := r.FindByAttr("data-slot", "popover-content")
	if content == nil {
		t.Fatalf("popover should open: %s", r.HTML())
	}
	cls := content.Attrs["class"]
	if !strings.Contains(cls, "bottom-full") || !strings.Contains(cls, "right-0") {
		t.Fatalf("side=top align=end classes missing: %s", cls)
	}

	r.KeyDown(content, "Escape")
	if r.FindByAttr("data-slot", "popover-content") != nil {
		t.Fatal("Escape should close the popover")
	}

	r.Click(r.Find("button"))
	r.Click(r.FindByAttr("data-slot", "popover-overlay"))
	if r.FindByAttr("data-slot", "popover-content") != nil {
		t.Fatal("outside click should close the popover")
	}
}

var picked []string

func dropdownApp() *g.Node {
	open, setOpen := g.UseState(false)
	return ui.Dropdown(
		ui.DropdownProps{Open: open, OnClose: func() { setOpen(false) }},
		ui.Button(ui.ButtonProps{OnClick: func(*g.Event) { setOpen(!open) }}, "actions"),
		ui.DropdownLabel("Event"),
		ui.DropdownItem(ui.DropdownItemProps{OnSelect: func() { picked = append(picked, "rename") }}, "Rename"),
		ui.DropdownItem(ui.DropdownItemProps{Disabled: true, OnSelect: func() { picked = append(picked, "nope") }}, "Disabled"),
		ui.DropdownSeparator(),
		ui.DropdownItem(ui.DropdownItemProps{OnSelect: func() { picked = append(picked, "delete") }, Class: "text-destructive"}, "Delete"),
	)
}

func TestDropdownSelectClosesMenu(t *testing.T) {
	picked = nil
	r := testdom.Mount(g.C0(dropdownApp))
	r.Click(r.Find("button")) // trigger
	if r.FindByAttr("data-slot", "dropdown-item") == nil {
		t.Fatalf("menu should open: %s", r.HTML())
	}
	if r.FindByAttr("data-slot", "dropdown-label") == nil || r.FindByAttr("data-slot", "dropdown-separator") == nil {
		t.Fatalf("label/separator missing: %s", r.HTML())
	}

	r.Click(r.FindText("Rename"))
	if len(picked) != 1 || picked[0] != "rename" {
		t.Fatalf("OnSelect not delivered: %v", picked)
	}
	if r.FindByAttr("data-slot", "dropdown-item") != nil {
		t.Fatal("selecting an item should close the menu")
	}

	// disabled items do nothing and keep the menu open
	r.Click(r.FindText("actions"))
	r.Click(r.FindText("Disabled"))
	if len(picked) != 1 {
		t.Fatalf("disabled item should not select: %v", picked)
	}
	if r.FindByAttr("data-slot", "dropdown-item") == nil {
		t.Fatal("disabled item should not close the menu")
	}

	// Escape closes via the menu's own keydown handler
	r.KeyDown(r.FindByAttr("data-slot", "popover-content"), "Escape")
	if r.FindByAttr("data-slot", "dropdown-item") != nil {
		t.Fatal("Escape should close the menu")
	}
}
