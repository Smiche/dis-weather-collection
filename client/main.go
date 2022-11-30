package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	get_config_view(&w, func(conf Config) {
		fmt.Println("Selected:", conf)
		show_main_view(&w)
		init_db_conn(conf)
	})
	// w.Resize(fyne.NewSize(600, 480))
	w.ShowAndRun()
}
