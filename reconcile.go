package grove

import (
	"github.com/gyoumi/grove/renderer"
)

// sameType reports whether old can be patched into new in place (as opposed
// to unmount + remount).
func sameType(old, new *Node) bool {
	if old.kind != new.kind || old.key != new.key {
		return false
	}
	switch old.kind {
	case kindElement:
		return old.tag == new.tag
	case kindComponent:
		return old.fnID == new.fnID && cheapEqual(old.fnKey, new.fnKey)
	}
	return true
}

// firstDOM returns the first real DOM node rendered by n, or nil if n
// currently renders nothing.
func firstDOM(n *Node) renderer.Node {
	if n == nil {
		return nil
	}
	switch n.kind {
	case kindElement, kindText:
		return n.dom
	case kindComponent:
		return firstDOM(n.rendered)
	case kindFragment:
		for _, c := range n.children {
			if d := firstDOM(c); d != nil {
				return d
			}
		}
		// kindPortal contributes no DOM here — its children live in the
		// portal root — so it falls through to nil.
	}
	return nil
}

// moveDOM re-inserts every top-level DOM node of n before anchor,
// preserving their relative order.
func (a *App) moveDOM(n *Node, parentDOM, anchor renderer.Node) {
	switch n.kind {
	case kindElement, kindText:
		a.r.InsertBefore(parentDOM, n.dom, anchor)
	case kindComponent:
		if n.rendered != nil {
			a.moveDOM(n.rendered, parentDOM, anchor)
		}
	case kindFragment:
		for _, c := range n.children {
			a.moveDOM(c, parentDOM, anchor)
		}
	}
}

// mount creates real DOM for n and inserts it under parentDOM before anchor.
func (a *App) mount(n *Node, parent *Node, pinst *instance, parentDOM, anchor renderer.Node) {
	n.parent = parent
	switch n.kind {
	case kindText:
		n.dom = a.r.CreateText(n.text)
		a.r.InsertBefore(parentDOM, n.dom, anchor)

	case kindElement:
		a.nextID++
		n.id = a.nextID
		a.byID[n.id] = n
		n.dom = a.r.CreateElement(n.tag, n.id)
		for k, v := range n.attrs {
			a.r.SetAttr(n.dom, k, v)
		}
		for k, v := range n.props {
			a.r.SetProp(n.dom, k, v)
		}
		for typ := range n.events {
			a.r.Listen(typ)
		}
		for _, c := range n.children {
			a.mount(c, n, pinst, n.dom, nil)
		}
		a.r.InsertBefore(parentDOM, n.dom, anchor)
		for _, ref := range n.refs {
			ref.Current = n.dom
		}

	case kindFragment:
		for _, c := range n.children {
			a.mount(c, n, pinst, parentDOM, anchor)
		}

	case kindPortal:
		// Children mount into the portal root (the app container), not the
		// caller's parentDOM, so they escape transformed ancestors.
		for _, c := range n.children {
			a.mount(c, n, pinst, a.r.PortalRoot(), nil)
		}

	case kindComponent:
		inst := newInstance(a, n, pinst)
		n.inst = inst
		rendered := inst.render()
		if rendered != nil {
			a.mount(rendered, n, inst, parentDOM, anchor)
		}
		n.rendered = rendered
	}
}

