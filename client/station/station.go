package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
	db "weather_client/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v5"
)

// Entry point.
func main() {
	a := app.New()
	w := a.NewWindow("Weather Station simulator")
	var conn *pgx.Conn

	// starts with the config selection view and once a config is picked db connection is initialized and view changes.
	db.Get_config_view(&w, func(conf db.Config) {
		fmt.Println("Selected:", conf)
		conn = db.Init_db_conn(conf)
		begin_simulation(&w, conn, conf)
	})
	container.New(layout.NewVBoxLayout())
	w.ShowAndRun()
}

// Displays the simulation UI and handles simulation logic.
func begin_simulation(window *fyne.Window, conn *pgx.Conn, conf db.Config) {
	// channel to stop the subroutine that adds measurements to db
	done := make(chan bool)

	// UI components
	connInfoString := fmt.Sprintf("Connected to: %s:%s", conn.Config().Host, conn.Config().Database)
	connInfoLabel := widget.NewLabel(connInfoString)

	currentValueLabel := widget.NewLabel("Value: 0")
	totalSentLabel := widget.NewLabel("Total sent: 0")
	measContainer := container.New(layout.NewHBoxLayout(), currentValueLabel, totalSentLabel)

	var stopButton *widget.Button
	stopButton = widget.NewButton("Stop", func() {
		done <- true
		stopButton.Disable()
	})

	// Populate the window with components.
	(*window).SetContent(container.New(layout.NewVBoxLayout(), connInfoLabel, measContainer, stopButton))

	// subroutine that executes an insert query every conf.Period milliseconds
	go func(curValueLabel *widget.Label, totValueLabel *widget.Label, dbConn *pgx.Conn, done chan bool) {
		// use MS duration specified in conf as a period
		duration := time.Millisecond * time.Duration(conf.Period)
		ticker := time.NewTicker(duration)

		count := 0
		// stop the timer when this func exits
		defer ticker.Stop()

		for {
			select {
			case <-done:
				fmt.Println("Done!")
				return
			case <-ticker.C:
				// random value 5-30
				val := rand.Float32()*25 + 5
				// Keep moving time by 1 second for every insert.
				valTime := conf.StartDate.Add(time.Duration(count) * time.Second)
				// Execute insert query.
				_, err := dbConn.Exec(context.Background(), "insert into measurement_local( device, value, \"time\", type, unit) values ($1, $2, $3, $4, $5) ", conf.DeviceId, val, valTime, conf.PhenomenonId, conf.UnitId)
				if err != nil {
					fmt.Println(err)
				}
				// Updating counter and UI
				count++
				currentValueLabel.SetText(fmt.Sprintf("Value: %f", val))
				totalSentLabel.SetText(fmt.Sprintf("Total sent: %d", count))
			}
		}
	}(currentValueLabel, totalSentLabel, conn, done)
}
