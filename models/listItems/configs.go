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
	MoreItemsLenInfo bool
	RenameGroupsView map[string]string
	GroupsView       bool
	ConfigsViewTheme ConfigsViewTheme
	ActionKeys       base.ActionKeys
}

func (configs *Configs) check() (warnings error) {
	if configs.CursorSymbol == "" {
		configs.CursorSymbol = ">"
	}
	if configs.MaxPageItems == 0 {
		configs.MaxPageItems = 20
	}
	return warnings
}
func WithBaseKeys(cfg *Configs) {
	baseKeys := base.GetBaseKeys()
	cfg.ActionKeys = baseKeys
}

// Configs view theme

type ConfigsViewTheme interface {
	Title(path []string, wp base.WindowParams) string
	ItemView(lism *ListItemModel, active bool, wp base.WindowParams) string
	Footer(page, allPages uint64, itemsGroups map[string]uint64, findingValue *string, wp base.WindowParams) string
}

// Base configs view handler

func WithBaseConfigsView(cfg *Configs) {
	baseCfgView := new(ConfigsViewHandler)
	cfg.ConfigsViewTheme = baseCfgView
}

type ConfigsViewHandler struct{}

func (cvh *ConfigsViewHandler) Title(path []string, wp base.WindowParams) string {
	var s string
	for i, name := range path {
		if i != 0 {
			s += " > "
		}
		s += name
	}
	s += "\n\n"
	return s
}

func (cvh *ConfigsViewHandler) ItemView(
	lim *ListItemModel, active bool,
	wp base.WindowParams,
) string {
	var s string
	if active {
		s += "> "
	} else {
		s += "  "
	}
	if lim.GetSelected() {
		s += "[*] "
	} else {
		s += "[ ] "
	}
	s += lim.GetName()
	if lim.GetGroup() != "" {
		s += fmt.Sprintf(" (%v)", lim.GetGroup())
	}
	s += "\n"
	return s
}

func (cvh *ConfigsViewHandler) Footer(
	page, allPages uint64,
	itemsGroups map[string]uint64,
	findingValue *string,
	wp base.WindowParams,
) string {
	var s string
	s += fmt.Sprintf("Page %v/%v.", page, allPages)
	for groupName, groupItems := range itemsGroups {
		s += fmt.Sprintf(" %v-%v", groupItems, groupName)
	}
	s += "/n/n"
	if findingValue != nil {
		s += fmt.Sprintf("Find: %v", findingValue)
	}
	return s
}
