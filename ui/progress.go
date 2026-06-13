package ui

import (
	"strconv"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// Progress renders a determinate progress bar. value is clamped to 0–100.
func Progress(value float64, class ...string) *g.Node {
	if value < 0 {
		value = 0
	}
	if value > 100 {
		value = 100
	}
	pct := strconv.FormatFloat(value, 'f', -1, 64)
	return g.Div(
		g.Class(style.CN("relative h-2 w-full overflow-hidden rounded-full bg-primary/20", []string(class))),
		g.Data("slot", "progress"),
		g.Role("progressbar"),
		g.Attr("aria-valuemin", "0"),
		g.Attr("aria-valuemax", "100"),
		g.Attr("aria-valuenow", pct),
		g.Div(
			g.Class("h-full rounded-full bg-primary transition-all"),
			g.Data("slot", "progress-bar"),
			g.Attr("style", "width: "+pct+"%"),
		),
	)
}
