package gallery

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/ui"
)

func dateTimeSection() *g.Node {
	return section("Date & time",
		g.Div(g.Class("flex flex-wrap items-start gap-4"),
			ui.Calendar(ui.CalendarProps{}),
			g.Div(g.Class("flex w-56 flex-col gap-3"),
				ui.DatePicker(ui.DatePickerProps{Placeholder: "Pick a date"}),
				ui.TimePicker(ui.TimePickerProps{Value: &ui.Time{Hour: 9, Minute: 30}}),
			),
		),
	)
}
