package ui

import (
	"slices"
	"strconv"
	"strings"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// Column describes one column of a DataTable over rows of type T.
type Column[T any] struct {
	Header string
	// Cell renders the cell content for a row (a string, node, or further
	// options — anything an element constructor accepts).
	Cell func(T) any
	// Less, when set, makes the column sortable; clicking its header sorts by
	// it (and toggles direction).
	Less  func(a, b T) bool
	Class string // applied to the header and cells
}

type DataTableProps[T any] struct {
	Columns []Column[T]
	Rows    []T
	// Filter, when set, shows a search box and keeps rows it returns true for.
	Filter   func(row T, query string) bool
	PageSize int    // rows per page; 0 shows them all
	Empty    string // message when nothing matches (default "No results.")
	Class    string
}

// DataTable renders rows in a sortable, optionally filterable and paginated
// table. It is generic over the row type, so cells and comparators are type
// safe:
//
//	ui.DataTable(ui.DataTableProps[Person]{
//	    Columns: []ui.Column[Person]{
//	        {Header: "Name", Cell: func(p Person) any { return p.Name },
//	            Less: func(a, b Person) bool { return a.Name < b.Name }},
//	    },
//	    Rows:   people,
//	    Filter: func(p Person, q string) bool { return strings.Contains(p.Name, q) },
//	})
func DataTable[T any](p DataTableProps[T]) *g.Node {
	return g.C(dataTableView[T], p)
}

func dataTableView[T any](p DataTableProps[T]) *g.Node {
	sortCol, setSortCol := g.UseState(-1)
	sortAsc, setSortAsc := g.UseState(true)
	query, setQuery := g.UseState("")
	page, setPage := g.UseState(0)

	rows := p.Rows
	if p.Filter != nil && strings.TrimSpace(query) != "" {
		kept := make([]T, 0, len(rows))
		for _, r := range rows {
			if p.Filter(r, query) {
				kept = append(kept, r)
			}
		}
		rows = kept
	}
	if sortCol >= 0 && sortCol < len(p.Columns) && p.Columns[sortCol].Less != nil {
		less := p.Columns[sortCol].Less
		rows = slices.Clone(rows)
		slices.SortStableFunc(rows, func(a, b T) int {
			switch {
			case less(a, b):
				return cmpIf(sortAsc, -1)
			case less(b, a):
				return cmpIf(sortAsc, 1)
			default:
				return 0
			}
		})
	}

	total := len(rows)
	pageCount := 1
	pg := 0
	pageRows := rows
	if p.PageSize > 0 {
		pageCount = max(1, (total+p.PageSize-1)/p.PageSize)
		pg = min(max(page, 0), pageCount-1)
		start := pg * p.PageSize
		pageRows = rows[start:min(start+p.PageSize, total)]
	}

	heads := make([]any, 0, len(p.Columns))
	for ci, col := range p.Columns {
		if col.Less == nil {
			heads = append(heads, TableHead(g.Class(col.Class), col.Header))
			continue
		}
		icon := "chevrons-up-down"
		if sortCol == ci {
			if sortAsc {
				icon = "chevron-up"
			} else {
				icon = "chevron-down"
			}
		}
		heads = append(heads, TableHead(g.Class(col.Class),
			g.Button(
				g.Type("button"),
				g.Class("-ml-1 inline-flex items-center gap-1 rounded px-1 hover:text-foreground"),
				g.Data("slot", "data-table-sort"),
				g.OnClick(func(*g.Event) {
					if sortCol == ci {
						setSortAsc(!sortAsc)
					} else {
						setSortCol(ci)
						setSortAsc(true)
					}
				}),
				col.Header,
				Icon(icon, "size-3.5 opacity-60"),
			),
		))
	}

	var bodyRows []any
	if len(pageRows) == 0 {
		empty := p.Empty
		if empty == "" {
			empty = "No results."
		}
		bodyRows = append(bodyRows, TableRow(
			TableCell(g.Attr("colspan", strconv.Itoa(len(p.Columns))), g.Class("h-24 text-center text-muted-foreground"), empty),
		))
	} else {
		for _, row := range pageRows {
			cells := make([]any, 0, len(p.Columns))
			for _, col := range p.Columns {
				cells = append(cells, TableCell(g.Class(col.Class), col.Cell(row)))
			}
			bodyRows = append(bodyRows, TableRow(cells...))
		}
	}

	parts := []any{
		g.Class(style.CN("flex w-full flex-col gap-3", p.Class)),
		g.Data("slot", "data-table"),
	}
	if p.Filter != nil {
		parts = append(parts, g.Div(g.Class("flex items-center"),
			Input(InputProps{
				Value: query, Placeholder: "Filter…", Class: "max-w-xs",
				OnInput: func(e *g.Event) { setQuery(e.Value()); setPage(0) },
			}),
		))
	}
	parts = append(parts, g.Div(g.Class("rounded-md border"),
		Table(
			TableHeader(TableRow(heads...)),
			TableBody(bodyRows...),
		),
	))
	if p.PageSize > 0 && pageCount > 1 {
		parts = append(parts, g.Div(g.Class("flex items-center justify-between"),
			g.Span(g.Class("text-sm text-muted-foreground"), g.Textf("Page %d of %d", pg+1, pageCount)),
			g.Div(g.Class("flex gap-2"),
				Button(ButtonProps{Variant: ButtonOutline, Size: ButtonSizeSm, Disabled: pg <= 0,
					OnClick: func(*g.Event) { setPage(pg - 1) }}, "Previous"),
				Button(ButtonProps{Variant: ButtonOutline, Size: ButtonSizeSm, Disabled: pg >= pageCount-1,
					OnClick: func(*g.Event) { setPage(pg + 1) }}, "Next"),
			),
		))
	}
	return g.Div(parts...)
}

func cmpIf(asc bool, v int) int {
	if asc {
		return v
	}
	return -v
}
