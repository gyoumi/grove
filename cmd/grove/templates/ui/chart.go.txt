package ui

import (
	"fmt"
	"strings"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type ChartType string

const (
	ChartLine ChartType = "line"
	ChartArea ChartType = "area"
	ChartBar  ChartType = "bar"
)

// ChartSeries is one data series. Color is any CSS color (a theme var works);
// it defaults to a slot from the built-in palette.
type ChartSeries struct {
	Label  string
	Color  string
	Values []float64
}

type ChartProps struct {
	Type   ChartType
	Series []ChartSeries
	Labels []string // x-axis labels (optional)
	Class  string   // sizes the plot area (default h-48)
}

// chartPalette supplies default series colors.
var chartPalette = []string{
	"var(--color-primary)", "#10b981", "#f59e0b", "#ec4899", "#6366f1",
}

// Chart draws line, area, or bar charts as inline SVG over the design
// system's colors. The plot stretches to its container (size it via Class);
// strokes stay crisp at any size.
func Chart(p ChartProps) *g.Node {
	lo, hi := chartRange(p.Series)

	// horizontal grid lines
	body := make([]any, 0, 8)
	body = append(body,
		g.Attr("viewBox", "0 0 100 100"),
		g.Attr("preserveAspectRatio", "none"),
		g.Class("size-full overflow-visible"),
	)
	for i := 1; i < 4; i++ {
		y := float64(i) / 4 * 100
		body = append(body, svgLine(0, y, 100, y, "stroke-border/60", "1"))
	}

	for si, s := range p.Series {
		color := s.Color
		if color == "" {
			color = chartPalette[si%len(chartPalette)]
		}
		switch p.Type {
		case ChartBar:
			body = append(body, barRects(s.Values, si, len(p.Series), lo, hi, color)...)
		case ChartArea:
			body = append(body,
				g.El("path", g.Attr("d", areaPath(s.Values, lo, hi)), g.Attr("fill", color),
					g.Attr("fill-opacity", "0.15"), g.Data("slot", "chart-series")),
				lineShape(s.Values, lo, hi, color),
			)
		default:
			body = append(body, lineShape(s.Values, lo, hi, color))
		}
	}

	plot := g.Div(
		g.Class(style.CN("h-48 w-full", p.Class)),
		g.Data("slot", "chart"),
		g.El("svg", body...),
	)
	if len(p.Labels) == 0 {
		return plot
	}

	labels := make([]any, 0, len(p.Labels)+1)
	labels = append(labels, g.Class("flex justify-between px-1 text-xs text-muted-foreground"))
	for _, l := range p.Labels {
		labels = append(labels, g.Span(l))
	}
	return g.Div(g.Class("flex w-full flex-col gap-1.5"), plot, g.Div(labels...))
}

func lineShape(values []float64, lo, hi float64, color string) *g.Node {
	return g.El("path",
		g.Attr("d", linePath(values, lo, hi)),
		g.Attr("fill", "none"),
		g.Attr("stroke", color),
		g.Attr("stroke-width", "2"),
		g.Attr("stroke-linecap", "round"),
		g.Attr("stroke-linejoin", "round"),
		g.Attr("vector-effect", "non-scaling-stroke"),
		g.Data("slot", "chart-series"),
	)
}

func svgLine(x1, y1, x2, y2 float64, class, width string) *g.Node {
	return g.El("line",
		g.Attr("x1", f2(x1)), g.Attr("y1", f2(y1)), g.Attr("x2", f2(x2)), g.Attr("y2", f2(y2)),
		g.Class(class), g.Attr("stroke", "currentColor"), g.Attr("stroke-width", width),
		g.Attr("vector-effect", "non-scaling-stroke"),
	)
}

func barRects(values []float64, si, nSeries int, lo, hi float64, color string) []any {
	if len(values) == 0 {
		return nil
	}
	slot := 100 / float64(len(values))
	groupW := slot * 0.7
	barW := groupW / float64(nSeries)
	out := make([]any, 0, len(values))
	for i, v := range values {
		y := pointY(v, lo, hi)
		x := float64(i)*slot + (slot-groupW)/2 + float64(si)*barW
		out = append(out, g.El("rect",
			g.Attr("x", f2(x)), g.Attr("y", f2(y)),
			g.Attr("width", f2(barW)), g.Attr("height", f2(94-y)),
			g.Attr("rx", "1"), g.Attr("fill", color),
			g.Data("slot", "chart-series"),
		))
	}
	return out
}

func chartRange(series []ChartSeries) (lo, hi float64) {
	lo, hi = 0, 0
	first := true
	for _, s := range series {
		for _, v := range s.Values {
			if first {
				lo, hi, first = v, v, false
				continue
			}
			lo = min(lo, v)
			hi = max(hi, v)
		}
	}
	if lo > 0 {
		lo = 0 // anchor the baseline at zero when all values are positive
	}
	return lo, hi
}

const chartPad = 6.0

func pointY(v, lo, hi float64) float64 {
	if hi == lo {
		return 50
	}
	return chartPad + (1-(v-lo)/(hi-lo))*(100-2*chartPad)
}

func linePath(values []float64, lo, hi float64) string {
	var b strings.Builder
	for i, v := range values {
		x := 0.0
		if len(values) > 1 {
			x = float64(i) / float64(len(values)-1) * 100
		}
		if i == 0 {
			fmt.Fprintf(&b, "M %s %s", f2(x), f2(pointY(v, lo, hi)))
		} else {
			fmt.Fprintf(&b, " L %s %s", f2(x), f2(pointY(v, lo, hi)))
		}
	}
	return b.String()
}

func areaPath(values []float64, lo, hi float64) string {
	if len(values) == 0 {
		return ""
	}
	return linePath(values, lo, hi) + " L 100 100 L 0 100 Z"
}

func f2(v float64) string { return fmt.Sprintf("%.2f", v) }
