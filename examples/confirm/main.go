package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/tea-models/models"
	"log"
	"os"
)

func main() {
	listModel, _ := models.NewListItemsModel(
		"Fields", false, false, "Base model", nil, 20, false,
	)
	confirmModel, _ := models.NewConfirmModel(
		"Do you want to do something?\n",
		listModel, func() {
			listModel.AddItem("Test field", "Value")
		})
	listModel.AddItem("Confirm field", confirmModel)
	p := tea.NewProgram(listModel)
	if _, err := p.Run(); err != nil {
		log.Fatalf("failed to run program: %v", err)
		os.Exit(1)
	}
}
