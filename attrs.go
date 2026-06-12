package grove

import (
	"strconv"
	"strings"
)

type attrOpt struct{ name, value string }

func (o attrOpt) Apply(n *Node) {
	if n.attrs == nil {
		n.attrs = map[string]string{}
	}
	n.attrs[o.name] = o.value
}

type classOpt string

func (o classOpt) Apply(n *Node) {
	if o == "" {
		return
	}
	if n.attrs == nil {
		n.attrs = map[string]string{}
	}
	if existing := n.attrs["class"]; existing != "" {
		n.attrs["class"] = existing + " " + string(o)
	} else {
		n.attrs["class"] = string(o)
	}
}

type propOpt struct {
	name  string
	value any
}

func (o propOpt) Apply(n *Node) {
	if n.props == nil {
		n.props = map[string]any{}
	}
	n.props[o.name] = o.value
}

type keyOpt string

func (o keyOpt) Apply(n *Node) { n.key = string(o) }

type noOpt struct{}

func (noOpt) Apply(*Node) {}

// Attr sets an arbitrary attribute.
func Attr(name, value string) Option { return attrOpt{name, value} }

// Prop sets an arbitrary DOM property (as opposed to an attribute).
func Prop(name string, value any) Option { return propOpt{name, value} }

// Class adds CSS classes; multiple Class options on one element accumulate.
func Class(classes ...string) Option { return classOpt(strings.Join(classes, " ")) }

// ClassIf adds the classes only when cond is true.
func ClassIf(cond bool, classes ...string) Option {
	if !cond {
		return noOpt{}
	}
	return Class(classes...)
}

// AttrIf sets the attribute only when cond is true.
func AttrIf(cond bool, name, value string) Option {
	if !cond {
		return noOpt{}
	}
	return attrOpt{name, value}
}

// OptIf applies the option only when cond is true — the Option counterpart
// of If for conditional children.
func OptIf(cond bool, o Option) Option {
	if !cond {
		return noOpt{}
	}
	return o
}

// Key sets the reconciliation key used to match this element across renders
// in dynamic lists.
func Key(key string) Option { return keyOpt(key) }

// Data sets a data-* attribute: Data("state", "open") → data-state="open".
func Data(name, value string) Option { return attrOpt{"data-" + name, value} }

// Aria sets an aria-* attribute: Aria("label", "Close") → aria-label.
func Aria(name, value string) Option { return attrOpt{"aria-" + name, value} }

// Style sets the inline style attribute.
func Style(css string) Option { return attrOpt{"style", css} }

func ID(v string) Option          { return attrOpt{"id", v} }
func Href(v string) Option        { return attrOpt{"href", v} }
func Src(v string) Option         { return attrOpt{"src", v} }
func Alt(v string) Option         { return attrOpt{"alt", v} }
func Title(v string) Option       { return attrOpt{"title", v} }
func Placeholder(v string) Option { return attrOpt{"placeholder", v} }
func Type(v string) Option        { return attrOpt{"type", v} }
func Name(v string) Option        { return attrOpt{"name", v} }
func For(v string) Option         { return attrOpt{"for", v} }
func Target(v string) Option      { return attrOpt{"target", v} }
func Rel(v string) Option         { return attrOpt{"rel", v} }
func Role(v string) Option        { return attrOpt{"role", v} }

func TabIndex(i int) Option { return attrOpt{"tabindex", strconv.Itoa(i)} }

// Disabled renders the boolean disabled attribute when true.
func Disabled(b bool) Option {
	if !b {
		return noOpt{}
	}
	return attrOpt{"disabled", ""}
}

// ReadOnly renders the boolean readonly attribute when true.
func ReadOnly(b bool) Option {
	if !b {
		return noOpt{}
	}
	return attrOpt{"readonly", ""}
}

// Required renders the boolean required attribute when true.
func Required(b bool) Option {
	if !b {
		return noOpt{}
	}
	return attrOpt{"required", ""}
}

// Value sets the value property — use this for controlled inputs. The
// reconciler re-syncs it on every render so the DOM can't drift from state.
func Value(v string) Option { return propOpt{"value", v} }

// Checked sets the checked property — use this for controlled checkboxes.
func Checked(b bool) Option { return propOpt{"checked", b} }

// DOMRef holds a handle to a real DOM node after mount; see BindRef.
type DOMRef = Ref[any]

type refOpt struct{ r *DOMRef }

func (o refOpt) Apply(n *Node) { n.ref = o.r }

// BindRef stores the element's platform DOM handle (a js.Value in the
// browser) in ref.Current after mount, and resets it to nil on unmount.
// Combine with UseRef: ref := g.UseRef[any](nil); g.Div(g.BindRef(ref)).
func BindRef(r *DOMRef) Option { return refOpt{r} }
