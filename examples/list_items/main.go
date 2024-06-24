package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/zovenor/logging/v2/prettyPrints"
	"github.com/zovenor/tea-models/models/listItems"
)

func updateF(lism *listItems.ListItemsModel, msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return lism, tea.Quit
		}
	}

	return nil, nil
}

func main() {
	prettyPrints.ClearTerminal()
	uf := updateF
	cfg := listItems.Configs{
		Name:             "Test app",
		SelectMode:       true,
		MaxPageItems:     20,
		FindMode:         true,
		ParentPath:       []string{"Main"},
		Parent:           nil,
		ShowIndexes:      true,
		DeletedMode:      true,
		UpdateFunc:       &uf,
		MoreItemsLenInfo: true,
	}
	lism, err := listItems.NewListItemsModel(&cfg, listItems.WithBaseKeys)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		lim := listItems.NewListItemModel()
		lim.SetName(fmt.Sprintf("Item %v", i+1))
		lim.SetValue(i)
		lism.AddItem(lim)
	}
	for i := 10; i < 20; i++ {
		lim := listItems.NewListItemModel()
		lim.SetName(fmt.Sprintf("Item %v", i+1))
		lim.SetValue(i)
		lim.SetGroup("needless")
		lism.AddItem(lim)
	}
	for i := 20; i < 30; i++ {
		lim := listItems.NewListItemModel()
		lim.SetName(fmt.Sprintf("Item %v", i+1))
		lim.SetValue(i)
		lim.SetGroup("ok")
		lism.AddItem(lim)
	}
	for i := 30; i < 41; i++ {
		lim := listItems.NewListItemModel()
		lim.SetName(fmt.Sprintf("Item %v", i+1))
		lim.SetValue(i)
		lism.AddItem(lim)
	}

	app := tea.NewProgram(lism)
	if _, err := app.Run(); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
}
