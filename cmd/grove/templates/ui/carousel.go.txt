package ui

import (
	"fmt"
	"strconv"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type CarouselProps struct {
	Dots  bool // show navigation dots
	Class string
}

// Carousel shows one slide at a time with previous/next controls (and
// optional dots). It tracks the current slide itself.
//
//	ui.Carousel(ui.CarouselProps{Dots: true}, slideA, slideB, slideC)
func Carousel(p CarouselProps, slides ...*g.Node) *g.Node {
	return g.C(carouselView, carouselArgs{p: p, slides: slides})
}

type carouselArgs struct {
	p      CarouselProps
	slides []*g.Node
}

func carouselView(a carouselArgs) *g.Node {
	index, setIndex := g.UseState(0)
	n := len(a.slides)
	if n == 0 {
		return g.Div(g.Data("slot", "carousel"))
	}
	if index > n-1 {
		index = n - 1
	}

	items := make([]any, 0, n)
	for _, s := range a.slides {
		items = append(items, g.Div(
			g.Class("min-w-0 shrink-0 grow-0 basis-full"),
			g.Data("slot", "carousel-item"),
			s,
		))
	}
	track := g.Div(
		g.Class("flex transition-transform duration-300 ease-out"),
		g.Data("slot", "carousel-track"),
		g.Attr("style", fmt.Sprintf("transform: translateX(-%d%%)", index*100)),
		items,
	)

	arrow := func(slot, icon string, disabled bool, onClick func()) *g.Node {
		side := "left-2"
		if slot == "carousel-next" {
			side = "right-2"
		}
		return g.Button(
			g.Type("button"),
			g.Class(style.CN("absolute top-1/2 z-10 inline-flex size-8 -translate-y-1/2 items-center justify-center rounded-full border bg-background text-foreground shadow transition-opacity hover:bg-accent disabled:opacity-40", side)),
			g.Data("slot", slot),
			g.Aria("label", slot),
			g.Disabled(disabled),
			g.OnClick(func(*g.Event) { onClick() }),
			Icon(icon, "size-4"),
		)
	}

	body := []any{
		g.Class(style.CN("relative", a.p.Class)),
		g.Data("slot", "carousel"),
		g.Role("region"),
		g.Div(g.Class("overflow-hidden rounded-lg"), track),
		arrow("carousel-previous", "chevron-left", index == 0, func() { setIndex(max(index-1, 0)) }),
		arrow("carousel-next", "chevron-right", index == n-1, func() { setIndex(min(index+1, n-1)) }),
	}

	if a.p.Dots {
		dots := make([]any, 0, n+1)
		dots = append(dots, g.Class("mt-3 flex items-center justify-center gap-1.5"), g.Data("slot", "carousel-dots"))
		for i := range n {
			dots = append(dots, g.Button(
				g.Type("button"),
				g.Class(style.CN("size-2 rounded-full transition-colors",
					map[string]bool{"bg-primary": i == index, "bg-muted-foreground/30 hover:bg-muted-foreground/60": i != index})),
				g.Data("slot", "carousel-dot"),
				g.Data("index", strconv.Itoa(i)),
				g.Aria("label", "Go to slide "+strconv.Itoa(i+1)),
				g.OnClick(func(*g.Event) { setIndex(i) }),
			))
		}
		body = append(body, g.Div(dots...))
	}

	return g.Div(body...)
}
