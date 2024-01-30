package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/tea-models/models"
	"log"
	"os"
)

func main() {
	items := []string{"Name", "Email", "Address"}
	model := models.NewListItemsModel("Fields", false, nil)
	for _, item := range items {
		model.AddItem(item, nil)
	}
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatalf("failed to run program: %v", err)
		os.Exit(1)
	}
}
