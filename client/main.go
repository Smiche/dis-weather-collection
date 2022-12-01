package main

import (
	"fmt"

	db "weather_client/db"

	"fyne.io/fyne/v2/app"
)

// Entry point for the main client.
func main() {
	a := app.New()
	w := a.NewWindow("WeatherDB Client")

	// Load config files and once a user picks a file init db connection and switch UI.
	db.Get_config_view(&w, func(conf db.Config) {
		fmt.Println("Selected:", conf)
		conn := db.Init_db_conn(conf)
		show_main_view(&w, conn)
	})

	// start the UI
	w.ShowAndRun()
}
