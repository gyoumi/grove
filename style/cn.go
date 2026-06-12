// Package style composes Tailwind class strings for component libraries:
// CN merges conditional classes and resolves utility conflicts (later
// classes win), and Variants expresses variant-based class maps.
package style

import (
	"fmt"
	"sort"
	"strings"
)

// CN builds a class string from strings, []string, and map[string]bool
// arguments (nil values are skipped), then resolves Tailwind utility
// conflicts so the last conflicting class wins:
//
//	CN("bg-blue-500 px-2", "bg-red-500")               // "px-2 bg-red-500"
//	CN("p-4", map[string]bool{"hidden": isHidden})     // conditional classes
//
// Conflict resolution covers the utility groups the ui components lean on;
// unknown classes are kept verbatim (deduplicated when identical).
func CN(args ...any) string {
	var toks []string
	var add func(a any)
	add = func(a any) {
		switch v := a.(type) {
		case nil:
		case string:
			toks = append(toks, strings.Fields(v)...)
		case []string:
			for _, s := range v {
				toks = append(toks, strings.Fields(s)...)
			}
		case map[string]bool:
			keys := make([]string, 0, len(v))
			for k := range v {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				if v[k] {
					toks = append(toks, strings.Fields(k)...)
				}
			}
		case []any:
			for _, x := range v {
				add(x)
			}
		default:
			panic(fmt.Sprintf("style.CN: unsupported argument type %T", a))
		}
	}
	for _, a := range args {
		add(a)
	}
	return Merge(toks...)
}

// Merge resolves Tailwind conflicts across the given class tokens, keeping
// the last token of each conflicting group.
func Merge(tokens ...string) string {
	type ref struct{ idx int }
	idx := map[string]*ref{}
	out := make([]string, 0, len(tokens))

	kill := func(key string) {
		if r, ok := idx[key]; ok {
			out[r.idx] = ""
			delete(idx, key)
		}
	}

	for _, tok := range tokens {
		pre, group := parseToken(tok)
		key := pre + "|" + group
		kill(key)
		for _, cg := range conflicts[group] {
			kill(pre + "|" + cg)
		}
		idx[key] = &ref{idx: len(out)}
		out = append(out, tok)
	}

	var b strings.Builder
	for _, t := range out {
		if t == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(t)
	}
	return b.String()
}

// parseToken splits a class into its variant prefix key (hover:, md:, ...)
// and the conflict group of its base utility.
func parseToken(tok string) (pre, group string) {
	t := tok
	if strings.HasPrefix(t, "!") {
		pre = "!"
		t = t[1:]
	}
	parts := splitVariants(t)
	base := parts[len(parts)-1]
	if len(parts) > 1 {
		variants := append([]string(nil), parts[:len(parts)-1]...)
		sort.Strings(variants)
		pre += strings.Join(variants, ":")
	}
	return pre, baseGroup(base)
}

