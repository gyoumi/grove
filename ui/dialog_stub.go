//go:build !(js && wasm)

package ui

import g "github.com/gyoumi/grove"

// Focus management needs a real browser; outside js/wasm (tests) these are
// no-ops.

func trapFocus(*g.Ref[any]) func() { return nil }

func cycleFocus(*g.Ref[any], *g.Event) {}
