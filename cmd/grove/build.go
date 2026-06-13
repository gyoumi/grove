package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runBuild(args []string) error {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	dir := fs.String("dir", ".", "app directory")
	tinygo := fs.Bool("tinygo", false, "compile with TinyGo (much smaller wasm)")
	fs.Parse(args)

	appDir, err := filepath.Abs(*dir)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(filepath.Join(appDir, "dist")); err != nil {
		return err
	}
	if err := buildAll(appDir, true, *tinygo); err != nil {
		return err
	}

	// index.html ships unchanged; the reload script is only injected by the
	// dev server at request time.
	src := filepath.Join(appDir, "index.html")
	if data, err := os.ReadFile(src); err == nil {
		if err := os.WriteFile(filepath.Join(appDir, "dist", "index.html"), data, 0o644); err != nil {
			return err
		}
	}

	wasm := filepath.Join(appDir, "dist", "app.wasm")
	if opt, err := exec.LookPath("wasm-opt"); err == nil {
		before := fileSize(wasm)
		cmd := exec.Command(opt, "-Oz",
			"--enable-bulk-memory", "--enable-nontrapping-float-to-int",
			wasm, "-o", wasm+".opt")
		if out, err := cmd.CombinedOutput(); err != nil {
			fmt.Fprintf(os.Stderr, "grove: wasm-opt skipped: %s\n", strings.TrimSpace(string(out)))
		} else if err := os.Rename(wasm+".opt", wasm); err == nil {
			fmt.Fprintf(os.Stderr, "grove: wasm-opt %s → %s\n", human(before), human(fileSize(wasm)))
		}
	}

	fmt.Printf("dist/app.wasm    %8s   (%s gzipped)\n", human(fileSize(wasm)), human(gzipSize(wasm)))
	css := filepath.Join(appDir, "dist", "styles.css")
	if fileSize(css) > 0 {
		fmt.Printf("dist/styles.css  %8s   (%s gzipped)\n", human(fileSize(css)), human(gzipSize(css)))
	}
	fmt.Println("serve the dist/ directory with any static file server (wasm needs the application/wasm content type, and gzip/brotli is worth enabling)")
	return nil
}
