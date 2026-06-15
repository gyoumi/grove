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

	display := name[strings.LastIndex(name, "/")+1:]
	mainGo := strings.ReplaceAll(initMainGo, "{{MODULE}}", name)
	mainGo = strings.ReplaceAll(mainGo, "{{APP}}", display)

	files := map[string]string{
		"go.mod":           gomod,
		"main.go":          mainGo,
		"index.html":       strings.ReplaceAll(initIndexHTML, "{{APP}}", display),
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

	// Vendor the full ui component set and the component gallery, so a new
	// app ships with every component and a browsable /components catalog.
	if err := vendorTemplates("ui", filepath.Join(name, "ui"), ""); err != nil {
		return err
	}
	if err := vendorTemplates("gallery", filepath.Join(name, "gallery"), name); err != nil {
		return err
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

	fmt.Printf("created %s/ with the full ui component set and a gallery\n\n", name)
	fmt.Printf("  cd %s\n", name)
	if *grovePath == "" {
		fmt.Println("  go get github.com/gyoumi/grove@latest")
	}
	fmt.Println("  grove serve")
	fmt.Println("\nedit ui/ freely (the components are yours); browse them all at /components")
	return nil
}

// vendorTemplates writes every embedded template under templates/<sub> into
// destDir as a .go file, replacing the {{MODULE}} placeholder with module
// when set (the gallery uses it to import the app's own ui package).
func vendorTemplates(sub, destDir, module string) error {
	entries, err := templates.ReadDir("templates/" + sub)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return err
	}
	for _, e := range entries {
		data, err := templates.ReadFile("templates/" + sub + "/" + e.Name())
		if err != nil {
			return err
		}
		if module != "" {
			data = []byte(strings.ReplaceAll(string(data), "{{MODULE}}", module))
		}
		dest := filepath.Join(destDir, strings.TrimSuffix(e.Name(), ".txt"))
		if err := os.WriteFile(dest, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

const initMainGo = `package main

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/dom"
	"github.com/gyoumi/grove/router"
	"{{MODULE}}/gallery"
	"{{MODULE}}/ui"
)

// App is the root: a small home page and a /components route showing the full
// grove component gallery. The header links between them and toggles dark
// mode. Replace Home with your app — the ui/ package and gallery/ are yours.
func App() *g.Node {
	dark, setDark := g.UseState(false)
	g.UseEffect(func() func() {
		dom.SetRootClass("dark", dark)
		return nil
	}, []any{dark})

	return g.Div(g.Class("mx-auto flex min-h-svh max-w-3xl flex-col gap-6 px-6 pb-6 text-foreground"),
		g.Header(g.Class("sticky top-0 z-40 -mx-6 flex items-center justify-between border-b border-border/60 bg-background/75 px-6 py-3 backdrop-blur-md"),
			router.Link("/", g.Class("text-xl font-bold tracking-tight no-underline"), "{{APP}}"),
			g.Div(g.Class("flex items-center gap-3"),
				router.Link("/components", g.Class("text-sm font-medium text-muted-foreground no-underline transition-colors hover:text-foreground"),
					"components"),
				ui.Tooltip(ui.TooltipProps{Label: "toggle dark mode"},
					ui.Switch(ui.SwitchProps{ID: "dark-mode", Checked: dark, OnChange: setDark}),
				),
			),
		),
		router.Routes(
			router.Route{Pattern: "/", Render: func(router.Params) *g.Node { return g.C0(Home) }},
			router.Route{Pattern: "/components", Render: func(router.Params) *g.Node { return g.C0(gallery.Page) }},
			router.Route{Pattern: "*", Render: func(router.Params) *g.Node {
				return g.P(g.Class("text-sm text-muted-foreground"), "That page doesn't exist.")
			}},
		),
		ui.Toaster(),
	)
}

// Home is the starter page. Edit it freely.
func Home() *g.Node {
	count, setCount := g.UseState(0)
	return g.Div(g.Class("flex flex-col items-center gap-5 py-16 text-center animate-rise"),
		g.H1(g.Class("text-4xl font-bold tracking-tight"), "{{APP}}"),
		g.P(g.Class("max-w-md text-balance text-muted-foreground"),
			"A grove app. Edit Home in main.go, or browse every component on the ",
			router.Link("/components", g.Class("font-medium text-primary underline-offset-4 hover:underline"), "components"),
			" page."),
		ui.Button(ui.ButtonProps{OnClick: func(*g.Event) { setCount(count + 1) }},
			g.Textf("count is %d", count)),
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
    <title>{{APP}}</title>
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

  /* A subtle entrance used by the gallery and starter page. */
  --animate-rise: rise 0.4s cubic-bezier(0.16, 1, 0.3, 1);

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

@keyframes rise { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }
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
