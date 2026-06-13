package ui

import (
	"time"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

// Date is a calendar day with no time-of-day or timezone — the unit the
// Calendar selects. The zero value is not a valid date; it doubles as
// "none" in props (a nil *Date) and as the value type for selections.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// Today is the current local date.
func Today() Date {
	n := time.Now()
	return Date{n.Year(), n.Month(), n.Day()}
}

// ParseDate reads a YYYY-MM-DD string; ok is false when it doesn't parse.
func ParseDate(s string) (d Date, ok bool) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, false
	}
	return Date{t.Year(), t.Month(), t.Day()}, true
}

// time normalizes the date to UTC midnight so arithmetic and comparison
// never depend on the host timezone.
func (d Date) time() time.Time {
	return time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.UTC)
}

// ISO formats the date as YYYY-MM-DD.
func (d Date) ISO() string { return d.time().Format("2006-01-02") }

// Before and After order two dates.
func (d Date) Before(o Date) bool { return d.time().Before(o.time()) }
func (d Date) After(o Date) bool  { return d.time().After(o.time()) }

// AddDays returns the date n days later (n may be negative).
func (d Date) AddDays(n int) Date {
	t := d.time().AddDate(0, 0, n)
	return Date{t.Year(), t.Month(), t.Day()}
}

// DaysBetween counts the days from a to b (b after a → positive).
func DaysBetween(a, b Date) int {
	return int(b.time().Sub(a.time()) / (24 * time.Hour))
}

// addMonths moves to the first of a month n away, ignoring the day so the
// month-grid navigation never overflows (e.g. Jan 31 → Feb).
func (d Date) addMonths(n int) Date {
	t := time.Date(d.Year, d.Month, 1, 0, 0, 0, 0, time.UTC).AddDate(0, n, 0)
	return Date{t.Year(), t.Month(), 1}
}

