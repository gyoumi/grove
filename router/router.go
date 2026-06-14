// Package router is a small client-side router for grove apps. In the
// browser it drives the real URL path with the History API (/event/42, no
// #), so links are clean and back/forward work; outside the browser (tests)
// it keeps the path in memory with the same API. Hosting a built app needs
// a fallback that serves index.html for unknown paths — grove serve does
// this for you.
//
//	router.Routes(
//	    router.Route{Pattern: "/", Render: home},
//	    router.Route{Pattern: "/event/:id", Render: showEvent},
//	    router.Route{Pattern: "*", Render: notFound},
//	)
package router

import (
	"strings"

	g "github.com/gyoumi/grove"
)

// Params holds the values captured by :name segments in a route pattern.
type Params map[string]string

// Route pairs a pattern with the component to render for it. Patterns are
// segment-wise: "/" matches the root, ":name" captures one segment, and a
// trailing "*" (or the bare pattern "*") matches anything.
type Route struct {
	Pattern string
	Render  func(Params) *g.Node
}

// Routes renders the first route whose pattern matches the current path,
// re-rendering whenever the path changes. Unmatched paths render nothing —
// add a "*" route for a not-found page.
func Routes(routes ...Route) *g.Node {
	return g.C(routesComponent, routes)
}

func routesComponent(routes []Route) *g.Node {
	path, setPath := g.UseState(Path())
	g.UseEffect(func() func() {
		return src.subscribe(func() { setPath(src.path()) })
	}, []any{})

	for _, r := range routes {
		if params, ok := match(r.Pattern, path); ok {
			return r.Render(params)
		}
	}
	return nil
}

// Path returns the current route path, "/" at minimum.
func Path() string { return src.path() }

// Navigate switches to the given path, e.g. router.Navigate("/event/42").
func Navigate(path string) { src.navigate(normalize(path)) }

// Link renders an anchor that navigates without a page load. The href is a
// clean path (/event/42), so links are shareable and right-click "open in
// new tab" works; the click is intercepted for client-side navigation.
// Children and options are passed through to the anchor element.
func Link(to string, args ...any) *g.Node {
	all := []any{
		g.Href(normalize(to)),
		g.OnClick(func(e *g.Event) {
			e.PreventDefault()
			Navigate(to)
		}),
	}
	return g.A(append(all, args...)...)
}

func normalize(p string) string {
	if p == "" {
		return "/"
	}
	if !strings.HasPrefix(p, "/") {
		return "/" + p
	}
	return p
}

// match reports whether pattern matches path, returning captured params.
func match(pattern, path string) (Params, bool) {
	if pattern == "*" {
		return Params{}, true
	}
	pat := splitPath(pattern)
	got := splitPath(path)
	params := Params{}
	for i, seg := range pat {
		if seg == "*" && i == len(pat)-1 {
			return params, true
		}
		if i >= len(got) {
			return nil, false
		}
		switch {
		case strings.HasPrefix(seg, ":"):
			params[seg[1:]] = got[i]
		case seg != got[i]:
			return nil, false
		}
	}
	if len(got) != len(pat) {
		return nil, false
	}
	return params, true
}

func splitPath(p string) []string {
	p = strings.Trim(normalize(p), "/")
	if p == "" {
		return nil
	}
	return strings.Split(p, "/")
}

// source abstracts where the path lives (URL hash vs memory).
type source interface {
	path() string
	navigate(p string)
	subscribe(fn func()) (unsubscribe func())
}

type subscribers struct {
	m      map[int]func()
	nextID int
}

func (s *subscribers) add(fn func()) func() {
	if s.m == nil {
		s.m = map[int]func(){}
	}
	s.nextID++
	id := s.nextID
	s.m[id] = fn
	return func() { delete(s.m, id) }
}

func (s *subscribers) notify() {
	for _, fn := range s.m {
		fn()
	}
}
