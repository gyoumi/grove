// Package testdom is an in-memory Renderer for testing grove components
// with plain go test — no browser needed. Mount a tree, fire events, and
// assert against the rendered HTML.
package testdom

import (
	"fmt"
	"html"
	"sort"
	"strings"

	grove "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/renderer"
)

// Elem is a fake DOM node. Tag == "" means a text node.
type Elem struct {
	Tag      string
	Text     string
	ID       int // grove element id
	Attrs    map[string]string
	Props    map[string]any
	Children []*Elem

	parent *Elem
}

// R is the fake renderer plus test conveniences.
type R struct {
	Container *Elem

	dispatch renderer.Dispatch
	pending  []func()
}

// Mount renders root into a fresh fake DOM and settles all scheduled work
// (including mount effects).
func Mount(root *grove.Node) *R {
	r := &R{Container: &Elem{Tag: "#container"}}
	grove.Mount(r, r.Container, root)
	r.Settle()
	return r
}

// --- renderer.Renderer implementation ---

func (r *R) SetDispatch(d renderer.Dispatch) { r.dispatch = d }

func (r *R) CreateElement(tag string, id int) renderer.Node {
	return &Elem{Tag: tag, ID: id, Attrs: map[string]string{}, Props: map[string]any{}}
}

func (r *R) CreateText(text string) renderer.Node { return &Elem{Text: text} }

func (r *R) SetText(n renderer.Node, text string) { n.(*Elem).Text = text }

func (r *R) SetAttr(n renderer.Node, name, value string) { n.(*Elem).Attrs[name] = value }

func (r *R) RemoveAttr(n renderer.Node, name string) { delete(n.(*Elem).Attrs, name) }

func (r *R) SetProp(n renderer.Node, name string, value any) {
	e := n.(*Elem)
	if value == nil {
		delete(e.Props, name)
		return
	}
	e.Props[name] = value
}

func (r *R) InsertBefore(parent, child, before renderer.Node) {
	p := parent.(*Elem)
	c := child.(*Elem)
	if c.parent != nil {
		c.parent.removeChild(c)
	}
	c.parent = p
	if before == nil {
		p.Children = append(p.Children, c)
		return
	}
	b := before.(*Elem)
	for i, k := range p.Children {
		if k == b {
			p.Children = append(p.Children[:i], append([]*Elem{c}, p.Children[i:]...)...)
			return
		}
	}
	panic(fmt.Sprintf("testdom: InsertBefore anchor <%s> is not a child of <%s>", b.Tag, p.Tag))
}

func (r *R) Remove(parent, child renderer.Node) {
	p := parent.(*Elem)
	c := child.(*Elem)
	p.removeChild(c)
	c.parent = nil
}

func (e *Elem) removeChild(c *Elem) {
	for i, k := range e.Children {
		if k == c {
			e.Children = append(e.Children[:i], e.Children[i+1:]...)
			return
		}
	}
}

func (r *R) Listen(event string) {}

func (r *R) Schedule(f func()) { r.pending = append(r.pending, f) }

// --- renderer.EventOps over map[string]any events ---

func walk(raw any, path []string) any {
	cur := raw
	for _, p := range path {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[p]
	}
	return cur
}

