package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/tea-models/models/confirm"
	"github.com/zovenor/tea-models/models/listItems"
)

func main() {
	listModel, _ := listItems.NewListItemsModel(&listItems.Configs{
		Name:       "Fields",
		SelectMode: false,
		ParentPath: "Base model",
		Parent:     nil,
	})
	confirmModel, _ := confirm.NewConfirmModel(
		"Do you want to do something?\n",
		listModel, func() {
			lim := listItems.NewListItemModel()
			lim.SetName("Test field")
			lim.SetValue("Value")
			listModel.AddItem(lim)
		})
	lim := listItems.NewListItemModel()
	lim.SetName("Confirm field")
	lim.SetValue(confirmModel)
	listModel.AddItem(lim)
	p := tea.NewProgram(listModel)
	if _, err := p.Run(); err != nil {
		log.Fatalf("failed to run program: %v", err)
		os.Exit(1)
	}
}
