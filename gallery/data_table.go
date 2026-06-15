package gallery

import (
	"fmt"
	"strings"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/ui"
)

// rosterRow is a sample row type for the DataTable demo.
type rosterRow struct {
	name string
	role string
	age  int
}

func dataTableSection() *g.Node {
	return section("Data table",
		ui.DataTable(ui.DataTableProps[rosterRow]{
			Columns: []ui.Column[rosterRow]{
				{Header: "Name", Cell: func(r rosterRow) any { return g.Strong(r.name) },
					Less: func(a, b rosterRow) bool { return a.name < b.name }},
				{Header: "Role", Cell: func(r rosterRow) any { return r.role }},
				{Header: "Age", Class: "text-right", Cell: func(r rosterRow) any { return fmt.Sprint(r.age) },
					Less: func(a, b rosterRow) bool { return a.age < b.age }},
			},
			Rows: []rosterRow{
				{"Ada Lovelace", "Engineer", 36},
				{"Grace Hopper", "Admiral", 85},
				{"Alan Turing", "Researcher", 41},
				{"Katherine Johnson", "Mathematician", 101},
				{"Linus Torvalds", "Maintainer", 54},
				{"Margaret Hamilton", "Director", 88},
			},
			Filter: func(r rosterRow, q string) bool {
				hay := strings.ToLower(fmt.Sprintf("%s %s %d", r.name, r.role, r.age))
				return strings.Contains(hay, strings.ToLower(q))
			},
			PageSize: 4,
		}),
	)
}
