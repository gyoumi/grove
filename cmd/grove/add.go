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
	"accordion":       {"accordion.go", "icon.go"},
	"alert":           {"alert.go"},
	"alert-dialog":    {"alert_dialog.go", "dialog_js.go", "dialog_stub.go"},
	"aspect-ratio":    {"aspect_ratio.go"},
	"avatar":          {"avatar.go"},
	"badge":           {"badge.go"},
	"breadcrumb":      {"breadcrumb.go", "icon.go"},
	"button":          {"button.go"},
	"button-group":    {"button_group.go"},
	"calendar":        {"calendar.go", "icon.go"},
	"card":            {"card.go"},
	"carousel":        {"carousel.go", "icon.go"},
	"chart":           {"chart.go"},
	"checkbox":        {"checkbox.go"},
	"collapsible":     {"collapsible.go"},
	"combobox":        {"combobox.go", "command.go", "icon.go", "popover.go", "popover_js.go", "popover_stub.go"},
	"command":         {"command.go", "icon.go"},
	"context-menu":    {"context_menu.go", "dropdown.go", "menu_js.go", "menu_stub.go", "popover.go", "popover_js.go", "popover_stub.go"},
	"data-table":      {"data_table.go", "table.go", "input.go", "button.go", "icon.go"},
	"date-picker":     {"date_picker.go", "calendar.go", "icon.go", "popover.go", "popover_js.go", "popover_stub.go"},
	"dialog":          {"dialog.go", "dialog_js.go", "dialog_stub.go"},
	"drawer":          {"drawer.go", "dialog_js.go", "dialog_stub.go"},
	"dropdown":        {"dropdown.go", "menu_js.go", "menu_stub.go", "popover.go", "popover_js.go", "popover_stub.go"},
	"empty":           {"empty.go"},
	"field":           {"field.go"},
	"form":            {"form.go", "field.go", "input.go"},
	"hover-card":      {"hover_card.go"},
	"icon":            {"icon.go"},
	"input":           {"input.go"},
	"input-group":     {"input_group.go"},
	"input-otp":       {"input_otp.go"},
	"item":            {"item.go"},
	"kbd":             {"kbd.go"},
	"label":           {"label.go"},
	"menubar":         {"menubar.go", "dropdown.go", "menu_js.go", "menu_stub.go", "popover.go", "popover_js.go", "popover_stub.go"},
	"native-select":   {"native_select.go", "icon.go"},
	"navigation-menu": {"navigation_menu.go", "icon.go"},
	"pagination":      {"pagination.go", "icon.go"},
	"popover":         {"popover.go", "popover_js.go", "popover_stub.go"},
	"progress":        {"progress.go"},
	"radio-group":     {"radio_group.go"},
	"resizable":       {"resizable.go", "resizable_js.go", "resizable_stub.go"},
	"scroll-area":     {"scroll_area.go"},
	"select":          {"select.go", "icon.go", "popover.go", "popover_js.go", "popover_stub.go"},
	"separator":       {"separator.go"},
	"sheet":           {"sheet.go", "dialog_js.go", "dialog_stub.go"},
	"sidebar":         {"sidebar.go", "icon.go"},
	"skeleton":        {"skeleton.go"},
	"slider":          {"slider.go"},
	"spinner":         {"spinner.go", "icon.go"},
	"switch":          {"switch.go"},
	"table":           {"table.go"},
	"tabs":            {"tabs.go"},
	"textarea":        {"textarea.go"},
	"time-picker":     {"time_picker.go", "icon.go", "popover.go", "popover_js.go", "popover_stub.go"},
	"toast":           {"toast.go", "toast_js.go", "toast_stub.go", "icon.go"},
	"toggle":          {"toggle.go"},
	"toggle-group":    {"toggle_group.go", "toggle.go"},
	"tooltip":         {"tooltip.go"},
	"typography":      {"typography.go"},
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
