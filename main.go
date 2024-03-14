package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"sncout/database"
	"sncout/models"
	"strconv"
	"time"
)

var a fyne.App
var w fyne.Window

var teamSearch *widget.Entry
var teamList *widget.List
var newTeam *widget.Button

var team *models.Team
var tabs *container.AppTabs

var data []*models.Team

var showed []*models.Team

func main() {
	database.Init()

	a = app.New()
	w = a.NewWindow("Sncout")

	teamSearch = widget.NewEntry()
	teamSearch.OnChanged = func(search string) {
		if search == "" {
			showed = data
			return
		}
		showed = nil
		num, err := strconv.Atoi(search)
		if err != nil {
			num = -1
		}
		for _, t := range data {
			if fuzzy.Match(t.Name, search) {
				showed = append(showed, t)
				fmt.Println("matched " + t.Name + " with " + search)
				continue
			}

			for _, m := range t.Matches {
				if fuzzy.Match(string(m.MatchType), search) || num == int(m.MatchNumber) {
					showed = append(showed, t)
					continue
				}
			}
		}
	}

	teamList = widget.NewList(
		func() int {
			return len(showed)
		},
		func() fyne.CanvasObject {
			return widget.NewButton("", nil)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Button).SetText(showed[i].Name)
			o.(*widget.Button).OnTapped = func() {
				team = showed[i]
			}
		})

	newTeam = widget.NewButtonWithIcon("Add", nil, createNewTeam)

	tabs = container.NewAppTabs(
		container.NewTabItemWithIcon("Teams", nil, container.New(layout.NewGridLayout(1), teamSearch, teamList, newTeam)))

	go refresh()
	go refreshData()

	w.SetContent(tabs)
	w.ShowAndRun()
}

func refresh() {
	teamList.Refresh()
	if team != nil {
		tabs = container.NewAppTabs(
			container.NewTabItemWithIcon("Teams", nil, container.New(layout.NewGridLayout(1), teamSearch, teamList, newTeam)),
			container.NewTabItemWithIcon("Team Information", nil, container.New(layout.NewGridLayout(1), widget.NewLabel(team.Name))),
			container.NewTabItemWithIcon("Matches", nil, widget.NewLabel("Matches")),
			container.NewTabItemWithIcon("Robot Information", nil, widget.NewLabel("Robot Information")))
	} else {
		tabs = container.NewAppTabs(
			container.NewTabItemWithIcon("Teams", nil, container.New(layout.NewGridLayout(1), teamSearch, teamList, newTeam)))
	}
}

func refreshData() {
	database.GetDB().Preload("Matches").Preload("Robot").Find(&data)
	time.Sleep(5 * time.Second)
}

func showTeam(team *models.Team) {
	tw := a.NewWindow("Team" + team.Name)
	tw.SetContent(widget.NewLabel(team.Name))
	tw.Show()
}

func createNewTeam() {
	tw := a.NewWindow("New Team")
	name := widget.NewEntry()

	form := widget.NewForm(&widget.FormItem{Text: "name", Widget: name})

	form.OnSubmit = func() {
		team = &models.Team{
			Name: name.Text,
		}
		database.GetDB().Save(team)
		refresh()
		tw.Hide()
		tw.Close()
	}

	tw.SetContent(form)
	tw.Show()
}
