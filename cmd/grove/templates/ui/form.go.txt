package ui

import (
	"maps"

	g "github.com/gyoumi/grove"
)

// Form holds a set of string field values and their validation errors. Create
// one with UseForm inside a component; bind inputs with Input and check rules
// with Validate.
type Form struct {
	values    map[string]string
	errors    map[string]string
	setValues func(map[string]string)
	setErrors func(map[string]string)
}

// UseForm initialises form state. initial seeds the field values.
func UseForm(initial map[string]string) *Form {
	values, setValues := g.UseState(maps.Clone(initial))
	if values == nil {
		values = map[string]string{}
	}
	errors, setErrors := g.UseState(map[string]string{})
	return &Form{values: values, errors: errors, setValues: setValues, setErrors: setErrors}
}

// Value and Error read a field's current value and error message.
func (f *Form) Value(name string) string { return f.values[name] }
func (f *Form) Error(name string) string { return f.errors[name] }

// Set updates a field's value (clearing its error).
func (f *Form) Set(name, value string) {
	v := maps.Clone(f.values)
	if v == nil {
		v = map[string]string{}
	}
	v[name] = value
	f.setValues(v)
	if f.errors[name] != "" {
		e := maps.Clone(f.errors)
		delete(e, name)
		f.setErrors(e)
	}
}

// Input returns the value and OnInput handler to wire a field to an Input.
func (f *Form) Input(name string) (value string, onInput func(*g.Event)) {
	return f.values[name], func(e *g.Event) { f.Set(name, e.Value()) }
}

// Validate runs each named rule against its field (a rule returns "" when the
// value is valid, or an error message). It records the errors and reports
// whether every field passed.
func (f *Form) Validate(rules map[string]func(string) string) bool {
	errs := map[string]string{}
	for name, rule := range rules {
		if msg := rule(f.values[name]); msg != "" {
			errs[name] = msg
		}
	}
	f.setErrors(errs)
	return len(errs) == 0
}

type FormFieldProps struct {
	Form        *Form
	Name        string
	Label       string
	Type        string // input type, e.g. "email", "password"
	Placeholder string
	Description string
}

// FormField is a labelled input bound to a Form field, showing the field's
// error (or description) beneath it.
func FormField(p FormFieldProps) *g.Node {
	value, onInput := p.Form.Input(p.Name)
	err := p.Form.Error(p.Name)

	inputClass := ""
	if err != "" {
		inputClass = "border-destructive focus-visible:ring-destructive"
	}
	parts := []any{
		FieldLabel(p.Name, p.Label),
		Input(InputProps{
			ID: p.Name, Type: p.Type, Value: value, Placeholder: p.Placeholder,
			OnInput: onInput, Class: inputClass,
		}),
	}
	switch {
	case err != "":
		parts = append(parts, FieldError(err))
	case p.Description != "":
		parts = append(parts, FieldDescription(p.Description))
	}
	return Field(parts...)
}
