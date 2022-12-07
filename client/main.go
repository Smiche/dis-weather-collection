package main

import (
	"fmt"

	db "weather_client/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var application fyne.App

// Entry point for the main client.
func main() {
	application = app.New()
	w := application.NewWindow("WeatherDB Client")

	// Load config files and once a user picks a file init db connection and switch UI.
	db.Get_config_view(&w, func(conf db.Config) {
		fmt.Println("Selected:", conf)
		conn := db.Init_db_conn(conf)
		show_main_view(&w, conn)
	})

	// start the UI
	w.ShowAndRun()
}
