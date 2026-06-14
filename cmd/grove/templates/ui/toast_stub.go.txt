//go:build !(js && wasm)

package ui

// Auto-dismiss needs a real timer; outside the browser (tests) toasts stay
// until dismissed explicitly.
func scheduleToastDismiss(id, ms int) {}
