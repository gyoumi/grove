package grove

import "github.com/gyoumi/grove/renderer"

// Event is passed to event handlers. It wraps the platform event (a
// js.Value in the browser) and offers typed accessors for the common cases;
// Str/Bool/Num walk arbitrary field paths for everything else.
type Event struct {
	// Type is the grove event name, e.g. "click" or "input".
	Type string
	// Raw is the platform event object, for escape-hatch interop.
	Raw any

	ops     renderer.EventOps
	stopped bool
}

// Value returns event.target.value — the current text of an input.
func (e *Event) Value() string { return e.ops.Str(e.Raw, "target", "value") }

// Checked returns event.target.checked.
func (e *Event) Checked() bool { return e.ops.Bool(e.Raw, "target", "checked") }

// Key returns the keyboard key for key events, e.g. "Enter" or "Escape".
func (e *Event) Key() string { return e.ops.Str(e.Raw, "key") }

// ShiftKey reports whether shift was held.
func (e *Event) ShiftKey() bool { return e.ops.Bool(e.Raw, "shiftKey") }

// Str reads an arbitrary string field path off the platform event,
// e.g. e.Str("target", "id").
func (e *Event) Str(path ...string) string { return e.ops.Str(e.Raw, path...) }

// Bool reads an arbitrary boolean field path off the platform event.
func (e *Event) Bool(path ...string) bool { return e.ops.Bool(e.Raw, path...) }

// Num reads an arbitrary numeric field path, e.g. e.Num("clientX").
func (e *Event) Num(path ...string) float64 { return e.ops.Num(e.Raw, path...) }

// PreventDefault cancels the browser's default action for the event.
func (e *Event) PreventDefault() { e.ops.PreventDefault(e.Raw) }

// StopPropagation stops the event from reaching handlers on ancestor
// elements (both grove handlers and native ones).
func (e *Event) StopPropagation() {
	e.stopped = true
	e.ops.StopPropagation(e.Raw)
}

type eventOpt struct {
	name    string
	handler func(*Event)
}

func (o eventOpt) Apply(n *Node) {
	if o.handler == nil {
		return
	}
	if n.events == nil {
		n.events = map[string]func(*Event){}
	}
	n.events[o.name] = o.handler
}

// On registers a handler for an arbitrary event type. Events are delivered
// via delegation, so only bubbling events reach handlers (focus and blur
// are translated to their bubbling focusin/focusout forms automatically).
func On(event string, h func(*Event)) Option { return eventOpt{event, h} }

func OnClick(h func(*Event)) Option     { return eventOpt{"click", h} }
func OnDblClick(h func(*Event)) Option  { return eventOpt{"dblclick", h} }
func OnInput(h func(*Event)) Option     { return eventOpt{"input", h} }
func OnChange(h func(*Event)) Option    { return eventOpt{"change", h} }
func OnSubmit(h func(*Event)) Option    { return eventOpt{"submit", h} }
func OnKeyDown(h func(*Event)) Option   { return eventOpt{"keydown", h} }
func OnKeyUp(h func(*Event)) Option     { return eventOpt{"keyup", h} }
func OnFocus(h func(*Event)) Option     { return eventOpt{"focus", h} }
func OnBlur(h func(*Event)) Option      { return eventOpt{"blur", h} }
func OnMouseDown(h func(*Event)) Option { return eventOpt{"mousedown", h} }
func OnMouseUp(h func(*Event)) Option   { return eventOpt{"mouseup", h} }
func OnMouseOver(h func(*Event)) Option { return eventOpt{"mouseover", h} }
func OnMouseOut(h func(*Event)) Option  { return eventOpt{"mouseout", h} }
