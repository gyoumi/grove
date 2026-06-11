package style

import "sort"

// Variants expresses a component's class variants, mirroring how shadcn
// defines them with class-variance-authority:
//
//	var buttonVariants = style.Variants{
//	    Base: "inline-flex items-center ...",
//	    Groups: map[string]map[string]string{
//	        "variant": {"default": "bg-primary ...", "outline": "border ..."},
//	        "size":    {"default": "h-9 px-4", "sm": "h-8 px-3"},
//	    },
//	    Defaults: map[string]string{"variant": "default", "size": "default"},
//	}
//
// Class picks the selected variant classes (falling back to Defaults) and
// merges them with any extra classes through CN.
type Variants struct {
	Base     string
	Groups   map[string]map[string]string
	Defaults map[string]string
}

// Class resolves the variant selection into a merged class string. selected
// may be nil or partial; extra arguments are appended and CN-merged, so
// caller-supplied classes override the variant's.
func (v Variants) Class(selected map[string]string, extra ...any) string {
	parts := []any{v.Base}
	names := make([]string, 0, len(v.Groups))
	for name := range v.Groups {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		val := selected[name]
		if val == "" {
			val = v.Defaults[name]
		}
		if cls := v.Groups[name][val]; cls != "" {
			parts = append(parts, cls)
		}
	}
	parts = append(parts, extra...)
	return CN(parts...)
}
