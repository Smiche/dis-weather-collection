package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func show_main_view(window *fyne.Window) {
	(*window).SetContent(widget.NewLabel("Hello World!"))
}
