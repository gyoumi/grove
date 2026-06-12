//go:build !(js && wasm)

package router

// memSource keeps the path in memory; used outside the browser (tests).
type memSource struct {
	p    string
	subs subscribers
}

var src source = &memSource{p: "/"}

func (m *memSource) path() string { return m.p }

func (m *memSource) navigate(p string) {
	if m.p == p {
		return
	}
	m.p = p
	m.subs.notify()
}

func (m *memSource) subscribe(fn func()) func() {
	return m.subs.add(fn)
}
