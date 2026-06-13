package ui_test

import (
	"testing"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/testdom"
	"github.com/gyoumi/grove/ui"
)

// dayInCurrentMonth returns day-of-month d in the month the pickers' inner
// calendar opens on, so tests can click a cell deterministically.
func dayInCurrentMonth(d int) ui.Date {
	t := ui.Today()
	return ui.Date{Year: t.Year, Month: t.Month, Day: d}
}

func singleDatePicker() *g.Node {
	val, setVal := g.UseState[*ui.Date](nil)
	return ui.DatePicker(ui.DatePickerProps{
		Value:    val,
		OnChange: func(d ui.Date) { dd := d; setVal(&dd) },
	})
}

func TestDatePickerSingleOpensAndPicks(t *testing.T) {
	r := testdom.Mount(g.C0(singleDatePicker))

	trig := r.FindByAttr("data-slot", "date-picker-trigger")
	if trig == nil || trig.Attrs["data-empty"] != "1" || trig.TextContent() == "" {
		t.Fatalf("trigger should start empty with a placeholder: %s", r.HTML())
	}
	if r.FindByAttr("data-slot", "calendar") != nil {
		t.Fatal("calendar should be closed until the trigger is clicked")
	}

	r.Click(trig)
	if r.FindByAttr("data-slot", "calendar") == nil {
		t.Fatalf("clicking the trigger should open the calendar: %s", r.HTML())
	}

	day := dayInCurrentMonth(1)
	r.Click(r.FindByAttr("data-day", day.ISO()))

	if r.FindByAttr("data-slot", "calendar") != nil {
		t.Fatal("picking a day should close the popover")
	}
	trig = r.FindByAttr("data-slot", "date-picker-trigger")
	if _, empty := trig.Attrs["data-empty"]; empty {
		t.Fatal("trigger should no longer be empty")
	}
	if got := trig.TextContent(); got == "" || !contains(got, day.ISO()) {
		t.Fatalf("trigger should show %s, got %q", day.ISO(), got)
	}
}

func rangeDatePicker() *g.Node {
	start, setStart := g.UseState[*ui.Date](nil)
	end, setEnd := g.UseState[*ui.Date](nil)
	return ui.DatePicker(ui.DatePickerProps{
		Mode:  ui.CalendarRange,
		Start: start, End: end,
		OnRange: func(s, e ui.Date) { ss, ee := s, e; setStart(&ss); setEnd(&ee) },
	})
}

func TestDatePickerRangeStaysOpenUntilBothEnds(t *testing.T) {
	r := testdom.Mount(g.C0(rangeDatePicker))
	r.Click(r.FindByAttr("data-slot", "date-picker-trigger"))

	first := dayInCurrentMonth(1)
	fifth := dayInCurrentMonth(5)

	r.Click(r.FindByAttr("data-day", first.ISO()))
	if r.FindByAttr("data-slot", "calendar") == nil {
		t.Fatal("popover should stay open after the first endpoint")
	}
	r.Click(r.FindByAttr("data-day", fifth.ISO()))
	if r.FindByAttr("data-slot", "calendar") != nil {
		t.Fatal("popover should close once the range is complete")
	}
	got := r.FindByAttr("data-slot", "date-picker-trigger").TextContent()
	if !contains(got, first.ISO()) || !contains(got, fifth.ISO()) {
		t.Fatalf("trigger should show the range, got %q", got)
	}
}

func timePicker() *g.Node {
	val, setVal := g.UseState[*ui.Time](nil)
	return ui.TimePicker(ui.TimePickerProps{
		Value:    val,
		OnChange: func(tm ui.Time) { t2 := tm; setVal(&t2) },
	})
}

func TestTimePickerPicksHourAndMinute(t *testing.T) {
	r := testdom.Mount(g.C0(timePicker))

	trig := r.FindByAttr("data-slot", "time-picker-trigger")
	if trig.Attrs["data-empty"] != "1" {
		t.Fatalf("time trigger should start empty: %s", r.HTML())
	}
	r.Click(trig)
	if r.FindByAttr("data-slot", "time-hours") == nil || r.FindByAttr("data-slot", "time-minutes") == nil {
		t.Fatalf("hour and minute columns should open: %s", r.HTML())
	}

	// default step 30 → minute column is 00 and 30 only
	if r.FindByAttr("data-minute", "30") == nil || r.FindByAttr("data-minute", "15") != nil {
		t.Fatalf("minute column should step by 30: %s", r.HTML())
	}

	r.Click(r.FindByAttr("data-hour", "9"))
	if r.FindByAttr("data-slot", "time-hours") == nil {
		t.Fatal("popover should stay open after picking an hour")
	}
	r.Click(r.FindByAttr("data-minute", "30"))

	trig = r.FindByAttr("data-slot", "time-picker-trigger")
	if got := trig.TextContent(); !contains(got, "09:30") {
		t.Fatalf("trigger should show 09:30, got %q", got)
	}
	// the chosen options are marked selected
	if r.FindByAttr("data-hour", "9").Attrs["data-selected"] != "1" {
		t.Fatal("picked hour should be marked selected")
	}
}

func TestTimeHelpers(t *testing.T) {
	tm := ui.TimeFromMinutes(9*60 + 30)
	if tm != (ui.Time{Hour: 9, Minute: 30}) || tm.String() != "09:30" || tm.Minutes() != 570 {
		t.Fatalf("Time helpers wrong: %+v %q %d", tm, tm.String(), tm.Minutes())
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
