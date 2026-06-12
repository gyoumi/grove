// Showcase of the ui package: every component on the default theme, with
// a dark-mode toggle and a working modal dialog.
// Run with: grove serve (from this directory).
package main

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/dom"
	"github.com/gyoumi/grove/ui"
)

func App() *g.Node {
	dark, setDark := g.UseState(false)
	open, setOpen := g.UseState(false)
	menuOpen, setMenuOpen := g.UseState(false)
	lastAction, setLastAction := g.UseState("none")
	name, setName := g.UseState("")
	agreed, setAgreed := g.UseState(false)

	g.UseEffect(func() func() {
		dom.SetRootClass("dark", dark)
		return nil
	}, []any{dark})

	darkLabel := "dark mode"
	if dark {
		darkLabel = "light mode"
	}

	return g.Div(g.Class("mx-auto flex min-h-svh max-w-2xl flex-col gap-8 p-8"),
		g.Header(g.Class("flex items-center justify-between"),
			g.Div(
				g.H1(g.Class("text-2xl font-semibold tracking-tight"), "grove/ui"),
				g.P(g.Class("text-sm text-muted-foreground"), "grove's component library"),
			),
			ui.Button(ui.ButtonProps{Variant: ui.ButtonOutline, OnClick: func(*g.Event) { setDark(!dark) }}, darkLabel),
		),

		section("Buttons",
			g.Div(g.Class("flex flex-wrap items-center gap-2"),
				ui.Button(ui.ButtonProps{}, "Default"),
				ui.Button(ui.ButtonProps{Variant: ui.ButtonSecondary}, "Secondary"),
				ui.Button(ui.ButtonProps{Variant: ui.ButtonDestructive}, "Destructive"),
				ui.Button(ui.ButtonProps{Variant: ui.ButtonOutline}, "Outline"),
				ui.Button(ui.ButtonProps{Variant: ui.ButtonGhost}, "Ghost"),
				ui.Button(ui.ButtonProps{Variant: ui.ButtonLink}, "Link"),
				ui.Button(ui.ButtonProps{Disabled: true}, "Disabled"),
				ui.Button(ui.ButtonProps{Size: ui.ButtonSizeSm, Variant: ui.ButtonSecondary}, "Small"),
				ui.Button(ui.ButtonProps{Size: ui.ButtonSizeLg}, "Large"),
			),
		),

		section("Badges",
			g.Div(g.Class("flex flex-wrap items-center gap-2"),
				ui.Badge(ui.BadgeProps{}, "Default"),
				ui.Badge(ui.BadgeProps{Variant: ui.BadgeSecondary}, "Secondary"),
				ui.Badge(ui.BadgeProps{Variant: ui.BadgeDestructive}, "Destructive"),
				ui.Badge(ui.BadgeProps{Variant: ui.BadgeOutline}, "Outline"),
			),
		),

		section("Card + form",
			ui.Card(
				ui.CardHeader(
					ui.CardTitle("Create account"),
					ui.CardDescription("Controlled inputs backed by UseState."),
				),
				ui.CardContent(g.Class("flex flex-col gap-4"),
					g.Div(g.Class("flex flex-col gap-2"),
						ui.Label("name", "Name"),
						ui.Input(ui.InputProps{
							ID:          "name",
							Value:       name,
							Placeholder: "Ada Lovelace",
							OnInput:     func(e *g.Event) { setName(e.Value()) },
						}),
					),
					g.Div(g.Class("flex items-center gap-2"),
						ui.Checkbox(ui.CheckboxProps{ID: "terms", Checked: agreed, OnChange: setAgreed}),
						ui.Label("terms", "I agree to the terms"),
					),
				),
				ui.CardFooter(g.Class("justify-between"),
					g.P(g.Class("text-sm text-muted-foreground"),
						g.If(name != "", g.Textf("hello, %s", name)),
					),
					ui.Button(ui.ButtonProps{Disabled: !agreed, OnClick: func(*g.Event) { setOpen(true) }}, "Sign up"),
				),
			),
		),

		section("Dropdown menu",
			g.Div(g.Class("flex items-center gap-4"),
				ui.Dropdown(ui.DropdownProps{Open: menuOpen, OnClose: func() { setMenuOpen(false) }},
					ui.Button(ui.ButtonProps{Variant: ui.ButtonOutline, OnClick: func(*g.Event) { setMenuOpen(!menuOpen) }}, "Actions"),
					ui.DropdownLabel("Demo actions"),
					ui.DropdownItem(ui.DropdownItemProps{OnSelect: func() { setLastAction("duplicate") }}, "Duplicate"),
					ui.DropdownItem(ui.DropdownItemProps{OnSelect: func() { setLastAction("archive") }}, "Archive"),
					ui.DropdownSeparator(),
					ui.DropdownItem(ui.DropdownItemProps{OnSelect: func() { setLastAction("delete") }, Class: "text-destructive"}, "Delete"),
				),
				g.P(g.Class("text-sm text-muted-foreground"), g.Textf("last action: %s", lastAction)),
			),
		),

		section("Alerts",
			g.Div(g.Class("flex flex-col gap-2"),
				ui.Alert(ui.AlertProps{},
					ui.AlertTitle("Heads up"),
					ui.AlertDescription("Everything here is rendered by Go compiled to WebAssembly."),
				),
				ui.Alert(ui.AlertProps{Variant: ui.AlertDestructive},
					ui.AlertTitle("Careful"),
					ui.AlertDescription("This variant uses the destructive theme colors."),
				),
			),
		),

		ui.Dialog(ui.DialogProps{Open: open, OnClose: func() { setOpen(false) }},
			ui.DialogHeader(
				ui.DialogTitle("Confirm sign up"),
				ui.DialogDescription(g.Textf("Create an account for %q?", name)),
			),
			ui.DialogFooter(
				ui.Button(ui.ButtonProps{Variant: ui.ButtonOutline, OnClick: func(*g.Event) { setOpen(false) }}, "Cancel"),
				ui.Button(ui.ButtonProps{OnClick: func(*g.Event) { setOpen(false) }}, "Confirm"),
			),
		),
	)
}

func section(title string, body *g.Node) *g.Node {
	return g.Section(g.Class("flex flex-col gap-3"),
		g.H2(g.Class("text-lg font-medium"), title),
		body,
		ui.Separator(false, g.Class("mt-2")),
	)
}

func main() {
	dom.Mount("#root", g.C0(App))
}
