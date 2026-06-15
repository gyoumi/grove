// Command gen copies the ui and gallery sources into cmd/grove/templates so
// `grove add` and `grove init` can embed them. Run from cmd/grove via go
// generate.
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	uiN := copyTree(filepath.Join("..", "..", "ui"), filepath.Join("templates", "ui"), nil)
	// The gallery imports grove's own ui package; in a scaffolded app that
	// becomes the app's vendored ui, so swap the import for a placeholder
	// that grove init fills in with the new module path.
	galleryN := copyTree(filepath.Join("..", "..", "gallery"), filepath.Join("templates", "gallery"), func(b []byte) []byte {
		return bytes.ReplaceAll(b, []byte("github.com/gyoumi/grove/ui"), []byte("{{MODULE}}/ui"))
	})
	fmt.Printf("gen: copied %d ui + %d gallery files into templates\n", uiN, galleryN)
}

// copyTree copies every non-test .go file from srcDir into dstDir with a .txt
// suffix (which keeps the templates out of the CLI's own build), optionally
// transforming the bytes first.
func copyTree(srcDir, dstDir string, transform func([]byte) []byte) int {
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		fail(err)
	}
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		fail(err)
	}
	n := 0
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") || name == "doc.go" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, name))
		if err != nil {
			fail(err)
		}
		if transform != nil {
			data = transform(data)
		}
		if err := os.WriteFile(filepath.Join(dstDir, name+".txt"), data, 0o644); err != nil {
			fail(err)
		}
		n++
	}
	return n
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "gen:", err)
	os.Exit(1)
}
