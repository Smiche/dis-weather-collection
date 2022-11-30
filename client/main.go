package main

import (
	"fmt"

	db "weather_client/db"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	db.Get_config_view(&w, func(conf db.Config) {
		fmt.Println("Selected:", conf)
		show_main_view(&w)
		db.Init_db_conn(conf)
	})
	// w.Resize(fyne.NewSize(600, 480))
	w.ShowAndRun()
}
