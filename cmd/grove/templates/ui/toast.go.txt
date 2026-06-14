package ui

import (
	"strconv"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type ToastVariant string

const (
	ToastDefault     ToastVariant = "default"
	ToastSuccess     ToastVariant = "success"
	ToastDestructive ToastVariant = "destructive"
)

type ToastOptions struct {
	Description string
	Variant     ToastVariant
	// Duration in ms before auto-dismiss; 0 uses the default (4s), negative
	// keeps it until dismissed.
	Duration int
}

type toastEntry struct {
	id    int
	title string
	opts  ToastOptions
}

// Toasts live in a small package-level store (wasm is single-threaded). A
// mounted Toaster subscribes and re-renders when it changes.
var (
	toastItems      []toastEntry
	toastNextID     int
	toastListeners  = map[int]func(){}
	toastListenerID int
)

func notifyToasts() {
	for _, fn := range toastListeners {
		fn()
	}
}

// Toast shows a transient notification. Call it from anywhere (an event
// handler, say); a Toaster mounted in the tree renders it.
func Toast(title string, opts ...ToastOptions) {
	o := ToastOptions{}
	if len(opts) > 0 {
		o = opts[0]
	}
	toastNextID++
	id := toastNextID
	toastItems = append(toastItems, toastEntry{id: id, title: title, opts: o})
	notifyToasts()

	d := o.Duration
	if d == 0 {
		d = 4000
	}
	if d > 0 {
		scheduleToastDismiss(id, d)
	}
}

func dismissToast(id int) {
	for i, t := range toastItems {
		if t.id == id {
			toastItems = append(toastItems[:i:i], toastItems[i+1:]...)
			notifyToasts()
			return
		}
	}
}

// DismissAllToasts clears every visible toast.
func DismissAllToasts() {
	if toastItems == nil {
		return
	}
	toastItems = nil
	notifyToasts()
}

// Toaster renders the toast viewport. Place it once near the app root; it
// portals its toasts to the mount container so they sit above everything.
func Toaster() *g.Node { return g.C0(toasterView) }

func toasterView() *g.Node {
	_, bump := g.UseReducer(func(n int, _ struct{}) int { return n + 1 }, 0)
	g.UseEffect(func() func() {
		toastListenerID++
		id := toastListenerID
		toastListeners[id] = func() { bump(struct{}{}) }
		return func() { delete(toastListeners, id) }
	}, []any{})

	if len(toastItems) == 0 {
		return nil
	}
	cards := make([]any, 0, len(toastItems))
	for _, t := range toastItems {
		cards = append(cards, toastCard(t).WithKey(strconv.Itoa(t.id)))
	}
	return g.Portal(g.Div(
		g.Class("pointer-events-none fixed bottom-4 right-4 z-[100] flex w-full max-w-sm flex-col gap-2"),
		g.Data("slot", "toaster"),
		cards,
	))
}

func toastCard(t toastEntry) *g.Node {
	chrome := "border bg-background text-foreground"
	var icon *g.Node
	switch t.opts.Variant {
	case ToastSuccess:
		chrome = "border border-emerald-500/40 bg-background"
		icon = Icon("circle-check", "size-5 shrink-0 text-emerald-500")
	case ToastDestructive:
		chrome = "border border-destructive/50 bg-background"
		icon = Icon("triangle-alert", "size-5 shrink-0 text-destructive")
	}

	body := []any{
		g.Class(style.CN("pointer-events-auto flex items-start gap-3 rounded-lg p-4 shadow-lg animate-slide-in-right", chrome)),
		g.Data("slot", "toast"),
		g.Data("variant", string(t.opts.Variant)),
		g.Role("status"),
	}
	if icon != nil {
		body = append(body, icon)
	}

	content := []any{g.Class("flex min-w-0 flex-1 flex-col gap-0.5"),
		g.Div(g.Class("text-sm font-medium"), t.title),
	}
	if t.opts.Description != "" {
		content = append(content, g.Div(g.Class("text-sm text-muted-foreground"), t.opts.Description))
	}

	body = append(body,
		g.Div(content...),
		g.Button(
			g.Type("button"),
			g.Class("shrink-0 rounded-md text-muted-foreground/70 transition-colors hover:text-foreground"),
			g.Data("slot", "toast-dismiss"),
			g.Attr("aria-label", "Dismiss"),
			g.OnClick(func(*g.Event) { dismissToast(t.id) }),
			Icon("x", "size-4"),
		),
	)
	return g.Div(body...)
}
