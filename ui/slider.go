package ui

import (
	"strconv"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

const sliderClass = "h-2 w-full cursor-pointer appearance-none rounded-full bg-primary/20 accent-primary disabled:cursor-not-allowed disabled:opacity-50"

type SliderProps struct {
	Value    float64
	Min      float64
	Max      float64 // defaults to 100 when Max==Min==0
	Step     float64 // defaults to 1
	Disabled bool
	OnChange func(float64)
	Class    string
}

// Slider is a controlled value slider rendered as a themed native range
// input, so dragging and keyboard control work everywhere without bespoke
// pointer math. OnChange fires with the new value as it moves.
func Slider(p SliderProps) *g.Node {
	min, max := p.Min, p.Max
	if min == 0 && max == 0 {
		max = 100
	}
	step := p.Step
	if step == 0 {
		step = 1
	}
	args := []any{
		g.Type("range"),
		g.Class(style.CN(sliderClass, p.Class)),
		g.Data("slot", "slider"),
		g.Attr("min", ftoa(min)),
		g.Attr("max", ftoa(max)),
		g.Attr("step", ftoa(step)),
		g.Value(ftoa(p.Value)),
		g.Disabled(p.Disabled),
		g.Attr("aria-valuenow", ftoa(p.Value)),
		g.OnInput(func(e *g.Event) {
			if p.OnChange != nil {
				if v, err := strconv.ParseFloat(e.Value(), 64); err == nil {
					p.OnChange(v)
				}
			}
		}),
	}
	return g.Input(args...)
}

func ftoa(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) }
