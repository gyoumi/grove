package grove

// instance is the persistent identity of a mounted component: its hook
// slots, its place in the tree, and its scheduling state.
type instance struct {
	app        *App
	node       *Node // current component vnode
	parentInst *instance
	depth      int

	hooks     []any
	cursor    int
	hookCount int // -1 until the first render fixes the count

	dirty     bool
	unmounted bool

	// context provider data (set every render by Context.Provider)
	ctxKey    *ctxKey
	ctxVal    any
	hasCtx    bool
	consumers map[*instance]struct{} // instances reading this provider
}

func newInstance(a *App, n *Node, parent *instance) *instance {
	depth := 0
	if parent != nil {
		depth = parent.depth + 1
	}
	return &instance{app: a, node: n, parentInst: parent, depth: depth, hookCount: -1}
}

// slot advances the hook cursor and reports whether this is a brand-new
// slot (first render reaching this hook call).
func (inst *instance) slot() (int, bool) {
	i := inst.cursor
	inst.cursor++
	if i > len(inst.hooks) {
		hookMismatch(inst, "a hook")
	}
	return i, i == len(inst.hooks)
}

// render invokes the component function with this instance as the hook
// target and returns the produced tree (which may be nil).
func (inst *instance) render() *Node {
	prev := currentInstance
	currentInstance = inst
	inst.cursor = 0
	defer func() { currentInstance = prev }()

	out := inst.node.fn()

	if inst.hookCount == -1 {
		inst.hookCount = inst.cursor
	} else if inst.cursor != inst.hookCount {
		hookMismatch(inst, "a hook")
	}
	inst.dirty = false
	return out
}

func (inst *instance) markDirty() {
	if inst.dirty || inst.unmounted {
		return
	}
	inst.dirty = true
	inst.app.enqueue(inst)
}

// runCleanups runs effect cleanups in hook order and releases context
// subscriptions; used on unmount.
func (inst *instance) runCleanups() {
	for _, h := range inst.hooks {
		switch cell := h.(type) {
		case *effectCell:
			if cell.cleanup != nil {
				cell.cleanup()
				cell.cleanup = nil
			}
		case *ctxCell:
			cell.release()
		}
	}
}
