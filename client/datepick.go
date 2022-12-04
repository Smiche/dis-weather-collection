package main

import (
	"time"

	xwidget "fyne.io/x/fyne/widget"
)

func date_picker(onSelected func(selected time.Time)) {
	w := application.NewWindow("Date picker")
	cal := xwidget.NewCalendar(time.Now(), func(selected time.Time) {
		onSelected(selected)
		w.Close()
	})
	w.SetContent(cal)
	w.Show()
}
