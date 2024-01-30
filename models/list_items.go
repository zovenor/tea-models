package models

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/zovenor/tea-models/models/base"
)

type ItemMsg *ListItemModel

type ListItemModel struct {
	name  string
	value interface{}

	selected bool
}

func NewListItemModel(name string, value interface{}) *ListItemModel {
	return &ListItemModel{
		name:  name,
		value: value,
	}
}

func NewListItemsModel(name string, selectedMode bool, parent tea.Model) *ListItemsModel {
	return &ListItemsModel{
		name:         name,
		selectedMode: selectedMode,
		parent:       parent,

		cursorSymbol: base.CURSOR_SYMBOL,
	}
}

type ListItemsModel struct {
	items        []*ListItemModel
	selectedMode bool
	name         string

	cursor       int
	cursorSymbol string
	err          error
	parent       tea.Model
}

func (lim *ListItemsModel) AddItem(name string, value interface{}) {
	lim.items = append(lim.items, NewListItemModel(name, value))
}

func (lim *ListItemsModel) SetCursorSymbol(cursorSymbol string) {
	lim.cursorSymbol = cursorSymbol
}
func (lim *ListItemsModel) Init() tea.Cmd {
	return nil
}

func (lim *ListItemsModel) View() string {
	var s string

	s += fmt.Sprintf("%v\n\n", lim.name)

	for i, item := range lim.items {

		if lim.cursor == i {
			s += lim.cursorSymbol + " "
		} else {
			s += base.RepeatSymbol(" ", len(lim.cursorSymbol)+1)
		}

		s += fmt.Sprintf("%v) %v\n", i+1, item.name)
	}

	if lim.err != nil {
		s += "\n"
		s += color.New(color.BgRed).Sprintf("error: %v", lim.err.Error())
		s += "\n"
	}

	return s
}

func (lim *ListItemsModel) getItemByIndex(index int) (*ListItemModel, error) {
	if index >= len(lim.items) {
		return nil, fmt.Errorf("index out of range")
	} else {
		return lim.items[index], nil
	}
}

func (lim *ListItemsModel) getParent() (tea.Model, error) {
	if lim.parent == nil {
		return nil, fmt.Errorf("parent is nil")
	}
	return lim.parent, nil
}
func (lim *ListItemsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	lim.err = nil

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case base.UP_KEY:
			if lim.cursor > 0 {
				lim.cursor--
			}
		case base.DOWN_KEY:
			if lim.cursor < len(lim.items)-1 {
				lim.cursor++
			}
		case base.FORWARD_KEY:
			neededItem, err := lim.getItemByIndex(lim.cursor)
			if err != nil {
				lim.err = err
				return lim, nil
			}
			switch value := neededItem.value.(type) {
			case tea.Model:
				return value, nil
			default:
				lim.err = fmt.Errorf("can not forward to this item")
				return lim, nil
			}
		case base.BACK_KEY:
			parent, err := lim.getParent()
			if err != nil {
				lim.err = err
				return lim, nil
			}
			return parent.Update(nil)
		case base.ENTER_KEY:
			parent, err := lim.getParent()
			if err != nil {
				lim.err = err
				return lim, nil
			}
			neededItem, err := lim.getItemByIndex(lim.cursor)
			if err != nil {
				lim.err = err
				return lim, nil
			}
			return parent.Update(ItemMsg(neededItem))
		case base.EXIT_KEY:
			return lim, tea.Quit
		}
	}

	return lim, nil
}
