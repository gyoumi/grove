package ui_test

import (
	"strings"
	"testing"

	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

func TestToastShowStackAndDismiss(t *testing.T) {
	ui.DismissAllToasts()
	defer ui.DismissAllToasts()

	r := testdom.Mount(ui.Toaster())
	if r.FindByAttr("data-slot", "toast") != nil {
		t.Fatal("no toasts should show initially")
	}

	ui.Toast("Saved", ui.ToastOptions{Description: "Your changes are saved.", Variant: ui.ToastSuccess})
	r.Settle()
	toast := r.FindByAttr("data-slot", "toast")
	if toast == nil {
		t.Fatalf("the toast should render after Toast(): %s", r.HTML())
	}
	if toast.Attrs["data-variant"] != "success" {
		t.Fatalf("variant should be recorded: %s", toast.HTML())
	}
	if !strings.Contains(toast.TextContent(), "Saved") || !strings.Contains(toast.TextContent(), "Your changes are saved.") {
		t.Fatalf("toast should show title + description: %s", toast.HTML())
	}
	if r.FindByAttr("data-icon", "circle-check") == nil {
		t.Fatalf("a success toast should show a check icon: %s", r.HTML())
	}

	// a second toast stacks alongside the first
	ui.Toast("Heads up")
	r.Settle()
	if r.FindText("Saved") == nil || r.FindText("Heads up") == nil {
		t.Fatalf("both toasts should be visible: %s", r.HTML())
	}

	// dismissing the first leaves the second
	r.Click(r.FindByAttr("data-slot", "toast-dismiss"))
	r.Settle()
	if r.FindText("Saved") != nil {
		t.Fatal("dismissing should remove the first toast")
	}
	if r.FindText("Heads up") == nil {
		t.Fatal("the second toast should remain")
	}

	// DismissAllToasts clears the rest
	ui.DismissAllToasts()
	r.Settle()
	if r.FindByAttr("data-slot", "toast") != nil {
		t.Fatalf("DismissAllToasts should clear everything: %s", r.HTML())
	}
}
