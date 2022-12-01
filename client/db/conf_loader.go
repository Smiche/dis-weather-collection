package db

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/magiconair/properties"
)

// Definition of the config properties file type.
type Config struct {
	Host         string    `properties:"host"`
	Port         int       `properties:"port,default=9000"`
	Username     string    `properties:"username"`
	Password     string    `properties:"password"`
	Database     string    `properties:"database"`
	DeviceId     int       `properties:"device_id"`
	PhenomenonId int       `properties:"phenomenon_id"`
	UnitId       int       `properties:"unit_id"`
	Filename     string    `properties:",default=filename"`
	StartDate    time.Time `properties:"start_date,layout=2006-01-02"`
	Period       int       `properties:"period_ms"`
}

func Get_config_view(window *fyne.Window, onConnect func(conf Config)) {
	// Load available .conf files in the starting directory.
	configs, confNames := Get_available_conf()

	// Making the UI layout and elements.
	text1 := canvas.NewText("Select a file from the available database connection configurations.", color.Black)
	content := container.New(layout.NewHBoxLayout(), text1)

	var c color.RGBA
	c.A = 0xFF
	c.R = 0xFA
	errText := canvas.NewText("Something went wrong when loading config.", c)
	errText.Hide()

	selectedValue := ""
	selectButton := widget.NewButton("Connect", func() {
		// When the button is clicked, find out which config was selected and call the supplied onConnect func with the config.
		idx := sort.Search(len(configs), func(i int) bool {
			fmt.Println(selectedValue, configs[i].Filename)
			return configs[i].Filename == selectedValue
		})
		if idx >= len(configs) {
			errText.Show()
		} else {
			onConnect(configs[idx])
		}
	})

	selectButton.Disable()

	dropdown := widget.NewSelect(confNames, func(value string) {
		if selectButton.Disabled() {
			selectButton.Enable()
		}
		selectedValue = value
	})

	// Populate the window with content. Vertical layout box and all the other elements.
	(*window).SetContent(container.New(layout.NewVBoxLayout(), content, dropdown, errText, selectButton))
}

// Gets a list of the available .conf files in the current directory. Also parses them.
func Get_available_conf() ([]Config, []string) {
	// get current working dir
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// read all files
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// get all .conf files
	filePaths := []string{}
	fileNames := []string{}
	for _, file := range files {
		if matched, err := filepath.Match("*.conf", file.Name()); matched && err == nil {
			filePaths = append(filePaths, filepath.Join(path, file.Name()))
			fileNames = append(fileNames, file.Name())
		}
	}

	configs := []Config{}

	// populate utility array that is used for the dropdown
	for _, confPath := range filePaths {
		configs = append(configs, Load_conf(confPath))
	}

	return configs, fileNames
}

// Parses a conf for given path.
func Load_conf(path string) Config {
	p := properties.MustLoadFile(path, properties.UTF8)

	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	cfg.Filename = filepath.Base(path)
	return cfg
}