func monthDays(m Date) int {
	// Day 0 of the next month is the last day of this one.
	return time.Date(m.Year, m.Month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func firstWeekdayOffset(m Date, weekStart time.Weekday) int {
	first := time.Date(m.Year, m.Month, 1, 0, 0, 0, 0, time.UTC).Weekday()
	return (int(first) - int(weekStart) + 7) % 7
}

// dowShort is indexed by time.Weekday (Sunday = 0).
var dowShort = [...]string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}

// CalendarMode picks single-day or date-range selection.
type CalendarMode int

const (
	CalendarSingle CalendarMode = iota
	CalendarRange
)

type CalendarProps struct {
	Mode CalendarMode

	// Single-day selection (Mode == CalendarSingle). Selected highlights a
	// day; OnSelect fires when one is clicked.
	Selected *Date
	OnSelect func(Date)

	// Range selection (Mode == CalendarRange). Start/End highlight the
	// committed range; OnRange fires (ordered) once the second endpoint is
	// clicked. The first click of a new range shows only that endpoint.
	Start, End *Date
	OnRange    func(start, end Date)

	// Month chooses the month first shown (Day ignored); the zero value
	// derives it from the selection, else today. Navigation is internal.
	Month Date

	// Min and Max, when set, disable days outside the (inclusive) bounds.
	Min, Max *Date

	// WeekStart is the leftmost weekday column (default Sunday).
	WeekStart time.Weekday

	Class string // extra classes for the calendar container
}

// classify reports how a day relates to the committed selection: a
// highlighted endpoint/day (selected) or strictly inside a range (inside).
func (p CalendarProps) classify(d Date) (selected, inside bool) {
	switch p.Mode {
	case CalendarRange:
		if (p.Start != nil && d == *p.Start) || (p.End != nil && d == *p.End) {
			return true, false
		}
		if p.Start != nil && p.End != nil && d.After(*p.Start) && d.Before(*p.End) {
			return false, true
		}
	default:
		if p.Selected != nil && d == *p.Selected {
			return true, false
		}
	}
	return false, false
}

// Calendar renders a themed month grid for picking a single date or a date
// range. It is fully controlled — selection lives in props and flows back
// through OnSelect/OnRange — except for the displayed month and an
// in-progress range endpoint, which it tracks itself.
func Calendar(p CalendarProps) *g.Node { return g.C(calendarView, p) }

// calPending is the first endpoint of a range mid-selection.
type calPending struct {
	anchor Date
	active bool
}

func calendarView(p CalendarProps) *g.Node {
	init := p.Month
	if init.Year == 0 {
		switch {
		case p.Mode == CalendarRange && p.Start != nil:
			init = *p.Start
		case p.Selected != nil:
			init = *p.Selected
		default:
			init = Today()
		}
	}
	month, setMonth := g.UseState(Date{init.Year, init.Month, 1})
	pend, setPend := g.UseState(calPending{})
	today := Today()

	selectDay := func(d Date) {
		if p.Mode != CalendarRange {
			if p.OnSelect != nil {
				p.OnSelect(d)
			}
			return
		}
		if !pend.active {
			setPend(calPending{anchor: d, active: true})
			return
		}
		s, e := pend.anchor, d
		if e.Before(s) {
			s, e = e, s
		}
		setPend(calPending{})
		if p.OnRange != nil {
			p.OnRange(s, e)
		}
	}

	header := g.Div(g.Class("flex items-center justify-between gap-1 pb-2"),
		calNav("cal-prev", "Previous month", "‹", func() { setMonth(month.addMonths(-1)) }),
		g.Div(g.Class("text-sm font-medium"), g.Data("slot", "cal-label"),
			g.Textf("%s %d", month.Month.String(), month.Year)),
		calNav("cal-next", "Next month", "›", func() { setMonth(month.addMonths(1)) }),
	)

	cells := make([]*g.Node, 0, 7+42)
	for i := range 7 {
		wd := time.Weekday((int(p.WeekStart) + i) % 7)
		cells = append(cells, g.Div(
			g.Class("flex size-9 items-center justify-center text-[0.8rem] font-normal text-muted-foreground"),
			dowShort[wd],
		))
	}
	for i := firstWeekdayOffset(month, p.WeekStart); i > 0; i-- {
		cells = append(cells, g.Div(g.Class("size-9")))
	}
	for day := 1; day <= monthDays(month); day++ {
		cells = append(cells, calDay(p, Date{month.Year, month.Month, day}, today, pend, selectDay))
	}

	return g.Div(
		g.Class(style.CN("inline-block select-none rounded-lg border bg-popover p-3 text-popover-foreground shadow-sm", p.Class)),
		g.Data("slot", "calendar"),
		header,
		g.Div(g.Class("grid grid-cols-7 gap-y-1"), cells),
	)
}

func calNav(slot, label, glyph string, onClick func()) *g.Node {
	return g.Button(
		g.Class("inline-flex size-7 items-center justify-center rounded-md border border-input bg-transparent text-base leading-none text-muted-foreground transition-colors hover:bg-accent hover:text-accent-foreground"),
		g.Type("button"),
		g.Data("slot", slot),
		g.Attr("aria-label", label),
		g.OnClick(func(*g.Event) { onClick() }),
		glyph,
	)
}

func calDay(p CalendarProps, d, today Date, pend calPending, onSelect func(Date)) *g.Node {
	disabled := (p.Min != nil && d.Before(*p.Min)) || (p.Max != nil && d.After(*p.Max))

	var selected, inside bool
	if pend.active {
		// While picking a range, show only the first endpoint.
		selected = d == pend.anchor
	} else {
		selected, inside = p.classify(d)
	}
	isToday := d == today

	cls := style.CN(
		"flex size-9 items-center justify-center rounded-md text-sm font-normal transition-colors",
		map[string]bool{
			"hover:bg-accent hover:text-accent-foreground":        !disabled && !selected && !inside,
			"bg-primary text-primary-foreground hover:bg-primary": selected,
			"rounded-none bg-accent text-accent-foreground":       inside,
			"text-muted-foreground opacity-40":                    disabled,
			"ring-1 ring-ring ring-inset":                         isToday && !selected && !inside,
		},
	)
	return g.Button(
		g.Class(cls),
		g.Type("button"),
		g.Data("slot", "cal-day"),
		g.Data("day", d.ISO()),
		g.AttrIf(selected, "data-selected", "1"),
		g.AttrIf(inside, "data-in-range", "1"),
		g.AttrIf(isToday, "data-today", "1"),
		g.Disabled(disabled),
		g.OnClick(func(*g.Event) {
			if disabled {
				return
			}
			onSelect(d)
		}),
		g.Textf("%d", d.Day),
	)
}
