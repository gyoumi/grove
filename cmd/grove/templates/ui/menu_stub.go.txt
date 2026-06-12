//go:build !(js && wasm)

package ui

import g "github.com/gyoumi/grove"

// Focus movement needs a real browser; outside js/wasm (tests) it is a
// no-op.
func menuFocusMove(*g.DOMRef, *g.Event, int) {}
