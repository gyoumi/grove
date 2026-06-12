// Package grove is a React-style UI framework for Go compiled to
// WebAssembly: function components render a virtual DOM, hooks hold state,
// and a reconciler diffs renders and applies minimal updates to the page.
package grove

import (
	"fmt"
	"reflect"
	"runtime"

	"github.com/gyoumi/grove/renderer"
)

type kind uint8

const (
	kindText kind = iota
	kindElement
	kindComponent
	kindFragment
)

// Node is a virtual DOM node: an element, a text node, a fragment, or a
// component invocation. Nodes are cheap one-shot descriptions — build a new
// tree every render and never reuse a *Node value in two places.
type Node struct {
	kind   kind
	tag    string
	text   string
	key    string
	attrs  map[string]string
	props  map[string]any
	events map[string]func(*Event)
	ref    *DOMRef

	children []*Node

	// component nodes
	fn      func() *Node
	fnID    uintptr // identity of the component function for reconciliation
	fnKey   any     // extra identity (used by context providers)
	fnProps any     // last props value, for Memo comparison
	memo    bool
	eq      func(old, new any) bool // custom Memo comparison (nil = ==)

	// runtime fields, populated when mounted
	id       int
	dom      renderer.Node
	inst     *instance
	parent   *Node
	rendered *Node
}

// WithKey sets the reconciliation key and returns the node, so component
// nodes in lists can be keyed: g.C(TodoItem, item).WithKey(item.ID).
// Elements can use the g.Key option instead.
func (n *Node) WithKey(key string) *Node {
	n.key = key
	return n
}

// Option configures an element node. Attribute, property, and event helpers
// (Class, Value, OnClick, ...) all return Options; custom Options work too.
type Option interface {
	Apply(n *Node)
}

// El builds an element node of the given tag. Arguments may be Options
// (attributes, properties, event handlers), *Node children, plain strings
// (turned into text nodes), []*Node slices, or nil (skipped — convenient
// for conditional children).
func El(tag string, args ...any) *Node {
	n := &Node{kind: kindElement, tag: tag}
	n.apply(args)
	return n
}

func (n *Node) apply(args []any) {
	for _, arg := range args {
		switch v := arg.(type) {
		case nil:
			// skipped: lets callers write g.If(cond, child) inline
		case *Node:
			if v != nil {
				n.children = append(n.children, v)
			}
		case string:
			n.children = append(n.children, Text(v))
		case []*Node:
			for _, c := range v {
				if c != nil {
					n.children = append(n.children, c)
				}
			}
		case []any:
			n.apply(v)
		case Option:
			v.Apply(n)
		default:
			panic(fmt.Sprintf("grove: invalid argument of type %T passed to element constructor", arg))
		}
	}
}

// Text creates a text node.
func Text(s string) *Node {
	return &Node{kind: kindText, text: s}
}

// Textf creates a text node from a format string.
func Textf(format string, a ...any) *Node {
	return Text(fmt.Sprintf(format, a...))
}

// Fragment groups children without introducing a wrapper element, like
// <>...</> in React.
func Fragment(args ...any) *Node {
	n := &Node{kind: kindFragment}
	n.apply(args)
	return n
}

// C0 wraps a parameterless function component as a node.
func C0(fn func() *Node) *Node {
	return &Node{kind: kindComponent, fn: fn, fnID: fnPtr(fn)}
}

// C wraps a function component and its props as a node. Like React, the
// component re-renders whenever its parent renders; the props value is
// whatever was passed on that render.
func C[P any](fn func(P) *Node, props P) *Node {
	return &Node{
		kind:    kindComponent,
		fn:      func() *Node { return fn(props) },
		fnID:    fnPtr(fn),
		fnProps: props,
	}
}

// Memo marks a component node to skip re-rendering when its props are
// unchanged (compared with ==; incomparable props such as funcs, slices,
// and maps never compare equal, so they defeat Memo — see MemoEq). State
// updates inside the component and context changes above it still apply.
//
//	g.Memo(g.C(Row, props))
//	g.Memo(g.C0(StaticHeader))   // never re-renders from the parent
//
// A skipped component keeps the event handlers from its last render, so
// handlers inside a Memo subtree should read changing data through a
// UseRef rather than closing over it.
func Memo(n *Node) *Node {
	if n.kind != kindComponent {
		panic("grove: Memo wraps component nodes (g.C/g.C0), not elements")
	}
	n.memo = true
	return n
}

// MemoEq is Memo with a custom props comparison, for props that contain
// slices or callbacks: compare the data that affects rendering and ignore
// the rest.
func MemoEq[P any](fn func(P) *Node, props P, eq func(old, new P) bool) *Node {
	n := C(fn, props)
	n.memo = true
	n.eq = func(a, b any) bool {
		ao, ok1 := a.(P)
		bo, ok2 := b.(P)
		return ok1 && ok2 && eq(ao, bo)
	}
	return n
}

func fnPtr(fn any) uintptr {
	return reflect.ValueOf(fn).Pointer()
}

func fnName(id uintptr) string {
	if f := runtime.FuncForPC(id); f != nil {
		return f.Name()
	}
	return "<unknown component>"
}

// If returns n when cond is true and nil otherwise. A nil child is skipped
// by element constructors, so this enables inline conditional rendering.
func If(cond bool, n *Node) *Node {
	if cond {
		return n
	}
	return nil
}

// IfElse returns then when cond is true, otherwise els.
func IfElse(cond bool, then, els *Node) *Node {
	if cond {
		return then
	}
	return els
}

// IfFn is like If but builds the node lazily, for children that are
// expensive (or invalid) to construct when the condition is false.
func IfFn(cond bool, build func() *Node) *Node {
	if cond {
		return build()
	}
	return nil
}

// Map renders a slice into child nodes. Give the children keys (g.Key for
// elements, WithKey for components) so list reordering reconciles correctly.
func Map[T any](items []T, f func(T) *Node) []*Node {
	out := make([]*Node, 0, len(items))
	for _, it := range items {
		if n := f(it); n != nil {
			out = append(out, n)
		}
	}
	return out
}

// MapIdx is Map with the index supplied as well.
func MapIdx[T any](items []T, f func(int, T) *Node) []*Node {
	out := make([]*Node, 0, len(items))
	for i, it := range items {
		if n := f(i, it); n != nil {
			out = append(out, n)
		}
	}
	return out
}
