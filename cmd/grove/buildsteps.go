package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// buildWasm compiles the app in dir to dist/app.wasm.
func buildWasm(dir string, release bool) error {
	args := []string{"build", "-o", filepath.Join("dist", "app.wasm")}
	if release {
		args = append(args, "-trimpath", "-ldflags=-s -w")
	}
	args = append(args, ".")
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build failed:\n%s", strings.TrimSpace(string(out)))
	}
	return nil
}

// copyWasmExec places the Go runtime's JS glue into dist/.
func copyWasmExec(dir string) error {
	goroot, err := exec.Command("go", "env", "GOROOT").Output()
	if err != nil {
		return fmt.Errorf("go env GOROOT: %w", err)
	}
	root := strings.TrimSpace(string(goroot))
	candidates := []string{
		filepath.Join(root, "lib", "wasm", "wasm_exec.js"),  // Go ≥ 1.24
		filepath.Join(root, "misc", "wasm", "wasm_exec.js"), // older Go
	}
	for _, c := range candidates {
		data, err := os.ReadFile(c)
		if err == nil {
			return os.WriteFile(filepath.Join(dir, "dist", "wasm_exec.js"), data, 0o644)
		}
	}
	return fmt.Errorf("wasm_exec.js not found under %s", root)
}

// buildCSS compiles styles/input.css to dist/styles.css with the Tailwind
// standalone CLI when available; otherwise it copies the file through so
// the app still loads (without utility classes).
func buildCSS(dir string, minify bool) error {
	input := filepath.Join(dir, "styles", "input.css")
	if _, err := os.Stat(input); err != nil {
		return nil // no stylesheet in this app
	}
	output := filepath.Join(dir, "dist", "styles.css")

	bin, err := findTailwind()
	if err != nil {
		fmt.Fprintf(os.Stderr, "grove: tailwind unavailable (%v); copying styles/input.css as-is\n", err)
		data, rerr := os.ReadFile(input)
		if rerr != nil {
			return rerr
		}
		return os.WriteFile(output, data, 0o644)
	}

	args := []string{"-i", filepath.Join("styles", "input.css"), "-o", filepath.Join("dist", "styles.css")}
	if minify {
		args = append(args, "--minify")
	}
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tailwind failed:\n%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func ensureDist(dir string) error {
	return os.MkdirAll(filepath.Join(dir, "dist"), 0o755)
}

// buildAll runs the full pipeline for serve/build.
func buildAll(dir string, release bool) error {
	if err := ensureDist(dir); err != nil {
		return err
	}
	if err := copyWasmExec(dir); err != nil {
		return err
	}
	if err := buildWasm(dir, release); err != nil {
		return err
	}
	return buildCSS(dir, release)
}

func fileSize(path string) int64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return fi.Size()
}

func gzipSize(path string) int64 {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	var n countingWriter
	zw := gzip.NewWriter(&n)
	if _, err := io.Copy(zw, f); err != nil {
		return 0
	}
	zw.Close()
	return n.n
}

type countingWriter struct{ n int64 }

func (w *countingWriter) Write(p []byte) (int, error) {
	w.n += int64(len(p))
	return len(p), nil
}

func human(n int64) string {
	switch {
	case n >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	case n >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(n)/(1<<10))
	}
	return fmt.Sprintf("%d B", n)
}
