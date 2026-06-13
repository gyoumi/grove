package ui

import g "github.com/gyoumi/grove"

// Field and its parts lay out a form control with its label, description, and
// error message:
//
//	ui.Field(
//	    ui.FieldLabel("email", "Email"),
//	    ui.Input(ui.InputProps{ID: "email"}),
//	    ui.FieldDescription("We'll never share it."),
//	    ui.FieldError("Required"),
//	)
func Field(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col gap-1.5"), g.Data("slot", "field")}, args...)
	return g.Div(all...)
}

// FieldGroup stacks several Fields with consistent spacing.
func FieldGroup(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col gap-4"), g.Data("slot", "field-group")}, args...)
	return g.Div(all...)
}

// FieldLabel labels the control with id forID ("" for none).
func FieldLabel(forID string, children ...any) *g.Node {
	all := []any{
		g.Class("text-sm font-medium leading-none"),
		g.Data("slot", "field-label"),
	}
	if forID != "" {
		all = append(all, g.For(forID))
	}
	return g.Label(append(all, children...)...)
}

func FieldDescription(children ...any) *g.Node {
	all := append([]any{g.Class("text-sm text-muted-foreground"), g.Data("slot", "field-description")}, children...)
	return g.El("p", all...)
}

func FieldError(children ...any) *g.Node {
	all := append([]any{g.Class("text-sm font-medium text-destructive"), g.Data("slot", "field-error")}, children...)
	return g.El("p", all...)
}

// FieldSet groups related fields under a legend.
func FieldSet(args ...any) *g.Node {
	all := append([]any{g.Class("flex flex-col gap-4 rounded-lg border p-4"), g.Data("slot", "field-set")}, args...)
	return g.El("fieldset", all...)
}

func FieldLegend(children ...any) *g.Node {
	all := append([]any{g.Class("px-1 text-sm font-medium"), g.Data("slot", "field-legend")}, children...)
	return g.El("legend", all...)
}
