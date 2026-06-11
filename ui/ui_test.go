package ui_test

import (
	"strings"
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func TestButtonVariants(t *testing.T) {
	r := testdom.Mount(ui.Button(ui.ButtonProps{Variant: ui.ButtonDestructive, Size: ui.ButtonSizeSm}, "Delete"))
	btn := r.Find("button")
	if btn == nil {
		t.Fatalf("no button rendered: %s", r.HTML())
	}
	cls := btn.Attrs["class"]
	for _, want := range []string{"bg-destructive", "h-8", "inline-flex"} {
		if !strings.Contains(cls, want) {
			t.Errorf("button class missing %q: %s", want, cls)
		}
	}
	if btn.Attrs["type"] != "button" {
		t.Errorf("default type should be button, got %q", btn.Attrs["type"])
	}
	if btn.TextContent() != "Delete" {
		t.Errorf("children not rendered: %s", btn.HTML())
	}
}

func TestButtonClassOverride(t *testing.T) {
	r := testdom.Mount(ui.Button(ui.ButtonProps{Class: "bg-accent w-full"}, "Go"))
	cls := r.Find("button").Attrs["class"]
	tokens := map[string]bool{}
	for _, tok := range strings.Fields(cls) {
		tokens[tok] = true
	}
	if tokens["bg-primary"] {
		t.Errorf("caller bg- class should override the variant's: %s", cls)
	}
	if !tokens["bg-accent"] || !tokens["w-full"] {
		t.Errorf("caller classes missing: %s", cls)
	}
	if !tokens["hover:bg-primary/90"] {
		t.Errorf("hover-scoped class should be untouched by a plain bg- override: %s", cls)
	}
}

func TestCardComposition(t *testing.T) {
	r := testdom.Mount(ui.Card(
		ui.CardHeader(ui.CardTitle("T"), ui.CardDescription("D")),
		ui.CardContent(g.P("body")),
		ui.CardFooter(g.Span("f")),
	))
	for _, slot := range []string{"card", "card-header", "card-title", "card-description", "card-content", "card-footer"} {
		if r.FindByAttr("data-slot", slot) == nil {
			t.Errorf("missing slot %s: %s", slot, r.HTML())
		}
	}
}

func TestControlledInputAndCheckbox(t *testing.T) {
	r := testdom.Mount(g.C0(formApp))

	in := r.FindByAttr("data-slot", "input")
	r.Input(in, "hello")
	if got := r.FindByAttr("data-slot", "echo").TextContent(); got != "hello" {
		t.Fatalf("input not controlled: %q", got)
	}

	cb := r.FindByAttr("data-slot", "checkbox")
	r.SetChecked(cb, true)
	if got := r.FindByAttr("data-slot", "checked").TextContent(); got != "true" {
		t.Fatalf("checkbox change not delivered: %q", got)
	}
}

func formApp() *g.Node {
	text, setText := g.UseState("")
	checked, setChecked := g.UseState(false)
	return g.Div(
		ui.Label("name", "Name"),
		ui.Input(ui.InputProps{ID: "name", Value: text, OnInput: func(e *g.Event) { setText(e.Value()) }}),
		ui.Checkbox(ui.CheckboxProps{Checked: checked, OnChange: setChecked}),
		g.Span(g.Data("slot", "echo"), text),
		g.Span(g.Data("slot", "checked"), g.Textf("%v", checked)),
	)
}

func dialogApp() *g.Node {
	open, setOpen := g.UseState(false)
	return g.Div(
		ui.Button(ui.ButtonProps{OnClick: func(*g.Event) { setOpen(true) }}, "Open"),
		ui.Dialog(ui.DialogProps{Open: open, OnClose: func() { setOpen(false) }},
			ui.DialogHeader(ui.DialogTitle("Hi")),
			ui.DialogFooter(ui.Button(ui.ButtonProps{}, "OK")),
		),
	)
}

func TestDialogOpenEscapeOverlay(t *testing.T) {
	r := testdom.Mount(g.C0(dialogApp))
	if r.FindByAttr("data-slot", "dialog-content") != nil {
		t.Fatal("dialog should start closed")
	}

	r.Click(r.FindByAttr("data-slot", "button"))
	content := r.FindByAttr("data-slot", "dialog-content")
	if content == nil {
		t.Fatalf("dialog should open: %s", r.HTML())
	}
	if content.Attrs["role"] != "dialog" || content.Attrs["aria-modal"] != "true" {
		t.Errorf("dialog aria attributes missing: %v", content.Attrs)
	}

	r.KeyDown(content, "Escape")
	if r.FindByAttr("data-slot", "dialog-content") != nil {
		t.Fatal("Escape should close the dialog")
	}

	r.Click(r.FindByAttr("data-slot", "button"))
	r.Click(r.FindByAttr("data-slot", "dialog-overlay"))
	if r.FindByAttr("data-slot", "dialog-content") != nil {
		t.Fatal("overlay click should close the dialog")
	}
}

func TestAlertAndBadgeAndSeparator(t *testing.T) {
	r := testdom.Mount(g.Div(
		ui.Alert(ui.AlertProps{Variant: ui.AlertDestructive},
			ui.AlertTitle("Error"), ui.AlertDescription("Something broke")),
		ui.Badge(ui.BadgeProps{Variant: ui.BadgeSecondary}, "New"),
		ui.Separator(false),
	))
	alert := r.FindByAttr("data-slot", "alert")
	if alert == nil || alert.Attrs["role"] != "alert" {
		t.Fatalf("alert missing: %s", r.HTML())
	}
	if !strings.Contains(alert.Attrs["class"], "text-destructive") {
		t.Errorf("alert variant classes missing: %s", alert.Attrs["class"])
	}
	if r.FindByAttr("data-slot", "badge") == nil || r.FindByAttr("data-slot", "separator") == nil {
		t.Fatalf("badge/separator missing: %s", r.HTML())
	}
}
