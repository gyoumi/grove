package ui

import (
	"strings"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// avatarPalette holds the backgrounds Avatar picks from by hashing the
// name, so a given person always gets the same color. Literal class names
// keep the Tailwind scanner aware of them.
var avatarPalette = []string{
	"bg-rose-500", "bg-orange-500", "bg-amber-500", "bg-lime-600",
	"bg-emerald-500", "bg-teal-500", "bg-sky-500", "bg-indigo-500",
	"bg-violet-500", "bg-fuchsia-500",
}

type AvatarSize string

const (
	AvatarSizeDefault AvatarSize = "default"
	AvatarSizeSm      AvatarSize = "sm"
	AvatarSizeLg      AvatarSize = "lg"
)

var avatarSizes = map[AvatarSize]string{
	AvatarSizeDefault: "size-8 text-xs",
	AvatarSizeSm:      "size-6 text-[10px]",
	AvatarSizeLg:      "size-10 text-sm",
}

type AvatarProps struct {
	// Name drives both the initials and the (stable) background color.
	Name  string
	Size  AvatarSize
	Class string
}

// Avatar renders a colored circle with the name's initials — up to two
// letters, from the first two words.
func Avatar(p AvatarProps, args ...any) *g.Node {
	size := avatarSizes[p.Size]
	if size == "" {
		size = avatarSizes[AvatarSizeDefault]
	}
	all := []any{
		g.Class(style.CN(
			"inline-flex shrink-0 select-none items-center justify-center rounded-full font-semibold uppercase text-white ring-2 ring-background",
			size, avatarPalette[nameHash(p.Name)%len(avatarPalette)], p.Class)),
		g.Data("slot", "avatar"),
		g.Title(p.Name),
		initials(p.Name),
	}
	return g.Span(append(all, args...)...)
}

func initials(name string) string {
	var out []rune
	for _, word := range strings.Fields(name) {
		out = append(out, []rune(word)[0])
		if len(out) == 2 {
			break
		}
	}
	if len(out) == 0 {
		return "?"
	}
	return string(out)
}

func nameHash(s string) int {
	h := 0
	for _, r := range s {
		h = h*31 + int(r)
	}
	if h < 0 {
		h = -h
	}
	return h
}
