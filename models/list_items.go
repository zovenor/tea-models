package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fatih/color"
	"github.com/zovenor/tea-models/models/base"
	"strings"
)

type ItemMsg *ListItemModel
type ItemsMsg []*ListItemModel

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

func NewListItemsModel(
	name string,
	selectMode bool,
	returnValue bool,
	parentPath string,
	parent tea.Model,
	maxItemsInPage int,
) (*ListItemsModel, error) {
	if maxItemsInPage < 1 {
		return nil, fmt.Errorf("maxItemsInPage should be more than 0")
	}
	return &ListItemsModel{
		name:        name,
		selectMode:  selectMode,
		returnValue: returnValue,
		parentPath:  parentPath,

		parent: parent,

		cursorSymbol:   base.CursorSymbol,
		maxItemsInPage: maxItemsInPage,
	}, nil
}

type ListItemsModel struct {
	items       []*ListItemModel
	selectMode  bool
	returnValue bool
	parentPath  string
	name        string

	cursor               int
	cursorSymbol         string
	err                  error
	parent               tea.Model
	viewListItemsIndexed []int
	findValue            string
	findModel            *textinput.Model
	findCursor           int
	maxItemsInPage       int
}

func (lim *ListItemsModel) AddItem(name string, value interface{}) {
	lim.viewListItemsIndexed = append(lim.viewListItemsIndexed, len(lim.items))
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

	if lim.parentPath != "" {
		s += fmt.Sprintf("%v > ", lim.parentPath)
	}

	s += fmt.Sprintf("%v\n\n", lim.name)

	for _, index := range lim.getPageItemsIndexes() {

		item, err := lim.getItemByIndex(index)
		if err != nil {
			continue
		}

		if lim.cursor == index {
			s += lim.cursorSymbol + " "
		} else {
			s += base.RepeatSymbol(" ", len(lim.cursorSymbol)+1)
		}

		if lim.selectMode {
			if item.selected {
				s += "[*] "
			} else {
				s += "[ ] "
			}
		}

		s += fmt.Sprintf("%v) %v\n", index+1, item.name)
	}

	if lim.findModel != nil {
		s += fmt.Sprintf("\n%v\n", lim.findModel.View())
	}

	if lim.err != nil {
		s += "\n"
		s += color.New(color.BgRed).Sprintf("error: %v", lim.err.Error())
		s += "\n"
		lim.err = nil
	}

	if lim.findModel == nil {
		s += base.GetHints(base.ExitKey, base.FindKey, base.SelectKey, base.EnterKey)
	} else {
		s += base.GetHints(base.ExitKey, base.EnterKey, base.CancelKey)
	}
	if lim.getPagesLen() > 0 {
		s += fmt.Sprintf("Page %v/%v", lim.getPageIndex()+1, lim.getPagesLen())
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

func (lim *ListItemsModel) filterByName(s string) {
	var neededIndexes []int
	for i, item := range lim.items {
		if strings.Contains(strings.ToLower(item.name), strings.ToLower(s)) {
			neededIndexes = append(neededIndexes, i)
		}
	}
	lim.viewListItemsIndexed = neededIndexes
}

func (lim *ListItemsModel) setCursorByFindCursor() {
	if len(lim.viewListItemsIndexed) == 0 {
		lim.findCursor = 0
	}
	if lim.findCursor >= len(lim.viewListItemsIndexed) {
		lim.findCursor = len(lim.viewListItemsIndexed) - 1
	}
	if lim.viewListItemsIndexed[lim.findCursor] < len(lim.items) {
		lim.cursor = lim.viewListItemsIndexed[lim.findCursor]
	} else {
		lim.setCursorByFindCursor()
	}
}

func (lim *ListItemsModel) nextIndex() {
	nextFindCursor := lim.findCursor + 1
	if nextFindCursor < len(lim.viewListItemsIndexed) {
		lim.findCursor = nextFindCursor
		lim.setCursorByFindCursor()
	} else {
		lim.setCursorByFindCursor()
	}
}

func (lim *ListItemsModel) getPagesLen() int {
	if len(lim.viewListItemsIndexed) < lim.maxItemsInPage {
		return 1
	}
	return len(lim.viewListItemsIndexed)/lim.maxItemsInPage + 1
}

func (lim *ListItemsModel) getPageIndex() int {
	if len(lim.viewListItemsIndexed) == 0 {
		return 0
	}
	return lim.findCursor / lim.maxItemsInPage
}

func (lim *ListItemsModel) getPageItemsIndexes() []int {
	if len(lim.viewListItemsIndexed) < lim.maxItemsInPage {
		return lim.viewListItemsIndexed
	} else if lim.getPageIndex()+1 < lim.getPagesLen() {
		return lim.viewListItemsIndexed[lim.getPageIndex()*lim.maxItemsInPage : (lim.getPageIndex()+1)*lim.maxItemsInPage]
	} else {
		return lim.viewListItemsIndexed[lim.getPageIndex()*lim.maxItemsInPage:]
	}
}

func (lim *ListItemsModel) lastIndex() {
	lastFindCursor := lim.findCursor - 1
	if lastFindCursor >= 0 {
		lim.findCursor = lastFindCursor
	}
	lim.setCursorByFindCursor()
}

func (lim *ListItemsModel) getSelectedItemsMsg() ItemsMsg {
	var items ItemsMsg
	for _, item := range lim.items {
		if item.selected {
			items = append(items, item)
		}
	}
	return items
}

func (lim *ListItemsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	if lim.findModel == nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case base.UpKey:
				lim.lastIndex()
			case base.DownKey:
				lim.nextIndex()
			case base.ForwardKey:
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
			case base.BackKey:
				parent, err := lim.getParent()
				if err != nil {
					lim.err = err
					return lim, nil
				}
				return parent.Update(nil)
			case base.EnterKey:
				if lim.returnValue {
					parent, err := lim.getParent()
					if err != nil {
						lim.err = err
						return lim, nil
					}
					if !lim.selectMode {
						neededItem, err := lim.getItemByIndex(lim.cursor)
						if err != nil {
							lim.err = err
							return lim, nil
						}
						return parent.Update(ItemMsg(neededItem))
					} else {
						neededItems := lim.getSelectedItemsMsg()
						return parent.Update(neededItems)
					}
				}
			case base.SelectKey:
				if lim.selectMode {
					neededItem, err := lim.getItemByIndex(lim.cursor)
					if err != nil {
						lim.err = err
						return lim, nil
					}
					neededItem.selected = !neededItem.selected
				}
			case base.FindKey:
				ti := textinput.New()
				ti.Placeholder = lim.findValue
				ti.SetValue(lim.findValue)
				ti.Focus()
				lim.findModel = &ti
			case base.ExitKey:
				return lim, tea.Quit
			}
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case base.EnterKey:
				lim.findValue = lim.findModel.Value()
				lim.filterByName(lim.findValue)
				lim.findModel = nil
				lim.setCursorByFindCursor()
			case base.CancelKey:
				lim.findModel = nil
			case base.ExitKey:
				return lim, tea.Quit
			default:
				var cmd tea.Cmd
				*lim.findModel, cmd = lim.findModel.Update(msg)
				return lim, cmd
			}
		}
	}

	return lim, nil
}
