package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v4"
)

// Handles the main view. Shown once user connects using a configuration.
func show_main_view(window *fyne.Window, conn *pgx.Conn) {

	(*window).SetContent(widget.NewLabel("Hello World!"))
}

func get_countries()
