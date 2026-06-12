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
// components. When the value changes, consumers are marked dirty directly,
// so they re-render even when a Memo boundary sits between them and the
// provider.
func (c *Context[T]) Provider(value T, children ...any) *Node {
	n := &Node{kind: kindComponent, fnKey: c.key}
	n.fn = func() *Node {
		inst := current()
		changed := inst.hasCtx && !cheapEqual(inst.ctxVal, value)
		inst.ctxKey = c.key
		inst.ctxVal = value
		inst.hasCtx = true
		if changed {
			for consumer := range inst.consumers {
				consumer.markDirty()
			}
		}
		return Fragment(children...)
	}
	// fnID stays 0 for all providers; fnKey distinguishes contexts during
	// reconciliation.
	return n
}

type ctxCell struct {
	owner    *instance
	provider *instance
}

// release deregisters the consumer; called on unmount.
func (cell *ctxCell) release() {
	if cell.provider != nil {
		delete(cell.provider.consumers, cell.owner)
	}
}

// UseContext returns the value of the nearest enclosing Provider for c, or
// the context's default when there is none.
func UseContext[T any](c *Context[T]) T {
	inst := current()
	i, fresh := inst.slot()
	if fresh {
		cell := &ctxCell{owner: inst}
		for a := inst.parentInst; a != nil; a = a.parentInst {
			if a.ctxKey == c.key {
				cell.provider = a
				if a.consumers == nil {
					a.consumers = map[*instance]struct{}{}
				}
				a.consumers[inst] = struct{}{}
				break
			}
		}
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
