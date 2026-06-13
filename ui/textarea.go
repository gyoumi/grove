package ui

import (
	"strconv"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

const textareaClass = "flex min-h-16 w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm transition-colors placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"

type TextareaProps struct {
	Value       string
	Placeholder string
	Class       string
	ID          string
	Rows        int
	Disabled    bool
	OnInput     func(*g.Event)
	OnChange    func(*g.Event)
}

// Textarea is a multi-line controlled text input. Its value re-syncs from
// Value on every render, like Input.
func Textarea(p TextareaProps) *g.Node {
	args := []any{
		g.Class(style.CN(textareaClass, p.Class)),
		g.Data("slot", "textarea"),
		g.Value(p.Value),
		g.Disabled(p.Disabled),
	}
	if p.ID != "" {
		args = append(args, g.ID(p.ID))
	}
	if p.Placeholder != "" {
		args = append(args, g.Attr("placeholder", p.Placeholder))
	}
	if p.Rows > 0 {
		args = append(args, g.Attr("rows", strconv.Itoa(p.Rows)))
	}
	if p.OnInput != nil {
		args = append(args, g.OnInput(p.OnInput))
	}
	if p.OnChange != nil {
		args = append(args, g.OnChange(p.OnChange))
	}
	return g.El("textarea", args...)
}
