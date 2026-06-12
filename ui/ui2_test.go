package ui_test

import (
	"strings"
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func TestAvatar(t *testing.T) {
	r := testdom.Mount(g.Div(
		ui.Avatar(ui.AvatarProps{Name: "Ada Lovelace"}),
		ui.Avatar(ui.AvatarProps{Name: "Ada Lovelace"}),
		ui.Avatar(ui.AvatarProps{Name: "Lin"}),
		ui.Avatar(ui.AvatarProps{Name: ""}),
	))
	avatars := r.FindAll("span")
	var slots []*testdom.Elem
	for _, s := range avatars {
		if s.Attrs["data-slot"] == "avatar" {
			slots = append(slots, s)
		}
	}
	if len(slots) != 4 {
		t.Fatalf("want 4 avatars: %s", r.HTML())
	}
	if got := slots[0].TextContent(); got != "AL" {
		t.Fatalf("initials = %q, want AL", got)
	}
	if slots[0].Attrs["class"] != slots[1].Attrs["class"] {
		t.Fatal("same name must give the same color")
	}
	if !strings.Contains(slots[0].Attrs["class"], "bg-") {
		t.Fatalf("no palette class: %s", slots[0].Attrs["class"])
	}
	if got := slots[3].TextContent(); got != "?" {
		t.Fatalf("empty name initials = %q", got)
	}
}

func switchApp() *g.Node {
	on, setOn := g.UseState(false)
	return g.Div(
		ui.Switch(ui.SwitchProps{ID: "sw", Checked: on, OnChange: setOn}),
		g.Span(g.Data("slot", "state"), g.Textf("%v", on)),
	)
}

func TestSwitch(t *testing.T) {
	r := testdom.Mount(g.C0(switchApp))
	sw := r.FindByAttr("data-slot", "switch")
	if sw.Attrs["role"] != "switch" || sw.Attrs["aria-checked"] != "false" {
		t.Fatalf("switch aria wrong: %v", sw.Attrs)
	}
	r.Click(sw)
	sw = r.FindByAttr("data-slot", "switch")
	if sw.Attrs["aria-checked"] != "true" || sw.Attrs["data-state"] != "checked" {
		t.Fatalf("switch did not toggle: %v", sw.Attrs)
	}
	if r.FindByAttr("data-slot", "state").TextContent() != "true" {
		t.Fatal("OnChange not delivered")
	}
	r.Click(sw)
	if r.FindByAttr("data-slot", "state").TextContent() != "false" {
		t.Fatal("switch should toggle back off")
	}
}

func TestTooltip(t *testing.T) {
	r := testdom.Mount(ui.Tooltip(ui.TooltipProps{Label: "delete this"},
		ui.Button(ui.ButtonProps{}, "x"),
	))
	bubble := r.FindByAttr("data-slot", "tooltip-bubble")
	if bubble == nil || bubble.Attrs["role"] != "tooltip" {
		t.Fatalf("tooltip bubble missing: %s", r.HTML())
	}
	if bubble.TextContent() != "delete this" {
		t.Fatalf("bubble text = %q", bubble.TextContent())
	}
	if !strings.Contains(bubble.Attrs["class"], "group-hover/tip:opacity-100") {
		t.Fatalf("hover reveal class missing: %s", bubble.Attrs["class"])
	}
	if r.Find("button") == nil {
		t.Fatal("trigger not rendered")
	}
}