// patch reconciles old into new. It returns the node that remains mounted
// (always new unless new is nil). parentDOM/anchor describe where the
// content lives, for mounts and replacements.
func (a *App) patch(old, new *Node, parent *Node, pinst *instance, parentDOM, anchor renderer.Node) *Node {
	if old == nil && new == nil {
		return nil
	}
	if old == nil {
		a.mount(new, parent, pinst, parentDOM, anchor)
		return new
	}
	if new == nil {
		a.unmount(old, parentDOM, true)
		return nil
	}
	if !sameType(old, new) {
		// Replace: mount the new node at the old one's position, then drop
		// the old. firstDOM(old) is a more precise anchor than the caller's
		// (the old content is still in the DOM at this point).
		repAnchor := firstDOM(old)
		if repAnchor == nil {
			repAnchor = anchor
		}
		a.mount(new, parent, pinst, parentDOM, repAnchor)
		a.unmount(old, parentDOM, true)
		return new
	}

	new.parent = parent
	switch new.kind {
	case kindText:
		new.dom = old.dom
		if new.text != old.text {
			a.r.SetText(new.dom, new.text)
		}

	case kindElement:
		new.dom = old.dom
		new.id = old.id
		a.byID[new.id] = new
		a.patchElement(old, new)
		new.children = a.patchChildren(old.children, new.children, new, pinst, new.dom, nil)

	case kindFragment:
		new.children = a.patchChildren(old.children, new.children, new, pinst, parentDOM, anchor)

	case kindPortal:
		new.children = a.patchChildren(old.children, new.children, new, pinst, a.r.PortalRoot(), nil)

	case kindComponent:
		inst := old.inst
		new.inst = inst
		inst.node = new
		if new.memo && !inst.dirty && memoEqual(old, new) {
			// Skip the render: adopt the old subtree unchanged.
			new.rendered = old.rendered
			if new.rendered != nil {
				new.rendered.parent = new
			}
			break
		}
		rendered := inst.render()
		new.rendered = a.patch(old.rendered, rendered, new, inst, parentDOM, anchor)
	}
	return new
}

func memoEqual(old, new *Node) bool {
	if new.eq != nil {
		return new.eq(old.fnProps, new.fnProps)
	}
	return cheapEqual(old.fnProps, new.fnProps)
}

func (a *App) patchElement(old, new *Node) {
	dom := new.dom
	for k, v := range new.attrs {
		if ov, ok := old.attrs[k]; !ok || ov != v {
			a.r.SetAttr(dom, k, v)
		}
	}
	for k := range old.attrs {
		if _, ok := new.attrs[k]; !ok {
			a.r.RemoveAttr(dom, k)
		}
	}
	for k, v := range new.props {
		// value and checked are always re-synced: the live DOM drifts as
		// the user types/clicks, and controlled inputs must win.
		if k == "value" || k == "checked" || !cheapEqual(old.props[k], v) {
			a.r.SetProp(dom, k, v)
		}
	}
	for k := range old.props {
		if _, ok := new.props[k]; !ok {
			a.r.SetProp(dom, k, nil)
		}
	}
	for typ := range new.events {
		if _, had := old.events[typ]; !had {
			a.r.Listen(typ)
		}
	}
	// Handlers themselves need no DOM work: dispatch reads new.events.

	kept := make(map[*DOMRef]bool, len(new.refs))
	for _, ref := range new.refs {
		kept[ref] = true
		ref.Current = dom
	}
	for _, ref := range old.refs {
		if !kept[ref] {
			ref.Current = nil
		}
	}
}

