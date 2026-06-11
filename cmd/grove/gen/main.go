// Command gen copies the ui package sources into cmd/grove/templates so
// `grove add` can embed them. Run from cmd/grove via go generate.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	srcDir := filepath.Join("..", "..", "ui")
	dstDir := filepath.Join("templates", "ui")
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
		// .txt suffix keeps the templates out of the CLI's own build.
		if err := os.WriteFile(filepath.Join(dstDir, name+".txt"), data, 0o644); err != nil {
			fail(err)
		}
		n++
	}
	fmt.Printf("gen: copied %d ui component files into %s\n", n, dstDir)
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "gen:", err)
	os.Exit(1)
}
