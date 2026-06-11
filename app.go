package grove

import (
	"sort"

	"github.com/gyoumi/grove/renderer"
)

// maxFlushIterations bounds render→effect→render cascades within one flush
// before grove assumes an unconditional state update loop.
const maxFlushIterations = 64

type effectRun struct {
	inst *instance
	cell *effectCell
}

// App is a mounted grove tree: the bridge between the reconciler and a
// Renderer, plus the update scheduler.
type App struct {
	r         renderer.Renderer
	container renderer.Node
	root      *Node

	byID   map[int]*Node
	nextID int

	dirtyQ    []*instance
	scheduled bool
	inFlush   bool
	effects   []effectRun
}

// Mount renders root into container using the given renderer and returns
// the running App. Browser apps normally call dom.Mount instead, which
// wraps this and blocks forever.
func Mount(r renderer.Renderer, container renderer.Node, root *Node) *App {
	a := &App{r: r, container: container, root: root, byID: map[int]*Node{}}
	r.SetDispatch(a.dispatch)
	a.mount(root, nil, nil, container, nil)
	a.flush() // run mount effects (and any renders they trigger)
	return a
}

func (a *App) enqueue(inst *instance) {
	a.dirtyQ = append(a.dirtyQ, inst)
	if !a.scheduled && !a.inFlush {
		a.scheduled = true
		a.r.Schedule(a.flush)
	}
}

func (a *App) queueEffect(inst *instance, c *effectCell) {
	if c.queued {
		return
	}
	c.queued = true
	a.effects = append(a.effects, effectRun{inst, c})
}

func (a *App) flush() {
	a.scheduled = false
	if a.inFlush {
		return
	}
	a.inFlush = true
	defer func() { a.inFlush = false }()

	for iter := 0; len(a.dirtyQ) > 0 || len(a.effects) > 0; iter++ {
		if iter >= maxFlushIterations {
			panic("grove: updates did not settle after " +
				"many render passes — a component is probably setting state unconditionally during render or in an effect without deps")
		}
		batch := a.dirtyQ
		a.dirtyQ = nil
		// Parents before children: a parent re-render also re-renders its
		// children, clearing their dirty flags so they're skipped below.
		sort.SliceStable(batch, func(i, j int) bool { return batch[i].depth < batch[j].depth })
		for _, inst := range batch {
			a.renderInstance(inst)
		}
		a.runEffects()
	}
}

func (a *App) renderInstance(inst *instance) {
	if inst.unmounted || !inst.dirty {
		return
	}
	n := inst.node
	parentDOM := a.hostParentDOM(n)
	anchor := a.nextAnchor(n)
	out := inst.render()
	n.rendered = a.patch(n.rendered, out, n, inst, parentDOM, anchor)
}

func (a *App) runEffects() {
	q := a.effects
	a.effects = nil
	for _, er := range q {
		er.cell.queued = false
		if er.inst.unmounted {
			continue
		}
		if er.cell.cleanup != nil {
			er.cell.cleanup()
			er.cell.cleanup = nil
		}
		if cleanup := er.cell.setup(); cleanup != nil {
			er.cell.cleanup = cleanup
		}
	}
}

// dispatch is installed into the renderer; it receives the nearest grove
// element for a platform event and bubbles it up the virtual tree.
func (a *App) dispatch(id int, event string, raw any) {
	n := a.byID[id]
	if n == nil {
		return
	}
	ev := &Event{Type: event, Raw: raw, ops: a.r}
	for cur := n; cur != nil && !ev.stopped; cur = cur.parent {
		if cur.kind == kindElement {
			if h := cur.events[event]; h != nil {
				h(ev)
			}
		}
	}
}

// hostParentDOM finds the real DOM node a component's output lives under.
func (a *App) hostParentDOM(n *Node) renderer.Node {
	for p := n.parent; p != nil; p = p.parent {
		if p.kind == kindElement {
			return p.dom
		}
	}
	return a.container
}

// nextAnchor finds the DOM node that follows n's content in its host
// parent, i.e. the InsertBefore anchor for re-rendering n in place.
func (a *App) nextAnchor(n *Node) renderer.Node {
	cur := n
	for {
		p := cur.parent
		if p == nil {
			return nil
		}
		if p.kind == kindComponent {
			cur = p
			continue
		}
		kids := p.children
		idx := -1
		for i, k := range kids {
			if k == cur {
				idx = i
				break
			}
		}
		if idx >= 0 {
			for _, sib := range kids[idx+1:] {
				if d := firstDOM(sib); d != nil {
					return d
				}
			}
		}
		if p.kind == kindElement {
			return nil // end of the host element's children
		}
		cur = p // fragment: keep ascending
	}
}
