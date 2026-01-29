package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/elastic/go-elasticsearch/v9"
)

type Configuration struct {
	esClient    *elasticsearch.Client
	searchIndex string

	app    fyne.App
	window fyne.Window

	input   *widget.Entry
	results *widget.Label
}

// buildUI creates the main user interface for the application.
func (uiConfig *Configuration) buildUI() {
	title := widget.NewLabel("Hybrid Search Demo")
	title.TextStyle.Bold = true
	title.Alignment = fyne.TextAlignCenter

	uiConfig.input = widget.NewEntry()
	uiConfig.input.SetPlaceHolder("Enter search query...")

	uiConfig.results = widget.NewLabel("")
	uiConfig.results.Wrapping = fyne.TextWrapWord

	scrollView := container.NewScroll(uiConfig.results)

	searchButton := widget.NewButton("Search", uiConfig.onSearch)

	searchWidget := container.NewVBox(uiConfig.input, searchButton)

	uiConfig.window.SetContent(container.NewBorder(
		container.NewVBox(title, searchWidget), nil, nil, nil,
		scrollView,
	))
}

// onSearch is the callback function for the search button.
func (uiConfig *Configuration) onSearch() {
	if uiConfig.input.Text == "" {
		uiConfig.results.SetText("Search terms cannot be empty. Try again!")
		return
	}
	uiConfig.results.SetText("Retrieving results...")

	queryResult, err := runQuery(uiConfig.esClient, uiConfig.searchIndex, strings.TrimSpace(uiConfig.input.Text))
	if err != nil {
		uiConfig.results.SetText(fmt.Sprintf("Error running query: %s", err))
		return
	}

	var resultData Result
	if err := json.Unmarshal([]byte(queryResult), &resultData); err != nil {
		uiConfig.results.SetText(fmt.Sprintf("Error parsing results: %s", err))
		return
	}

	uiConfig.results.SetText(showUIResults(resultData))
}
