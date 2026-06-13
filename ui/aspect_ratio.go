package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// AspectRatio constrains its content to a fixed ratio. Pass a Tailwind aspect
// utility (aspect-video, aspect-square, or aspect-[4/3]); content fills it.
//
//	ui.AspectRatio("aspect-video", ui.Img(...))
func AspectRatio(ratio string, children ...any) *g.Node {
	args := []any{
		g.Class(style.CN("relative w-full overflow-hidden [&>*]:size-full [&>*]:object-cover", ratio)),
		g.Data("slot", "aspect-ratio"),
	}
	return g.Div(append(args, children...)...)
}
