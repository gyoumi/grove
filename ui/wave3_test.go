package ui_test

import (
	"strings"
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func selectApp() *g.Node {
	v, set := g.UseState("")
	return ui.Select(ui.SelectProps{
		Value: v, Placeholder: "Pick a fruit", OnChange: set,
		Options: []ui.SelectOption{{Value: "a", Label: "Apple"}, {Value: "b", Label: "Banana"}},
	})
}

func TestSelectOpensAndPicks(t *testing.T) {
	r := testdom.Mount(g.C0(selectApp))
	trig := r.FindByAttr("data-slot", "select-trigger")
	if trig == nil || trig.Attrs["data-empty"] != "1" {
		t.Fatalf("select should start empty: %s", r.HTML())
	}
	if r.FindByAttr("data-slot", "select-item") != nil {
		t.Fatal("options should be closed until the trigger is clicked")
	}
	r.Click(trig)
	if r.FindByAttr("data-slot", "select-item") == nil || r.FindText("Banana") == nil {
		t.Fatalf("options should open: %s", r.HTML())
	}
	r.Click(r.FindText("Apple"))
	if r.FindByAttr("data-slot", "select-item") != nil {
		t.Fatal("picking should close the listbox")
	}
	trig = r.FindByAttr("data-slot", "select-trigger")
	if _, empty := trig.Attrs["data-empty"]; empty || !strings.Contains(trig.TextContent(), "Apple") {
		t.Fatalf("trigger should show the choice: %s", trig.HTML())
	}
}

func alertApp() *g.Node {
	open, setOpen := g.UseState(true)
	return ui.AlertDialog(ui.AlertDialogProps{Open: open, OnClose: func() { setOpen(false) }},
		ui.AlertDialogHeader(ui.AlertDialogTitle("Delete?")),
		ui.AlertDialogFooter(
			ui.AlertDialogCancel(func(*g.Event) { setOpen(false) }, "Cancel"),
			ui.AlertDialogAction(func(*g.Event) { setOpen(false) }, "Confirm"),
		),
	)
}

func TestAlertDialogDoesNotDismissOnOverlay(t *testing.T) {
	r := testdom.Mount(g.C0(alertApp))
	if r.FindByAttr("data-slot", "alert-dialog-content") == nil {
		t.Fatalf("alert dialog should be open: %s", r.HTML())
	}
	// clicking the overlay must NOT close an alert dialog
	r.Click(r.FindByAttr("data-slot", "alert-dialog-overlay"))
	if r.FindByAttr("data-slot", "alert-dialog-content") == nil {
		t.Fatal("overlay click should not dismiss an alert dialog")
	}
	// the cancel action does close it
	r.Click(r.FindByAttr("data-slot", "alert-dialog-cancel"))
	if r.FindByAttr("data-slot", "alert-dialog-content") != nil {
		t.Fatalf("cancel should close the alert dialog: %s", r.HTML())
	}
}

var sheetSetOpen func(bool)

func sheetExitApp() *g.Node {
	open, set := g.UseState(true)
	sheetSetOpen = set
	return ui.Sheet(ui.SheetProps{Open: open, Side: ui.SheetRight, OnClose: func() { set(false) }},
		ui.SheetHeader(ui.SheetTitle("Filters")))
}

func TestSheetSlidesOutThenUnmounts(t *testing.T) {
	r := testdom.Mount(g.C0(sheetExitApp))
	if r.FindByAttr("data-slot", "sheet-content") == nil {
		t.Fatalf("sheet should be open: %s", r.HTML())
	}

	// Closing keeps the panel mounted, in a slide-out state, until the
	// animation ends.
	sheetSetOpen(false)
	r.Settle()
	content := r.FindByAttr("data-slot", "sheet-content")
	if content == nil {
		t.Fatal("sheet should stay mounted during the close-out animation")
	}
	if content.Attrs["data-state"] != "closing" || !strings.Contains(content.Attrs["class"], "animate-slide-out-right") {
		t.Fatalf("sheet should be sliding out: %s", content.HTML())
	}

	// The slide-out animation finishing unmounts it.
	r.Fire(content, "animationend", map[string]any{"animationName": "slide-out-right"})
	if r.FindByAttr("data-slot", "sheet-content") != nil {
		t.Fatalf("sheet should unmount once it has slid out: %s", r.HTML())
	}
}

func TestSheetSidePlacement(t *testing.T) {
	r := testdom.Mount(ui.Sheet(ui.SheetProps{Open: true, Side: ui.SheetLeft},
		ui.SheetHeader(ui.SheetTitle("Filters")),
	))
	content := r.FindByAttr("data-slot", "sheet-content")
	if content == nil || content.Attrs["data-side"] != "left" {
		t.Fatalf("sheet should record its side: %s", r.HTML())
	}
	if !strings.Contains(content.Attrs["class"], "left-0") {
		t.Fatalf("left sheet should be anchored left: %s", content.Attrs["class"])
	}
}

func menubarApp() *g.Node {
	picked := g.UseRef("")
	return ui.Menubar(ui.MenubarProps{},
		ui.MenubarMenu{Label: "File", Items: []any{
			ui.DropdownItem(ui.DropdownItemProps{OnSelect: func() { picked.Current = "new" }}, "New"),
		}},
	)
}

func TestMenubarOpensMenu(t *testing.T) {
	r := testdom.Mount(g.C0(menubarApp))
	if r.FindByAttr("data-slot", "dropdown-item") != nil {
		t.Fatal("menu should start closed")
	}
	r.Click(r.FindText("File"))
	if r.FindText("New") == nil {
		t.Fatalf("clicking File should open its menu: %s", r.HTML())
	}
	// selecting an item closes the menu
	r.Click(r.FindText("New"))
	if r.FindByAttr("data-slot", "dropdown-item") != nil {
		t.Fatal("selecting should close the menu")
	}
}
