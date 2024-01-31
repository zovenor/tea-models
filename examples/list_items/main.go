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
	model, _ := models.NewListItemsModel(
		"Fields", false, false, "Base model", nil, 20, true,
	)
	for _, item := range items {
		model.AddItem(item, nil)
	}
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatalf("failed to run program: %v", err)
		os.Exit(1)
	}
}
