package grove

import "fmt"

// currentInstance is the component instance currently rendering. Go wasm in
// the browser is single-threaded, so a package global is safe; hooks must
// not be called from goroutines.
var currentInstance *instance

func current() *instance {
	if currentInstance == nil {
		panic("grove: hooks can only be called while a component renders (inside a function passed to C or C0)")
	}
	return currentInstance
}

func hookMismatch(inst *instance, hook string) {
	panic(fmt.Sprintf(
		"grove: %s called in a different order than the previous render of %s — hooks must run unconditionally, in the same order, on every render (rules of hooks)",
		hook, fnName(inst.node.fnID)))
}

// cheapEqual compares two values with ==, treating incomparable types
// (slices, maps, funcs) as never equal.
func cheapEqual(a, b any) (eq bool) {
	defer func() {
		if recover() != nil {
			eq = false
		}
	}()
	return a == b
}

// depsEqual reports whether two dependency lists are shallowly equal.
func depsEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !cheapEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

type stateCell[T any] struct {
	val T
	set func(T)
}

// UseState declares a piece of component state. It returns the current
// value and a setter; calling the setter re-renders the component (batched
// with other updates in the same tick). Setting an equal value is a no-op.
//
// Like React, the returned value is a snapshot: handlers capture the value
// from their render. Use UseReducer when an update must derive from the
// latest state.
func UseState[T any](initial T) (T, func(T)) {
	return UseStateLazy(func() T { return initial })
}

// UseStateLazy is UseState with the initial value computed only on the
// first render.
func UseStateLazy[T any](initial func() T) (T, func(T)) {
	inst := current()
	i, fresh := inst.slot()
	if fresh {
		c := &stateCell[T]{val: initial()}
		c.set = func(v T) {
			if inst.unmounted {
				return
			}
			if cheapEqual(c.val, v) {
				return
			}
			c.val = v
			inst.markDirty()
		}
		inst.hooks = append(inst.hooks, c)
	}
	c, ok := inst.hooks[i].(*stateCell[T])
	if !ok {
		hookMismatch(inst, "UseState")
	}
	return c.val, c.set
}

type reducerCell[S, A any] struct {
	val      S
	reducer  func(S, A) S
	dispatch func(A)
}

// UseReducer declares state advanced by a reducer: dispatch(action) computes
// the next state from the latest state, so updates compose correctly even
// when several fire in one tick.
func UseReducer[S, A any](reducer func(S, A) S, initial S) (S, func(A)) {
	inst := current()
	i, fresh := inst.slot()
	if fresh {
		c := &reducerCell[S, A]{val: initial, reducer: reducer}
		c.dispatch = func(action A) {
			if inst.unmounted {
				return
			}
			next := c.reducer(c.val, action)
			if cheapEqual(c.val, next) {
				return
			}
			c.val = next
			inst.markDirty()
		}
		inst.hooks = append(inst.hooks, c)
	}
	c, ok := inst.hooks[i].(*reducerCell[S, A])
	if !ok {
		hookMismatch(inst, "UseReducer")
	}
	c.reducer = reducer
	return c.val, c.dispatch
}

type effectCell struct {
	setup   func() func()
	cleanup func()
	deps    []any
	queued  bool
}

// UseEffect runs setup after the render is committed to the DOM. setup may
// return a cleanup function, run before the next setup and on unmount.
//
// deps controls when setup re-runs, mirroring React:
//
//	UseEffect(fn, nil)          // after every render
//	UseEffect(fn, []any{})      // once, on mount
//	UseEffect(fn, []any{a, b})  // whenever a or b changes (shallow ==)
//
// Don't put funcs, slices, or maps in deps — they never compare equal, so
// the effect would run every render.
func UseEffect(setup func() func(), deps []any) {
	inst := current()
	i, fresh := inst.slot()
	if fresh {
		c := &effectCell{setup: setup, deps: deps}
		inst.hooks = append(inst.hooks, c)
		inst.app.queueEffect(inst, c)
		return
	}
	c, ok := inst.hooks[i].(*effectCell)
	if !ok {
		hookMismatch(inst, "UseEffect")
	}
	c.setup = setup
	changed := deps == nil || !depsEqual(c.deps, deps)
	c.deps = deps
	if changed {
		inst.app.queueEffect(inst, c)
	}
}

type memoCell[T any] struct {
	val  T
	deps []any
}

// UseMemo caches the result of compute, recomputing only when deps change
// (same deps semantics as UseEffect).
func UseMemo[T any](compute func() T, deps []any) T {
	inst := current()
	i, fresh := inst.slot()
	if fresh {
		c := &memoCell[T]{val: compute(), deps: deps}
		inst.hooks = append(inst.hooks, c)
		return c.val
	}
	c, ok := inst.hooks[i].(*memoCell[T])
	if !ok {
		hookMismatch(inst, "UseMemo")
	}
	if deps == nil || !depsEqual(c.deps, deps) {
		c.val = compute()
		c.deps = deps
	}
	return c.val
}

// UseCallback caches a function value across renders until deps change.
// Provided for React familiarity; in grove it is exactly UseMemo over fn.
func UseCallback[F any](fn F, deps []any) F {
	return UseMemo(func() F { return fn }, deps)
}

// Ref is a mutable box that persists across renders without triggering
// re-renders when written.
type Ref[T any] struct {
	Current T
}

// UseRef returns the same *Ref on every render of the component.
func UseRef[T any](initial T) *Ref[T] {
	inst := current()
	i, fresh := inst.slot()
	if fresh {
		inst.hooks = append(inst.hooks, &Ref[T]{Current: initial})
	}
	r, ok := inst.hooks[i].(*Ref[T])
	if !ok {
		hookMismatch(inst, "UseRef")
	}
	return r
}
