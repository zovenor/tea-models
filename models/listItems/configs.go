package listItems

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zovenor/tea-models/models/base"
)

const BaseMaxPageItems uint16 = 20

type Configs struct {
	Name             string
	SelectMode       bool
	FindMode         bool
	ParentPath       string
	Parent           tea.Model
	MaxPageItems     uint16
	ShowIndexes      bool
	MapArgs          map[string]interface{}
	UpdateFunc       *func(*ListItemsModel, tea.Msg) (tea.Model, tea.Cmd)
	CursorSymbol     string
	DeletedMode      bool
	Keys             []base.Key
	MoreItemsLenInfo bool
	RenameGroupsView map[string]string
	GroupsView       bool
}

func (configs *Configs) check() (warnings []error) {
	warnings = make([]error, 0)
	if configs.MaxPageItems == 0 {
		warnings = append(warnings, fmt.Errorf("max page items equal zero, replaced to base value -> 20"))
	}
	if configs.CursorSymbol == "" {
		configs.CursorSymbol = ">"
	}
	return warnings
}
