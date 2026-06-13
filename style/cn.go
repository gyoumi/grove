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
// and the conflict group of its base utility. Important markers — leading !
// (v3) or trailing ! (v4) — scope the class like a variant, so important
// and non-important utilities never displace each other.
func parseToken(tok string) (pre, group string) {
	t := tok
	imp := ""
	if strings.HasPrefix(t, "!") {
		imp = "!"
		t = t[1:]
	}
	parts := splitVariants(t)
	base := parts[len(parts)-1]
	if strings.HasSuffix(base, "!") {
		imp = "!"
		base = base[:len(base)-1]
	}
	pre = imp
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
	"inline-table": "display", "table-caption": "display", "table-cell": "display",
	"table-column": "display", "table-column-group": "display",
	"table-footer-group": "display", "table-header-group": "display",
	"table-row": "display", "table-row-group": "display", "list-item": "display",

	"static": "position", "fixed": "position", "absolute": "position",
	"relative": "position", "sticky": "position",

	"visible": "visibility", "invisible": "visibility", "collapse": "visibility",

	"isolate": "isolation", "isolation-auto": "isolation",
	"box-border": "box-sizing", "box-content": "box-sizing",
	"container": "container",

	"uppercase": "text-transform", "lowercase": "text-transform",
	"capitalize": "text-transform", "normal-case": "text-transform",

	"underline": "text-decoration", "overline": "text-decoration",
	"line-through": "text-decoration", "no-underline": "text-decoration",

	"italic": "font-style", "not-italic": "font-style",
	"antialiased": "font-smoothing", "subpixel-antialiased": "font-smoothing",

	// font-variant-numeric composes; normal-nums resets the whole property
	"normal-nums": "fvn",
	"ordinal":     "fvn-ordinal", "slashed-zero": "fvn-slashed",
	"lining-nums": "fvn-figure", "oldstyle-nums": "fvn-figure",
	"proportional-nums": "fvn-spacing", "tabular-nums": "fvn-spacing",
	"diagonal-fractions": "fvn-fraction", "stacked-fractions": "fvn-fraction",

	"truncate": "truncate",
	"sr-only":  "sr", "not-sr-only": "sr",

	"wrap-normal": "overflow-wrap", "wrap-break-word": "overflow-wrap",
	"wrap-anywhere": "overflow-wrap",

	"table-auto": "table-layout", "table-fixed": "table-layout",

	"border": "border-w-",
	"shadow": "shadow", "shadow-inner": "shadow", "shadow-none": "shadow",
	"inset-ring": "inset-ring-w", "inset-shadow": "inset-shadow-size",
	"rounded":    "rounded",
	"transition": "transition",
	"transform":  "transform", "transform-gpu": "transform",
	"transform-cpu": "transform", "transform-none": "transform",
	"grow": "grow", "shrink": "shrink",
	"ring": "ring-w", "outline": "outline-style", "outline-none": "outline-style",
	"resize": "resize", "filter": "filter", "underline-offset": "underline-offset",

	"grayscale": "grayscale", "invert": "invert", "sepia": "sepia",

	"divide-x-reverse": "divide-x-reverse", "divide-y-reverse": "divide-y-reverse",
	"space-x-reverse": "space-x-reverse", "space-y-reverse": "space-y-reverse",
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

	// inset-ring/inset-shadow are box-shadow utilities, not position; they
	// must classify before the inset- position family.
	{prefix: "inset-ring-", fn: insetRingGroup},
	{prefix: "inset-shadow-", fn: insetShadowGroup},
	{prefix: "inset-x-", group: "inset-x"}, {prefix: "inset-y-", group: "inset-y"},
	{prefix: "inset-", group: "inset"},
	{prefix: "top-", group: "top"}, {prefix: "right-", group: "right"},
	{prefix: "bottom-", group: "bottom"}, {prefix: "left-", group: "left"},
	{prefix: "start-", group: "start"}, {prefix: "end-", group: "end"},
	{prefix: "z-", group: "z"},
	{prefix: "float-", group: "float"}, {prefix: "clear-", group: "clear"},
	{prefix: "box-decoration-", group: "box-decoration"},

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
	// col-auto / col-7 set the grid-column shorthand, same as col-span
	{prefix: "col-", group: "col-span"}, {prefix: "row-", group: "row-span"},
	{prefix: "auto-cols-", group: "auto-cols"}, {prefix: "auto-rows-", group: "auto-rows"},

	{prefix: "gap-x-", group: "gap-x"}, {prefix: "gap-y-", group: "gap-y"},
	{prefix: "gap-", group: "gap"},

	{prefix: "justify-items-", group: "justify-items"},
	{prefix: "justify-self-", group: "justify-self"},
	{prefix: "justify-", group: "justify"},
	{prefix: "place-items-", group: "place-items"},
	{prefix: "place-content-", group: "place-content"},
	{prefix: "place-self-", group: "place-self"},
	{prefix: "items-", group: "items"}, {prefix: "content-", fn: contentGroup},
	{prefix: "self-", group: "self"},

	{prefix: "font-stretch-", group: "font-stretch"},
	{prefix: "font-", fn: fontGroup},
	{prefix: "text-", fn: textGroup},
	{prefix: "leading-", group: "leading"}, {prefix: "tracking-", group: "tracking"},
	{prefix: "whitespace-", group: "whitespace"}, {prefix: "break-", fn: breakGroup},
	{prefix: "hyphens-", group: "hyphens"},
	{prefix: "indent-", group: "indent"}, {prefix: "align-", group: "align"},
	{prefix: "list-", fn: listGroup},
	{prefix: "decoration-", fn: decorationGroup},
	{prefix: "underline-offset-", group: "underline-offset"},
	{prefix: "line-clamp-", group: "line-clamp"},

	{prefix: "bg-", fn: bgGroup},
	{prefix: "from-", fn: gradientGroup("from")}, {prefix: "via-", fn: gradientGroup("via")},
	{prefix: "to-", fn: gradientGroup("to")},

	{prefix: "border-", fn: borderGroup},
	{prefix: "divide-x", group: "divide-x"}, {prefix: "divide-y", group: "divide-y"},
	{prefix: "divide-", fn: divideGroup},
	{prefix: "rounded-", fn: roundedGroup},

	{prefix: "shadow-", fn: shadowGroup},
	{prefix: "opacity-", group: "opacity"},
	{prefix: "mix-blend-", group: "mix-blend"},
	{prefix: "ring-", fn: ringGroup},
	{prefix: "outline-", fn: outlineGroup},

	{prefix: "overflow-x-", group: "overflow-x"}, {prefix: "overflow-y-", group: "overflow-y"},
	{prefix: "overflow-", group: "overflow"},
	{prefix: "overscroll-x-", group: "overscroll-x"}, {prefix: "overscroll-y-", group: "overscroll-y"},
	{prefix: "overscroll-", group: "overscroll"},
	{prefix: "object-", fn: objectGroup},
	{prefix: "aspect-", group: "aspect"}, {prefix: "columns-", group: "columns"},

	{prefix: "cursor-", group: "cursor"}, {prefix: "select-", group: "select"},
	{prefix: "pointer-events-", group: "pointer-events"},
	{prefix: "touch-", fn: touchGroup}, {prefix: "scroll-", fn: scrollGroup},
	{prefix: "snap-", fn: snapGroup}, {prefix: "will-change-", group: "will-change"},
	{prefix: "appearance-", group: "appearance"}, {prefix: "resize-", group: "resize"},
	{prefix: "caption-", group: "caption"},
	{prefix: "forced-color-adjust-", group: "forced-color-adjust"},

	{prefix: "transition-", group: "transition"},
	{prefix: "duration-", group: "duration"}, {prefix: "ease-", group: "ease"},
	{prefix: "delay-", group: "delay"}, {prefix: "animate-", group: "animate"},

	{prefix: "scale-x-", group: "scale-x"}, {prefix: "scale-y-", group: "scale-y"},
	{prefix: "scale-", group: "scale"},
	{prefix: "rotate-x-", group: "rotate-x"}, {prefix: "rotate-y-", group: "rotate-y"},
	{prefix: "rotate-z-", group: "rotate-z"},
	{prefix: "rotate-", group: "rotate"},
	{prefix: "translate-x-", group: "translate-x"}, {prefix: "translate-y-", group: "translate-y"},
	{prefix: "translate-", group: "translate"},
	{prefix: "skew-x-", group: "skew-x"}, {prefix: "skew-y-", group: "skew-y"},
	{prefix: "skew-", group: "skew"}, {prefix: "origin-", group: "origin"},
	{prefix: "perspective-origin-", group: "perspective-origin"},
	{prefix: "perspective-", group: "perspective"},

	{prefix: "fill-", group: "fill"},
	{prefix: "stroke-", fn: strokeGroup},
	{prefix: "accent-", group: "accent"}, {prefix: "caret-", group: "caret"},

	{prefix: "blur-", group: "blur"}, {prefix: "brightness-", group: "brightness"},
	{prefix: "contrast-", group: "contrast"}, {prefix: "saturate-", group: "saturate"},
	{prefix: "grayscale-", group: "grayscale"}, {prefix: "invert-", group: "invert"},
	{prefix: "sepia-", group: "sepia"}, {prefix: "hue-rotate-", group: "hue-rotate"},
	{prefix: "drop-shadow-", fn: dropShadowGroup},
	{prefix: "backdrop-", fn: backdropGroup},
	{prefix: "mask-", fn: maskGroup},

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

	"inset":   {"inset-x", "inset-y", "top", "right", "bottom", "left", "start", "end"},
	"inset-x": {"left", "right", "start", "end"},
	"inset-y": {"top", "bottom"},

	"size":       {"w", "h"},
	"gap":        {"gap-x", "gap-y"},
	"overflow":   {"overflow-x", "overflow-y"},
	"overscroll": {"overscroll-x", "overscroll-y"},
	"flex":       {"grow", "shrink", "basis"},

	"text-size": {"leading"},

	// font-variant-numeric: normal-nums resets every axis, and any axis
	// utility cancels a normal-nums reset
	"fvn":          {"fvn-ordinal", "fvn-slashed", "fvn-figure", "fvn-spacing", "fvn-fraction"},
	"fvn-ordinal":  {"fvn"},
	"fvn-slashed":  {"fvn"},
	"fvn-figure":   {"fvn"},
	"fvn-spacing":  {"fvn"},
	"fvn-fraction": {"fvn"},

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

	"border-spacing": {"border-spacing-x", "border-spacing-y"},

	"translate": {"translate-x", "translate-y"},
	"scale":     {"scale-x", "scale-y"},
	"skew":      {"skew-x", "skew-y"},

	// col-span/col-auto/col-7 set the grid-column shorthand, overriding
	// the start/end longhands (likewise rows)
	"col-span": {"col-start", "col-end"},
	"row-span": {"row-start", "row-end"},

	"scroll-m":  {"scroll-mx", "scroll-my", "scroll-mt", "scroll-mr", "scroll-mb", "scroll-ml", "scroll-ms", "scroll-me"},
	"scroll-mx": {"scroll-ml", "scroll-mr", "scroll-ms", "scroll-me"},
	"scroll-my": {"scroll-mt", "scroll-mb"},
	"scroll-p":  {"scroll-px", "scroll-py", "scroll-pt", "scroll-pr", "scroll-pb", "scroll-pl", "scroll-ps", "scroll-pe"},
	"scroll-px": {"scroll-pl", "scroll-pr", "scroll-ps", "scroll-pe"},
	"scroll-py": {"scroll-pt", "scroll-pb"},

	// touch-auto/none/manipulation reset the pan/pinch axes, and any axis
	// utility cancels the reset
	"touch":    {"touch-x", "touch-y", "touch-pz"},
	"touch-x":  {"touch"},
	"touch-y":  {"touch"},
	"touch-pz": {"touch"},
}

