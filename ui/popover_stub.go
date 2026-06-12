//go:build !(js && wasm)

package ui

import g "github.com/gyoumi/grove"

// Collision measurement needs a real browser; outside js/wasm (tests) the
// requested side is kept.
func measureFlip(*g.DOMRef, PopoverSide) (PopoverSide, bool) { return "", false }