// splitVariants splits on top-level colons, ignoring colons inside
// brackets/parens (arbitrary values and selectors).
func splitVariants(s string) []string {
	var parts []string
	depth, start := 0, 0
	for i, r := range s {
		switch r {
		case '[', '(':
			depth++
		case ']', ')':
			depth--
		case ':':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	return append(parts, s[start:])
}

func set(items ...string) map[string]bool {
	m := make(map[string]bool, len(items))
	for _, s := range items {
		m[s] = true
	}
	return m
}

var exactGroups = map[string]string{
	"block": "display", "inline-block": "display", "inline": "display",
	"flex": "display", "inline-flex": "display", "grid": "display",
	"inline-grid": "display", "table": "display", "contents": "display",
	"hidden": "display", "flow-root": "display",

	"static": "position", "fixed": "position", "absolute": "position",
	"relative": "position", "sticky": "position",

	"visible": "visibility", "invisible": "visibility", "collapse": "visibility",

	"uppercase": "text-transform", "lowercase": "text-transform",
	"capitalize": "text-transform", "normal-case": "text-transform",

	"underline": "text-decoration", "overline": "text-decoration",
	"line-through": "text-decoration", "no-underline": "text-decoration",

	"italic": "font-style", "not-italic": "font-style",

	"truncate": "truncate",
	"sr-only":  "sr", "not-sr-only": "sr",

	"border": "border-w-",
	"shadow": "shadow", "shadow-inner": "shadow", "shadow-none": "shadow",
	"rounded":    "rounded",
	"transition": "transition",
	"transform":  "transform",
	"grow":       "grow", "shrink": "shrink",
	"ring": "ring-w", "outline": "outline-style", "outline-none": "outline-style",
	"resize": "resize", "filter": "filter", "underline-offset": "underline-offset",
}

type prefixRule struct {
	prefix string
	group  string              // fixed group, or
	fn     func(string) string // classifier on the remainder
}

// Longest prefixes first within each family so e.g. px- matches before p-.
var prefixRules = []prefixRule{
	{prefix: "px-", group: "pad-x"}, {prefix: "py-", group: "pad-y"},
	{prefix: "pt-", group: "pad-t"}, {prefix: "pr-", group: "pad-r"},
	{prefix: "pb-", group: "pad-b"}, {prefix: "pl-", group: "pad-l"},
	{prefix: "ps-", group: "pad-s"}, {prefix: "pe-", group: "pad-e"},
	{prefix: "p-", group: "pad"},

	{prefix: "mx-", group: "margin-x"}, {prefix: "my-", group: "margin-y"},
	{prefix: "mt-", group: "margin-t"}, {prefix: "mr-", group: "margin-r"},
	{prefix: "mb-", group: "margin-b"}, {prefix: "ml-", group: "margin-l"},
	{prefix: "ms-", group: "margin-s"}, {prefix: "me-", group: "margin-e"},
	{prefix: "m-", group: "margin"},

	{prefix: "space-x-", group: "space-x"}, {prefix: "space-y-", group: "space-y"},

	{prefix: "inset-x-", group: "inset-x"}, {prefix: "inset-y-", group: "inset-y"},
	{prefix: "inset-", group: "inset"},
	{prefix: "top-", group: "top"}, {prefix: "right-", group: "right"},
	{prefix: "bottom-", group: "bottom"}, {prefix: "left-", group: "left"},
	{prefix: "z-", group: "z"},

	{prefix: "size-", group: "size"},
	{prefix: "min-w-", group: "min-w"}, {prefix: "min-h-", group: "min-h"},
	{prefix: "max-w-", group: "max-w"}, {prefix: "max-h-", group: "max-h"},
	{prefix: "w-", group: "w"}, {prefix: "h-", group: "h"},

	{prefix: "basis-", group: "basis"}, {prefix: "grow-", group: "grow"},
	{prefix: "shrink-", group: "shrink"}, {prefix: "order-", group: "order"},
	{prefix: "flex-", fn: flexGroup},

	{prefix: "grid-cols-", group: "grid-cols"}, {prefix: "grid-rows-", group: "grid-rows"},
	{prefix: "grid-flow-", group: "grid-flow"},
	{prefix: "col-span-", group: "col-span"}, {prefix: "col-start-", group: "col-start"},
	{prefix: "col-end-", group: "col-end"},
	{prefix: "row-span-", group: "row-span"}, {prefix: "row-start-", group: "row-start"},
	{prefix: "row-end-", group: "row-end"},
	{prefix: "col-", group: "col"}, {prefix: "row-", group: "row"},
	{prefix: "auto-cols-", group: "auto-cols"}, {prefix: "auto-rows-", group: "auto-rows"},

	{prefix: "gap-x-", group: "gap-x"}, {prefix: "gap-y-", group: "gap-y"},
	{prefix: "gap-", group: "gap"},

	{prefix: "justify-items-", group: "justify-items"},
	{prefix: "justify-self-", group: "justify-self"},
	{prefix: "justify-", group: "justify"},
	{prefix: "place-items-", group: "place-items"},
	{prefix: "place-content-", group: "place-content"},
	{prefix: "place-self-", group: "place-self"},
	{prefix: "items-", group: "items"}, {prefix: "content-", group: "content"},
	{prefix: "self-", group: "self"},

	{prefix: "font-", fn: fontGroup},
	{prefix: "text-", fn: textGroup},
	{prefix: "leading-", group: "leading"}, {prefix: "tracking-", group: "tracking"},
	{prefix: "whitespace-", group: "whitespace"}, {prefix: "break-", group: "break"},
	{prefix: "indent-", group: "indent"}, {prefix: "align-", group: "align"},
	{prefix: "list-", group: "list"},
	{prefix: "decoration-", fn: decorationGroup},
	{prefix: "underline-offset-", group: "underline-offset"},
	{prefix: "line-clamp-", group: "line-clamp"},

	{prefix: "bg-", fn: bgGroup},
	{prefix: "from-", group: "gradient-from"}, {prefix: "via-", group: "gradient-via"},
	{prefix: "to-", group: "gradient-to"},

	{prefix: "border-", fn: borderGroup},
	{prefix: "divide-x", group: "divide-x"}, {prefix: "divide-y", group: "divide-y"},
	{prefix: "divide-", group: "divide-color"},
	{prefix: "rounded-", fn: roundedGroup},

	{prefix: "shadow-", group: "shadow"},
	{prefix: "opacity-", group: "opacity"},
	{prefix: "ring-", fn: ringGroup},
	{prefix: "outline-", fn: outlineGroup},

	{prefix: "overflow-x-", group: "overflow-x"}, {prefix: "overflow-y-", group: "overflow-y"},
	{prefix: "overflow-", group: "overflow"},
	{prefix: "overscroll-", group: "overscroll"},
	{prefix: "object-", fn: objectGroup},
	{prefix: "aspect-", group: "aspect"}, {prefix: "columns-", group: "columns"},

	{prefix: "cursor-", group: "cursor"}, {prefix: "select-", group: "select"},
	{prefix: "pointer-events-", group: "pointer-events"},
	{prefix: "touch-", group: "touch"}, {prefix: "scroll-", group: "scroll"},
	{prefix: "snap-", group: "snap"}, {prefix: "will-change-", group: "will-change"},
	{prefix: "appearance-", group: "appearance"}, {prefix: "resize-", group: "resize"},

	{prefix: "transition-", group: "transition"},
	{prefix: "duration-", group: "duration"}, {prefix: "ease-", group: "ease"},
	{prefix: "delay-", group: "delay"}, {prefix: "animate-", group: "animate"},

	{prefix: "scale-x-", group: "scale-x"}, {prefix: "scale-y-", group: "scale-y"},
	{prefix: "scale-", group: "scale"},
	{prefix: "rotate-", group: "rotate"},
	{prefix: "translate-x-", group: "translate-x"}, {prefix: "translate-y-", group: "translate-y"},
	{prefix: "translate-", group: "translate"},
	{prefix: "skew-", group: "skew"}, {prefix: "origin-", group: "origin"},

	{prefix: "fill-", group: "fill"},
	{prefix: "stroke-", fn: strokeGroup},
	{prefix: "accent-", group: "accent"}, {prefix: "caret-", group: "caret"},

	{prefix: "blur-", group: "blur"}, {prefix: "brightness-", group: "brightness"},
	{prefix: "contrast-", group: "contrast"}, {prefix: "saturate-", group: "saturate"},
	{prefix: "drop-shadow-", group: "drop-shadow"},
	{prefix: "backdrop-", group: "backdrop"},

	{prefix: "field-sizing-", group: "field-sizing"},
}

func baseGroup(base string) string {
	s := strings.TrimPrefix(base, "-")
	if g, ok := exactGroups[s]; ok {
		return g
	}
	// arbitrary properties: [mask-type:luminance] conflicts per property
	if strings.HasPrefix(s, "[") {
		if i := strings.IndexByte(s, ':'); i > 1 {
			return "arb:" + s[1:i]
		}
	}
	for _, rule := range prefixRules {
		if rest, ok := strings.CutPrefix(s, rule.prefix); ok {
			if rule.fn != nil {
				return rule.fn(rest)
			}
			return rule.group
		}
	}
	return "tok:" + s
}

// preSlash strips a trailing /modifier (text-sm/6, bg-black/50) so value
// classification sees the stem; slashes inside brackets are preserved.
func preSlash(s string) string {
	depth := 0
	for i, r := range s {
		switch r {
		case '[', '(':
			depth++
		case ']', ')':
			depth--
		case '/':
			if depth == 0 {
				return s[:i]
			}
		}
	}
	return s
}

// conflicts lists, per group, the more specific groups a class overrides
// (p-4 overrides px-2; the reverse is not true).
var conflicts = map[string][]string{
	"pad":   {"pad-x", "pad-y", "pad-t", "pad-r", "pad-b", "pad-l", "pad-s", "pad-e"},
	"pad-x": {"pad-l", "pad-r", "pad-s", "pad-e"},
	"pad-y": {"pad-t", "pad-b"},

	"margin":   {"margin-x", "margin-y", "margin-t", "margin-r", "margin-b", "margin-l", "margin-s", "margin-e"},
	"margin-x": {"margin-l", "margin-r", "margin-s", "margin-e"},
	"margin-y": {"margin-t", "margin-b"},

	"inset":   {"inset-x", "inset-y", "top", "right", "bottom", "left"},
	"inset-x": {"left", "right"},
	"inset-y": {"top", "bottom"},

	"size":     {"w", "h"},
	"gap":      {"gap-x", "gap-y"},
	"overflow": {"overflow-x", "overflow-y"},
	"flex":     {"grow", "shrink", "basis"},

	"text-size": {"leading"},

	"rounded": {"rounded-t", "rounded-r", "rounded-b", "rounded-l",
		"rounded-tl", "rounded-tr", "rounded-br", "rounded-bl",
		"rounded-s", "rounded-e", "rounded-ss", "rounded-se", "rounded-es", "rounded-ee"},
	"rounded-t": {"rounded-tl", "rounded-tr"},
	"rounded-r": {"rounded-tr", "rounded-br"},
	"rounded-b": {"rounded-bl", "rounded-br"},
	"rounded-l": {"rounded-tl", "rounded-bl"},

	"border-w-":  {"border-w-t", "border-w-r", "border-w-b", "border-w-l", "border-w-x", "border-w-y", "border-w-s", "border-w-e"},
	"border-w-x": {"border-w-l", "border-w-r"},
	"border-w-y": {"border-w-t", "border-w-b"},

	"border-color-":  {"border-color-t", "border-color-r", "border-color-b", "border-color-l", "border-color-x", "border-color-y", "border-color-s", "border-color-e"},
	"border-color-x": {"border-color-l", "border-color-r"},
	"border-color-y": {"border-color-t", "border-color-b"},

	"translate": {"translate-x", "translate-y"},
	"scale":     {"scale-x", "scale-y"},

	// col-auto / row-auto reset the whole grid-column/row shorthand
	"col": {"col-span", "col-start", "col-end"},
	"row": {"row-span", "row-start", "row-end"},
}

var (
	textSizes = set("xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl", "5xl", "6xl", "7xl", "8xl", "9xl")
	textAlign = set("left", "center", "right", "justify", "start", "end")

	fontWeights = set("thin", "extralight", "light", "normal", "medium", "semibold", "bold", "extrabold", "black")

	borderSides     = set("t", "r", "b", "l", "x", "y", "s", "e")
	roundedCorners  = set("t", "r", "b", "l", "tl", "tr", "br", "bl", "s", "e", "ss", "se", "es", "ee")
	decorationStyle = set("solid", "double", "dotted", "dashed", "wavy")
)

func looksColor(s string) bool {
	for _, p := range []string{"#", "rgb", "hsl", "oklch", "oklab", "hwb", "lab(", "lch(", "color", "var(--color"} {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

// widthLike reports whether a utility remainder denotes a width/size value
// rather than a color: "", numbers, px, or non-color arbitrary values.
func widthLike(rest string) bool {
	if rest == "" || rest == "px" {
		return true
	}
	if c := rest[0]; c >= '0' && c <= '9' {
		return true
	}
	if strings.HasPrefix(rest, "[") {
		return !looksColor(strings.Trim(rest, "[]"))
	}
	return false
}

func flexGroup(rest string) string {
	switch rest {
	case "row", "row-reverse", "col", "col-reverse":
		return "flex-direction"
	case "wrap", "nowrap", "wrap-reverse":
		return "flex-wrap"
	}
	return "flex"
}

func fontGroup(rest string) string {
	if fontWeights[rest] {
		return "font-weight"
	}
	return "font-family"
}

func textGroup(rest string) string {
	if textSizes[preSlash(rest)] {
		// the /modifier form (text-sm/6) sets size and line-height at once
		return "text-size"
	}
	if textAlign[rest] {
		return "text-align"
	}
	switch rest {
	case "wrap", "nowrap", "balance", "pretty":
		return "text-wrap"
	case "ellipsis", "clip":
		return "text-overflow"
	}
	if strings.HasPrefix(rest, "[") && !looksColor(strings.Trim(rest, "[]")) {
		return "text-size"
	}
	return "text-color"
}

func bgGroup(rest string) string {
	if rest == "none" || strings.HasPrefix(rest, "gradient-") ||
		strings.HasPrefix(rest, "linear-") || strings.HasPrefix(rest, "radial-") ||
		strings.HasPrefix(rest, "conic-") {
		return "bg-image"
	}
	switch rest {
	case "cover", "contain", "auto":
		return "bg-size"
	case "fixed", "local", "scroll":
		return "bg-attachment"
	case "repeat", "no-repeat", "repeat-x", "repeat-y", "repeat-round", "repeat-space":
		return "bg-repeat"
	case "center", "top", "bottom", "left", "right",
		"top-left", "top-right", "bottom-left", "bottom-right",
		"left-top", "left-bottom", "right-top", "right-bottom":
		return "bg-position"
	}
	if strings.HasPrefix(rest, "clip-") {
		return "bg-clip"
	}
	if strings.HasPrefix(rest, "origin-") {
		return "bg-origin"
	}
	return "bg-color"
}

func borderGroup(rest string) string {
	switch rest {
	case "solid", "dashed", "dotted", "double", "hidden", "none":
		return "border-style"
	case "collapse", "separate":
		return "border-collapse"
	}
	if strings.HasPrefix(rest, "spacing") {
		return "border-spacing"
	}
	side, rem := "", rest
	if i := strings.IndexByte(rest, '-'); i > 0 && borderSides[rest[:i]] {
		side, rem = rest[:i], rest[i+1:]
	} else if borderSides[rest] {
		side, rem = rest, ""
	}
	if widthLike(rem) {
		return "border-w-" + side
	}
	return "border-color-" + side
}

func roundedGroup(rest string) string {
	if i := strings.IndexByte(rest, '-'); i > 0 && roundedCorners[rest[:i]] {
		return "rounded-" + rest[:i]
	}
	if roundedCorners[rest] {
		return "rounded-" + rest
	}
	return "rounded"
}

func ringGroup(rest string) string {
	if rest == "inset" {
		return "ring-inset"
	}
	if after, ok := strings.CutPrefix(rest, "offset-"); ok {
		if widthLike(after) {
			return "ring-offset-w"
		}
		return "ring-offset-color"
	}
	if widthLike(rest) {
		return "ring-w"
	}
	return "ring-color"
}

func outlineGroup(rest string) string {
	switch rest {
	case "none", "hidden", "solid", "dashed", "dotted", "double":
		return "outline-style"
	}
	if strings.HasPrefix(rest, "offset-") {
		return "outline-offset"
	}
	if widthLike(rest) {
		return "outline-w"
	}
	return "outline-color"
}

func objectGroup(rest string) string {
	switch rest {
	case "contain", "cover", "fill", "none", "scale-down":
		return "object-fit"
	}
	return "object-position"
}

func strokeGroup(rest string) string {
	if widthLike(rest) {
		return "stroke-w"
	}
	return "stroke-color"
}

func decorationGroup(rest string) string {
	if decorationStyle[rest] {
		return "decoration-style"
	}
	if widthLike(rest) {
		return "decoration-w"
	}
	return "decoration-color"
}
