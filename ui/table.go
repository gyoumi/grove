package ui

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// Table and its parts compose a themed data table:
//
//	ui.Table(
//	    ui.TableHeader(ui.TableRow(ui.TableHead("Name"), ui.TableHead("Email"))),
//	    ui.TableBody(ui.TableRow(ui.TableCell("Ada"), ui.TableCell("ada@x.io"))),
//	)
//
// Table wraps itself in a horizontally scrollable container.
func Table(args ...any) *g.Node {
	return g.Div(
		g.Class("relative w-full overflow-auto"),
		g.Data("slot", "table-container"),
		tslot("table", "w-full caption-bottom text-sm", "table", args),
	)
}

func TableHeader(args ...any) *g.Node {
	return tslot("thead", "[&_tr]:border-b", "table-header", args)
}
func TableBody(args ...any) *g.Node {
	return tslot("tbody", "[&_tr:last-child]:border-0", "table-body", args)
}
func TableFooter(args ...any) *g.Node {
	return tslot("tfoot", "border-t bg-muted/50 font-medium", "table-footer", args)
}
func TableRow(args ...any) *g.Node {
	return tslot("tr", "border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted", "table-row", args)
}
func TableHead(args ...any) *g.Node {
	return tslot("th", "h-10 px-2 text-left align-middle font-medium text-muted-foreground", "table-head", args)
}
func TableCell(args ...any) *g.Node {
	return tslot("td", "p-2 align-middle", "table-cell", args)
}
func TableCaption(args ...any) *g.Node {
	return tslot("caption", "mt-4 text-sm text-muted-foreground", "table-caption", args)
}

func tslot(tag, base, slot string, args []any) *g.Node {
	all := append([]any{g.Class(style.CN(base)), g.Data("slot", slot)}, args...)
	return g.El(tag, all...)
}
