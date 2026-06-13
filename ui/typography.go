package ui

import g "github.com/gyoumi/grove"

// Typography helpers apply the design system's text styles to common prose
// elements, so headings and copy look consistent without hand-written
// classes. Each takes the usual element args (text, nodes, extra g.Class).

func TypographyH1(args ...any) *g.Node {
	return g.El("h1", prepend("scroll-m-20 text-4xl font-extrabold tracking-tight lg:text-5xl", args)...)
}
func TypographyH2(args ...any) *g.Node {
	return g.El("h2", prepend("scroll-m-20 border-b pb-2 text-3xl font-semibold tracking-tight first:mt-0", args)...)
}
func TypographyH3(args ...any) *g.Node {
	return g.El("h3", prepend("scroll-m-20 text-2xl font-semibold tracking-tight", args)...)
}
func TypographyH4(args ...any) *g.Node {
	return g.El("h4", prepend("scroll-m-20 text-xl font-semibold tracking-tight", args)...)
}
func TypographyP(args ...any) *g.Node {
	return g.El("p", prepend("leading-7 [&:not(:first-child)]:mt-6", args)...)
}
func TypographyLead(args ...any) *g.Node {
	return g.El("p", prepend("text-xl text-muted-foreground", args)...)
}
func TypographyLarge(args ...any) *g.Node {
	return g.Div(prepend("text-lg font-semibold", args)...)
}
func TypographySmall(args ...any) *g.Node {
	return g.El("small", prepend("text-sm font-medium leading-none", args)...)
}
func TypographyMuted(args ...any) *g.Node {
	return g.El("p", prepend("text-sm text-muted-foreground", args)...)
}
func TypographyBlockquote(args ...any) *g.Node {
	return g.El("blockquote", prepend("mt-6 border-l-2 pl-6 italic", args)...)
}
func TypographyInlineCode(args ...any) *g.Node {
	return g.El("code", prepend("relative rounded bg-muted px-[0.3rem] py-[0.2rem] font-mono text-sm font-semibold", args)...)
}
func TypographyList(args ...any) *g.Node {
	return g.El("ul", prepend("my-6 ml-6 list-disc [&>li]:mt-2", args)...)
}

func prepend(base string, args []any) []any {
	return append([]any{g.Class(base), g.Data("slot", "typography")}, args...)
}
