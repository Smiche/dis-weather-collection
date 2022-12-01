package main

import (
	"context"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/jackc/pgx/v5"
)

type Station struct {
	ID             int32
	Name           string
	Number         int32
	OrganizationId int32
	Type           string
	Latitude       float32
	Longitude      float32
	Altitude       float32
	City           string
	Country        string
}

// Handles the main view. Shown once user connects using a configuration.
func show_main_view(window *fyne.Window, conn *pgx.Conn) {
	stations, err := get_stations(conn)
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

	tabs := container.NewAppTabs(
		container.NewTabItem("Stations", container.New(layout.NewVBoxLayout(), queryText, table)),
		container.NewTabItem("Local Query", widget.NewLabel("World!")),
		container.NewTabItem("Global Query", widget.NewLabel("World!")),
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationTop)

	(*window).SetContent(tabs)
}

func get_stations(conn *pgx.Conn) ([]Station, error) {
	rows, err := conn.Query(context.Background(), "select * from station_all")
	if err != nil {
		fmt.Println(err)
	}

	stations, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Station])
	fmt.Println(stations)
	return stations, err
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
