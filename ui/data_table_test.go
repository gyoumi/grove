package ui_test

import (
	"strconv"
	"strings"
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

type person struct {
	name string
	age  int
}

func dataTableApp() *g.Node {
	return ui.DataTable(ui.DataTableProps[person]{
		Columns: []ui.Column[person]{
			{Header: "Name", Cell: func(p person) any { return p.name },
				Less: func(a, b person) bool { return a.name < b.name }},
			{Header: "Age", Cell: func(p person) any { return strconv.Itoa(p.age) },
				Less: func(a, b person) bool { return a.age < b.age }},
		},
		Rows:     []person{{"Cara", 40}, {"Ada", 30}, {"Bob", 35}},
		Filter:   func(p person, q string) bool { return strings.Contains(strings.ToLower(p.name), strings.ToLower(q)) },
		PageSize: 2,
	})
}

func firstCell(r *testdom.R) string {
	tds := r.FindAll("td")
	if len(tds) == 0 {
		return ""
	}
	return tds[0].TextContent()
}

func TestDataTableSortFilterPaginate(t *testing.T) {
	r := testdom.Mount(g.C0(dataTableApp))

	// paginated to two rows per page
	if r.FindText("Page 1 of 2") == nil {
		t.Fatalf("should paginate: %s", r.HTML())
	}
	if n := len(r.FindAll("td")); n != 4 { // 2 rows × 2 columns
		t.Fatalf("page 1 should show 2 rows (4 cells), got %d", n)
	}

	// sort by Name ascending → Ada first
	r.Click(r.FindText("Name"))
	if got := firstCell(r); got != "Ada" {
		t.Fatalf("ascending sort should start with Ada, got %q", got)
	}
	// clicking again flips to descending → Cara first
	r.Click(r.FindText("Name"))
	if got := firstCell(r); got != "Cara" {
		t.Fatalf("descending sort should start with Cara, got %q", got)
	}

	// next page (in descending order Cara, Bob, Ada → page 2 is Ada)
	r.Click(r.FindText("Next"))
	if r.FindText("Page 2 of 2") == nil || firstCell(r) != "Ada" {
		t.Fatalf("page 2 should show the last row: %s", r.HTML())
	}

	// filtering narrows the rows and resets to page 1
	r.Input(r.FindByAttr("data-slot", "input"), "bo")
	if n := len(r.FindAll("td")); n != 2 || firstCell(r) != "Bob" {
		t.Fatalf("filter 'bo' should leave just Bob: %s", r.HTML())
	}

	// a filter with no matches shows the empty state
	r.Input(r.FindByAttr("data-slot", "input"), "zzz")
	if r.FindText("No results.") == nil {
		t.Fatalf("no matches should show the empty state: %s", r.HTML())
	}
}
