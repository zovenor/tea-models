package listItems

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/zovenor/tea-models/models/base"
)

const BaseMaxPageItems uint16 = 20

type Configs struct {
	Name             string
	SelectMode       bool
	FindMode         bool
	ParentPath       []string
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
	Footer(page, allPages uint64, itemsGroups []ItemsGroup, findingValue *string, wp base.WindowParams) string
}

// Base configs view handler

func WithBaseConfigsView(cfg *Configs) {
	baseCfgView := new(ConfigsViewHandler)
	cfg.ConfigsViewTheme = baseCfgView
}

type ConfigsViewHandler struct{}

func (cvh *ConfigsViewHandler) Title(path []string, wp base.WindowParams) string {
	var s string
	s += strings.Join(path, " ► ")
	s += "\n\n"
	return s
}

func (cvh *ConfigsViewHandler) ItemView(
	lim *ListItemModel, active bool,
	wp base.WindowParams,
) string {
	var s string
	if active {
		s += "► "
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

type ItemsGroup struct {
	Name  string
	Total uint64
}

func (cvh *ConfigsViewHandler) Footer(
	page, allPages uint64,
	itemsGroups []ItemsGroup,
	findingValue *string,
	wp base.WindowParams,
) string {
	var s string = "\n"
	s += fmt.Sprintf("Page %v/%v.", page, allPages)
	for _, groupItems := range itemsGroups {
		s += fmt.Sprintf(" %v-%v", groupItems.Total, groupItems.Name)
	}
	s += "\n\n"
	if findingValue != nil {
		s += fmt.Sprintf("Find: %v", *findingValue)
	}
	return s
}
