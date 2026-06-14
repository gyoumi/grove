package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const reloadScript = `<script>(function(){var es=new EventSource("/__grove__/reload");es.onmessage=function(){location.reload()}})()</script>`

func runServe(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.Int("port", 8080, "port to listen on")
	dir := fs.String("dir", ".", "app directory")
	tinygo := fs.Bool("tinygo", false, "compile with TinyGo (much smaller wasm)")
	fs.Parse(args)

	appDir, err := filepath.Abs(*dir)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(appDir, "index.html")); err != nil {
		return fmt.Errorf("%s has no index.html — is this a grove app? (try grove init)", appDir)
	}

	hub := &sseHub{clients: map[chan struct{}]bool{}}

	fmt.Fprintf(os.Stderr, "grove: building %s\n", appDir)
	if err := buildAll(appDir, false, *tinygo); err != nil {
		// Keep serving so the browser reloads once the code compiles again.
		fmt.Fprintf(os.Stderr, "grove: %v\n", err)
	} else {
		printSizes(appDir)
	}

	go watch(appDir, func() {
		start := time.Now()
		if err := buildAll(appDir, false, *tinygo); err != nil {
			fmt.Fprintf(os.Stderr, "grove: %v\n", err)
			return
		}
		fmt.Fprintf(os.Stderr, "grove: rebuilt in %s, reloading\n", time.Since(start).Round(time.Millisecond))
		hub.broadcast()
	})

	mux := http.NewServeMux()
	mux.HandleFunc("/__grove__/reload", hub.serve)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" || path == "/index.html" {
			serveIndex(w, appDir)
			return
		}
		if strings.Contains(path, "..") {
			http.NotFound(w, r)
			return
		}
		// dist/ first (build outputs), then the app dir (static assets).
		for _, base := range []string{filepath.Join(appDir, "dist"), appDir} {
			full := filepath.Join(base, filepath.FromSlash(path))
			if fi, err := os.Stat(full); err == nil && !fi.IsDir() {
				http.ServeFile(w, r, full)
				return
			}
		}
		// SPA fallback: a path with no file extension is a client-side route
		// (/components, /event/42), so serve index.html and let the router
		// handle it. Missing assets (.wasm, .css, …) still 404.
		if filepath.Ext(path) == "" {
			serveIndex(w, appDir)
			return
		}
		http.NotFound(w, r)
	})

	addr := fmt.Sprintf("localhost:%d", *port)
	fmt.Fprintf(os.Stderr, "grove: serving on http://%s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func serveIndex(w http.ResponseWriter, appDir string) {
	data, err := os.ReadFile(filepath.Join(appDir, "index.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	html := string(data)
	if i := strings.LastIndex(strings.ToLower(html), "</body>"); i >= 0 {
		html = html[:i] + reloadScript + html[i:]
	} else {
		html += reloadScript
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func printSizes(appDir string) {
	wasm := filepath.Join(appDir, "dist", "app.wasm")
	fmt.Fprintf(os.Stderr, "grove: app.wasm %s (%s gzipped)\n", human(fileSize(wasm)), human(gzipSize(wasm)))
}

// watch polls for changes to source files and runs rebuild after each
// change settles. Polling keeps the CLI dependency-free and works the same
// on every platform and filesystem.
func watch(appDir string, rebuild func()) {
	last := snapshot(appDir)
	for {
		time.Sleep(400 * time.Millisecond)
		cur := snapshot(appDir)
		if cur != last {
			last = cur
			rebuild()
		}
	}
}

// snapshot fingerprints source file paths, sizes, and mtimes.
func snapshot(appDir string) string {
	var b strings.Builder
	filepath.WalkDir(appDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		name := d.Name()
		if d.IsDir() {
			if name == "dist" || name == ".git" || strings.HasPrefix(name, ".") && path != appDir {
				return filepath.SkipDir
			}
			return nil
		}
		switch filepath.Ext(name) {
		case ".go", ".css", ".html", ".mod":
			if fi, err := d.Info(); err == nil {
				fmt.Fprintf(&b, "%s|%d|%d\n", path, fi.Size(), fi.ModTime().UnixNano())
			}
		}
		return nil
	})
	return b.String()
}

type sseHub struct {
	mu      sync.Mutex
	clients map[chan struct{}]bool
}

func (h *sseHub) broadcast() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for c := range h.clients {
		select {
		case c <- struct{}{}:
		default:
		}
	}
}

func (h *sseHub) serve(w http.ResponseWriter, r *http.Request) {
	fl, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	c := make(chan struct{}, 1)
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.clients, c)
		h.mu.Unlock()
	}()

	fmt.Fprint(w, ": connected\n\n")
	fl.Flush()
	for {
		select {
		case <-c:
			fmt.Fprint(w, "data: reload\n\n")
			fl.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
