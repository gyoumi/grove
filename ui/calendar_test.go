package ui_test

import (
	"testing"
	"time"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

var june2026 = ui.Date{Year: 2026, Month: time.June, Day: 1}

func TestCalendarRendersAndNavigates(t *testing.T) {
	r := testdom.Mount(ui.Calendar(ui.CalendarProps{Month: june2026}))

	if lbl := r.FindByAttr("data-slot", "cal-label"); lbl == nil || lbl.TextContent() != "June 2026" {
		t.Fatalf("month label wrong: %s", r.HTML())
	}
	// June has 30 days; July 1 must not be in this month's grid.
	if r.FindByAttr("data-day", "2026-06-30") == nil || r.FindByAttr("data-day", "2026-07-01") != nil {
		t.Fatalf("June grid should be 1..30: %s", r.HTML())
	}

	r.Click(r.FindByAttr("data-slot", "cal-next"))
	if lbl := r.FindByAttr("data-slot", "cal-label"); lbl == nil || lbl.TextContent() != "July 2026" {
		t.Fatalf("next month not shown: %s", r.HTML())
	}
	if r.FindByAttr("data-day", "2026-07-31") == nil {
		t.Fatalf("July should have 31 days: %s", r.HTML())
	}

	// back across the year boundary twice from July → May
	r.Click(r.FindByAttr("data-slot", "cal-prev"))
	r.Click(r.FindByAttr("data-slot", "cal-prev"))
	if lbl := r.FindByAttr("data-slot", "cal-label"); lbl.TextContent() != "May 2026" {
		t.Fatalf("prev nav wrong: %s", lbl.TextContent())
	}
}

func singleCal(onSelect func(ui.Date)) func() *g.Node {
	return func() *g.Node {
		sel, setSel := g.UseState[*ui.Date](nil)
		return ui.Calendar(ui.CalendarProps{
			Month: june2026, Mode: ui.CalendarSingle, Selected: sel,
			OnSelect: func(d ui.Date) { dd := d; setSel(&dd); onSelect(d) },
		})
	}
}

func TestCalendarSingleSelect(t *testing.T) {
	var got ui.Date
	r := testdom.Mount(g.C0(singleCal(func(d ui.Date) { got = d })))

	r.Click(r.FindByAttr("data-day", "2026-06-10"))
	if got != (ui.Date{Year: 2026, Month: time.June, Day: 10}) {
		t.Fatalf("OnSelect got %+v", got)
	}
	cell := r.FindByAttr("data-day", "2026-06-10")
	if cell.Attrs["data-selected"] != "1" {
		t.Fatalf("selected day should be marked: %s", cell.HTML())
	}
	// selecting another day moves the highlight
	r.Click(r.FindByAttr("data-day", "2026-06-20"))
	if r.FindByAttr("data-day", "2026-06-10").Attrs["data-selected"] == "1" {
		t.Fatal("old selection should clear")
	}
	if r.FindByAttr("data-day", "2026-06-20").Attrs["data-selected"] != "1" {
		t.Fatal("new selection should be marked")
	}
}

func rangeCal(onRange func(s, e ui.Date)) func() *g.Node {
	return func() *g.Node {
		start, setStart := g.UseState[*ui.Date](nil)
		end, setEnd := g.UseState[*ui.Date](nil)
		return ui.Calendar(ui.CalendarProps{
			Month: june2026, Mode: ui.CalendarRange, Start: start, End: end,
			OnRange: func(s, e ui.Date) {
				ss, ee := s, e
				setStart(&ss)
				setEnd(&ee)
				if onRange != nil {
					onRange(s, e)
				}
			},
		})
	}
}

func TestCalendarRangeSelect(t *testing.T) {
	var gotS, gotE ui.Date
	r := testdom.Mount(g.C0(rangeCal(func(s, e ui.Date) { gotS, gotE = s, e })))

	// first click sets the anchor and highlights only it
	r.Click(r.FindByAttr("data-day", "2026-06-10"))
	if r.FindByAttr("data-day", "2026-06-10").Attrs["data-selected"] != "1" {
		t.Fatalf("anchor should highlight: %s", r.HTML())
	}
	if r.FindByAttr("data-day", "2026-06-12").Attrs["data-in-range"] == "1" {
		t.Fatal("nothing should be in-range before the second click")
	}

	// second click commits the ordered range
	r.Click(r.FindByAttr("data-day", "2026-06-14"))
	if gotS != (ui.Date{Year: 2026, Month: time.June, Day: 10}) ||
		gotE != (ui.Date{Year: 2026, Month: time.June, Day: 14}) {
		t.Fatalf("OnRange got %+v..%+v", gotS, gotE)
	}
	for _, day := range []string{"2026-06-10", "2026-06-14"} {
		if r.FindByAttr("data-day", day).Attrs["data-selected"] != "1" {
			t.Fatalf("endpoint %s should be selected", day)
		}
	}
	if r.FindByAttr("data-day", "2026-06-12").Attrs["data-in-range"] != "1" {
		t.Fatalf("interior day should be in-range: %s", r.HTML())
	}
	if got := ui.DaysBetween(gotS, gotE); got != 4 {
		t.Fatalf("DaysBetween = %d, want 4", got)
	}
}

func TestCalendarRangeReversedOrders(t *testing.T) {
	var gotS, gotE ui.Date
	r := testdom.Mount(g.C0(rangeCal(func(s, e ui.Date) { gotS, gotE = s, e })))
	r.Click(r.FindByAttr("data-day", "2026-06-20")) // later first
	r.Click(r.FindByAttr("data-day", "2026-06-15")) // earlier second
	if gotS != (ui.Date{Year: 2026, Month: time.June, Day: 15}) ||
		gotE != (ui.Date{Year: 2026, Month: time.June, Day: 20}) {
		t.Fatalf("range should be ordered, got %+v..%+v", gotS, gotE)
	}
}

func TestCalendarMinDisablesPastDays(t *testing.T) {
	min := ui.Date{Year: 2026, Month: time.June, Day: 10}
	selected := 0
	app := func() *g.Node {
		return ui.Calendar(ui.CalendarProps{
			Month: june2026, Min: &min, OnSelect: func(ui.Date) { selected++ },
		})
	}
	r := testdom.Mount(g.C0(app))

	early := r.FindByAttr("data-day", "2026-06-05")
	if _, ok := early.Attrs["disabled"]; !ok {
		t.Fatalf("days before Min should be disabled: %s", early.HTML())
	}
	r.Click(early)
	if selected != 0 {
		t.Fatal("clicking a disabled day must not select it")
	}
	r.Click(r.FindByAttr("data-day", "2026-06-15"))
	if selected != 1 {
		t.Fatal("a day on/after Min should select")
	}
}

func TestParseDateAndISORoundTrip(t *testing.T) {
	d, ok := ui.ParseDate("2026-06-13")
	if !ok || d != (ui.Date{Year: 2026, Month: time.June, Day: 13}) {
		t.Fatalf("ParseDate: %+v ok=%v", d, ok)
	}
	if d.ISO() != "2026-06-13" {
		t.Fatalf("ISO: %s", d.ISO())
	}
	if _, ok := ui.ParseDate("nonsense"); ok {
		t.Fatal("bad input should not parse")
	}
	if got := d.AddDays(3); got != (ui.Date{Year: 2026, Month: time.June, Day: 16}) {
		t.Fatalf("AddDays: %+v", got)
	}
}
