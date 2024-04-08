package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/tea-models/models"
	"log"
	"os"
)

func main() {
	items := []string{
		"Name", "Email", "Address", "Name", "Email", "Address", "Name", "Email", "Address", "Name", "Email", "Address",
		"Name", "Email", "Address", "Name", "Email", "Address", "Name", "Email", "Address", "Name", "Email", "Address",
	}
	listItemsConf := models.ListItemsConf{
		Name:           "Fields",
		SelectMode:     false,
		ReturnValue:    false,
		ParentPath:     "Base model",
		Parent:         nil,
		MaxItemsInPage: 20,
		Indexes:        true,
	}
	model, _ := models.NewListItemsModel(listItemsConf)
	for _, item := range items {
		model.AddItem(item, nil)
	}
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatalf("failed to run program: %v", err)
		os.Exit(1)
	}
}
