package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"time"
	db "weather_client/db"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v5"
	chart "github.com/wcharczuk/go-chart/v2"
)

// Creates all elements for the local view
func Local_tab(localContainer *fyne.Container, conn *pgx.Conn) {
	var startTime time.Time
	var endTime time.Time
	var selectedStation int

	// label, button and container for start time picker
	curSDateLabel := widget.NewLabel("")
	sDateButton := widget.NewButton("Start Date", func() {
		date_picker(func(selected time.Time) {
			curSDateLabel.SetText(selected.Format("2006-01-02"))
			startTime = selected
		})
	})
	sDateCont := container.New(layout.NewHBoxLayout(), sDateButton, curSDateLabel)

	// label, button and container for end time picker
	curEDateLabel := widget.NewLabel("")
	eDateButton := widget.NewButton("End Date", func() {
		date_picker(func(selected time.Time) {
			curEDateLabel.SetText(selected.Format("2006-01-02"))
			endTime = selected
		})
	})
	eDateCont := container.New(layout.NewHBoxLayout(), eDateButton, curEDateLabel)

	// query local stations and show them in the dropdown selection
	stations, _ := db.Get_local_stations(conn)
	var stationNames []string

	for _, station := range stations {
		stationNames = append(stationNames, station.Name)
	}

	dropdown := widget.NewSelect(stationNames, func(value string) {
		for _, station := range stations {
			if station.Name == value {
				selectedStation = int(station.ID)
			}
		}
	})

	noResultsLabel := widget.NewLabel("No results found for the query.")
	noResultsLabel.Hide()

	// queryRawData := false
	// rawDataTick := widget.NewCheck("Optional", func(value bool) {
	// 	queryRawData = value
	// })

	imageContainer := container.New(layout.NewCenterLayout())

	queryButton := widget.NewButton("Query", func() {
		measurements := db.Query_local_data(conn, selectedStation, startTime, endTime)
		if len(measurements) > 0 {
			noResultsLabel.Hide()
			imageContainer.Show()
			local_view_chart(imageContainer, measurements)
		} else {
			imageContainer.Hide()
			noResultsLabel.Show()
		}
	})

	// label that shows the query
	queryText := canvas.NewText("Query: select * from meas_min_max_day_local where station_info=$1 and time >= $2 and time <= $3 order by time ASC", color.Black)
	queryText.TextStyle = fyne.TextStyle{Italic: true}
	queryText.Alignment = fyne.TextAlignLeading

	// add all elements to the view
	localContainer.Add(sDateCont)
	localContainer.Add(eDateCont)
	localContainer.Add(dropdown)
	localContainer.Add(queryButton)
	localContainer.Add(noResultsLabel)
	localContainer.Add(imageContainer)
	localContainer.Add(queryText)
}

func local_view_chart(cont *fyne.Container, measurements []db.MeasurementMinMax) {
	var xValues []time.Time
	var yAvg []float64
	var yMin []float64
	var yMax []float64
	unit := ""

	// get unit from first item
	if len(measurements) > 0 {
		unit = measurements[0].Unit
	}

	// flatten measurements to values for the chart
	for _, meas := range measurements {
		xValues = append(xValues, meas.Time)
		yAvg = append(yAvg, meas.Avg)
		yMin = append(yMin, meas.Min)
		yMax = append(yMax, meas.Max)
	}

	// create the chart
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: "Time",
		},
		YAxis: chart.YAxis{
			Name: fmt.Sprintf("Max, Avg, Min %s", unit),
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xValues,
				YValues: yMin,
				Style: chart.Style{
					StrokeColor: chart.ColorGreen,
				},
			},
			chart.TimeSeries{
				XValues: xValues,
				YValues: yMax,
				Style: chart.Style{
					StrokeColor: chart.ColorRed,
				},
			},
			chart.TimeSeries{
				XValues: xValues,
				YValues: yAvg,
				Style: chart.Style{
					StrokeColor: chart.ColorBlue,
				},
			},
		},
	}

	// Render the chart to an image buffer
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		log.Fatal(err)
	}

	image, err := png.Decode(bytes.NewReader(buffer.Bytes()))
	if err != nil {
		log.Fatal(err)
	}

	// Make the fyneio component from the image and add it to the screen.
	canvasImage := canvas.NewImageFromImage(image)
	canvasImage.FillMode = canvas.ImageFillOriginal
	// empty the image container and add the updated chart
	cont.RemoveAll()
	cont.Add(canvasImage)
}
