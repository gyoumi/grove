package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	grovePath := fs.String("grove", "", "path to a local grove checkout (adds a replace directive)")
	pos := parseMixed(fs, args)
	if len(pos) != 1 {
		return fmt.Errorf("usage: grove init <app-name> [--grove <path>]")
	}
	name := pos[0]
	if _, err := os.Stat(name); err == nil {
		return fmt.Errorf("directory %s already exists", name)
	}

	gomod := fmt.Sprintf("module %s\n\ngo 1.25\n\nrequire github.com/gyoumi/grove v0.0.0\n", name)
	if *grovePath != "" {
		abs, err := filepath.Abs(*grovePath)
		if err != nil {
			return err
		}
		gomod += fmt.Sprintf("\nreplace github.com/gyoumi/grove => %s\n", abs)
	}

	files := map[string]string{
		"go.mod":           gomod,
		"main.go":          initMainGo,
		"index.html":       strings.ReplaceAll(initIndexHTML, "{{name}}", name),
		"styles/input.css": themeCSS,
		".gitignore":       "dist/\n",
	}
	for rel, content := range files {
		path := filepath.Join(name, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
	}

	// Resolve the module graph now so grove serve works immediately.
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = name
	if out, err := tidy.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "grove: go mod tidy failed:\n%s\n", strings.TrimSpace(string(out)))
		if *grovePath == "" {
			fmt.Fprintln(os.Stderr, "grove: run `go get github.com/gyoumi/grove@latest` inside the app, then `go mod tidy`")
		}
	}

	fmt.Printf("created %s/\n\n", name)
	fmt.Printf("  cd %s\n", name)
	if *grovePath == "" {
		fmt.Println("  go get github.com/gyoumi/grove@latest")
	}
	fmt.Println("  grove serve")
	return nil
}

const initMainGo = `package main

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/dom"
)

func App() *g.Node {
	count, setCount := g.UseState(0)
	return g.Div(g.Class("flex min-h-svh flex-col items-center justify-center gap-4 bg-background text-foreground"),
		g.H1(g.Class("text-3xl font-semibold tracking-tight"), "grove"),
		g.Button(
			g.Class("inline-flex h-9 items-center justify-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground shadow hover:bg-primary/90"),
			g.OnClick(func(*g.Event) { setCount(count + 1) }),
			g.Textf("count is %d", count),
		),
		g.P(g.Class("text-sm text-muted-foreground"), "edit main.go and save to reload"),
	)
}

func main() {
	dom.Mount("#root", g.C0(App))
}
`

const initIndexHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>{{name}}</title>
    <link rel="stylesheet" href="/styles.css" />
  </head>
  <body>
    <div id="root"></div>
    <script src="/wasm_exec.js"></script>
    <script>
      const go = new Go();
      WebAssembly.instantiateStreaming(fetch("/app.wasm"), go.importObject).then(
        (result) => go.run(result.instance),
      );
    </script>
  </body>
</html>
`

// themeCSS is Tailwind v4 plus grove's CSS-variable design system: every
// color and radius the ui components use is a variable here, so restyling
// an app (including dark mode) means editing variables, not components.
const themeCSS = `@import "tailwindcss";
@source "../**/*.go";

@custom-variant dark (&:is(.dark *));

:root {
  --radius: 0.625rem;
  --background: oklch(1 0 0);
  --foreground: oklch(0.145 0 0);
  --card: oklch(1 0 0);
  --card-foreground: oklch(0.145 0 0);
  --popover: oklch(1 0 0);
  --popover-foreground: oklch(0.145 0 0);
  --primary: oklch(0.205 0 0);
  --primary-foreground: oklch(0.985 0 0);
  --secondary: oklch(0.97 0 0);
  --secondary-foreground: oklch(0.205 0 0);
  --muted: oklch(0.97 0 0);
  --muted-foreground: oklch(0.556 0 0);
  --accent: oklch(0.97 0 0);
  --accent-foreground: oklch(0.205 0 0);
  --destructive: oklch(0.577 0.245 27.325);
  --destructive-foreground: oklch(0.985 0 0);
  --border: oklch(0.922 0 0);
  --input: oklch(0.922 0 0);
  --ring: oklch(0.708 0 0);
}

