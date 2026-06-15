package gallery

import (
	"strings"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/ui"
)

func formSection() *g.Node {
	return section("Form", g.C0(formDemo))
}

// formDemo is a validated sign-up form built on UseForm.
func formDemo() *g.Node {
	f := ui.UseForm(map[string]string{})
	submit := func(*g.Event) {
		ok := f.Validate(map[string]func(string) string{
			"name": func(v string) string {
				if strings.TrimSpace(v) == "" {
					return "Name is required."
				}
				return ""
			},
			"email": func(v string) string {
				if !strings.Contains(v, "@") {
					return "Enter a valid email."
				}
				return ""
			},
		})
		if ok {
			ui.Toast("Submitted", ui.ToastOptions{Description: "The form is valid.", Variant: ui.ToastSuccess})
		}
	}
	return g.Div(g.Class("max-w-sm"),
		ui.FieldGroup(
			ui.FormField(ui.FormFieldProps{Form: f, Name: "name", Label: "Name", Placeholder: "Ada Lovelace"}),
			ui.FormField(ui.FormFieldProps{Form: f, Name: "email", Label: "Email", Type: "email", Placeholder: "ada@x.io", Description: "We'll never share it."}),
			ui.Button(ui.ButtonProps{OnClick: submit}, "Sign up"),
		),
	)
}
