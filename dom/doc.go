// Package dom renders grove trees into the real browser DOM via syscall/js.
// It only builds for GOOS=js GOARCH=wasm; this file exists so the package
// is still visible (as documentation) to other platforms' builds.
//
// Typical usage:
//
//	func main() {
//		dom.Mount("#root", g.C0(App))
//	}
//
// Mount blocks forever to keep the Go program (and its event handlers)
// alive.
package dom
