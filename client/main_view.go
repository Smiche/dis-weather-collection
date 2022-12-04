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

// Handles the main view. Shown once user connects using a configuration.
func show_main_view(window *fyne.Window, conn *pgx.Conn) {
	stations, err := db.Get_stations(conn)
	if err != nil {
		return
	}

	// Table root component
	table := container.New(layout.NewHBoxLayout())

	// Adding each column
	for i := 0; i < 5; i++ {
		table.Add(container.New(layout.NewVBoxLayout()))
	}
	// Headers
	table.Objects[0].(*fyne.Container).Add(get_header("#"))
	table.Objects[1].(*fyne.Container).Add(get_header("Country"))
	table.Objects[2].(*fyne.Container).Add(get_header("Name"))
	table.Objects[3].(*fyne.Container).Add(get_header("Type"))
	table.Objects[4].(*fyne.Container).Add(get_header("City"))

	// Rows
	for _, station := range stations {
		lNumber := get_cell(fmt.Sprint(station.Number))
		lCountry := get_cell(station.Country)
		lName := get_cell(station.Name)
		lType := get_cell(station.Type)
		lCity := get_cell(station.City)

		table.Objects[0].(*fyne.Container).Add(lNumber)
		table.Objects[1].(*fyne.Container).Add(lCountry)
		table.Objects[2].(*fyne.Container).Add(lName)
		table.Objects[3].(*fyne.Container).Add(lType)
		table.Objects[4].(*fyne.Container).Add(lCity)
	}

	queryText := canvas.NewText("Query: \"select * from station_all\"", color.Black)
	queryText.Alignment = fyne.TextAlignLeading
	queryText.TextStyle = fyne.TextStyle{Italic: true}

	localContainer := container.New(layout.NewVBoxLayout())

	tabs := container.NewAppTabs(
		container.NewTabItem("Stations", container.New(layout.NewVBoxLayout(), queryText, table)),
		container.NewTabItem("Local Query", localContainer),
		container.NewTabItem("Global Query", widget.NewLabel("World!")),
	)

	tabs.OnSelected = func(ti *container.TabItem) {
		// Track tab switching and render
		if ti.Text == "Local Query" {
			localContainer.RemoveAll()
			Local_view(localContainer, conn)
		}
	}

	tabs.SetTabLocation(container.TabLocationTop)

	(*window).SetContent(tabs)
}

// generates a canvas.Text with header style
func get_header(name string) *canvas.Text {
	text := canvas.NewText(name, color.Black)
	text.Alignment = fyne.TextAlignLeading
	text.TextStyle = fyne.TextStyle{Bold: true}
	return text
}

// generates a canvas.Text for the cells
func get_cell(name string) *canvas.Text {
	text := canvas.NewText(name, color.Black)
	text.Alignment = fyne.TextAlignLeading
	text.TextStyle = fyne.TextStyle{}
	return text
}

// Draws a chart plot on the Local tab
func Local_view(localContainer *fyne.Container, conn *pgx.Conn) {
	curSDateLabel := widget.NewLabel("")
	sDateButton := widget.NewButton("Start Date", func() {
		date_picker(func(selected time.Time) {
			curSDateLabel.SetText(selected.Format("2006-01-02"))
		})
	})

	sDateCont := container.New(layout.NewHBoxLayout(), sDateButton, curSDateLabel)

	curEDateLabel := widget.NewLabel("")
	eDateButton := widget.NewButton("End Date", func() {
		date_picker(func(selected time.Time) {
			curEDateLabel.SetText(selected.Format("2006-01-02"))
		})
	})

	eDateCont := container.New(layout.NewHBoxLayout(), eDateButton, curEDateLabel)

	stations, _ := db.Get_local_stations(conn)
	var stationNames []string

	for _, station := range stations {
		stationNames = append(stationNames, station.Name)
	}

	dropdown := widget.NewSelect(stationNames, func(value string) {
		for _, station := range stations {
			if station.Name == value {
				fmt.Print("Selected:")
			}
		}
	})

	localContainer.Add(sDateCont)
	localContainer.Add(eDateCont)
	localContainer.Add(dropdown)

	local_view_chart(localContainer, db.Query_local_data(conn, 1))
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
	imageContainer := container.New(layout.NewCenterLayout(), canvasImage)
	cont.Add(imageContainer)
}
