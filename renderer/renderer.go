// Package renderer defines the platform interface the grove reconciler
// renders through. The browser implementation lives in the dom package;
// an in-memory implementation for tests lives in testdom.
package renderer

// Node is an opaque handle to a platform-level DOM node. Each Renderer
// implementation supplies its own concrete type and only ever receives
// handles it produced itself.
type Node any

// Dispatch delivers a platform event to the reconciler. id is the grove id
// of the nearest grove-created element at or above the event target, event
// is the grove event name (e.g. "click", "focus"), and raw is the platform
// event object.
type Dispatch func(id int, event string, raw any)

// Renderer is the set of primitive DOM operations the reconciler needs.
type Renderer interface {
	// SetDispatch installs the callback used to deliver events. Called once
	// at mount, before any Listen call.
	SetDispatch(d Dispatch)

	// PortalRoot is the node that Portal children mount under (the app's
	// mount container), so they sit outside any transformed app content and
	// position relative to the viewport.
	PortalRoot() Node

	CreateElement(tag string, id int) Node
	CreateText(text string) Node
	SetText(n Node, text string)
	SetAttr(n Node, name, value string)
	RemoveAttr(n Node, name string)
	// SetProp sets a DOM property (as opposed to an attribute) — used for
	// value, checked, and friends. A nil value clears the property.
	SetProp(n Node, name string, value any)
	// InsertBefore inserts (or moves) child under parent, before the given
	// sibling. A nil before appends at the end.
	InsertBefore(parent, child, before Node)
	Remove(parent, child Node)

	// Listen makes sure events of the given grove event type are delivered
	// to the Dispatch callback. Idempotent; listeners are delegated, so this
	// is per event type, not per element.
	Listen(event string)

	// Schedule runs f asynchronously after the current call stack unwinds
	// (a microtask in the browser). The reconciler uses it to batch renders.
	Schedule(f func())

	// Flush is called once per commit, after a render pass has applied all
	// its DOM mutations and before effects run, so effects observe a
	// finished DOM. Renderers that mutate the DOM immediately leave this a
	// no-op; a batched renderer applies its accumulated op buffer here.
	Flush()

	EventOps
}

// EventOps gives the reconciler typed access into platform event objects so
// the core package never needs to know their concrete representation.
type EventOps interface {
	PreventDefault(raw any)
	StopPropagation(raw any)
	// Str walks the given field path (e.g. "target", "value") and returns
	// the value as a string, or "" if absent.
	Str(raw any, path ...string) string
	Bool(raw any, path ...string) bool
	Num(raw any, path ...string) float64
}
