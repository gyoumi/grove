package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// iconShape is one drawing primitive of an icon (a path, circle, line, …).
type iconShape struct {
	tag   string
	attrs [][2]string
}

// icons defines each icon as stroke primitives on a 24×24 grid. They are
// stored as data — not prebuilt nodes — so every Icon call gets fresh node
// identities for the reconciler.
var icons = map[string][]iconShape{
	"chevron-down":  {{"path", [][2]string{{"d", "m6 9 6 6 6-6"}}}},
	"chevron-up":    {{"path", [][2]string{{"d", "m18 15-6-6-6 6"}}}},
	"chevron-left":  {{"path", [][2]string{{"d", "m15 18-6-6 6-6"}}}},
	"chevron-right": {{"path", [][2]string{{"d", "m9 18 6-6-6-6"}}}},
	"chevrons-up-down": {
		{"path", [][2]string{{"d", "m7 15 5 5 5-5"}}},
		{"path", [][2]string{{"d", "m7 9 5-5 5 5"}}},
	},
	"check":       {{"path", [][2]string{{"d", "M20 6 9 17l-5-5"}}}},
	"check-small": {{"path", [][2]string{{"d", "M20 6 9 17l-5-5"}}}},
	"minus":       {{"path", [][2]string{{"d", "M5 12h14"}}}},
	"plus":        {{"path", [][2]string{{"d", "M5 12h14"}}}, {"path", [][2]string{{"d", "M12 5v14"}}}},
	"x": {
		{"path", [][2]string{{"d", "M18 6 6 18"}}},
		{"path", [][2]string{{"d", "m6 6 12 12"}}},
	},
	"search": {
		{"circle", [][2]string{{"cx", "11"}, {"cy", "11"}, {"r", "8"}}},
		{"path", [][2]string{{"d", "m21 21-4.3-4.3"}}},
	},
	"calendar": {
		{"path", [][2]string{{"d", "M8 2v4"}}},
		{"path", [][2]string{{"d", "M16 2v4"}}},
		{"rect", [][2]string{{"width", "18"}, {"height", "18"}, {"x", "3"}, {"y", "4"}, {"rx", "2"}}},
		{"path", [][2]string{{"d", "M3 10h18"}}},
	},
	"clock": {
		{"circle", [][2]string{{"cx", "12"}, {"cy", "12"}, {"r", "10"}}},
		{"path", [][2]string{{"d", "M12 6v6l4 2"}}},
	},
	"arrow-left":  {{"path", [][2]string{{"d", "m12 19-7-7 7-7"}}}, {"path", [][2]string{{"d", "M19 12H5"}}}},
	"arrow-right": {{"path", [][2]string{{"d", "M5 12h14"}}}, {"path", [][2]string{{"d", "m12 5 7 7-7 7"}}}},
	"more-horizontal": {
		{"circle", [][2]string{{"cx", "12"}, {"cy", "12"}, {"r", "1"}}},
		{"circle", [][2]string{{"cx", "19"}, {"cy", "12"}, {"r", "1"}}},
		{"circle", [][2]string{{"cx", "5"}, {"cy", "12"}, {"r", "1"}}},
	},
	"more-vertical": {
		{"circle", [][2]string{{"cx", "12"}, {"cy", "12"}, {"r", "1"}}},
		{"circle", [][2]string{{"cx", "12"}, {"cy", "5"}, {"r", "1"}}},
		{"circle", [][2]string{{"cx", "12"}, {"cy", "19"}, {"r", "1"}}},
	},
	"circle":       {{"circle", [][2]string{{"cx", "12"}, {"cy", "12"}, {"r", "10"}}}},
	"circle-check": {{"circle", [][2]string{{"cx", "12"}, {"cy", "12"}, {"r", "10"}}}, {"path", [][2]string{{"d", "m9 12 2 2 4-4"}}}},
	"info": {
		{"circle", [][2]string{{"cx", "12"}, {"cy", "12"}, {"r", "10"}}},
		{"path", [][2]string{{"d", "M12 16v-4"}}},
		{"path", [][2]string{{"d", "M12 8h.01"}}},
	},
	"triangle-alert": {
		{"path", [][2]string{{"d", "m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"}}},
		{"path", [][2]string{{"d", "M12 9v4"}}},
		{"path", [][2]string{{"d", "M12 17h.01"}}},
	},
	"panel-left": {
		{"rect", [][2]string{{"width", "18"}, {"height", "18"}, {"x", "3"}, {"y", "3"}, {"rx", "2"}}},
		{"path", [][2]string{{"d", "M9 3v18"}}},
	},
	"grip-vertical": {
		{"circle", [][2]string{{"cx", "9"}, {"cy", "12"}, {"r", "1"}}},
		{"circle", [][2]string{{"cx", "9"}, {"cy", "5"}, {"r", "1"}}},
		{"circle", [][2]string{{"cx", "9"}, {"cy", "19"}, {"r", "1"}}},
		{"circle", [][2]string{{"cx", "15"}, {"cy", "12"}, {"r", "1"}}},
		{"circle", [][2]string{{"cx", "15"}, {"cy", "5"}, {"r", "1"}}},
		{"circle", [][2]string{{"cx", "15"}, {"cy", "19"}, {"r", "1"}}},
	},
	"loader": {{"path", [][2]string{{"d", "M21 12a9 9 0 1 1-6.219-8.56"}}}},
}

// Icon renders a stroke icon as inline SVG that inherits the current text
// color and is sized by classes (default size-4). Extra classes resize or
// recolor it: ui.Icon("check", "size-5 text-emerald-500"). Unknown names
// render an empty (still valid) svg.
func Icon(name string, class ...string) *g.Node {
	args := []any{
		g.Attr("viewBox", "0 0 24 24"),
		g.Attr("fill", "none"),
		g.Attr("stroke", "currentColor"),
		g.Attr("stroke-width", "2"),
		g.Attr("stroke-linecap", "round"),
		g.Attr("stroke-linejoin", "round"),
		g.Class(style.CN("inline-block size-4 shrink-0", []string(class))),
		g.Data("slot", "icon"),
		g.Data("icon", name),
		g.Attr("aria-hidden", "true"),
	}
	for _, sh := range icons[name] {
		shapeArgs := make([]any, 0, len(sh.attrs))
		for _, a := range sh.attrs {
			shapeArgs = append(shapeArgs, g.Attr(a[0], a[1]))
		}
		args = append(args, g.El(sh.tag, shapeArgs...))
	}
	return g.El("svg", args...)
}