func (r *R) Str(raw any, path ...string) string {
	switch v := walk(raw, path).(type) {
	case nil:
		return ""
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

func (r *R) Bool(raw any, path ...string) bool {
	b, _ := walk(raw, path).(bool)
	return b
}

func (r *R) Num(raw any, path ...string) float64 {
	switch v := walk(raw, path).(type) {
	case float64:
		return v
	case int:
		return float64(v)
	}
	return 0
}

func (r *R) PreventDefault(raw any) {
	if m, ok := raw.(map[string]any); ok {
		m["defaultPrevented"] = true
	}
}

func (r *R) StopPropagation(raw any) {}

// --- test helpers ---

// Settle runs scheduled flushes until none remain.
func (r *R) Settle() {
	for i := 0; len(r.pending) > 0; i++ {
		if i > 1000 {
			panic("testdom: scheduled work never settled")
		}
		q := r.pending
		r.pending = nil
		for _, f := range q {
			f()
		}
	}
}

// Fire dispatches an event of the given type on el and settles. raw may be
// nil; a "target" entry reflecting el's value/checked props is filled in if
// absent.
func (r *R) Fire(el *Elem, event string, raw map[string]any) {
	if r.dispatch == nil {
		panic("testdom: Fire before Mount")
	}
	if el == nil {
		panic("testdom: Fire on nil element (selector found nothing?)")
	}
	if raw == nil {
		raw = map[string]any{}
	}
	if _, ok := raw["target"]; !ok {
		target := map[string]any{}
		if v, ok := el.Props["value"]; ok {
			target["value"] = v
		}
		if c, ok := el.Props["checked"]; ok {
			target["checked"] = c
		}
		raw["target"] = target
	}
	r.dispatch(el.ID, event, raw)
	r.Settle()
}

// Click fires a click event on el.
func (r *R) Click(el *Elem) { r.Fire(el, "click", nil) }

// Input simulates typing: sets el's value prop and fires an input event.
func (r *R) Input(el *Elem, value string) {
	el.Props["value"] = value
	r.Fire(el, "input", map[string]any{"target": map[string]any{"value": value}})
}

// SetChecked simulates toggling a checkbox and fires a change event.
func (r *R) SetChecked(el *Elem, checked bool) {
	el.Props["checked"] = checked
	r.Fire(el, "change", map[string]any{"target": map[string]any{"checked": checked}})
}

// KeyDown fires a keydown event with the given key (e.g. "Enter").
func (r *R) KeyDown(el *Elem, key string) {
	r.Fire(el, "keydown", map[string]any{"key": key})
}

// Find returns the first element with the given tag, depth-first, or nil.
func (r *R) Find(tag string) *Elem {
	return r.Container.find(func(e *Elem) bool { return e.Tag == tag })
}

// FindAll returns all elements with the given tag in document order.
func (r *R) FindAll(tag string) []*Elem {
	var out []*Elem
	r.Container.each(func(e *Elem) {
		if e.Tag == tag {
			out = append(out, e)
		}
	})
	return out
}

// FindByAttr returns the first element whose attribute matches, or nil.
func (r *R) FindByAttr(name, value string) *Elem {
	return r.Container.find(func(e *Elem) bool { return e.Attrs[name] == value })
}

// FindText returns the deepest element whose text content contains s — the
// element you'd click on, not the page wrapper that also happens to contain
// the text. Ties go to the earliest match in document order.
func (r *R) FindText(s string) *Elem {
	var deepest func(e *Elem) *Elem
	deepest = func(e *Elem) *Elem {
		for _, c := range e.Children {
			if m := deepest(c); m != nil {
				return m
			}
		}
		if e.Tag != "" && e != r.Container && strings.Contains(e.TextContent(), s) {
			return e
		}
		return nil
	}
	return deepest(r.Container)
}

func (e *Elem) find(pred func(*Elem) bool) *Elem {
	var found *Elem
	e.each(func(k *Elem) {
		if found == nil && k != e && pred(k) {
			found = k
		}
	})
	return found
}

func (e *Elem) each(f func(*Elem)) {
	f(e)
	for _, c := range e.Children {
		c.each(f)
	}
}

// TextContent returns the concatenated text of e and its descendants.
func (e *Elem) TextContent() string {
	if e.Tag == "" {
		return e.Text
	}
	var b strings.Builder
	for _, c := range e.Children {
		b.WriteString(c.TextContent())
	}
	return b.String()
}

var voidTags = map[string]bool{
	"br": true, "hr": true, "img": true, "input": true, "meta": true,
	"link": true, "source": true, "col": true, "wbr": true,
}

// HTML renders the mounted tree as an HTML string for snapshot assertions.
// Attributes are sorted; value/checked/disabled props render as attributes.
func (r *R) HTML() string {
	var b strings.Builder
	for _, c := range r.Container.Children {
		c.writeHTML(&b)
	}
	return b.String()
}

// HTML renders this element (and subtree) as HTML.
func (e *Elem) HTML() string {
	var b strings.Builder
	e.writeHTML(&b)
	return b.String()
}

func (e *Elem) writeHTML(b *strings.Builder) {
	if e.Tag == "" {
		b.WriteString(html.EscapeString(e.Text))
		return
	}
	b.WriteByte('<')
	b.WriteString(e.Tag)

	type kv struct{ k, v string }
	var pairs []kv
	for k, v := range e.Attrs {
		pairs = append(pairs, kv{k, v})
	}
	for k, v := range e.Props {
		switch val := v.(type) {
		case bool:
			if val {
				pairs = append(pairs, kv{k, ""})
			}
		default:
			pairs = append(pairs, kv{k, fmt.Sprint(v)})
		}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].k < pairs[j].k })
	for _, p := range pairs {
		b.WriteByte(' ')
		b.WriteString(p.k)
		if !(p.v == "" && isBoolish(p.k)) {
			b.WriteString(`="`)
			b.WriteString(html.EscapeString(p.v))
			b.WriteByte('"')
		}
	}

	if voidTags[e.Tag] && len(e.Children) == 0 {
		b.WriteString("/>")
		return
	}
	b.WriteByte('>')
	for _, c := range e.Children {
		c.writeHTML(b)
	}
	b.WriteString("</")
	b.WriteString(e.Tag)
	b.WriteByte('>')
}

func isBoolish(name string) bool {
	switch name {
	case "disabled", "checked", "readonly", "required", "selected", "open":
		return true
	}
	return false
}
