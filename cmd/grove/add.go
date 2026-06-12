package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:generate go run ./gen

//go:embed templates
var templates embed.FS

// components maps a component name to the source files it needs. The files
// are byte-for-byte the grove ui package sources (see go:generate above),
// declared as package ui so they drop into an app's ui/ directory.
var components = map[string][]string{
	"alert":     {"alert.go"},
	"avatar":    {"avatar.go"},
	"badge":     {"badge.go"},
	"button":    {"button.go"},
	"card":      {"card.go"},
	"checkbox":  {"checkbox.go"},
	"dialog":    {"dialog.go", "dialog_js.go", "dialog_stub.go"},
	"input":     {"input.go"},
	"label":     {"label.go"},
	"separator": {"separator.go"},
	"switch":    {"switch.go"},
	"tooltip":   {"tooltip.go"},
}

func componentNames() []string {
	names := make([]string, 0, len(components))
	for n := range components {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func runAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	list := fs.Bool("list", false, "list available components")
	dir := fs.String("dir", ".", "app directory")
	force := fs.Bool("force", false, "overwrite existing files")
	names := parseMixed(fs, args)

	if *list {
		fmt.Println(strings.Join(componentNames(), "\n"))
		return nil
	}
	if len(names) == 0 {
		return fmt.Errorf("usage: grove add <component> — available: %s", strings.Join(componentNames(), ", "))
	}

	uiDir := filepath.Join(*dir, "ui")
	if err := os.MkdirAll(uiDir, 0o755); err != nil {
		return err
	}

	for _, name := range names {
		files, ok := components[strings.ToLower(name)]
		if !ok {
			return fmt.Errorf("unknown component %q — available: %s", name, strings.Join(componentNames(), ", "))
		}
		for _, f := range files {
			data, err := templates.ReadFile("templates/ui/" + f + ".txt")
			if err != nil {
				return fmt.Errorf("template for %s missing (rebuild the CLI after go generate): %w", f, err)
			}
			dest := filepath.Join(uiDir, f)
			if _, err := os.Stat(dest); err == nil && !*force {
				fmt.Fprintf(os.Stderr, "grove: %s exists, skipping (use -force to overwrite)\n", dest)
				continue
			}
			if err := os.WriteFile(dest, data, 0o644); err != nil {
				return err
			}
			fmt.Printf("added %s\n", dest)
		}
	}
	fmt.Println("components are yours now — edit them freely (they import github.com/gyoumi/grove/style)")
	return nil
}