.dark {
  --background: oklch(0.145 0 0);
  --foreground: oklch(0.985 0 0);
  --card: oklch(0.205 0 0);
  --card-foreground: oklch(0.985 0 0);
  --popover: oklch(0.205 0 0);
  --popover-foreground: oklch(0.985 0 0);
  --primary: oklch(0.922 0 0);
  --primary-foreground: oklch(0.205 0 0);
  --secondary: oklch(0.269 0 0);
  --secondary-foreground: oklch(0.985 0 0);
  --muted: oklch(0.269 0 0);
  --muted-foreground: oklch(0.708 0 0);
  --accent: oklch(0.269 0 0);
  --accent-foreground: oklch(0.985 0 0);
  --destructive: oklch(0.704 0.191 22.216);
  --destructive-foreground: oklch(0.985 0 0);
  --border: oklch(1 0 0 / 10%);
  --input: oklch(1 0 0 / 15%);
  --ring: oklch(0.556 0 0);
}

@theme inline {
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-card: var(--card);
  --color-card-foreground: var(--card-foreground);
  --color-popover: var(--popover);
  --color-popover-foreground: var(--popover-foreground);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
  --color-secondary: var(--secondary);
  --color-secondary-foreground: var(--secondary-foreground);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-foreground);
  --color-accent: var(--accent);
  --color-accent-foreground: var(--accent-foreground);
  --color-destructive: var(--destructive);
  --color-destructive-foreground: var(--destructive-foreground);
  --color-border: var(--border);
  --color-input: var(--input);
  --color-ring: var(--ring);
  --radius-sm: calc(var(--radius) - 4px);
  --radius-md: calc(var(--radius) - 2px);
  --radius-lg: var(--radius);
  --radius-xl: calc(var(--radius) + 4px);

  /* Overlay enter/leave animations used by the ui Dialog/Sheet/Drawer. */
  --animate-overlay-in: overlay-in 0.2s ease;
  --animate-overlay-out: overlay-out 0.2s ease forwards;
  --animate-dialog-in: dialog-in 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  --animate-slide-in-right: slide-in-right 0.3s cubic-bezier(0.32, 0.72, 0, 1);
  --animate-slide-in-left: slide-in-left 0.3s cubic-bezier(0.32, 0.72, 0, 1);
  --animate-slide-in-top: slide-in-top 0.3s cubic-bezier(0.32, 0.72, 0, 1);
  --animate-slide-in-bottom: slide-in-bottom 0.3s cubic-bezier(0.32, 0.72, 0, 1);
  --animate-slide-out-right: slide-out-right 0.25s cubic-bezier(0.32, 0.72, 0, 1) forwards;
  --animate-slide-out-left: slide-out-left 0.25s cubic-bezier(0.32, 0.72, 0, 1) forwards;
  --animate-slide-out-top: slide-out-top 0.25s cubic-bezier(0.32, 0.72, 0, 1) forwards;
  --animate-slide-out-bottom: slide-out-bottom 0.25s cubic-bezier(0.32, 0.72, 0, 1) forwards;
}

@keyframes overlay-in { from { opacity: 0; } to { opacity: 1; } }
@keyframes overlay-out { from { opacity: 1; } to { opacity: 0; } }
/* The dialog centers with the translate property, so it only scales/fades. */
@keyframes dialog-in {
  from { opacity: 0; transform: scale(0.96); }
  to { opacity: 1; transform: scale(1); }
}
@keyframes slide-in-right { from { transform: translateX(100%); } to { transform: translateX(0); } }
@keyframes slide-in-left { from { transform: translateX(-100%); } to { transform: translateX(0); } }
@keyframes slide-in-top { from { transform: translateY(-100%); } to { transform: translateY(0); } }
@keyframes slide-in-bottom { from { transform: translateY(100%); } to { transform: translateY(0); } }
@keyframes slide-out-right { to { transform: translateX(100%); } }
@keyframes slide-out-left { to { transform: translateX(-100%); } }
@keyframes slide-out-top { to { transform: translateY(-100%); } }
@keyframes slide-out-bottom { to { transform: translateY(100%); } }

@layer base {
  * {
    border-color: var(--color-border);
  }
  body {
    background-color: var(--color-background);
    color: var(--color-foreground);
  }
}
`
