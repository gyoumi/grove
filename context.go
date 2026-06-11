package grove

// ctxKey is the identity of a context; every NewContext call allocates a
// distinct one.
type ctxKey struct{ _ byte }

// Context carries a value down the component tree without prop drilling,
// like React context. Create one at package level with NewContext, wrap a
// subtree with Provider, and read it anywhere below with UseContext.
type Context[T any] struct {
	def T
	key *ctxKey
}

// NewContext creates a context with a default value, returned by UseContext
// when no Provider is above the caller.
func NewContext[T any](def T) *Context[T] {
	return &Context[T]{def: def, key: &ctxKey{}}
}

// Provider makes value available to UseContext in all descendant
// components. Children re-render with the new value whenever the provider's
// parent re-renders (grove has no render-skipping memo yet, so descendants
// always see the latest value).
func (c *Context[T]) Provider(value T, children ...any) *Node {
	n := &Node{kind: kindComponent, fnKey: c.key}
	n.fn = func() *Node {
		inst := current()
		inst.ctxKey = c.key
		inst.ctxVal = value
		return Fragment(children...)
	}
	// fnID stays 0 for all providers; fnKey distinguishes contexts during
	// reconciliation.
	return n
}

type ctxCell struct {
	provider *instance
	resolved bool
}

// UseContext returns the value of the nearest enclosing Provider for c, or
// the context's default when there is none.
func UseContext[T any](c *Context[T]) T {
	inst := current()
	i, fresh := inst.slot()
	if fresh {
		cell := &ctxCell{}
		for a := inst.parentInst; a != nil; a = a.parentInst {
			if a.ctxKey == c.key {
				cell.provider = a
				break
			}
		}
		cell.resolved = true
		inst.hooks = append(inst.hooks, cell)
	}
	cell, ok := inst.hooks[i].(*ctxCell)
	if !ok {
		hookMismatch(inst, "UseContext")
	}
	if cell.provider == nil {
		return c.def
	}
	return cell.provider.ctxVal.(T)
}
