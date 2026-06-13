//go:build !(js && wasm)

package ui

import g "github.com/gyoumi/grove"

// Collision measurement needs a real browser; outside js/wasm (tests) the
// requested side is kept, the panel never shifts, and there is no viewport
// to resize.
func measureFlip(*g.DOMRef, PopoverSide) (PopoverSide, bool) { return "", false }
func measureShift(*g.DOMRef, PopoverSide, int) (int, bool)   { return 0, false }
func onViewportResize(func()) func()                         { return func() {} }