var (
	textSizes = set("xs", "sm", "base", "lg", "xl", "2xl", "3xl", "4xl", "5xl", "6xl", "7xl", "8xl", "9xl")
	textAlign = set("left", "center", "right", "justify", "start", "end")

	fontWeights = set("thin", "extralight", "light", "normal", "medium", "semibold", "bold", "extrabold", "black")

	borderSides     = set("t", "r", "b", "l", "x", "y", "s", "e")
	roundedCorners  = set("t", "r", "b", "l", "tl", "tr", "br", "bl", "s", "e", "ss", "se", "es", "ee")
	decorationStyle = set("solid", "double", "dotted", "dashed", "wavy")

	shadowSizes = set("2xs", "xs", "sm", "md", "lg", "xl", "2xl", "none", "inner")
	spacingEnds = set("", "x", "y", "t", "r", "b", "l", "s", "e")
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
	if after, ok := strings.CutPrefix(rest, "shadow-"); ok {
		if shadowLike(after) {
			return "text-shadow"
		}
		return "text-shadow-color"
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
	if strings.HasPrefix(rest, "blend-") {
		return "bg-blend"
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
	if strings.HasPrefix(rest, "spacing-") {
		switch {
		case strings.HasPrefix(rest, "spacing-x-"):
			return "border-spacing-x"
		case strings.HasPrefix(rest, "spacing-y-"):
			return "border-spacing-y"
		}
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

// shadowLike reports whether a shadow utility remainder denotes a shadow
// size/shape (shadow-lg, shadow-[0_1px_2px]) rather than a shadow color.
func shadowLike(rest string) bool {
	r := preSlash(rest)
	if shadowSizes[r] {
		return true
	}
	if strings.HasPrefix(r, "[") {
		return !looksColor(strings.Trim(r, "[]"))
	}
	return false
}

func shadowGroup(rest string) string {
	if shadowLike(rest) {
		return "shadow"
	}
	return "shadow-color"
}

func dropShadowGroup(rest string) string {
	if shadowLike(rest) {
		return "drop-shadow"
	}
	return "drop-shadow-color"
}

func insetShadowGroup(rest string) string {
	if shadowLike(rest) {
		return "inset-shadow-size"
	}
	return "inset-shadow-color"
}

func insetRingGroup(rest string) string {
	if widthLike(rest) {
		return "inset-ring-w"
	}
	return "inset-ring-color"
}

// breakGroup separates the four CSS properties behind break-*: page break
// control (break-after/before/inside) and word breaking (break-all, ...).
func breakGroup(rest string) string {
	for _, sub := range []string{"after-", "before-", "inside-"} {
		if strings.HasPrefix(rest, sub) {
			return "break-" + sub[:len(sub)-1]
		}
	}
	return "word-break"
}

func listGroup(rest string) string {
	switch rest {
	case "inside", "outside":
		return "list-position"
	}
	if strings.HasPrefix(rest, "image-") {
		return "list-image"
	}
	return "list-type"
}

// contentGroup: content-none / content-[...] set the CSS content property;
// everything else (content-center, content-between, ...) is align-content.
func contentGroup(rest string) string {
	if rest == "none" || strings.HasPrefix(rest, "[") {
		return "content"
	}
	return "align-content"
}

// gradientGroup tells gradient stop positions (from-10%, from-[25%]) apart
// from stop colors (from-red-500); each stop has its own pair of groups.
func gradientGroup(stop string) func(string) string {
	return func(rest string) string {
		if strings.HasSuffix(strings.Trim(preSlash(rest), "[]"), "%") {
			return "gradient-" + stop + "-pos"
		}
		return "gradient-" + stop
	}
}

func divideGroup(rest string) string {
	switch rest {
	case "solid", "dashed", "dotted", "double", "none":
		return "divide-style"
	}
	return "divide-color"
}

var backdropFilters = []string{
	"blur", "brightness", "contrast", "grayscale", "hue-rotate", "invert",
	"opacity", "saturate", "sepia",
}

func backdropGroup(rest string) string {
	for _, f := range backdropFilters {
		if rest == f || strings.HasPrefix(rest, f+"-") {
			return "backdrop-" + f
		}
	}
	return "backdrop"
}

// scrollGroup splits scroll-* into behavior (scroll-smooth), the scroll
// margin family, and the scroll padding family, mirroring m-/p- sides.
func scrollGroup(rest string) string {
	switch rest {
	case "auto", "smooth":
		return "scroll-behavior"
	}
	if seg, _, ok := strings.Cut(rest, "-"); ok &&
		(seg != "" && (seg[0] == 'm' || seg[0] == 'p') && spacingEnds[seg[1:]]) {
		return "scroll-" + seg
	}
	return "scroll"
}

func snapGroup(rest string) string {
	switch rest {
	case "start", "end", "center", "align-none":
		return "snap-align"
	case "normal", "always":
		return "snap-stop"
	case "none", "x", "y", "both":
		return "snap-type"
	case "mandatory", "proximity":
		return "snap-strictness"
	}
	return "snap"
}

func touchGroup(rest string) string {
	switch rest {
	case "pan-x", "pan-left", "pan-right":
		return "touch-x"
	case "pan-y", "pan-up", "pan-down":
		return "touch-y"
	case "pinch-zoom":
		return "touch-pz"
	}
	return "touch" // auto, none, manipulation reset the whole property
}

// maskGroup covers the mask-* family per underlying property; anything not
// recognized as clip/origin/mode/type/composite/size/repeat/position is the
// mask image itself.
func maskGroup(rest string) string {
	switch {
	case strings.HasPrefix(rest, "clip-"), rest == "no-clip":
		return "mask-clip"
	case strings.HasPrefix(rest, "origin-"):
		return "mask-origin"
	case strings.HasPrefix(rest, "type-"):
		return "mask-type"
	}
	switch rest {
	case "alpha", "luminance", "match":
		return "mask-mode"
	case "add", "subtract", "intersect", "exclude":
		return "mask-composite"
	case "cover", "contain", "auto":
		return "mask-size"
	case "center", "top", "bottom", "left", "right",
		"top-left", "top-right", "bottom-left", "bottom-right":
		return "mask-position"
	}
	if strings.HasPrefix(rest, "repeat") || rest == "no-repeat" {
		return "mask-repeat"
	}
	return "mask-image"
}