// patchChildren reconciles two child lists. tail is the DOM node that
// follows the list in parentDOM (nil means the list ends the parent).
func (a *App) patchChildren(oldKids, newKids []*Node, parent *Node, pinst *instance, parentDOM, tail renderer.Node) []*Node {
	// Index old children: keyed by key, unkeyed by position.
	var keyed map[string]*Node
	var unkeyed []*Node
	for _, o := range oldKids {
		if o.key != "" {
			if keyed == nil {
				keyed = map[string]*Node{}
			}
			keyed[o.key] = o
		} else {
			unkeyed = append(unkeyed, o)
		}
	}
	oldIndex := make(map[*Node]int, len(oldKids))
	for i, o := range oldKids {
		oldIndex[o] = i
	}

	// Match new children to reusable old ones (same key/position AND same
	// type — a type change means remount, handled as no-match).
	pairs := make([]*Node, len(newKids))
	matched := make(map[*Node]bool, len(oldKids))
	ui := 0
	for i, nk := range newKids {
		if nk.key != "" {
			if o := keyed[nk.key]; o != nil && sameType(o, nk) {
				pairs[i] = o
				matched[o] = true
				delete(keyed, nk.key)
			}
		} else if ui < len(unkeyed) {
			o := unkeyed[ui]
			ui++
			if sameType(o, nk) {
				pairs[i] = o
				matched[o] = true
			}
		}
	}

	// Drop old children that found no match.
	for _, o := range oldKids {
		if !matched[o] {
			a.unmount(o, parentDOM, true)
		}
	}

	// Decide which kept nodes can stay put: the longest run of kept nodes
	// whose old indices already increase in the new order needs no DOM
	// moves, so only the rest move. That run is the longest increasing
	// subsequence of the kept nodes' old indices — minimizing moves rather
	// than the greedy "stay if not smaller than the next" heuristic, which
	// could move N-1 nodes for a single-element rotation.
	var keptNew, keptOld []int
	for i, o := range pairs {
		if o != nil {
			keptNew = append(keptNew, i)
			keptOld = append(keptOld, oldIndex[o])
		}
	}
	stay := make([]bool, len(newKids))
	for _, k := range longestIncreasingSubsequence(keptOld) {
		stay[keptNew[k]] = true
	}

	// Place right-to-left so the anchor for each child is always known.
	anchor := tail
	for i := len(newKids) - 1; i >= 0; i-- {
		nk := newKids[i]
		o := pairs[i]
		switch {
		case o == nil:
			a.mount(nk, parent, pinst, parentDOM, anchor)
		case stay[i]:
			a.patch(o, nk, parent, pinst, parentDOM, anchor)
		default:
			a.patch(o, nk, parent, pinst, parentDOM, anchor)
			a.moveDOM(nk, parentDOM, anchor)
		}
		if fd := firstDOM(nk); fd != nil {
			anchor = fd
		}
	}
	return newKids
}

// longestIncreasingSubsequence returns the indices into nums of a longest
// strictly increasing subsequence (old indices are distinct). Standard
// patience-sorting in O(n log n) with predecessor links for reconstruction.
func longestIncreasingSubsequence(nums []int) []int {
	if len(nums) == 0 {
		return nil
	}
	tails := make([]int, 0, len(nums)) // tails[k]: index of the smallest tail of an increasing run of length k+1
	prev := make([]int, len(nums))     // prev[i]: predecessor of i in its run
	for i, x := range nums {
		prev[i] = -1
		lo, hi := 0, len(tails)
		for lo < hi {
			mid := (lo + hi) / 2
			if nums[tails[mid]] < x {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		if lo > 0 {
			prev[i] = tails[lo-1]
		}
		if lo == len(tails) {
			tails = append(tails, i)
		} else {
			tails[lo] = i
		}
	}
	out := make([]int, len(tails))
	k := tails[len(tails)-1]
	for j := len(tails) - 1; j >= 0; j-- {
		out[j] = k
		k = prev[k]
	}
	return out
}

// unmount tears down n: effect cleanups (children first, like React), event
// registry entries, refs, and — at the top level — DOM removal.
func (a *App) unmount(n *Node, parentDOM renderer.Node, removeDOM bool) {
	switch n.kind {
	case kindText:
		if removeDOM {
			a.r.Remove(parentDOM, n.dom)
		}

	case kindElement:
		for _, c := range n.children {
			a.unmount(c, n.dom, false)
		}
		for _, ref := range n.refs {
			ref.Current = nil
		}
		delete(a.byID, n.id)
		if removeDOM {
			a.r.Remove(parentDOM, n.dom)
		}

	case kindFragment:
		for _, c := range n.children {
			a.unmount(c, parentDOM, removeDOM)
		}

	case kindPortal:
		// Portal children live in the portal root, not in any unmounting
		// ancestor's DOM subtree, so they must always be removed explicitly.
		for _, c := range n.children {
			a.unmount(c, a.r.PortalRoot(), true)
		}

	case kindComponent:
		inst := n.inst
		inst.unmounted = true
		if n.rendered != nil {
			a.unmount(n.rendered, parentDOM, removeDOM)
		}
		inst.runCleanups()
	}
}
